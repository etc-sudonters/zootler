package rng

import "math/rand"

type Xoshiro256PP [4]uint64

func Xoshiro256PPSource64(x *Xoshiro256PP) rand.Source64 {
	return &randAdapter{x}
}

func NewXoshiro256PPFromU64(seed uint64) *Xoshiro256PP {
	x := Xoshiro256PP([4]uint64{0, 0, 0, 0})
	for i := range x {
		seed = SplitMix64(seed)
		x[i] = seed
	}
	return &x
}

func (x *Xoshiro256PP) NextUint64() uint64 {
	result := rotl(x[0]+x[3], 23) + x[0]
	t := x[1] << 27
	x[2] ^= x[0]
	x[3] ^= x[1]
	x[1] ^= x[2]
	x[0] ^= x[3]
	x[2] ^= t
	x[3] = rotl(x[3], 45)
	return result
}

func rotl(x uint64, k int) uint64 {
	return (x << k) | (x >> (64 - k))
}

func SplitMix64(state uint64) uint64 {
	state += 0x9e3779b97f4a7c15
	z := (state ^ (state >> 30)) * 0xbf58476d1ce4e5b9
	z = (z ^ (z >> 27)) * 0x94d049bb133111eb
	return z ^ (z >> 31)
}

type randAdapter struct {
	x *Xoshiro256PP
}

func (a *randAdapter) Int63() int64 {
	n := a.x.NextUint64()
	return int64(n << 1 >> 1)
}

func (a *randAdapter) Seed(seed int64) {
	a.x = NewXoshiro256PPFromU64(uint64(seed))
}

func (a *randAdapter) Uint64() uint64 {
	return a.x.NextUint64()
}
