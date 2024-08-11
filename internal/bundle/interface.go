package bundle

import (
	"fmt"
	"sudonters/zootler/internal/table"

	"github.com/etc-sudonters/substrate/reiterate"
	"github.com/etc-sudonters/substrate/skelly/bitset"
)

// the pointer is invalid until MoveNext is called
// the pointer changes after each MoveNext
// the pointer becomes invalid after MoveNext returns false -- including the first call
type Interface interface {
	reiterate.Iterator[*table.RowTuple]
	Len() int
}

type Empty struct{}

func (e Empty) MoveNext() bool           { return false }
func (e Empty) Current() *table.RowTuple { return nil }
func (e Empty) Len() int                 { return 0 }

func Many(fill bitset.Bitset64, columns table.Columns) Interface {
	return &many{
		len:  fill.Len(),
		iter: Iter(fill),
		r:    nil,
		cols: columns,
	}
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

	return &single{r: tup}, nil
}

type many struct {
	len  int
	iter reiterate.Iterator[uint64]
	r    *table.RowTuple
	cols table.Columns
}

func (r many) MoveNext() bool {
	return r.iter.MoveNext()
}

func (r *many) Current() *table.RowTuple {
	if r.r == nil {
		r.r = new(table.RowTuple)
		r.r.Values = make(table.Values, len(r.cols))
		r.r.Cols = make(table.ColumnMetas, len(r.cols))

		for i, col := range r.cols {
			r.r.Cols[i].Id = col.Id()
			r.r.Cols[i].T = col.Type()
		}
	}

	thisId := table.RowId(r.iter.Current())
	vt := r.r
	vt.Id = thisId
	for i, col := range r.cols {
		vt.Values[i] = col.Column().Get(thisId)
	}

	return vt
}

func (r many) Len() int {
	return r.len
}

type single struct {
	r     table.RowTuple
	moved bool
}

func (s *single) MoveNext() bool {
	if s.moved {
		return false
	}

	s.moved = true
	return s.moved
}

func (s single) Current() *table.RowTuple {
	if s.moved {
		return &s.r
	}
	return nil
}

func (s single) Len() int {
	return 1
}
