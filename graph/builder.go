package graph

import (
	"errors"
)

var ErrNodeNotFound = errors.New("node not found")
var ErrOriginNotFound = errors.New("origin not found")
var ErrDestNotFound = errors.New("destination not found")

type Builder struct {
	G Directed
}

type EdgeMap map[Origination][]Destination

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

func (b *Builder) AddEdge(o Origination, d Destination) error {

	if origins, exists := b.G.origins[o]; exists {
		b.G.origins[o] = append(origins, d)
	} else {
		return ErrOriginNotFound
	}

	if destinations, exists := b.G.dests[d]; exists {
		b.G.dests[d] = append(destinations, o)
	} else {
		return ErrDestNotFound
	}

	return nil
}

func (b *Builder) AddEdges(e EdgeMap) error {
	for o, neighbors := range e {
		for _, d := range neighbors {
			if err := b.AddEdge(o, d); err != nil {
				return err
			}
		}
	}
	return nil
}
