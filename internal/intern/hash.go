package intern

type Handle[T any] uint32

func NewInterner[T comparable]() HashIntern[T] {
	return HashIntern[T]{
		interned: make(map[T]Handle[T], 16),
	}
}

func NewInternerF[T comparable](f func(T) T) HashInternF[T] {
	return HashInternF[T]{
		f: f, HashIntern: NewInterner[T](),
	}
}

type HashIntern[T comparable] struct {
	interned map[T]Handle[T]
}

func (c HashIntern[T]) Intern(t T) Handle[T] {
	if handle, exists := c.interned[t]; exists {
		return Handle[T](handle)
	}

	handle := Handle[T](len(c.interned))
	c.interned[t] = handle
	return handle
}

func (c HashIntern[T]) All(yield func(Handle[T], T) bool) {
	for thing, handle := range c.interned {
		if !yield(handle, thing) {
			return
		}
	}
}

func (c HashIntern[T]) Len() int {
	return len(c.interned)
}

type HashInternF[T comparable] struct {
	HashIntern[T]
	f func(T) T
}

func (c HashInternF[T]) Intern(t T) Handle[T] {
	return c.HashIntern.Intern(c.f(t))
}
