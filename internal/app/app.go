package app

import (
	"context"
	"reflect"
	"sudonters/zootler/internal/query"

	"github.com/etc-sudonters/substrate/mirrors"
)

type ApplicationShuttingDown struct{}

func (_ ApplicationShuttingDown) Error() string {
	return "Application shutting down"
}

type Zootlr struct {
	ctx       context.Context
	cancel    context.CancelCauseFunc
	reason    error
	storage   query.Engine
	resources map[reflect.Type]any
}

func (z *Zootlr) Run(cmd AppCmd) error {
	return cmd(z)
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

func(z *Zootlr) AddResource(t any) {
    z.resources[reflect.TypeOf(t)] = t
}

type SetupFunc func(*Zootlr) error

func Setup(sc SetupApp) SetupFunc {
	return func(z *Zootlr) error {
		return sc.Setup(z)
	}
}

func AddResource[T any](res T) SetupFunc {
	return func(z *Zootlr) error {
		z.resources[mirrors.TypeOf[T]()] = res
		return nil
	}
}

type Resource[T any] struct {
	Res T
}

func GetResource[T any](z *Zootlr) *Resource[T] {
	resource, exists := z.resources[mirrors.TypeOf[T]()]
	if exists {
		return &Resource[T]{Res: resource.(T)}
	}

	return nil
}

func New(ctx context.Context, ops ...SetupFunc) (*Zootlr, error) {
	var z Zootlr
	var engineError error
	z.resources = make(map[reflect.Type]any)
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
	Setup(*Zootlr) error
}

type AppCmd func(*Zootlr) error
