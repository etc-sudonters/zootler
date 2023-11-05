package reitertools

type Index[T any] struct {
	Index int
	Elem  T
}

func Enumerate[T any](i Iterator[T]) Iterator[Index[T]] {
	return EnumerateFrom(i, 0)
}

func EnumerateFrom[T any](i Iterator[T], startAt int) Iterator[Index[T]] {
	return Map(i, enumerate[T](startAt))
}

func enumerate[T any](current int) func(T) Index[T] {
	return func(t T) Index[T] {
		this := Index[T]{
			Index: current,
			Elem:  t,
		}
		current++
		return this
	}
}
