package jiro

import (
	"fmt"
	"sudonters/zootler/internal/components"
	"sudonters/zootler/internal/entities"
	"sudonters/zootler/internal/entity"
	"sudonters/zootler/internal/table"

	"github.com/etc-sudonters/substrate/skelly/graph"
	"github.com/etc-sudonters/substrate/slipup"
)

type worldnode struct {
	Events    map[string]string `json:"events"`
	Exits     map[string]string `json:"exits"`
	Locations map[string]string `json:"locations"`

	RegionName string `json:"region_name"`

	AltHint string `json:"alt_hint"`
	Hint    string `json:"hint"`

	Dungeon    string `json:"dungeon"`
	IsBossRoom bool   `json:"is_boss_room"`

	Savewarp   string `json:"savewarp"`
	Scene      string `json:"scene"`
	TimePasses bool   `json:"time_passes"`
}

type nodeedge struct {
	name components.Name
	rule components.RawLogic
	kind edgekind
}

func (wn worldnode) Edges(yield func(components.Name, nodeedge) bool) {
	var ne nodeedge
	every := func(kind edgekind, edges map[string]string) bool {
		ne.kind = kind
		for dest, rule := range edges {
			ne.name = components.Name(fmt.Sprintf("%s -> %s", wn.RegionName, dest))
			ne.rule = components.RawLogic(rule)
			if !yield(components.Name(dest), ne) {
				return false
			}
		}
		return true
	}

	pairs := []struct {
		k edgekind
		m map[string]string
	}{
		{event, wn.Events},
		{exity, wn.Exits},
		{check, wn.Locations},
	}

	for _, pair := range pairs {
		if !every(pair.k, pair.m) {
			return
		}
	}
}

func (wn worldnode) EntityName() components.Name {
	return components.Name(wn.RegionName)
}

func (wn worldnode) AsComponents() table.Values {
	var vt table.Values
	assign := func(v table.Value) {
		vt = append(vt, v)
	}
	if wn.Hint != "" {
		assign(components.HintRegion{
			Name: wn.Hint,
			Alt:  wn.AltHint,
		})
	}

	if wn.Dungeon != "" {
		assign(components.Dungeon(wn.Dungeon))
	}

	if wn.IsBossRoom {
		assign(components.BossRoom{})
	}

	if wn.Scene != "" {
		assign(components.Scene(wn.Scene))
	}

	if wn.Savewarp != "" {
		assign(components.SavewarpName(wn.Savewarp))
	}

	if wn.TimePasses {
		assign(components.TimePasses{})
	}

	return vt
}

type loadstate struct {
	locs entities.Map[entities.Location]
	edge entities.Map[entities.Edge]
	grph graph.Builder
}

type edgekind uint8

const (
	check edgekind = iota
	event
	exity
)

func (e edgekind) component() table.Value {
	switch e {
	case check:
		return components.CheckEdge{}
	case event:
		return components.EventEdge{}
	case exity:
		return components.ExitEdge{}
	default:
		panic("unknown edge kind")
	}
}

func (l *loadstate) load(wn worldnode) error {
	origin, originErr := l.locs.Entity(components.Name(wn.RegionName))
	paniconerr(originErr)
	if err := origin.AddComponents(wn.AsComponents()); err != nil {
		paniconerr(err)
	}

	for destName, nodeEdge := range wn.Edges {
		dest, destErr := l.locs.Entity(destName)
		paniconerr(destErr)
		_, edgeErr := l.connect(nodeEdge.name, origin, dest, nodeEdge.kind, nodeEdge.rule)
		paniconerr(edgeErr)
	}

	return nil
}

func (l *loadstate) connect(name components.Name, origin, dest entities.Location, kind edgekind, rule components.RawLogic) (entities.Edge, error) {
	edge, edgeErr := l.edge.Entity(name)
	if edgeErr != nil {
		return entities.Edge{}, slipup.Describef(edgeErr, "edge %s", name)
	}

	comps := table.Values{
		rule, kind.component(),
		components.Connection{Origin: entity.Model(origin.Id()), Dest: entity.Model(dest.Id())},
	}

	if err := edge.AddComponents(comps); err != nil {
		return edge, slipup.Describef(err, "while adding components to '%s'", name)
	}

	return edge, nil
}

func paniconerr(e error) {
	if e != nil {
		panic(e)
	}
}
