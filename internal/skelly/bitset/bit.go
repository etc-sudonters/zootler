package bitset

import (
	"math/bits"
)

func Buckets32(i uint32) int {
	return int(i / 32)
}

func BitIndex32(i uint32) uint32 {
	return 1 << (i % 32)
}

func New32(i int) Bitset32 {
	var b Bitset32
	b.buckets = make([]uint32, i)
	return b
}

func Create32(members ...uint32) Bitset32 {
	b := New32(0)
	for _, m := range members {
		b.Set(m)
	}
	return b
}

func WithBucketsFor32(i uint32) Bitset32 {
	return New32(Buckets32(i))
}

func FromRaw32(parts []uint32) Bitset32 {
	var b Bitset32
	b.buckets = parts
	return b
}

func ToRawParts32(b Bitset32) []uint32 {
	ret := make([]uint32, len(b.buckets))
	copy(ret, b.buckets)
	return ret
}

type Bitset32 struct {
	buckets []uint32
}

func IsEmpty32(b Bitset32) bool {
	for i := range b.buckets {
		if b.buckets[i] != 0 {
			return false
		}
	}
	return true
}

func Copy32(b Bitset32) Bitset32 {
	var n Bitset32
	n.buckets = make([]uint32, len(b.buckets))
	copy(n.buckets, b.buckets)
	return n
}

func (this *Bitset32) resize(bucket int) {
	if bucket < len(this.buckets) {
		return
	}

	buckets := make([]uint32, bucket+1)
	copy(buckets, this.buckets)
	this.buckets = buckets
}

func (this *Bitset32) Set(i uint32) bool {
	idx := Buckets32(i)
	bit := BitIndex32(i)
	this.resize(idx)
	bucket := this.buckets[idx]
	this.buckets[idx] = bucket | bit
	return bucket&bit == 0
}

func (this Bitset32) Unset(i uint32) {
	bucket := Buckets32(i)

	if bucket >= len(this.buckets) {
		return
	}

	this.buckets[bucket] &= ^BitIndex32(i)
}

func (this Bitset32) IsSet(i uint32) bool {
	bucket := Buckets32(i)
	if bucket >= len(this.buckets) {
		return false
	}

	bit := BitIndex32(i)
	return bit == (bit & this.buckets[bucket])
}

func (this Bitset32) Complement() Bitset32 {
	n := Copy32(this)
	for i, bits := range n.buckets {
		n.buckets[i] = ^bits
	}
	return n
}

func (this Bitset32) Intersect(n Bitset32) Bitset32 {
	buckets := min(len(this.buckets), len(n.buckets))
	r := Bitset32{}
	r.buckets = make([]uint32, buckets)

	for i := range r.buckets {
		r.buckets[i] = this.buckets[i] & n.buckets[i]
	}

	return r
}

func (this Bitset32) Union(n Bitset32) Bitset32 {
	if len(n.buckets) > len(this.buckets) {
		this, n = n, this
	}

	ret := Copy32(this)
	for i, bits := range n.buckets {
		ret.buckets[i] |= bits
	}

	return ret
}

func (this Bitset32) Difference(n Bitset32) Bitset32 {
	buckets := make([]uint32, max(len(this.buckets), len(n.buckets)))
	copy(buckets, n.buckets)
	n.buckets = buckets
	return this.Intersect(n.Complement())
}

func (this Bitset32) Eq(n Bitset32) bool {
	hi, lo := this.buckets, n.buckets

	if len(lo) > len(hi) {
		hi, lo = lo, hi
		for _, bucket := range hi[len(lo):] {
			if bucket != 0 {
				return false
			}
		}
	}

	for i := range lo {
		if lo[i] != hi[i] {
			return false
		}
	}

	return true
}

func (this Bitset32) Len() int {
	var count int

	for _, bucket := range this.buckets {
		for ; bucket != 0; bucket &= bucket - 1 {
			count++
		}
	}

	return count
}

func (this Bitset32) Elems() []uint32 {
	var elems []uint32

	for k, bucket := range this.buckets {
		k := uint32(k)
		for bucket != 0 {
			tz := uint32(bits.TrailingZeros32(bucket))
			elems = append(elems, k*32+tz)
			bucket ^= (1 << tz)
		}
	}

	return elems
}

func (this Bitset32) Pop() uint32 {
	for k, bucket := range this.buckets {
		if bucket == 0 {
			continue
		}

		tz := uint32(bits.TrailingZeros32(bucket))
		this.buckets[k] = bucket ^ (1 << tz)
		return uint32(k*32) + tz
	}
	return 0
}

func (this Bitset32) IsEmpty() bool {
	for _, bucket := range this.buckets {
		if bucket != 0 {
			return false
		}
	}
	return true

}

func (this Bitset32) IsSuperSetOf(other Bitset32) bool {
	intersect := this.Intersect(other)
	return intersect.Eq(other)
}

func (this Bitset32) IsSubSetOf(other Bitset32) bool {
	return other.IsSuperSetOf(this)
}
