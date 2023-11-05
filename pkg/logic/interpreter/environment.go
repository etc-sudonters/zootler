package interpreter

type Environment struct {
	parent *Environment
	values map[string]Value
}

func NewEnv() Environment {
	var e Environment
	e.values = make(map[string]Value)
	return e
}

func (e Environment) Get(name string) (Value, bool) {
	v, ok := e.values[name]
	if !ok && e.parent != nil {
		v, ok = e.parent.Get(name)
	}

	return v, ok
}

func (e Environment) Set(name string, v Value) {
	e.values[name] = v
}

func (e Environment) Enclosed() Environment {
	inner := NewEnv()
	inner.parent = &e
	return inner
}
