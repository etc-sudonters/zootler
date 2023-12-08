package table

import (
	"errors"
	"fmt"
)

var _ Table = (*table)(nil)

type Table interface {
	InsertColumn(col ColumnData) (ColumnId, error)
	RetrieveColumns(cols ...ColumnId) ([]ColumnData, error)
}

type table struct {
	/*
	 *	columns are 1 indexed but we handle that complexity internally
	 *	this is to detect accidental creation of zero values columnids
	 */
	cols []ColumnData
}

func (tbl *table) InsertColumn(col ColumnData) (ColumnId, error) {
	id := ColumnId(len(tbl.cols) + 1)
	col.id = id
	tbl.cols = append(tbl.cols, col)
	return id, nil
}

func (tbl table) RetrieveColumns(cols ...ColumnId) ([]ColumnData, error) {
	if len(cols) == 0 {
		return nil, errors.New("no columns requested")
	}

	columns := make([]ColumnData, len(cols))

	for i, id := range cols {
		id := id - 1
		if int(id) >= len(tbl.cols) || 0 >= id {
			return nil, fmt.Errorf("invalid column %q requested", id)
		}

		columns[i] = tbl.cols[id]
	}

	return columns, nil
}
