package graph

import (
	"github.com/etc-sudonters/rando/set"
)

type Node int

type Destination Node
type Origination Node

func WithCapacity(c int) Model {
	return Model{
		nodes:   make([]Neighbors, 0, c),
		inEdges: make(map[Node]Neighbors, c),
	}
}

type Model struct {
	nodes   []Neighbors
	inEdges map[Node]Neighbors
}

func (g Model) Predecessors(n Node) []Origination {
	if !g.canNodeExist(n) {
		return nil
	}

	originations := make([]Origination, 0, len(g.inEdges[n]))

	for origin := range g.inEdges[n] {
		originations = append(originations, Origination(origin))
	}

	return originations
}

func (g Model) Successors(n Node) []Destination {
	if !g.canNodeExist(n) {
		return nil
	}

	destinations := make([]Destination, 0, len(g.nodes[n]))

	for dest := range g.nodes[n] {
		destinations = append(destinations, Destination(dest))
	}

	return destinations
}

func (g Model) canNodeExist(n Node) bool {
	actualIdx := int(n)
	return 0 <= actualIdx && actualIdx < len(g.nodes)
}

type Neighbors set.Hash[Node]

func (n Neighbors) Add(i Node) {
	(set.Hash[Node])(n).Add(i)
}
