package world

import (
	"context"
	"sync"

	"sudonters/zootler/internal/graph"
	"sudonters/zootler/internal/set"
	"sudonters/zootler/pkg/entity"
	"sudonters/zootler/pkg/logic"
)

type edge struct {
	origin graph.Origination
	dest   graph.Destination
}

type Id int

type World struct {
	Id        Id
	Entities  Pool
	Graph     graph.Directed
	edgeCache map[edge]entity.Model
	nodeCache map[graph.Node]entity.Model
}

func (w *World) FindReachableWorld(ctx context.Context) (set.Hash[graph.Node], error) {
	reachable := set.New[graph.Node]()

	bfs := graph.BreadthFirst[graph.Destination]{
		Selector: &RulesAwareSelector[graph.Destination]{
			w, graph.Successors,
		},
		Visitor: &graph.VisitSet{S: reachable},
	}

	spawns, err := w.Entities.Query([]entity.Selector{entity.With[logic.Spawn]{}})
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	for _, e := range spawns {
		e := e
		wg.Add(1)
		go func() {
			defer wg.Done()
			bfs.Walk(ctx, w.Graph, graph.Node(e.Model()))
		}()
	}

	wg.Wait()
	return reachable, nil
}
