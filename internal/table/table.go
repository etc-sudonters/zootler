package table

import (
	"fmt"
	"reflect"

	"github.com/etc-sudonters/substrate/mirrors"
	"github.com/etc-sudonters/substrate/reiterate"
	"github.com/etc-sudonters/substrate/skelly/bitset"
	"github.com/etc-sudonters/substrate/stageleft"
)

type IteratorWithLength[T any] interface {
	reiterate.Iterator[T]
	Len() int
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
	panic("not implemented") // TODO: Implement
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
