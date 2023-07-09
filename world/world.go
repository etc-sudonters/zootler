package world

import (
	"context"
	"fmt"

	"github.com/etc-sudonters/rando/entity"
	"github.com/etc-sudonters/rando/graph"
	"github.com/etc-sudonters/rando/set"
)

type edge struct {
	origin graph.Origination
	dest   graph.Destination
}

type Id int

type World struct {
	Id        Id
	Entities  Pool
	Graph     graph.Builder
	edgeCache map[edge]entity.Model
	nodeCache map[graph.Node]entity.Model
}

func (w *World) AddNode(name entity.TagName) *entity.View {
	if w.nodeCache == nil {
		w.nodeCache = make(map[graph.Node]entity.Model)
	}

	node := w.Graph.AddNode()
	ent := w.Entities.AddNode(name, node)
	w.nodeCache[node] = ent.Id
	return ent
}

func (w *World) AddEdge(name entity.TagName, o graph.Origination, d graph.Destination) (*entity.View, error) {
	if w.edgeCache == nil {
		w.edgeCache = make(map[edge]entity.Model)
	}

	ent := w.Entities.AddEdge(name, o, d)
	err := w.Graph.AddEdge(o, d)

	w.edgeCache[edge{o, d}] = ent.Id

	return ent, err
}

func (w *World) FindReachableWorld(ctx context.Context) set.Hash[entity.Model] {
	reachable := make(set.Hash[entity.Model])

	bfs := graph.BreadthFirst{
		Selector: &RulesAwareSelector{
			w,
			graph.Successors,
		},
		Visitor: graph.VisitorFunc(func(c context.Context, n graph.Node) error {
			entity, ok := w.nodeCache[n]
			if !ok {
				return fmt.Errorf("unknown node %d", n)
			}

			// tagging the node and edge components with reachable
			// would be much more useful actually, but requires
			// a more sophiscated storage system
			reachable.Add(entity)
			return nil
		}),
	}

	spawns := w.Entities.Query(
		entity.Tagged(SpawnComponent),
		entity.Tagged(NodeComponent),
	)

	for _, s := range spawns {
		spawn, _ := s.Get(NodeComponent)
		bfs.Walk(ctx, w.Graph.G, spawn.(graph.Node))
	}

	return reachable
}
