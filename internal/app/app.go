package app

import (
	"context"
	"math/rand/v2"
	"reflect"
	"sudonters/zootler/internal/app"
	"sudonters/zootler/internal/generation"
	"sudonters/zootler/internal/query"
	"sudonters/zootler/internal/settings"

	"github.com/etc-sudonters/substrate/mirrors"
	"github.com/etc-sudonters/substrate/slipup"
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

// allows caller to bring their own entropy or invoke some spooky action at a distance
type RngFactory interface {
	Create(seed uint64) *rand.Rand
	Seed() (seed uint64)
}

type GenerationFactory interface {
	New(settings.ZootrSettings) (*generation.Generation, error)
}

type defaultGenerationFactory Zootlr

func (d *defaultGenerationFactory) New(settings settings.ZootrSettings) (*generation.Generation, error) {
	g := d.new()
	s := g.Seed
	z := *app.Zootlr(d)
	rng := GetResource[RngFactory](z)

	if s.Settings.Seed == 0 {
		s.Settings.Seed = rng.Res.Seed()
	}
	if s.Settings.Worlds == 0 {
		s.Settings.Worlds = 1
	}
	g.Rngesus = rng.Res.Create(s.Settings.Seed)
	s.Worlds = make([]generation.WorldBuilder, s.Settings.Worlds)

	if err := d.compileEdges(s); err != nil {
		return nil, slipup.Describe(err, "while creating generation")
	}

	return g, nil
}

func (d *defaultGenerationFactory) new() *generation.Generation {
	g := new(generation.Generation)
	z := *app.Zootlr(d)

	g.Seed = new(generation.SeedBuilder)
	g.Ctx, g.Cancel = context.WithCancelCause(z.ctx)

	return g
}

func (d *defaultGenerationFactory) compileEdges(sb *generation.SeedBuilder) error {
	return nil
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

func (z *Zootlr) AddResource(t any) {
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
