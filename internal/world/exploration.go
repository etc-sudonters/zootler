package world

import (
	"github.com/etc-sudonters/substrate/skelly/bitset"
	"github.com/etc-sudonters/substrate/skelly/graph"
)

type tracker[T any] interface {
	Push(T)
	Pop() (T, error)
	Len() int
}

type Exploration[TLoc any, TEdge any] struct {
	graph   *Graph[TLoc, TEdge]
	tracker tracker[graph.Node]
	visited bitset.Bitset32
}

func (expl *Exploration[TLoc, TEdge]) Accept(n graph.Destination) {
	expl.tracker.Push(graph.Node(n))
}

func (expl *Exploration[TLoc, TEdge]) Walk(yield func(destedge[TLoc, TEdge]) bool) {
	for expl.tracker.Len() > 0 {
		node, _ := expl.tracker.Pop()
		if !bitset.Set(&expl.visited, node) {
			continue
		}
		for _, successor := range expl.graph.Successors(node) {
			if !yield(successor) {
				return
			}
		}
	}
}
