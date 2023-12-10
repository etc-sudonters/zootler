package table

import (
	"github.com/etc-sudonters/substrate/skelly/bitset"
)

type Table struct {
	Cols []ColumnData
	Rows []bitset.Bitset64
}

func (tbl *Table) CreateColumn(b colbuilder) (ColumnId, error) {
	col := b.build(ColumnId(len(tbl.Cols)))
	tbl.Cols = append(tbl.Cols, col)
	return col.id, nil
}

func (tbl *Table) InsertRow() RowId {
	id := RowId(len(tbl.Rows))
	tbl.Rows = append(tbl.Rows, bitset.New(-1))
	return id
}
