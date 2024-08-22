package runtime

import (
	"context"
	"github.com/etc-sudonters/substrate/slipup"
)

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

var decl declMarker = declMarker{}

type FuncNamespace struct {
	funcs map[string]Function
}

func (n *FuncNamespace) IsFunc(name string) bool {
	_, found := n.funcs[name]
	return found
}

func (n *FuncNamespace) GetFunc(name string) (Function, error) {
	f, found := n.funcs[name]
	if !found {
		return nil, ErrUnboundName
	}
	if _, isTombstone := f.(declMarker); isTombstone {
		return nil, slipup.Describef(ErrUnboundName, "function declared but not defined: '%s'", name)
	}
	return f, nil
}

func (n *FuncNamespace) DeclFunction(name string) error {
	n.funcs[name] = decl
	return nil
}

func (n *FuncNamespace) AddFunc(name string, f Function) {
	n.funcs[name] = f
}

func NewFuncNamespace() *FuncNamespace {
	return &FuncNamespace{
		funcs: make(map[string]Function, 512),
	}
}

type declMarker struct{}

func (_ declMarker) Arity() int {
	panic("not implemented")
}
func (_ declMarker) Run(_ context.Context, _ *VM, _ Values) (Value, error) {
	panic("not implemented")
}
