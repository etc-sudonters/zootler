package table

type Predicate interface{}
type Selection interface{}

type Query struct {
	Load  []Selection
	Where []Predicate
}

type RowTuple struct {
	Row    RowId
	Values []Value
}

type ChangeTuple struct {
	Row    RowId
	Values []struct {
		Id    ColumnId
		Value Value
	}
}
