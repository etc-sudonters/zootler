package world

import (
	"github.com/etc-sudonters/substrate/skelly/bitset"
	"github.com/etc-sudonters/substrate/skelly/graph"
	"github.com/etc-sudonters/substrate/skelly/queue"
	"github.com/etc-sudonters/substrate/skelly/stack"
)

type Root graph.Origination

func NewGraph[TLoc any, TEdge any](g graph.Directed) Graph[TLoc, TEdge] {
	return Graph[TLoc, TEdge]{
		graph: g,
		locs:  make(map[graph.Node]TLoc, g.NodeCount()),
		edges: make(map[graph.Edge]TEdge, g.NodeCount()*4),
	}
}

type Graph[TLoc any, TEdge any] struct {
	graph graph.Directed
	locs  map[graph.Node]TLoc
	edges map[graph.Edge]TEdge
}

func (g *Graph[TLoc, TEdge]) SetNodeAttributes(n graph.Node, attr TLoc) {
	g.locs[n] = attr
}

func (g *Graph[TLoc, TEdge]) SetEdgeAttributes(o graph.Origination, d graph.Destination, attr TEdge) {
	g.edges[graph.Edge{o, d}] = attr
}

func (g *Graph[TLoc, TEdge]) NodeAttributes(n graph.Node) (TLoc, bool) {
	attrs, exists := g.locs[n]
	return attrs, exists
}

func (g *Graph[TLoc, TEdge]) EdgeAttributes(o graph.Origination, d graph.Destination) (TEdge, bool) {
	attrs, exists := g.edges[graph.Edge{o, d}]
	return attrs, exists
}

type destedge[TLoc any, TEdge any] struct {
	Destination      graph.Destination
	DestinationAttrs TLoc
	EdgeAttrs        TEdge
	Origination      graph.Origination
}

func (g *Graph[TLoc, TEdge]) Successors(n graph.Node) []destedge[TLoc, TEdge] {
	var dests []destedge[TLoc, TEdge]
	here := graph.Origination(n)
	succs, noSuccs := g.graph.Successors(n)
	if noSuccs == nil {
		for dest := range bitset.Iter64T[graph.Destination](succs).All {
			loc, _ := g.NodeAttributes(graph.Node(dest))
			edge, _ := g.EdgeAttributes(here, dest)

			dests = append(dests, destedge[TLoc, TEdge]{
				Origination:      here,
				Destination:      dest,
				DestinationAttrs: loc,
				EdgeAttrs:        edge,
			})
		}
	}
	return dests
}

func (g *Graph[TLoc, TEdge]) BFS(root Root) Exploration[TLoc, TEdge] {
	expl := Exploration[TLoc, TEdge]{
		graph:   g,
		visited: bitset.New(bitset.Buckets(uint64(g.graph.NodeCount()))),
		tracker: queue.Make[graph.Node](0, 16),
	}
	expl.tracker.Push(graph.Node(root))
	return expl
}

func (g *Graph[TLoc, TEdge]) DFS(root Root) Exploration[TLoc, TEdge] {
	expl := Exploration[TLoc, TEdge]{
		graph:   g,
		visited: bitset.New(bitset.Buckets(uint64(g.graph.NodeCount()))),
		tracker: stack.Make[graph.Node](0, 16),
	}
	expl.tracker.Push(graph.Node(root))
	return expl
}
