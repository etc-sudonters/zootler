package columns

import (
	"sudonters/zootler/internal/skelly/bitset32"
	"sudonters/zootler/internal/table"
)

func NewBit(singleton table.Value) *Bit {
	return &Bit{t: singleton, members: &bitset32.Bitset32{}}
}

/*
 * Column backed by a bitset, consequently rows stored in this column do not
 * express unique values. Instead the presence of a row is handled by a
 * singleton value set on the column.
 */
type Bit struct {
	t       table.Value
	members *bitset32.Bitset32
}

func (m Bit) Get(e table.RowId) table.Value {
	if m.members.IsSet(uint32(e)) {
		return m.t
	}
	return nil
}

func (m *Bit) Set(e table.RowId, c table.Value) {
	m.members.Set(uint32(e))
}

func (m *Bit) Unset(e table.RowId) {
	m.members.Unset(uint32(e))
}

func (m *Bit) ScanFor(c table.Value) bitset32.Bitset32 {
	return bitset32.Copy32(*m.members)
}

func (m *Bit) Len() int {
	return m.members.Len()
}

func BitColumnOf[T any]() *table.ColumnBuilder {
	var t T
	return BitColumnUsing(t)
}

func BitColumnUsing[T any](t T) *table.ColumnBuilder {
	return table.BuildColumnOf[T](NewBit(t))
}
