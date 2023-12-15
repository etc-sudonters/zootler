package table

import (
	"github.com/etc-sudonters/substrate/skelly/bitset"
)

type ColumnIds []ColumnId
type Columns []ColumnData
type Values []Value

type RowTuple struct {
	Id RowId
	ValueTuple
}

type ValueTuple struct {
	Cols   ColumnIds
	Values Values
}

func New() *Table {
	return &Table{
		Cols: make([]ColumnData, 0),
		Rows: make([]*bitset.Bitset64, 0),
	}
}

type Table struct {
	Cols []ColumnData
	Rows []*bitset.Bitset64
}

func (tbl *Table) CreateColumn(b *ColumnBuilder) (ColumnData, error) {
	col := b.build(ColumnId(len(tbl.Cols)))
	tbl.Cols = append(tbl.Cols, col)
	return col, nil
}

func (tbl *Table) InsertRow() RowId {
	id := RowId(len(tbl.Rows))
	tbl.Rows = append(tbl.Rows, &bitset.Bitset64{})
	return id
}

func (tbl *Table) SetValue(r RowId, c ColumnId, v Value) error {
	col := tbl.Cols[c]
	col.column.Set(r, v)
	row := tbl.Rows[r]
	row.Set(uint64(c))
	return nil
}

func (tbl *Table) UnsetValue(r RowId, c ColumnId) error {
	col := tbl.Cols[c]
	col.column.Unset(r)
	row := tbl.Rows[r]
	row.Unset(uint64(c))
	return nil
}
