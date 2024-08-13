package vm

import "sudonters/zootler/internal/rules/bytecode"

type ExecutionEnvironment struct {
	identifiers map[string]bytecode.Value
	parent      *ExecutionEnvironment
}

func NewEnv() *ExecutionEnvironment {
	return &ExecutionEnvironment{
		identifiers: make(map[string]bytecode.Value),
	}
}

func (e *ExecutionEnvironment) ChildScope() *ExecutionEnvironment {
	return &ExecutionEnvironment{
		identifiers: make(map[string]bytecode.Value),
		parent:      e,
	}
}

func (e *ExecutionEnvironment) Set(name string, v bytecode.Value) {
	e.identifiers[name] = v
}

func (e *ExecutionEnvironment) Lookup(name string) (bytecode.Value, bool) {
	if v, found := e.identifiers[name]; found {
		return v, true
	}

	if e.parent != nil {
		return e.Lookup(name)
	}

	return bytecode.NullValue(), false
}
