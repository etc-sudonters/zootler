package columns

import (
	"reflect"
	"sudonters/zootler/internal/table"

	"github.com/etc-sudonters/substrate/skelly/bitset"
)

func NewSliceColumn(id table.ColumnId, entityBuckets int) *SliceColumn {
	r := new(SliceColumn)
	r.id = id
	r.components = make([]table.Value, 0)
	r.members = bitset.New(entityBuckets)
	return r
}

// this should be used for components that belong to a majority of entities
type SliceColumn struct {
	id         table.ColumnId
	typ        reflect.Type
	components []table.Value
	members    bitset.Bitset64
}

func (row *SliceColumn) Set(e table.RowId, c table.Value) {
	row.ensureSize(int(e))
	row.components[e] = c
	row.members.Set(int(e))
}

func (row *SliceColumn) Unset(e table.RowId) {
	if len(row.components) < int(e) {
		return
	}

	row.components[e] = nil
	row.members.Clear(int(e))
}

func (row SliceColumn) Get(e table.RowId) table.Value {
	if !row.members.Test(int(e)) {
		return nil
	}

	return row.components[e]
}

func (row *SliceColumn) ensureSize(n int) {
	if len(row.components) > n {
		return
	}

	expaded := make([]table.Value, n+1, n*2)
	copy(expaded, row.components)
	row.components = expaded
}
