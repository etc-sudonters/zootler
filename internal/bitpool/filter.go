package bitpool

import (
	"sudonters/zootler/internal/entity"

	"github.com/etc-sudonters/substrate/skelly/bitset"
)

type filter struct {
	i bitset.Bitset64
	e bitset.Bitset64
}

func (f *filter) init(k int) {
	f.i = bitset.New(k)
	f.e = bitset.New(k)
}

func (f filter) include(t entity.ComponentId) {
	f.i.Set(int(t))
}

func (f filter) exclude(t entity.ComponentId) {
	f.e.Set(int(t))
}

func (f filter) test(b bitview) bool {
	if !bitset.IsEmpty(f.i) && !b.comps.Intersect(f.i).Eq(f.i) {
		return false
	}

	if !b.comps.Difference(f.e).Eq(b.comps) {
		return false
	}

	return true
}
