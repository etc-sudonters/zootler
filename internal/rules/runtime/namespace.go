package runtime

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
	return f, nil
}

func (n *FuncNamespace) AddFunc(name string, f Function) {
	n.funcs[name] = f
}

func NewFuncNamespace() *FuncNamespace {
	return &FuncNamespace{
		funcs: make(map[string]Function, 256),
	}
}
