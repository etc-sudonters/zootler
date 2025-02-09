package bitset32

type u32 interface {
	~uint32
}

func CreateT[T u32](members ...T) Bitset {
	b := New(0)
	for _, m := range members {
		b.Set(uint32(m))
	}
	return b
}

func Set[T u32](b *Bitset, t T) bool {
	return b.Set(uint32(t))
}

func Unset[T u32](b *Bitset, t T) {
	b.Unset(uint32(t))
}

func IsSet[T u32](b *Bitset, t T) bool {
	return b.IsSet(uint32(t))
}

func Intersects(this, other Bitset) bool {
	return !this.Intersect(other).IsEmpty()
}
