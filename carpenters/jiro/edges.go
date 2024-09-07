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
