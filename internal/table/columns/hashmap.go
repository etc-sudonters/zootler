package columns

import (
	"reflect"
	"sudonters/zootler/internal/table"

	"github.com/etc-sudonters/substrate/skelly/bitset"
)

func NewMap() *Map {
	return &Map{make(map[table.RowId]table.Value)}
}

type Map struct {
	entities map[table.RowId]table.Value
}

func (s *Map) Get(e table.RowId) table.Value {
	return s.entities[e]
}

func (s *Map) Set(e table.RowId, c table.Value) {
	s.entities[e] = c
}

func (s *Map) Unset(e table.RowId) {
	delete(s.entities, e)
}

func (s *Map) ScanFor(v table.Value) (b bitset.Bitset64) {
	for id, value := range s.entities {
		if reflect.DeepEqual(v, value) {
			b.Set(uint64(id))
		}
	}

	return
}

func (s *Map) Len() int {
	return len(s.entities)
}

func HashMapColumn[T any]() *table.ColumnBuilder {
	return table.BuildColumnOf[T](NewMap())
}
