package graph

import (
	"fmt"
)

// an entity in the graph
type Node int

// an entity in the graph that is the origination for an edge
type Origination Node

// an entity in the graph that is the destination for an edge
type Destination Node

// describes a graph as the set of originating edges
type OriginationMap map[Origination][]Destination

// describes a graph as the set of terminating edges
type destinationMap map[Destination][]Origination

// allows specifying direction when interacting with a Directed
type Direction interface {
	Origination | Destination
}

func (n Node) String() string {
	return fmt.Sprintf("Node{%d}", n)
}

func (d Destination) String() string {
	return fmt.Sprintf("Destination{%d}", d)
}

func (o Origination) String() string {
	return fmt.Sprintf("Origination{%d}", o)
}

func WithCapacity(c int) Directed {
	return Directed{
		origins: make(OriginationMap, c),
		dests:   make(destinationMap, c),
	}
}

func FromOriginationMap(src OriginationMap) Directed {
	b := Builder{WithCapacity(len(src))}
	b.AddEdges(src)
	return b.G
}

// adjanceny list, edges maintain insertion order, nodes do not
// do not construct directly, use a provided ctor
// direct usage of Directed is readonly
type Directed struct {
	origins map[Origination][]Destination
	dests   map[Destination][]Origination
}

// given Node n, find all other nodes that point at it
func (g Directed) Predecessors(n Node) ([]Origination, error) {
	l := copyEdgeList(g.dests[Destination(n)])
	if l == nil {
		return l, ErrDestNotFound
	}
	return l, nil
}

// given Node n, find all nodes that it points at
func (g Directed) Successors(n Node) ([]Destination, error) {
	l := copyEdgeList(g.origins[Origination(n)])
	if l == nil {
		return l, ErrOriginNotFound
	}
	return l, nil
}

func copyEdgeList[T Direction](src []T) []T {
	if src == nil {
		return nil
	}
	dst := make([]T, len(src))
	copy(dst, src)
	return dst
}
