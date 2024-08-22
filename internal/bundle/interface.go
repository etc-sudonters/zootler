package bundle

import (
	"fmt"
	"sudonters/zootler/internal/table"

	"github.com/etc-sudonters/substrate/skelly/bitset"
)

type RowIter = func(table.RowId, table.ValueTuple) bool

type Interface interface {
	All(RowIter)
	Len() int
}

type Empty struct{}

func (e Empty) All(RowIter) {}
func (e Empty) Len() int    { return 0 }

func Many(fill bitset.Bitset64, columns table.Columns) Interface {
	return &many{fill, columns}
}

func Single(fill bitset.Bitset64, columns table.Columns) (Interface, error) {
	if fill.Len() != 1 {
		return nil, fmt.Errorf("%w: had %d", ErrExpectSingleRow, fill.Len())
	}

	var tup table.RowTuple
	tup.Cols = make(table.ColumnMetas, len(columns))
	tup.Values = make(table.Values, len(columns))
	tup.Id = table.RowId(fill.Elems()[0])

	for i, c := range columns {
		tup.Cols[i].Id = c.Id()
		tup.Cols[i].T = c.Type()
		tup.Values[i] = c.Column().Get(tup.Id)
	}

	s := single(tup)
	return &s, nil
}

type many struct {
	fill    bitset.Bitset64
	columns table.Columns
}

func (r *many) All(yield RowIter) {
	var vt table.ValueTuple
	vt.Values = make(table.Values, len(r.columns))
	vt.Cols = make(table.ColumnMetas, len(r.columns))
	for i, col := range r.columns {
		vt.Cols[i].Id = col.Id()
		vt.Cols[i].T = col.Type()
	}

	biter := bitset.Iter64(r.fill)

	for rowId := range biter.All {
		for i, col := range r.columns {
			vt.Values[i] = col.Column().Get(table.RowId(rowId))
		}

		if !yield(table.RowId(rowId), vt) {
			return
		}
	}

}

func (r many) Len() int {
	return r.fill.Len()
}

type single table.RowTuple

func (s *single) All(yield RowIter) {
	yield(s.Id, s.ValueTuple)
}

func (s single) Len() int {
	return 1
}
