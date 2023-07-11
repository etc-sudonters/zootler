package testutils

import (
	"context"
	"testing"
)

func CreateContextFrom(t *testing.T) (context.Context, context.CancelCauseFunc) {
	var cancel context.CancelFunc
	ctx := context.Background()

	if deadline, ok := t.Deadline(); ok {
		ctx, cancel = context.WithDeadline(ctx, deadline)
	}

	ctx, cancelCause := context.WithCancelCause(ctx)

	return ctx, func(e error) {
		cancelCause(e)
		if cancel != nil {
			cancel()
		}
	}
}
