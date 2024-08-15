package runtime

import (
	"context"
	"sudonters/zootler/internal/slipup"
)

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
