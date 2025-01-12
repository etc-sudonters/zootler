package objects

type BuildTable func(*Table)

func TableFrom(builder TableBuilder) BuildTable {
	return func(tbl *Table) {
		tbl.pointers = builder.pointers
		tbl.builtins = builder.funcs
		tbl.constants = make([]Object, len(builder.constants))
		for obj, index := range builder.constants {
			tbl.constants[index] = obj
		}
	}
}

func TableWithBuiltIns(funcs BuiltInFunctions) BuildTable {
	return func(tbl *Table) {
		tbl.builtins = funcs
	}
}

func NewTable(build ...BuildTable) Table {
	var tbl Table
	for i := range build {
		build[i](&tbl)
	}
	return tbl
}

type Table struct {
	constants []Object
	pointers  []Ptr
	builtins  BuiltInFunctions
}

func (this *Table) Constant(index Index) Object {
	return this.constants[int(index)]
}

func (this *Table) Pointer(index Index) Ptr {
	return this.pointers[int(index)]
}

func (this *Table) BuiltIn(index Index) *BuiltInFunction {
	return &this.builtins[int(index)]
}
