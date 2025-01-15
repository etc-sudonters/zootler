package bitset

type new32 interface {
	~uint32
}

func CreateT32[T new32](members ...T) Bitset32 {
	b := New32(0)
	for _, m := range members {
		b.Set(uint32(m))
	}
	return b
}

func Set32[T new32](b *Bitset32, t T) bool {
	return b.Set(uint32(t))
}

func Unset32[T new32](b *Bitset32, t T) {
	b.Unset(uint32(t))
}

func IsSet32[T new32](b *Bitset32, t T) bool {
	return b.IsSet(uint32(t))
}
