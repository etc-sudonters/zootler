package jiro

import (
	"sudonters/zootler/internal/components"
	"sudonters/zootler/internal/table"
)

type nodeedge struct {
	name components.Name
	rule components.RawLogic
	kind edgekind
}

type edgekind uint8

const (
	check edgekind = iota
	event
	exity
)

func (e edgekind) ascomponents() table.Values {
	var vt table.Values
	switch e {
	case check:
		vt = append(vt, components.CheckEdge{})
	case event:
		vt = append(vt, components.EventEdge{}, components.CollectableGameToken{})
	case exity:
		vt = append(vt, components.ExitEdge{})
	default:
		panic("unknown edge kind")
	}

	return vt
}
