package jiro

import (
	"fmt"
	"sudonters/zootler/internal/components"
	"sudonters/zootler/internal/table"
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
