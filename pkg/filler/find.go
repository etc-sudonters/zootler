package filler

import (
	"context"
	"sync"

	"sudonters/zootler/internal/entity"
	"sudonters/zootler/pkg/logic"
	"sudonters/zootler/pkg/world"

	"github.com/etc-sudonters/substrate/skelly/graph"
	set "github.com/etc-sudonters/substrate/skelly/set/hash"
)

type setVisitor struct {
	s set.Hash[graph.Node]
}

func (s setVisitor) Visit(_ context.Context, g graph.Node) error {
	s.s.Add(g)
	return nil
}

func FindReachableWorld(ctx context.Context, w *world.World) (set.Hash[graph.Node], error) {
	reachable := set.New[graph.Node]()

	bfs := graph.BreadthFirst[graph.Destination]{
		Selector: &RulesAwareSelector[graph.Destination]{
			w, graph.Successors,
		},
		Visitor: setVisitor{reachable},
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
