package magicbean

import (
	"sudonters/zootler/zecs"

	"github.com/etc-sudonters/substrate/skelly/graph"
)

type ExplorableEdge struct {
	Kind   EdgeKind
	Entity zecs.Entity
	Rule   CompiledRule
}

type ExplorableWorld struct {
	Graph graph.Directed
	Edges map[Transit]ExplorableEdge
}
