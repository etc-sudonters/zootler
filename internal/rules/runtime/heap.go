package runtime

// we have heap at home
type VmHeap struct {
	Funcs map[string]Function
}
