package graph

import (
	"context"

	"github.com/etc-sudonters/rando/set"
)

type Visitor interface {
	Visit(context.Context, Node) error
}

type VisitorFunc func(context.Context, Node) error

func (v VisitorFunc) Visit(c context.Context, g Node) error {
	return v(c, g)
}

type Selector interface {
	Select(context.Context, Model, Node) (Neighbors, error)
}

type SelectorFunc func(context.Context, Model, Node) (Neighbors, error)

func (s SelectorFunc) Select(ctx context.Context, g Model, n Node) (Neighbors, error) {
	return s(ctx, g, n)
}

type Walker interface {
	Walk(context.Context, Model, Node) error
}

type BreadthFirst struct {
	Selector Selector
	Visitor  Visitor
}

func (b BreadthFirst) Walk(ctx context.Context, g Model, r Node) error {
	q := []Node{r}
	seen := set.New[Node]()
	seen.Add(r)

	for len(q) > 0 {
		if err := ctx.Err(); err != nil {
			return err
		}

		head := q[0]
		q = q[1:]

		err := b.Visitor.Visit(ctx, head)
		if err != nil {
			return err
		}

		neighbors, err := b.Selector.Select(ctx, g, head)
		if err != nil {
			return err
		}

		for s := range neighbors {
			if !seen.Exists(s) {
				seen.Add(s)
				q = append(q, s)
			}
		}
	}

	return nil
}

var Successors SelectorFunc = func(_ context.Context, g Model, n Node) (Neighbors, error) {
	neighbors := make(Neighbors)

	for _, s := range g.Successors(n) {
		neighbors.Add(Node(s))
	}

	return neighbors, nil
}
