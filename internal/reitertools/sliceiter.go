package reitertools

func SliceIter[E any, T ~[]E](t T) *sliceIter[E, T] {
	return &sliceIter[E, T]{t, -1}
}

func SubsliceIter[E any, T ~[]E](t T, idx int) *sliceIter[E, T] {
	return &sliceIter[E, T]{t, idx - 1}
}

type sliceIter[E any, T ~[]E] struct {
	t   T
	idx int
}

func (s sliceIter[E, T]) Index() int { return s.idx }

func (s *sliceIter[E, T]) MoveNext() bool {
	if s.idx+1 >= len(s.t) {
		return false
	}

	s.idx++
	return true
}

func (s *sliceIter[E, T]) Current() E {
	return s.t[s.idx]
}
