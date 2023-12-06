package table

import "reflect"

type RowId uint64
type ColumnId uint64
type Value interface{}

type ColumnFactory func() Column

// core column interface
type Column interface {
	Get(e RowId) Value
	Set(e RowId, c Value)
	Unset(e RowId)
}

type ColumnData struct {
	column Column
	typ    reflect.Type
	id     ColumnId
}
