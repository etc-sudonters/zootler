package main

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"
	"sudonters/zootler/internal/bundle"
	"sudonters/zootler/internal/query"
	"sudonters/zootler/internal/slipup"
	"sudonters/zootler/internal/table"
	"sudonters/zootler/pkg/logic"
	"sudonters/zootler/pkg/world/components"

	"github.com/etc-sudonters/substrate/dontio"
	"github.com/etc-sudonters/substrate/mirrors"
	"github.com/etc-sudonters/substrate/skelly/graph"
)

type WorldFileLoader struct {
	Path, Helpers string
	IncludeMQ     bool
}

func (w WorldFileLoader) Configure(ctx context.Context, e query.Engine) error {
	stdio, stdErr := dontio.StdFromContext(ctx)
	if stdErr != nil {
		return stdErr
	}
	std := std{stdio}
	if err := w.helpers(e); err != nil {
		return slipup.Trace(err, "loading helpers")
	}

	std.WriteLineOut("reading dir  '%s'", w.Path)
	directory, dirErr := os.ReadDir(w.Path)
	if dirErr != nil {
		return slipup.Trace(dirErr, w.Path)
	}

	locTbl, locErr := CreateLocationMap(ctx, e)
	if locErr != nil {
		return slipup.Trace(locErr, "creating location id table")
	}

	for _, entry := range directory {
		if !IsFile(entry) {
			continue
		}

		path := path.Join(w.Path, entry.Name())
		std.WriteLineOut("reading file '%s'", path)
		if err := w.logicFile(ctx, path, locTbl, e); err != nil {
			return slipup.Trace(err, path)
		}
	}

	return nil
}

func (w WorldFileLoader) helpers(e query.Engine) error {
	helpers := make(map[string]string, 256)

	for name, code := range helpers {
		e.InsertRow(components.Name(name), logic.RawRule(code), components.Helper{})
	}
	return nil
}

func (w WorldFileLoader) logicFile(
	ctx context.Context, path string, locations *LocationMap, e query.Engine,
) error {
	stdio, stdErr := dontio.StdFromContext(ctx)
	std := std{stdio}
	if stdErr != nil {
		return stdErr
	}
	isMq := strings.Contains(path, "mq")

	if isMq && !w.IncludeMQ {
		return nil
	}

	rawLocs, readErr := ReadJsonFile[[]WorldFileLocation](path)
	if readErr != nil {
		return slipup.Trace(readErr, path)
	}

	for _, raw := range rawLocs {
		here, buildErr := locations.Build(components.Name(raw.Name))
		if buildErr != nil {
			return slipup.TraceMsg(buildErr, "building %s", raw.Name)
		}

		values := table.Values{
			components.HintRegion{
				Name: raw.HintRegion,
				Alt:  raw.HintRegionAlt,
			},
		}

		if raw.BossRoom {
			values = append(values, components.BossRoom{})
		}

		if raw.Dungeon != "" {
			values = append(values, components.Dungeon(raw.Dungeon))
		}

		if isMq {
			values = append(values, components.MasterQuest{})
		}

		if err := here.add(values); err != nil {
			return slipup.TraceMsg(err, "while populating %s", raw.Name)
		}

		for destination, rule := range raw.Exits {
			std.WriteLineOut("linking '%s -> %s'", here.name, destination)
			linkErr := here.linkTo(components.Name(destination), components.RawLogic{Rule: rule})
			if linkErr != nil {
				return linkErr
			}
		}
	}

	return nil
}

type WorldFileLocation struct {
	BossRoom      bool              `json:"is_boss_room"`
	DoesTimePass  bool              `json:"time_passes"`
	Dungeon       string            `json:"dungeon"`
	HintRegion    string            `json:"hint"`     // hint zone
	HintRegionAlt string            `json:"alt_hint"` // more specific location -- restricted to gannon's chambers?
	Name          string            `json:"region_name"`
	Savewarp      string            `json:"savewarp"` // effectively another exit
	Scene         string            `json:"scene"`    // similar to hint?
	Events        map[string]string `json:"events"`
	Exits         map[string]string `json:"exits"`
	Locations     map[string]string `json:"locations"`
}

type locationBuilder struct {
	parent *LocationMap
	id     table.RowId
	name   components.Name
}

func (l locationBuilder) add(v table.Values) error {
	return l.parent.e.SetValues(l.id, v)
}

func (l locationBuilder) linkTo(name components.Name, rule components.RawLogic) error {
	destination, linkErr := l.parent.Build(name)
	if linkErr != nil {
		return slipup.TraceMsg(linkErr, "while linking '%s -> %s'", l.name, name)
	}

	edgeName := fmt.Sprintf("%s -> %s", l.name, name)
	_, edgeCreateErr := l.parent.e.InsertRow(edgeName, rule, components.Edge{
		Origin: graph.Origination(l.id),
		Dest:   graph.Destination(destination.id),
	})
	if edgeCreateErr != nil {
		return slipup.TraceMsg(edgeCreateErr, "while creating edge %s", edgeName)
	}
	return nil
}

func CreateLocationMap(ctx context.Context, e query.Engine) (*LocationMap, error) {
	stdio, stdErr := dontio.StdFromContext(ctx)
	std := std{stdio}
	if stdErr != nil {
		return nil, stdErr
	}
	query := e.CreateQuery()
	query.Load(mirrors.TypeOf[components.Name]())
	query.Exists(mirrors.TypeOf[components.Location]())
	rows, qErr := e.Retrieve(query)
	if qErr != nil {
		return nil, qErr
	}

	std.WriteLineOut("retrieved '%d' named locations", rows.Len())
	locations, loadErr := bundle.ToMap(rows, func(r *table.RowTuple) (string, table.RowId, error) {
		model := r.Id
		name, ok := r.Values[0].(components.Name)
		if !ok {
			return "", r.Id, fmt.Errorf("could not cast row %+v to 'components.Name'", r)
		}
		return normalize(name), model, nil
	})

	if loadErr != nil {
		return nil, loadErr
	}

	return &LocationMap{
		e:         e,
		locations: locations,
	}, nil
}

type LocationMap struct {
	e         query.Engine
	locations map[string]table.RowId
}

func (l *LocationMap) Build(name components.Name) (*locationBuilder, error) {
	var b locationBuilder
	b.name = name
	b.parent = l
	normalized := normalize(name)
	if id, ok := l.locations[normalized]; ok {
		b.id = id
		return &b, nil
	}

	row, insertErr := l.e.InsertRow(name, components.Location{})
	if insertErr != nil {
		return nil, insertErr
	}
	b.id = row
	return &b, nil
}
