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
	dirEntries, dirErr := os.ReadDir(w.Path)
	if dirErr != nil {
		return slipup.Trace(dirErr, w.Path)
	}

	locTbl, locErr := CreateLocationIds(ctx, e)
	if locErr != nil {
		return slipup.Trace(locErr, "creating location id table")
	}

	for _, dentry := range dirEntries {
		if !IsFile(dentry) {
			continue
		}

		path := path.Join(w.Path, dentry.Name())
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
	ctx context.Context, path string, locations *LocationIds, e query.Engine,
) error {
	isMq := strings.Contains(path, "mq")

	if isMq && !w.IncludeMQ {
		return nil
	}

	rawLocs, readErr := ReadJsonFile[[]WorldFileLocation](path)
	if readErr != nil {
		return slipup.Trace(readErr, path)
	}

	for _, raw := range rawLocs {
		here, mintErr := locations.Mint(raw.Name)
		if mintErr != nil {
			return slipup.TraceMsg(mintErr, "minting %s", raw.Name)
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

		if err := e.SetValues(here, values); err != nil {
			return slipup.TraceMsg(err, "while populating %s", raw.Name)
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

func CreateLocationIds(ctx context.Context, e query.Engine) (*LocationIds, error) {
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

	return &LocationIds{
		e:         e,
		locations: locations,
	}, nil
}

type LocationIds struct {
	e         query.Engine
	locations map[string]table.RowId
}

func (l LocationIds) Mint(name string) (table.RowId, error) {
	normalized := normalize(name)
	if id, ok := l.locations[normalized]; ok {
		return id, nil
	}

	row, insertErr := l.e.InsertRow(components.Name(name), components.Location{})
	if insertErr != nil {
		return table.INVALID_ROWID, insertErr
	}

	return row, nil
}
