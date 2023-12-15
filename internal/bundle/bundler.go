package bundle

import (
	"errors"
	"fmt"
	"sudonters/zootler/internal/iters"
	"sudonters/zootler/internal/table"

	"github.com/etc-sudonters/substrate/reiterate"
	"github.com/etc-sudonters/substrate/skelly/bitset"
)

var ErrExpectSingleRow = errors.New("expected exactly 1 row")

// the pointer is invalid until the first MoveNext
// the pointer changes after each MoveNext
// implementations are free to define behavior when the iterator is exhausted
type Interface interface {
	reiterate.Iterator[*table.RowTuple]
	Len() int
}

type Fill struct {
	bitset.Bitset64
}

type Bundler func(fill Fill, columns table.Columns) Interface

func RowOrdered(fill Fill, columns table.Columns) Interface {
	return &rowOrdered{
		fill: fill,
		iter: iters.Bitset64(fill.Bitset64),
		r:    nil,
		cols: columns,
	}
}

func SingleRow(fill Fill, columns table.Columns) (Interface, error) {
	if fill.Len() != 1 {
		return nil, fmt.Errorf("%w: had %d", ErrExpectSingleRow, fill.Len())
	}

	var tup table.RowTuple
	tup.Cols = make(table.ColumnIds, len(columns))
	tup.Values = make(table.Values, len(columns))
	tup.Id = table.RowId(fill.Elems()[0])

	for i, c := range columns {
		tup.Cols[i] = c.Id()
		tup.Values[i] = c.Column().Get(tup.Id)
	}

	return &single{r: tup}, nil
}

type rowOrdered struct {
	fill Fill
	iter reiterate.Iterator[int]
	r    *table.RowTuple
	cols table.Columns
}

func (r rowOrdered) MoveNext() bool {
	return r.iter.MoveNext()
}

func (r *rowOrdered) Current() *table.RowTuple {
	if r.r == nil {
		r.r = new(table.RowTuple)
		r.r.Id = table.RowId(r.iter.Current())
		r.r.Values = make(table.Values, len(r.cols))
		r.r.Cols = make(table.ColumnIds, len(r.cols))

		for i, col := range r.cols {
			r.r.Cols[i] = col.Id()
		}
	}

	vt := r.r
	for i, col := range r.cols {
		vt.Values[i] = col.Column().Get(vt.Id)
	}

	return vt
}

func (r rowOrdered) Len() int {
	return r.fill.Len()
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
