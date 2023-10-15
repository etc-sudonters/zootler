package graph

import (
	"context"
	"errors"

	"sudonters/zootler/internal/queue"
	"sudonters/zootler/internal/set"
	"sudonters/zootler/internal/stack"
)

type (
	// used to provide human readable diagnostic output
	DebugFunc func(string, ...any)

	// calls F on current node, selected and err results from S
	DebugSelector[T Direction] struct {
		F DebugFunc
		S Selector[T]
	}

	// calls F on current node, error from V
	DebugVisitor struct {
		F DebugFunc
		V Visitor
	}

	// provides current node to every Visitor in slice
	VisitorArray []Visitor

	// records visited nodes in a queue.Q[Node]
	VisitQueue struct {
		Q queue.Q[Node]
	}

	// records visited nodes in a stack.S[Node]
	VisitStack struct {
		S stack.S[Node]
	}

	// records visited nodes in a set.Hash[Node]
	VisitSet struct {
		S set.Hash[Node]
	}
)

func (v VisitorArray) Visit(ctx context.Context, node Node) error {
	var errs []error

	for i := range v {
		if err := v[i].Visit(ctx, node); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

func (d DebugVisitor) Visit(ctx context.Context, node Node) error {
	d.F("visiting node %s", node)
	err := d.V.Visit(ctx, node)
	if err != nil {
		d.F("error on visit to %s: %s", node, err)
	}
	return err
}

func (q *VisitQueue) Visit(_ context.Context, n Node) error {
	q.Q = q.Q.Push(n)
	return nil
}

func (s *VisitStack) Visit(_ context.Context, n Node) error {
	s.S = s.S.Push(n)
	return nil
}

func (s *VisitSet) Visit(_ context.Context, n Node) error {
	if s.S == nil {
		s.S = set.New[Node]()
	}
	s.S.Add(n)
	return nil
}

func (d DebugSelector[T]) Select(g Directed, n Node) ([]T, error) {
	d.F("selecting from %s", n)
	selected, err := d.S.Select(g, n)
	if err != nil {
		d.F("error while selecting: %s", err)
	}

	d.F("selected %+v", selected)

	return selected, err
}
