package columns

import "sudonters/zootler/internal/table"

func NewMap() Map {
	return Map{make(map[table.RowId]table.Value)}
}

type Map struct {
	entities map[table.RowId]table.Value
}

func (s Map) Get(e table.RowId) table.Value {
	return s.entities[e]
}

func (s Map) Set(e table.RowId, c table.Value) {
	s.entities[e] = c
}

func (s Map) Unset(e table.RowId) {
	delete(s.entities, e)
}

func BuildHashColumn[T any]() *table.ColumnBuilder {
	return table.BuildColumnOf[T](NewMap())
}
