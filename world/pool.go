package world

import (
	"github.com/etc-sudonters/rando/entity"
	"github.com/etc-sudonters/rando/graph"
)

const (
	WorldIdComponent    entity.TagName = "origin-world"
	NodeComponent       entity.TagName = "node"
	EdgeComponent       entity.TagName = "edge"
	OriginNodeComponent entity.TagName = "origin-node"
	DestNodeComponent   entity.TagName = "dest-node"
	TokenComponent      entity.TagName = "token"
)

type Pool struct {
	World Id
	*entity.Pool
}

func (p Pool) Add(t entity.TagName) *entity.View {
	return p.Pool.Add(t).Use(OriginWorldArchetype(p.World))
}

func (p Pool) AddNode(t entity.TagName, n graph.Node) *entity.View {
	return p.Add(t).Use(NodeArchetype(n))
}

func (p Pool) AddEdge(t entity.TagName, o graph.Origination, d graph.Destination) *entity.View {
	return p.Add(t).Use(EdgeArchetype(o, d))
}

func (p Pool) AddToken(t entity.TagName) *entity.View {
	return p.Add(t).Use(entity.ArchetypeFunc(TokenArchetype))
}
