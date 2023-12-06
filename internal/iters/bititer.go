package iters

import (
	"math/bits"

	"github.com/etc-sudonters/substrate/reiterate"
	"github.com/etc-sudonters/substrate/skelly/bitset"
)

func Bitset64(b bitset.Bitset64) reiterate.Iterator[int] {
	return &bitset64{bitset.ToRawParts(b), 0, -1}
}

type bitset64 struct {
	parts   []uint64
	current uint64
	partIdx int
}

func (b *bitset64) MoveNext() bool {
	if b.current == 0 {
		b.partIdx++
		if b.partIdx >= len(b.parts) {
			return false
		}

		b.current = b.parts[b.partIdx]
		return true
	}

	b.current ^= (1 << bits.TrailingZeros64(b.current))
	return true
}

func (b *bitset64) Current() int {
	return b.partIdx*64 + bits.TrailingZeros64(b.current)
}
