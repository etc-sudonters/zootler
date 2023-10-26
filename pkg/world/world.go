package world

import (
	"sudonters/zootler/internal/entity"

	"github.com/etc-sudonters/substrate/skelly/graph"
)

type Edge struct {
	Origin      graph.Origination
	Destination graph.Destination
}

type World struct {
	Entities  Pool
	Graph     graph.Directed
	edgeCache map[Edge]entity.Model
	nodeCache map[graph.Node]entity.Model
}

func (w World) EdgeEntity(e Edge) (entity.Model, bool) {
	m, ok := w.edgeCache[e]
	return m, ok
}
