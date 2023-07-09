package world

import (
	"context"

	"github.com/etc-sudonters/rando/graph"
	"github.com/etc-sudonters/rando/logic"
)

type RulesAwareSelector struct {
	World *World
	graph.Selector
}

func (s *RulesAwareSelector) Select(ctx context.Context, g graph.Model, n graph.Node) (graph.Neighbors, error) {
	neighbors := make(graph.Neighbors)
	candidates, err := s.Selector.Select(ctx, g, n)

	if err != nil {
		return neighbors, err
	}

	if len(candidates) == 0 {
		return neighbors, nil
	}

	for c := range candidates {
		edge, _ := s.World.edgeCache[edge{graph.Origination(n), graph.Destination(c)}]
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
