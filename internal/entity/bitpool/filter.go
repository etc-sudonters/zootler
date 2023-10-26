package bitpool

import (
	"sudonters/zootler/internal/entity"

	"github.com/etc-sudonters/substrate/skelly/set/bits"
)

type filter struct {
	i bits.Bitset64
	e bits.Bitset64
}

func (f *filter) init(k int) {
	f.i = bits.New(k)
	f.e = bits.New(k)
}

func (f filter) include(t entity.ComponentId) {
	f.i.Set(int(t))
}

func (f filter) exclude(t entity.ComponentId) {
	f.e.Set(int(t))
}

func (f filter) test(b bitview) bool {
	if !bits.IsEmpty(f.i) && !b.comps.Intersect(f.i).Eq(f.i) {
		return false
	}

	if !b.comps.Difference(f.e).Eq(b.comps) {
		return false
	}

	return true
}
