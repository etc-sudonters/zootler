package runtime

import "context"

type Function interface {
	Arity() int
	Run(context.Context, *VM, Values) (Value, error)
}

func NewCompiledFunc(arity int, chunk *Chunk, env *ExecutionEnvironment) *CompiledFunc {
	return &CompiledFunc{arity, chunk, env}
}

type CompiledFunc struct {
	arity int
	chunk *Chunk
	env   *ExecutionEnvironment
}

func (c *CompiledFunc) Arity() int {
	return c.arity
}

func (c *CompiledFunc) Run(ctx context.Context, vm *VM, values Values) (Value, error) {
	return vm.RunCompiledFunc(ctx, c, values)
}

func NewSimpleNativeFunc(arity int, call func(context.Context, *VM, Values) (Value, error)) SimpleNativeFunc {
	return SimpleNativeFunc{arity, call}
}

type SimpleNativeFunc struct {
	arity int
	call  func(context.Context, *VM, Values) (Value, error)
}

func (n SimpleNativeFunc) Arity() int {
	return n.arity
}

func (n SimpleNativeFunc) Run(ctx context.Context, vm *VM, values Values) (Value, error) {
	return n.call(ctx, vm, values)
}
