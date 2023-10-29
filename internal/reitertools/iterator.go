package reitertools

type Iterator[E any] interface {
	MoveNext() bool
	Current() E
}

func Filter[E any](i Iterator[E], f func(E) bool) Iterator[E] {
	return filter[E]{i, f}
}

func Flatten[T any, U any](src Iterator[T], f func(T) Iterator[U]) Iterator[U] {
	return &flattener[T, U]{src, f, nil}
}

type flattener[T any, U any] struct {
	src Iterator[T]
	f   func(T) Iterator[U]
	sub Iterator[U]
}

func (f *flattener[T, U]) MoveNext() bool {
	if f.sub != nil && f.sub.MoveNext() {
		return true
	}

	for f.src.MoveNext() {
		f.sub = f.f(f.src.Current())
		if f.sub.MoveNext() {
			return true
		}
	}

	return false
}

func (f *flattener[T, U]) Current() U {
	return f.sub.Current()
}

type filter[E any] struct {
	i Iterator[E]
	f func(E) bool
}

func (f filter[E]) MoveNext() bool {
	for {
		if !f.i.MoveNext() {
			return false
		}

		if f.f(f.i.Current()) {
			return true
		}
	}
}

func (f filter[E]) Current() E {
	return f.i.Current()
}
