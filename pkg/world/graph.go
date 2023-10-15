package world

import (
	"sudonters/zootler/internal/graph"
	"sudonters/zootler/pkg/logic"
)

type RulesAwareSelector[T graph.Direction] struct {
	W *World
	S graph.Selector[T]
}

func edgeTo(n graph.Node, other interface{}) edge {
	switch t := other.(type) {
	case graph.Origination:
		return edge{t, graph.Destination(n)}
	case graph.Destination:
		return edge{graph.Origination(n), t}
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
		edge := s.W.edgeCache[edgeTo(n, c)]
		s.W.Entities.Get(edge, &rule)

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
