package reitertools

type Iterator[E any] interface {
	MoveNext() bool
	Current() E
	Index() int
}

func Map[T any, U any](src Iterator[T], f func(T, int) U) Iterator[U] {
	return &mapiter[T, U]{src, f}
}

func Filter[E any](i Iterator[E], f func(E, int) bool) Iterator[E] {
	return &filter[E]{i, f}
}

func Flatten[T any, U any](src Iterator[T], f func(T) Iterator[U]) Iterator[U] {
	return &flattener[T, U]{src, f, nil, -1}
}

type flattener[T any, U any] struct {
	src Iterator[T]
	f   func(T) Iterator[U]
	sub Iterator[U]
	i   int
}

func (f *flattener[T, U]) MoveNext() bool {
	if f.sub != nil && f.sub.MoveNext() {
		f.i++
		return true
	}

	for f.src.MoveNext() {
		f.sub = f.f(f.src.Current())
		if f.sub.MoveNext() {
			f.i++
			return true
		}
	}

	return false
}

func (f flattener[T, U]) Index() int {
	return f.i
}

func (f *flattener[T, U]) Current() U {
	return f.sub.Current()
}

type filter[E any] struct {
	i Iterator[E]
	f func(E, int) bool
}

func (f *filter[E]) MoveNext() bool {
	for {
		if !f.i.MoveNext() {
			return false
		}

		if f.f(f.i.Current(), f.i.Index()) {
			return true
		}
	}
}

func (f filter[E]) Current() E {
	return f.i.Current()
}

func (f filter[E]) Index() int {
	return f.i.Index()
}

type mapiter[T any, U any] struct {
	src Iterator[T]
	f   func(T, int) U
}

func (m *mapiter[T, U]) MoveNext() bool {
	return m.src.MoveNext()
}

func (m mapiter[T, U]) Current() U {
	return m.f(m.src.Current(), m.src.Index())
}

func (m mapiter[T, U]) Index() int {
	return m.src.Index()
}
