package columns

import (
	"sudonters/zootler/internal/table"

	"github.com/etc-sudonters/substrate/skelly/bitset"
)

func NewBit(singleton table.Value) *Bit {
	return &Bit{t: singleton, members: &bitset.Bitset64{}}
}

/*
 * Column backed by a bitset, consequently rows stored in this column do not
 * express unique values. Instead the presence of a row is handled by a
 * singleton value set on the column.
 */
type Bit struct {
	t       table.Value
	members *bitset.Bitset64
}

func (m Bit) Get(e table.RowId) table.Value {
	if m.members.IsSet(uint64(e)) {
		return m.t
	}
	return nil
}

func (m *Bit) Set(e table.RowId, c table.Value) {
	m.members.Set(uint64(e))
}

func (m *Bit) Unset(e table.RowId) {
	m.members.Unset(uint64(e))
}

func BuildMarkerColumn[T any]() *table.ColumnBuilder {
	var t T
	return table.BuildColumnOf[T](NewBit(t))
}

func BuildSingletonColumn[T any](t T) *table.ColumnBuilder {
	return table.BuildColumnOf[T](NewBit(t))
}
