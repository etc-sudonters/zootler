package filler

import (
	"errors"
	"fmt"
	"sudonters/zootler/internal/entity"
	"sudonters/zootler/pkg/logic"
	"sudonters/zootler/pkg/world"

	"github.com/etc-sudonters/substrate/skelly/graph"
)

type RulesAwareSelector[T graph.Direction] struct {
	W            *world.World
	S            graph.Selector[T]
	Unfullfilled []T
}

func (s RulesAwareSelector[T]) Select(g graph.Directed, n graph.Node) ([]T, error) {
	s.Unfullfilled = nil
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
		edge, _ := s.W.Edge(world.Edge{
			Origination: entity.Model(n),
			Destination: entity.Model(c),
		})
		if err := edge.Get(&rule); err != nil {
			if errors.Is(err, entity.ErrNotLoaded) {
				panic(fmt.Errorf("edge %d exists without a rule: %w", edge.Model(), err))
			}
			return nil, err
		}

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
		} else {
			s.Unfullfilled = append(s.Unfullfilled, c)
		}
	}

	return accessibleNeighbors, nil
}
