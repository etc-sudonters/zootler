package columns

import (
	"sudonters/zootler/internal/table"

	"github.com/etc-sudonters/substrate/skelly/bitset"
)

func NewBit[T any](k int) Bit[T] {
	return Bit[T]{members: bitset.New(k)}
}

// these columns do not carry data but group entities into subgroups
type Bit[T table.Value] struct {
	t       T
	members bitset.Bitset64
}

func (m Bit[T]) Get(e table.RowId) table.Value {
	if m.members.Test(int(e)) {
		return m.t
	}
	return nil
}

func (m Bit[T]) Set(e table.RowId, c table.Value) {
	m.members.Set(int(e))
}

func (m Bit[T]) Unset(e table.RowId) {
	m.members.Clear(int(e))
}
