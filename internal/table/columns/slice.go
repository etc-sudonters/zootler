package columns

import (
	"fmt"
	"reflect"
	"sudonters/zootler/internal/table"

	"github.com/etc-sudonters/substrate/skelly/bitset"
)

func SliceColumn[T any]() *table.ColumnBuilder {
	return table.BuildColumnOf[T](NewSlice())
}

func SizedSliceColumn[T any](size uint) *table.ColumnBuilder {
	return table.BuildColumnOf[T](SizedSlice(size))
}

func NewSlice() *Slice {
	r := new(Slice)
	r.components = make([]table.Value, 0)
	r.members = &bitset.Bitset64{}
	return r
}

func SizedSlice(size uint) *Slice {
	bitset := bitset.WithBucketsFor(uint64(size))
	r := new(Slice)
	r.components = make([]table.Value, size)
	r.members = &bitset
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
		n = 8
	}

	expanded := make([]table.Value, n+1, n*2)
	copy(expanded, row.components)
	row.components = expanded
}

func (row Slice) ScanFor(v table.Value) bitset.Bitset64 {
	density := float64(row.members.Len()) / float64(len(row.components))
	if density > 0.6 {
		return row.scanValues(v)
	} else {
		return row.scanMembers(v)
	}
}

func (row Slice) Len() int {
	return row.members.Len()
}

func (row Slice) scanMembers(v table.Value) (b bitset.Bitset64) {
    bititer := bitset.Iter64(*row.members)
	for id := range bititer.All {
		value := row.components[id]
		if value == nil {
			panic(fmt.Errorf("nil value indexed in slice row at %d", id))
		}

		if reflect.DeepEqual(v, value) {
			b.Set(uint64(id))
		}
	}

	return
}

func (row Slice) scanValues(v table.Value) (b bitset.Bitset64) {
	for id, value := range row.components {
		if value == nil || !reflect.DeepEqual(v, value) {
			continue
		}

		b.Set(uint64(id))
	}

	return
}
