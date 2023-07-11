package graphutils

import (
	"context"
	"errors"
	"testing"

	"github.com/etc-sudonters/zootler/graph"
)

type VisitorArray []graph.Visitor
type DebugVisitor struct {
	*testing.T
}

func (v VisitorArray) Visit(ctx context.Context, node graph.Node) error {
	var errs []error

	for i := range v {
		if err := v[i].Visit(ctx, node); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

func (t *DebugVisitor) Visit(_ context.Context, node graph.Node) error {
	t.Logf("visiting node %s", node)
	return nil
}
