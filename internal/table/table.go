package table

import (
	"fmt"
	"reflect"
	"sudonters/zootler/internal/iters"

	"github.com/etc-sudonters/substrate/mirrors"
	"github.com/etc-sudonters/substrate/reiterate"
	"github.com/etc-sudonters/substrate/skelly/bitset"
	"github.com/etc-sudonters/substrate/stageleft"
)

type IteratorWithLength[T any] interface {
	reiterate.Iterator[T]
	Len() int
}

type ResultIter struct {
	tbl     *table
	columns []ColumnId
	rows    reiterate.Iterator[int]
	len     int
}

func (r ResultIter) Len() int {
	return r.len
}

func (r ResultIter) MoveNext() bool {
	return r.rows.MoveNext()
}

func (r ResultIter) Current() RowTuple {
	current := RowId(r.rows.Current())
	values := make([]Value, len(r.columns))

	for i, col := range r.columns {
		values[i] = r.tbl.columns[col].column.Get(current)
	}

	return RowTuple{
		Row:    current,
		Values: values,
	}
}

type Table interface {
	Select(q Query) IteratorWithLength[RowTuple]
	Set(r RowTuple)
	Unset(r RowTuple)
}

type table struct {
	columns []ColumnData
	rows    []bitset.Bitset64
	typemap mirrors.TypeMap
}

func (t *table) Select(q Query) IteratorWithLength[RowTuple] {
	rows := bitset.New(bitset.Buckets(len(t.rows))).Complement()

	for _, pred := range q.Where {
		predCol := t.columns[pred.ColumnId()]
		idx := pred.DecideIndex(predCol)
		rows = rows.Intersect(idx.Index.Get(nil))
	}

	return ResultIter{
		tbl: t,
		columns: reiterate.Map[Selection, []Selection, ColumnId, []ColumnId](q.Load, func(s Selection) ColumnId {
			return s.Column
		}),
		rows: iters.Bitset64(rows),
		len:  rows.Len(),
	}
}

func (t *table) Set(r RowTuple) {
	for i, id := range t.intoColumnIds(r.Values) {
		t.columns[id].column.Set(r.Row, r.Values[i])
	}
}

func (t *table) Unset(r RowTuple) {
	for _, id := range t.intoColumnIds(r.Values) {
		t.columns[id].column.Unset(r.Row)
	}
}

func (t table) intoColumnIds(vs []Value) []ColumnId {
	ids := make([]ColumnId, len(vs))

	for i, v := range vs {
		typ := reflect.TypeOf(v)
		id, ok := t.typemap[typ]
		if !ok {
			panic(
				stageleft.AttachExitCode(
					fmt.Errorf("unknown column type: %s", typ.Name()),
					stageleft.ExitCode(99),
				))
		}

		ids[i] = ColumnId(id)
	}

	return ids
}
