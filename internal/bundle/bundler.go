package bundle

import (
	"sudonters/zootler/internal/iters"
	"sudonters/zootler/internal/table"

	"github.com/etc-sudonters/substrate/reiterate"
	"github.com/etc-sudonters/substrate/skelly/bitset"
)

type RowTuple struct {
	Id table.RowId
	ValueTuple
}

type ValueTuple struct {
	Cols   ColumnIds
	Values Values
}

// the pointer is invalid until the first MoveNext
// the pointer changes after each MoveNext
// implementations are free to define behavior when the iterator is exhausted
type Interface reiterate.Iterator[*RowTuple]

type Fill struct {
	bitset.Bitset64
}
type ColumnIds []table.ColumnId
type Columns []table.ColumnData
type Values []table.Value

type Bundler func(fill Fill, columns Columns) Interface

func DefaultBundle(fill Fill, columns Columns) Interface {
	return &rowOrderedBundle{
		iter: iters.Bitset64(fill.Bitset64),
		r:    nil,
		cols: columns,
	}
}

func SingleRowBundle(fill Fill, columns Columns) Interface {
	var tup RowTuple
	tup.Cols = make(ColumnIds, len(columns))
	tup.Values = make(Values, len(columns))
	tup.Id = table.RowId(fill.Elems()[0])

	for i, c := range columns {
		tup.Cols[i] = c.Id()
		tup.Values[i] = c.Column().Get(tup.Id)
	}

	return &singleRowBundle{r: tup}
}

type rowOrderedBundle struct {
	iter reiterate.Iterator[int]
	r    *RowTuple
	cols Columns
}

func (r rowOrderedBundle) MoveNext() bool {
	return r.iter.MoveNext()
}

func (r *rowOrderedBundle) Current() *RowTuple {
	if r.r == nil {
		r.r = new(RowTuple)
		r.r.Id = table.RowId(r.iter.Current())
		r.r.Values = make(Values, len(r.cols))
		r.r.Cols = make(ColumnIds, len(r.cols))

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

type singleRowBundle struct {
	r     RowTuple
	moved bool
}

func (s *singleRowBundle) MoveNext() bool {
	if s.moved {
		return false
	}

	s.moved = true
	return s.moved
}

func (s singleRowBundle) Current() *RowTuple {
	if s.moved {
		return &s.r
	}
	return nil
}
