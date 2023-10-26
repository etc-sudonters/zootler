package filler

import (
	"sudonters/zootler/pkg/logic"
	"sudonters/zootler/pkg/world"

	"github.com/etc-sudonters/substrate/skelly/graph"
)

type RulesAwareSelector[T graph.Direction] struct {
	W *world.World
	S graph.Selector[T]
}

func edgeTo(n graph.Node, other interface{}) world.Edge {
	switch t := other.(type) {
	case graph.Origination:
		return world.Edge{t, graph.Destination(n)}
	case graph.Destination:
		return world.Edge{graph.Origination(n), t}
	}

	panic("unreachable")

}

func (s RulesAwareSelector[T]) Select(g graph.Directed, n graph.Node) ([]T, error) {
	candidates, err := s.S.Select(g, n)

	if err != nil {
		return nil, err
	}

	if len(candidates) == 0 {
		return nil, nil
	}

	accessibleNeighbors := make([]T, 0, len(candidates))

	var rule logic.Rule

	for _, c := range candidates {
		edge, _ := s.W.EdgeEntity(edgeTo(n, c))
		s.W.Entities.Get(edge, []interface{}{&rule})

		if rule == nil {
			accessibleNeighbors = append(accessibleNeighbors, c)
			continue
		}

		fulfillment, err := rule.Fulfill(s.W.Entities)

		if err != nil {
			return nil, err
		}

		if fulfillment {
			accessibleNeighbors = append(accessibleNeighbors, c)
		}
	}

	return accessibleNeighbors, nil
}
