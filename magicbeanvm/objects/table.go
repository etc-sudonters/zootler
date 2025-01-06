package objects

type Index uint16

const maxObjects int = 0xFFFF

func NewTableBuilder() TableBuilder {
	builtins := make(frozen[string], len(builtInIndex))
	for name := range builtInIndex {
		builtins[name] = builtInIndex[name]
	}

	tbl := TableBuilder{
		constants: make(tracker[Object], 256),
		names:     make(tracker[string], 512),
		builtins:  builtins,
	}

	return tbl
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

type frozen[V comparable] map[V]Index

func (this frozen[V]) track(tracking V) Index {
	if index, exists := this[tracking]; exists {
		return index
	}
	panic("cannot add new items to frozen tracker")
}

type TableBuilder struct {
	constants tracker[Object]
	builtins  frozen[string]
	names     tracker[string]
}

func (this *TableBuilder) Constant(constant Object) Index {
	return this.constants.track(constant)
}

func (this *TableBuilder) Name(name string) Index {
	return this.names.track(name)
}

func (this *TableBuilder) BuiltIn(name string) Index {
	return this.builtins.track(name)
}
