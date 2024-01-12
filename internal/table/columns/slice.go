package columns

import (
	"reflect"
	"sudonters/zootler/internal/table"

	"github.com/etc-sudonters/substrate/skelly/bitset"
)

func BuildSliceColumn[T any]() *table.ColumnBuilder {
	return table.BuildColumnOf[T](NewSlice())
}

func NewSlice() *Slice {
	r := new(Slice)
	r.components = make([]table.Value, 0)
	r.members = &bitset.Bitset64{}
	return r
}

/*
 * Column backed by a slice. The slice is grown only when values for a row are
 * inserted into this column. The size of the slice is equal to the largest row
 * with a value in the column. This column does not attempt to compensate for
 * empty rows in the slice, they are left `nil`.
 */
type Slice struct {
	id         table.ColumnId
	typ        reflect.Type
	components []table.Value
	members    *bitset.Bitset64
}

func (row *Slice) Set(e table.RowId, c table.Value) {
	row.ensureSize(int(e))
	row.components[e] = c
	row.members.Set(uint64(e))
}

func (row *Slice) Unset(e table.RowId) {
	if len(row.components) < int(e) {
		return
	}

	row.components[e] = nil
	row.members.Unset(uint64(e))
}

func (row Slice) Get(e table.RowId) table.Value {
	if !row.members.IsSet(uint64(e)) {
		return nil
	}

	return row.components[e]
}

func (row *Slice) ensureSize(n int) {
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
