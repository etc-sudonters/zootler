package runtime

import "context"

type Function interface {
	Arity() int
	Run(context.Context, *VM, Values) (Value, error)
}

type CompiledFuncValue struct {
	arity int
	chunk *Chunk
	env   *ExecutionEnvironment
}

func (c *CompiledFuncValue) Arity() int {
	return c.arity
}

func (c *CompiledFuncValue) Run(ctx context.Context, vm *VM, values Values) (Value, error) {
	return vm.RunCompiledFunc(ctx, c, values)
}

type NativeFuncValue struct {
	arity int
	call  func(context.Context, *VM, Values) (Value, error)
}

func (n *NativeFuncValue) Arity() int {
	return n.arity
}

func (n *NativeFuncValue) Run(ctx context.Context, vm *VM, values Values) (Value, error) {
	return n.call(ctx, vm, values)
}
