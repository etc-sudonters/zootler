package table

import (
	"github.com/etc-sudonters/substrate/skelly/bitset"
)

type ColumnIds []ColumnId
type Columns []ColumnData
type Values []Value
type Row = bitset.Bitset64
type Rows []*Row

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
		Cols:    make([]ColumnData, 0),
		Rows:    make(Rows, 0),
		indexes: make(map[ColumnId]Index, 0),
	}
}

type Table struct {
	Cols    []ColumnData
	Rows    Rows
	indexes map[ColumnId]Index
}

func (tbl *Table) Lookup(c ColumnId, v Value) bitset.Bitset64 {
	if idx, ok := tbl.indexes[c]; ok {
		return idx.Rows(v)
	}
	return tbl.Cols[c].column.ScanFor(v)
}

func (tbl *Table) CreateColumn(b *ColumnBuilder) (ColumnData, error) {
	id := ColumnId(len(tbl.Cols))
	col := b.build(id)
	tbl.Cols = append(tbl.Cols, col)
	if b.index != nil {
		tbl.indexes[id] = b.index
	}
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
	if idx, ok := tbl.indexes[c]; ok {
		idx.Set(r, v)
	}
	return nil
}

func (tbl *Table) UnsetValue(r RowId, c ColumnId) error {
	col := tbl.Cols[c]
	col.column.Unset(r)
	row := tbl.Rows[r]
	row.Unset(uint64(c))
	if idx, ok := tbl.indexes[c]; ok {
		idx.Unset(r)
	}
	return nil
}
