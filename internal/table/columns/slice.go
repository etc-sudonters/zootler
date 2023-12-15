package columns

import (
	"reflect"
	"sudonters/zootler/internal/table"

	"github.com/etc-sudonters/substrate/skelly/bitset"
)

func NewSliceColumn() *SliceColumn {
	r := new(SliceColumn)
	r.components = make([]table.Value, 0)
	r.members = &bitset.Bitset64{}
	return r
}

type SliceColumn struct {
	id         table.ColumnId
	typ        reflect.Type
	components []table.Value
	members    *bitset.Bitset64
}

func (row *SliceColumn) Set(e table.RowId, c table.Value) {
	row.ensureSize(int(e))
	row.components[e] = c
	row.members.Set(uint64(e))
}

func (row *SliceColumn) Unset(e table.RowId) {
	if len(row.components) < int(e) {
		return
	}

	row.components[e] = nil
	row.members.IsSet(uint64(e))
}

func (row SliceColumn) Get(e table.RowId) table.Value {
	if !row.members.IsSet(uint64(e)) {
		return nil
	}

	return row.components[e]
}

func (row *SliceColumn) ensureSize(n int) {
	if len(row.components) > n {
		return
	}

	if n == 0 {
		n = 1
	}

	expaded := make([]table.Value, n+1, n*2)
	copy(expaded, row.components)
	row.components = expaded
}
