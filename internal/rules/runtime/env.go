package runtime

type ExecutionEnvironment struct {
	identifiers map[string]Value
	parent      *ExecutionEnvironment
}

func NewEnv() *ExecutionEnvironment {
	return &ExecutionEnvironment{
		identifiers: make(map[string]Value),
	}
}

func (e *ExecutionEnvironment) ChildScope() *ExecutionEnvironment {
	return &ExecutionEnvironment{
		identifiers: make(map[string]Value),
		parent:      e,
	}
}

func (e *ExecutionEnvironment) Set(name string, v Value) {
	e.identifiers[name] = v
}

func (e *ExecutionEnvironment) Lookup(name string) (Value, bool) {
	if v, found := e.identifiers[name]; found {
		return v, true
	}

	if e.parent != nil {
		return e.Lookup(name)
	}

	return NullValue(), false
}
