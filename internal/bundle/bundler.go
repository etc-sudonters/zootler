package bundle

import (
	"errors"
	"fmt"
	"math/bits"
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

func Filled(b bitset.Bitset64) Fill {
	return Fill{b}
}

func ToMap[TKey comparable, TValue any](i Interface, f func(*table.RowTuple) (TKey, TValue, error)) (map[TKey]TValue, error) {
	m := make(map[TKey]TValue, i.Len())
	for i.MoveNext() {
		key, value, err := f(i.Current())
		if err != nil {
			return nil, err
		}

		m[key] = value
	}
	return m, nil
}

type Bundler func(fill Fill, columns table.Columns) Interface

func RowOrdered(fill Fill, columns table.Columns) Interface {
	return &rowOrdered{
		fill: fill,
		iter: Iter(fill.Bitset64),
		r:    nil,
		cols: columns,
	}
}

func SingleRow(fill Fill, columns table.Columns) (Interface, error) {
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

type rowOrdered struct {
	fill Fill
	iter reiterate.Iterator[uint64]
	r    *table.RowTuple
	cols table.Columns
}

func (r rowOrdered) MoveNext() bool {
	return r.iter.MoveNext()
}

func (r *rowOrdered) Current() *table.RowTuple {
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

func Iter(b bitset.Bitset64) reiterate.Iterator[uint64] {
	return &iter{bitset.ToRawParts(b), 0, -1}
}

type iter struct {
	parts   []uint64
	current uint64
	partIdx int
}

func (b *iter) MoveNext() bool {
	if b.partIdx >= len(b.parts) {
		return false
	}

	for b.current == 0 {
		b.partIdx++
		if b.partIdx >= len(b.parts) {
			return false
		}

		candidate := b.parts[b.partIdx]
		if candidate != 0 {
			b.current = candidate
			return true
		}
	}

	b.current ^= (1 << bits.TrailingZeros64(b.current))
	return true
}

func (b *iter) Current() uint64 {
	return uint64(b.partIdx)*64 + uint64(bits.TrailingZeros64(b.current))
}
