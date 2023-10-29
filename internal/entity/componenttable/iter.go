package componenttable

import (
	"reflect"
	"sudonters/zootler/internal/entity"
)

func (t *Table) Rows() *tableIter {
	return &tableIter{t, 0}
}

type tableIter struct {
	t   *Table
	idx int
}

type RowData struct {
	r *Row
}

func (r RowData) Capacity() int {
	return len(r.r.components)
}

func (r RowData) Len() int {
	return r.r.members.Len()
}

func (r RowData) Id() entity.ComponentId {
	return r.r.id
}

func (r RowData) Type() reflect.Type {
	return r.r.typ
}

func (t *tableIter) MoveNext() bool {
	if 1+t.idx >= len(t.t.rows) {
		return false
	}

	t.idx++
	return true
}

func (t *tableIter) Current() RowData {
	return RowData{t.t.rows[t.idx]}
}
