package graph

import (
	"errors"

	"github.com/etc-sudonters/zootler/internal/bag"
)

var ErrNodeNotFound = errors.New("node not found")
var ErrOriginNotFound = errors.New("origin not found")
var ErrDestNotFound = errors.New("destination not found")

// adds nodes and edges to a Directed
type Builder struct {
	G Directed
}

// adds a node to the graph if it doesn't exist already
func (b *Builder) AddNode(n Node) {
	if _, exists := b.G.dests[Destination(n)]; !exists {
		b.G.dests[Destination(n)] = []Origination{}
		b.G.origins[Origination(n)] = []Destination{}
	}
}

func (b *Builder) AddNodes(ns []Node) {
	for _, n := range ns {
		b.AddNode(n)
	}
}

// connects o -> d, if either node doesn't exist they are created
// edges maintain insertion order, duplicate edges are not added
func (b *Builder) AddEdge(o Origination, d Destination) error {
	b.AddNode(Node(o))
	b.AddNode(Node(d))
	if !bag.Contains(d, b.G.origins[o]) {
		b.G.origins[o] = append(b.G.origins[o], d)
	}

	if !bag.Contains(o, b.G.dests[d]) {
		b.G.dests[d] = append(b.G.dests[d], o)
	}
	return nil
}

func (b *Builder) AddEdges(e OriginationMap) error {
	for o, neighbors := range e {
		for _, d := range neighbors {
			if err := b.AddEdge(o, d); err != nil {
				return err
			}
		}
	}
	return nil
}
