package app

import (
	"context"
	"sudonters/zootler/internal/query"
)

type ApplicationShuttingDown struct{}

func (_ ApplicationShuttingDown) Error() string {
	return "Application shutting down"
}

type Zootlr struct {
	ctx     context.Context
	cancel  context.CancelCauseFunc
	reason  error
	storage query.Engine
}

func (z *Zootlr) Engine() query.Engine {
	return z.storage
}

func (z *Zootlr) Ctx() context.Context {
	return z.ctx
}

func (z *Zootlr) Shutdown() {
	z.cancel(ApplicationShuttingDown{})
}

func (z *Zootlr) Error(err error) {
	if z.reason != nil {
		z.reason = err
		z.cancel(z.reason)
	}
}

type SetupFunc func(*Zootlr) error

func Setup(sc SetupApp) SetupFunc {
	return func(z *Zootlr) error {
		return sc.Setup(z.Ctx(), z.Engine())
	}
}

func New(ctx context.Context, ops ...SetupFunc) (*Zootlr, error) {
	var z Zootlr
	var engineError error
	z.storage, engineError = query.NewEngine()
	if engineError != nil {
		return nil, engineError
	}

	z.ctx, z.cancel = context.WithCancelCause(ctx)
	for i := range ops {
		if err := ops[i](&z); err != nil {
			return nil, err
		}
	}
	return &z, nil
}

type SetupApp interface {
	Setup(context.Context, query.Engine) error
}
