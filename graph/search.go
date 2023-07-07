package graph

import (
	"context"

	"github.com/etc-sudonters/rando/set"
)

type Traverser interface {
	Traverse(context.Context, Model, Node) error
}

type BreadthFirstSuccessors struct {
	Strategy func(Node, []Destination) ([]Destination, error)
}

func (b BreadthFirstSuccessors) Traverse(ctx context.Context, g Model, r Node) error {
	q := []Node{r}
	seen := set.New[Node]()
	seen.Add(r)

	for len(q) > 0 {
		head := q[0]
		q = q[1:]

		successors, err := b.Strategy(head, g.Successors(head))
		if err != nil {
			return err
		}

		for _, s := range successors {
			s := Node(s)
			if !seen.Exists(s) {
				seen.Add(s)
				q = append(q, s)
			}
		}
	}

	return nil
}
