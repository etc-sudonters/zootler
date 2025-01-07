package objects

type BuildTable func(*Table)

func TableFrom(builder TableBuilder) BuildTable {
	return func(tbl *Table) {
		tbl.constants = make([]Object, len(builder.constants))
		for obj, index := range builder.constants {
			tbl.constants[index] = obj
		}
	}
}

func TableWithBuiltIns(funcs BuiltInFunctionTable) BuildTable {
	return func(tbl *Table) {
		tbl.BuiltIns = NewBuiltins(funcs)
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
	BuiltIns  BuiltInFunctions
}

func (this *Table) Constant(index Index) Object {
	return this.constants[int(index)]
}

func (this *Table) BuiltIn(index Index) *BuiltInFunc {
	return this.BuiltIns.Get(index)
}
