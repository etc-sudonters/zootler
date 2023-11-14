package componenttable

import (
	"reflect"
	"sudonters/zootler/internal/entity"

	"github.com/etc-sudonters/substrate/reiterate"
)

func (t *Table) Rows() reiterate.Iterator[RowData] {
	return &tableIter{t, 0}
}

type tableIter struct {
	t   *Table
	idx int
}

func (t tableIter) Index() int { return t.idx }

type RowData struct {
	r *Row
}

func (r RowData) Get(e entity.Model) entity.Component {
	return r.r.Get(e)
}

func (r RowData) Components() reiterate.Iterator[RowEntry] {
	return r.r.Components()
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

type RowEntry struct {
	Entity    entity.Model
	Component entity.Component
}
