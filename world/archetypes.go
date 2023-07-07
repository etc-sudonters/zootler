package world

import (
	"github.com/etc-sudonters/rando/entity"
	"github.com/etc-sudonters/rando/graph"
)

var TokenArchetype = entity.ArchetypeTag(TokenComponent)

func OriginWorldArchetype(w Id) entity.ArchetypeFunc {
	return func(t entity.Tags) {
		t.Apply(WorldIdComponent, w)
	}
}

func NodeArchetype(n graph.Node) entity.ArchetypeFunc {
	return func(t entity.Tags) {
		t.Apply(NodeComponent, n)
	}
}

func EdgeArchetype(o graph.Origination, d graph.Destination) entity.ArchetypeFunc {
	return func(t entity.Tags) {
		t.Apply(OriginNodeComponent, o).
			Apply(DestNodeComponent, d).
			Use(entity.ArchetypeTag(EdgeComponent))
	}
}
