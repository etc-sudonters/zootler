package columns

import (
	"sudonters/zootler/internal/table"

	"github.com/etc-sudonters/substrate/skelly/bitset"
)

/*
 * Provides strongly typed ids and membership
 */
type RowId[T ~uint64] struct {
	t       table.Value
	members *bitset.Bitset64
}

func (m RowId[T]) Get(e table.RowId) table.Value {
	if m.members.IsSet(uint64(e)) {
		return T(e)
	}
	return nil
}

func (m *RowId[T]) Set(e table.RowId, c table.Value) {
	m.members.Set(uint64(e))
}

func (m *RowId[T]) Unset(e table.RowId) {
	m.members.Unset(uint64(e))
}

func (m *RowId[T]) ScanFor(c table.Value) bitset.Bitset64 {
	return bitset.Copy(*m.members)
}
