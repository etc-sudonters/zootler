package intern

type Handle[T any] uint32

func NewInterner[T comparable]() HashIntern[T] {
	return HashIntern[T]{
		interned: make(map[T]Handle[T], 16),
	}
}

func NewInternerF[T any, U comparable](f func(T) U) HashInternF[T, U] {
	return HashInternF[T, U]{
		f:          f,
		HashIntern: NewInterner[U](),
	}
}

type HashIntern[T comparable] struct {
	interned map[T]Handle[T]
}

func (c HashIntern[T]) Intern(t T) Handle[T] {
	if handle, exists := c.interned[t]; exists {
		return Handle[T](handle)
	}

	handle := Handle[T](len(c.interned) + 1)
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

type HashInternF[T any, U comparable] struct {
	HashIntern[U]
	f func(T) U
}

func (c HashInternF[T, U]) Intern(t T) Handle[U] {
	return c.HashIntern.Intern(c.f(t))
}

func (c HashIntern[T]) IntoTable() []T {
	ts := make([]T, len(c.interned))

	for t, idx := range c.interned {
		ts[idx-1] = t
	}

	return ts
}
