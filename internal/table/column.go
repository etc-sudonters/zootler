package table

import (
	"fmt"
	"reflect"

	"github.com/etc-sudonters/substrate/mirrors"
	"github.com/etc-sudonters/substrate/skelly/bitset"
)

type RowId uint64
type ColumnId uint64
type Value interface{}

type ColumnFactory func() Column

// core column interface
type Column interface {
	Get(e RowId) Value
	Set(e RowId, c Value)
	Unset(e RowId)
	ScanFor(Value) bitset.Bitset64
}

type ColumnMetadata interface {
	Type() reflect.Type
	Id() ColumnId
}

type ColumnData struct {
	id     ColumnId
	typ    reflect.Type
	column Column
}

func (c ColumnData) Column() Column {
	return c.column
}

func (c ColumnData) Type() reflect.Type {
	return c.typ
}

func (c ColumnData) Id() ColumnId {
	return c.id
}

func BuildColumn(col Column, typ reflect.Type) *ColumnBuilder {
	if col == nil {
		panic("nil column")
	}

	if typ == nil {
		panic("nil type information")
	}

	b := new(ColumnBuilder)
	b.column = col
	b.typ = typ
	return b
}

func BuildColumnOf[T Value](col Column) *ColumnBuilder {
	return BuildColumn(col, mirrors.TypeOf[T]())
}

type ColumnBuilder struct {
	typ    reflect.Type
	column Column
	index  Index
}

func (c *ColumnBuilder) Index(i Index) *ColumnBuilder {
	if c.index != nil {
		// let me know when to support multiple indexes
		panic(fmt.Errorf("column already indexed: %s", c.typ.Name()))
	}
	c.index = i
	return c
}

func (c ColumnBuilder) Type() reflect.Type {
	return c.typ
}

func (c ColumnBuilder) build(id ColumnId) ColumnData {
	return ColumnData{
		id:     id,
		typ:    c.typ,
		column: c.column,
	}
}
