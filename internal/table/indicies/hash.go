package indicies

import (
	"sudonters/zootler/internal/table"

	"github.com/etc-sudonters/substrate/skelly/bitset"
	"github.com/etc-sudonters/substrate/skelly/hashset"
)

type HashIndex[T comparable] struct {
	f func(table.Value) T

	entities map[T]hashset.Hash[table.RowId]
}

func (h HashIndex[T]) Set(e table.RowId, c table.Value) {
	h.entities[h.f(c)].Add(e)
}

func (h HashIndex[T]) Unset(e table.RowId, c table.Value) {
	delete(h.entities[h.f(c)], e)
}

func (h HashIndex[T]) Matches(c table.Value) int {
	return len(h.entities[h.f(c)])
}

func (h HashIndex[T]) Get(c table.Value) bitset.Bitset64 {
	var b bitset.Bitset64

	for e := range h.entities[h.f(c)] {
		b.Set(int(e))
	}

	return b
}
