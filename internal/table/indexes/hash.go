package indexes

import (
	"sudonters/zootler/internal/table"

	"github.com/etc-sudonters/substrate/skelly/bitset"
)

type hashbitmap[T comparable] map[T]bitset.Bitset64

func (h hashbitmap[T]) isset(key T, which uint64) bool {
	members := h.membersetfor(key)
	return members.IsSet(which)
}

func (h hashbitmap[T]) set(key T, which uint64) {
	members := h.membersetfor(key)
	members.Set(which)
	h[key] = members
}

func (h hashbitmap[T]) unset(which uint64) {
	for key, members := range h {
		if members.IsSet(which) {
			(&members).Unset(which)
			h[key] = members
			break
		}
	}
}

func (h hashbitmap[T]) membersetfor(key T) bitset.Bitset64 {
	members, ok := h[key]
	if !ok {
		members = bitset.Bitset64{}
	}

	return members
}

type HashIndex[T comparable] struct {
	members hashbitmap[T]
	hasher  TableHashingFunc[T]
}

func CreateHashIndex[TComponent any, TIndex comparable](f HashingFunc[TComponent, TIndex]) *HashIndex[TIndex] {
	return NewHashIndex(TableHasherFrom(f))
}

func NewHashIndex[T comparable](f TableHashingFunc[T]) *HashIndex[T] {
	return &HashIndex[T]{
		members: make(hashbitmap[T], 8),
		hasher:  f,
	}
}

func (h *HashIndex[T]) Set(r table.RowId, v table.Value) {
	idx, ok := h.hasher(v)
	if !ok {
		return
	}
	h.members.set(idx, uint64(r))
}

func (h *HashIndex[T]) Unset(r table.RowId) {
	h.members.unset(uint64(r))
}

// this bitset is intersected / & / AND'd
func (h *HashIndex[T]) Rows(v table.Value) bitset.Bitset64 {
	idx, ok := h.hasher(v)
	if !ok {
		return bitset.Bitset64{}
	}
	return bitset.Copy(h.members.membersetfor(idx))
}
