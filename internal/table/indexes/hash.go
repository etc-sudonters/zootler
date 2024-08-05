package indexes

import (
	"sudonters/zootler/internal/table"

	"github.com/etc-sudonters/substrate/mirrors"
	"github.com/etc-sudonters/substrate/skelly/bitset"
)

type IndexHasher[T comparable] func(table.Value) (T, bool)

func HashIndexFrom[TValue any, TIndex comparable](f func(TValue) (TIndex, bool)) *HashMapIndex[TIndex] {
	return NewHashMapIndex(
		func(v table.Value) (TIndex, bool) {
			if concrete, casted := v.(TValue); casted {
				return f(concrete)
			}
			return mirrors.Empty[TIndex](), false
		})
}

func NewHashMapIndex[T comparable](i IndexHasher[T]) *HashMapIndex[T] {
	return &HashMapIndex[T]{
		members: make(map[T]table.RowId, 8),
		hasher:  i,
	}
}

type HashMapIndex[TIndex comparable] struct {
	members map[TIndex]table.RowId
	hasher  IndexHasher[TIndex]
}

func (h *HashMapIndex[TIndex]) Set(e table.RowId, c table.Value) {
	idx, ok := h.hasher(c)
	if !ok {
		return
	}

	h.members[idx] = e
}

func (h *HashMapIndex[TIndex]) Unset(e table.RowId) {
	var key TIndex
	var found bool
	for k, row := range h.members {
		if row == e {
			key = k
			found = true
			break
		}
	}

	if found {
		delete(h.members, key)
	}
}

func (h *HashMapIndex[TIndex]) Rows(c table.Value) (b bitset.Bitset64) {
	idx, hashed := h.hasher(c)
	if hashed {
		entity, exists := h.members[idx]
		if exists {
			b.Set(uint64(entity))
		}
	}

	return
}
