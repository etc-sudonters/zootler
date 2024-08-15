package main

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"
	"sudonters/zootler/internal/bundle"
	"sudonters/zootler/internal/components"
	"sudonters/zootler/internal/entity"
	"sudonters/zootler/internal/query"
	"sudonters/zootler/internal/rules/parser"
	"sudonters/zootler/internal/rules/runtime"
	"sudonters/zootler/internal/slipup"
	"sudonters/zootler/internal/table"
)

type WorldFileLoader struct {
	Path, Helpers string
	IncludeMQ     bool
	Compiler      runtime.Compiler
}

func (w WorldFileLoader) Setup(ctx context.Context, e query.Engine) error {
	WriteLineOut(ctx, "reading dir  '%s'", w.Path)
	directory, dirErr := os.ReadDir(w.Path)
	if dirErr != nil {
		return slipup.Describe(dirErr, w.Path)
	}

	locTbl, locErr := CreateLocationMap(ctx, e)
	if locErr != nil {
		return slipup.Describe(locErr, "creating location id table")
	}

	for _, entry := range directory {
		if !IsFile(entry) {
			continue
		}

		path := path.Join(w.Path, entry.Name())
		if err := w.logicFile(ctx, path, locTbl); err != nil {
			return slipup.Describe(err, path)
		}
	}

	return nil
}

func (w WorldFileLoader) helpers() error {
	helpers := make(map[string]string, 256)
	for decl, body := range helpers {
		f, err := parser.ParseFunctionDecl(decl, body)
		if err != nil {
			return slipup.Describef(err, "while parsing function decl '%s'", decl)
		}

		if compileErr := w.Compiler.CompileFunctionDecl(f); compileErr != nil {
			return slipup.Describef(err, "while compiling function decl '%s'", decl)
		}
	}
	return nil
}

func (w WorldFileLoader) logicFile(
	_ context.Context, path string, locations *LocationMap,
) error {
	isMq := strings.Contains(path, "mq")

	if isMq && !w.IncludeMQ {
		return nil
	}

	rawLocs, readErr := ReadJsonFile[[]WorldFileLocation](path)
	if readErr != nil {
		return slipup.Describe(readErr, path)
	}

	for _, raw := range rawLocs {
		here, buildErr := locations.Build(components.Name(raw.Name))
		if buildErr != nil {
			return slipup.Describef(buildErr, "building %s", raw.Name)
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

		if raw.DoesTimePass {
			values = append(values, components.TimePasses{})
		}

		if raw.Savewarp != "" {
			values = append(values, components.SavewarpName(raw.Savewarp))
		}

		if isMq {
			values = append(values, components.MasterQuest{})
		}

		if err := here.add(values); err != nil {
			return slipup.Describef(err, "while populating %s", raw.Name)
		}

		for destination, rule := range raw.Exits {
			linkErr := here.linkTo(components.Name(destination), components.RawLogic{Rule: rule}, components.ExitEdge{})
			if linkErr != nil {
				return slipup.Describef(linkErr, "while linking '%s -> %s'", here.name, destination)
			}
		}

		for check, rule := range raw.Locations {
			linkErr := here.linkTo(components.Name(check), components.RawLogic{Rule: rule}, components.CheckEdge{})
			if linkErr != nil {
				return slipup.Describef(linkErr, "while linking '%s -> %s'", here.name, check)
			}
		}

		for event, rule := range raw.Events {
			linkErr := here.linkTo(components.Name(event), components.RawLogic{Rule: rule}, components.EventEdge{})
			if linkErr != nil {
				return slipup.Describef(linkErr, "while linking '%s -> %s'", here.name, event)
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

func (l locationBuilder) linkTo(name components.Name, rule components.RawLogic, vs ...table.Value) error {
	edgeName := fmt.Sprintf("%s -> %s", l.name, name)
	destination, linkErr := l.parent.Build(name)
	if linkErr != nil {
		return slipup.Describef(linkErr, "while linking '%s'", edgeName)
	}

	edge, edgeCreateErr := l.parent.e.InsertRow(components.Name(edgeName), rule, components.Edge{
		Origin: entity.Model(l.id),
		Dest:   entity.Model(destination.id),
	})
	if edgeCreateErr != nil {
		return slipup.Describef(edgeCreateErr, "while creating edge %s", edgeName)
	}

	if len(vs) != 0 {
		if additionalErr := l.parent.e.SetValues(edge, table.Values(vs)); additionalErr != nil {
			return slipup.Describef(additionalErr, "while customizing '%s'", edgeName)
		}
	}
	return nil
}

func CreateLocationMap(ctx context.Context, e query.Engine) (*LocationMap, error) {
	q := e.CreateQuery()
	q.Load(query.MustAsColumnId[components.Name](e))
	q.Exists(query.MustAsColumnId[components.Location](e))
	rows, qErr := e.Retrieve(q)
	if qErr != nil {
		return nil, qErr
	}

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
