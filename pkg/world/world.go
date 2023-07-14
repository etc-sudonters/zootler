package world

import (
	"context"
	"sync"

	"github.com/etc-sudonters/zootler/internal/datastructures/graph"
	"github.com/etc-sudonters/zootler/internal/datastructures/set"
	"github.com/etc-sudonters/zootler/pkg/entity"
	"github.com/etc-sudonters/zootler/pkg/logic"
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

func (w *World) AddNode(name logic.Name) (entity.View, error) {
	ent, err := w.Entities.Create(name)
	if err != nil {
		return nil, err
	}
	ent.Add(logic.Node{})
	w.Graph.AddNode(graph.Node(ent.Model()))
	return ent, nil
}

func (w *World) AddEdge(name logic.Name, o graph.Origination, d graph.Destination) (entity.View, error) {
	if w.edgeCache == nil {
		w.edgeCache = make(map[edge]entity.Model)
	}

	if _, ok := w.edgeCache[edge{o, d}]; ok {
		return nil, nil
	}

	if err := w.Graph.AddEdge(o, d); err != nil {
		return nil, err
	}

	ent, err := w.Entities.Create(name)
	ent.Add(logic.Edge{
		Destination: entity.Model(d),
		Origination: entity.Model(o),
	})

	w.edgeCache[edge{o, d}] = ent.Model()

	return ent, err
}

func (w *World) FindReachableWorld(ctx context.Context) (set.Hash[graph.Node], error) {
	reachable := set.New[graph.Node]()

	bfs := graph.BreadthFirst[graph.Destination]{
		Selector: &RulesAwareSelector[graph.Destination]{
			w, graph.Successors,
		},
		Visitor: &graph.VisitSet{S: reachable},
	}

	spawns, err := w.Entities.Query(logic.Spawn{})
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	for _, e := range spawns {
		e := e
		wg.Add(1)
		go func() {
			defer wg.Done()
			bfs.Walk(ctx, w.Graph.G, graph.Node(e.Model()))
		}()
	}

	wg.Wait()
	return reachable, nil
}
