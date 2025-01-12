package objects

import "fmt"

type Index uint16

const maxObjects int = 0xFFFF

func NewTableBuilder() TableBuilder {
	return TableBuilder{
		constants:    make(tracker[Object], 256),
		names:        make(tracker[string], 512),
		builtinNames: make(tracker[string], 16),
		funcs:        make(BuiltInFunctions, 8),
	}
}

type tracker[V comparable] map[V]Index

func (this tracker[V]) track(tracking V) Index {
	if index, exists := this[tracking]; exists {
		return index
	}

	size := len(this)
	if size > maxObjects {
		panic("too many constants")
	}

	index := Index(size)
	this[tracking] = index
	return index
}

type TableBuilder struct {
	constants    tracker[Object]
	builtinNames tracker[string]
	names        tracker[string]
	funcs        BuiltInFunctions
	pointers     []Ptr
}

func (this *TableBuilder) AddConstant(constant Object) Index {
	return this.constants.track(constant)
}

func (this *TableBuilder) GetPointerFor(name string) Index {
	index, exists := this.names[name]
	if !exists {
		panic(fmt.Errorf("pointer for %q not declared", name))
	}

	return index
}

func (this *TableBuilder) AddPointer(name string, ptr Ptr) Index {
	idx := this.names.track(name)
	this.pointers = resize(this.pointers, int(idx))
	this.pointers[int(idx)] = ptr
	return idx
}

func (this *TableBuilder) GetBuiltIn(name string) Index {
	index, exists := this.builtinNames[name]
	if !exists {
		panic(fmt.Errorf("built in %q not declared", name))
	}
	return index
}

func (this *TableBuilder) AddBuiltIn(builtIn BuiltInFunction) Index {
	idx := this.builtinNames.track(builtIn.Name)
	this.funcs = resize(this.funcs, int(idx))
	this.funcs[int(idx)] = builtIn
	return idx
}

func resize[T any, TS ~[]T](arr TS, size int) TS {
	if size < len(arr) {
		return arr
	}

	if size < 31 {
		size = 32
	} else {
		size = size * 2
	}

	grown := make(TS, size)
	copy(grown, arr)
	return grown
}
