package graph

import (
	"context"
	"errors"

	"github.com/etc-sudonters/zootler/set"
)

type Visitor interface {
	Visit(context.Context, Node) error
}

type VisitorFunc func(context.Context, Node) error

func (v VisitorFunc) Visit(c context.Context, g Node) error {
	return v(c, g)
}

type Selector[T DirectionConstraint] interface {
	Select(Directed, Node) (Neighbors[T], error)
}

type SelectorFunc[T DirectionConstraint] func(Directed, Node) (Neighbors[T], error)

func (s SelectorFunc[T]) Select(g Directed, n Node) (Neighbors[T], error) {
	return s(g, n)
}

type Walker[T DirectionConstraint] interface {
	Walk(context.Context, Directed, Node) error
}

type BreadthFirst[T DirectionConstraint] struct {
	Selector[T]
	Visitor
}

func (b BreadthFirst[T]) Walk(ctx context.Context, g Directed, r Node) error {
	q := queue[Node]{r}
	seen := set.New[Node]()
	seen.Add(r)

	var node Node
	for len(q) > 0 {
		if err := ctx.Err(); err != nil {
			return err
		}

		node, q, _ = q.Pop()

		if err := b.Visitor.Visit(ctx, node); err != nil {
			return err
		}

		neighbors, err := b.Selector.Select(g, node)
		if err != nil {
			return err
		}

		for neighbor := range neighbors {
			neighbor := Node(neighbor)
			if !seen.Exists(neighbor) {
				seen.Add(neighbor)
				q = q.Push(neighbor)
			}
		}
	}

	return nil
}

type DepthFirst[T DirectionConstraint] struct {
	Visitor
	Selector[T]
}

func (d DepthFirst[T]) Walk(ctx context.Context, g Directed, r Node) error {
	s := stack[Node]{r}
	seen := set.New[Node]()

	var node Node
	for len(s) > 0 {
		if err := ctx.Err(); err != nil {
			return err
		}

		node, s, _ = s.Pop()

		if !seen.Exists(node) {
			if err := d.Visitor.Visit(ctx, Node(node)); err != nil {
				return err
			}

			seen.Add(node)
			neighbors, err := d.Selector.Select(g, Node(node))
			if err != nil {
				return err
			}

			for neighbor := range neighbors {
				s = s.Push(Node(neighbor))
			}
		}
	}

	return nil
}

var Successors SelectorFunc[Destination] = successors

func successors(g Directed, n Node) (Neighbors[Destination], error) {
	neighbors := make(Neighbors[Destination])

	for _, s := range g.Successors(Origination(n)) {
		neighbors.Add(s)
	}

	return neighbors, nil
}

type queue[T any] []T

func (q queue[T]) Push(t T) queue[T] {
	return append(q, t)
}

func (q queue[T]) Pop() (T, queue[T], error) {
	var t T
	if len(q) == 0 {
		return t, nil, errors.New("empty queue")
	}

	t, q = q[0], q[1:]

	return t, q, nil
}

type stack[T any] []T

func (s stack[T]) Push(t T) stack[T] {
	return append([]T{t}, s...)
}

func (s stack[T]) Pop() (T, []T, error) {
	var t T
	if len(s) == 0 {
		return t, nil, errors.New("empty stack")
	}

	t, s = s[0], s[1:]
	return t, s, nil
}
