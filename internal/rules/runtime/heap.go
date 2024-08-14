package runtime

type VmMemory struct {
	funcs map[string]Function
}

func (v *VmMemory) GetFunc(name string) (Function, error) {
	f, found := v.funcs[name]
	if !found {
		return nil, ErrUnboundName
	}
	return f, nil
}

func (v *VmMemory) AddFunc(name string, f Function) {
	v.funcs[name] = f
}

func NewVmMem() *VmMemory {
	return &VmMemory{
		funcs: make(map[string]Function, 256),
	}
}
