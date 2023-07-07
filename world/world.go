package world

import (
	"github.com/etc-sudonters/rando/entity"
	"github.com/etc-sudonters/rando/graph"
)

type Id int

type World struct {
	Id       Id
	Entities Pool
	Graph    graph.Builder
}

func (w World) AddNode(name entity.TagName) *entity.View {
	return w.Entities.AddNode(name, w.Graph.AddNode())
}

func (w World) AddEdge(name entity.TagName, o graph.Origination, d graph.Destination) (*entity.View, error) {
	ent := w.Entities.AddEdge(name, o, d)
	err := w.Graph.AddEdge(o, d)

	return ent, err
}
