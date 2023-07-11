package world

import (
	"github.com/etc-sudonters/zootler/graph"
	"github.com/etc-sudonters/zootler/logic"
)

type RulesAwareSelector[T graph.DirectionConstraint] struct {
	World *World
	graph.Selector[T]
}

func (s *RulesAwareSelector[T]) Select(g graph.Directed, n graph.Node) (graph.Neighbors[T], error) {
	neighbors := make(graph.Neighbors[T])
	candidates, err := s.Selector.Select(g, n)

	if err != nil {
		return neighbors, err
	}

	if len(candidates) == 0 {
		return neighbors, nil
	}

	for c := range candidates {
		edge, _ := s.World.edgeCache[edge{graph.Origination(n), c}]
		view := s.World.Entities.Get(edge)
		component, _ := view.Get(RuleComponent)
		rule, hasRules := component.(logic.Rule)

		if !hasRules {
			neighbors.Add(c)
			continue
		}

		fulfillment, err := rule.Fulfill(s.World.Entities)

		if err != nil {
			return neighbors, err
		}

		if fulfillment {
			neighbors.Add(c)
		}
	}

	return neighbors, nil
}
