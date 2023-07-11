package graph

import (
	"fmt"
)

type Node int
type Origination Node
type Destination Node

type NodeConstraint interface {
	Node
}

type DirectionConstraint interface {
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
		origins: make(map[Origination][]Destination, c),
		dests:   make(map[Destination][]Origination, c),
	}
}

type Directed struct {
	origins map[Origination][]Destination
	dests   map[Destination][]Origination
}

func (g Directed) Predecessors(n Destination) []Origination {
	return copyEdgeList(g.dests[n])
}

func (g Directed) Successors(n Origination) []Destination {
	return copyEdgeList(g.origins[n])
}

func copyEdgeList[T DirectionConstraint](src []T) []T {
	if src == nil {
		return nil
	}
	dst := make([]T, len(src))
	copy(dst, src)
	return dst
}
