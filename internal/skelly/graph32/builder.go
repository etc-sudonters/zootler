package graph32

import (
	"sudonters/zootler/internal/skelly/bitset32"
)

type Builder struct {
	Directed
}

func (b *Builder) AddNode(n Node) {
	if _, exists := b.origins[n]; !exists {
		b.origins[n] = bitset32.Bitset{}
	}
}

func (b *Builder) AddNodes(ns []Node) {
	for _, n := range ns {
		b.AddNode(n)
	}
}

func (b *Builder) AddEdge(origin, dest Node) {
	b.AddNode(origin)
	b.AddNode(dest)
	origins := b.origins[origin]
	bitset32.Set(&origins, dest)
	b.origins[origin] = origins
}
