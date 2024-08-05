package main

import (
	"fmt"
	"sudonters/zootler/internal/bundle"
	"sudonters/zootler/internal/entity"
	"sudonters/zootler/internal/query"
	"sudonters/zootler/internal/table"
	"sudonters/zootler/pkg/world/components"

	"github.com/etc-sudonters/substrate/mirrors"
	"github.com/etc-sudonters/substrate/skelly/graph"
)

type WorldGraphFileLocation struct {
	AltHint      string            `json:"alt_hint"` // more specific location -- restricted to gannon's chambers?
	BossRoom     bool              `json:"is_boss_room"`
	DoesTimePass bool              `json"time_passes"`
	Dungeon      string            `json:"dungeon"`
	Hint         string            `json:"hint"` // hint zone
	Region       string            `json:"region_name"`
	Savewarp     string            `json:"savewarp"` // effectively another exit
	Scene        string            `json:"scene"`    // similar to hint?
	Events       map[string]string `json:"events"`
	Exits        map[string]string `json:"exits"`
	Locations    map[string]string `json:"locations"`
}

type WorldGraphLoader struct {
	Helpers, Path string
	IncludeMQ     bool
}

func (w WorldGraphLoader) Load() error {
	return nil
}

type EdgeRule struct {
	Origin   graph.Origination
	Dest     graph.Destination
	Bytecode []uint8
}

func EdgeRuleIndexer(e EdgeRule) (string, bool) {
	return fmt.Sprintf("%d-%d", e.Origin, e.Dest), true
}

func CreateLocationIds(e query.Engine) (*LocationIds, error) {
	query := e.CreateQuery()
	query.Exists(mirrors.TypeOf[components.Location]())
	query.Load(mirrors.TypeOf[components.Name]())
	rows, qErr := e.Retrieve(query)
	if qErr != nil {
		return nil, qErr
	}

	locations, loadErr := bundle.ToMap(rows, func(r *table.RowTuple) (string, entity.Model, error) {
		model, name := r.Id, r.Values[1].(components.Name)
		return normalize(name), entity.Model(model), nil
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
	locations map[string]entity.Model
}

func (l LocationIds) Mint(name string) (entity.Model, error) {
	normalized := normalize(name)
	if id, ok := l.locations[normalized]; ok {
		return id, nil
	}

	row, insertErr := l.e.InsertRow(components.Name(name))
	if insertErr != nil {
		return entity.INVALID_ENTITY, insertErr
	}

	return entity.Model(row), nil
}
