package queue

import "errors"

type Q[T any] []T

func From[E any, T ~[]E](src T) Q[E] {
	q := make(Q[E], len(src))
	copy(q, src)
	return q
}

func (q Q[T]) Push(t T) Q[T] {
	return append(q, t)
}

func (q Q[T]) Pop() (T, Q[T], error) {
	var t T
	if len(q) == 0 {
		return t, q, errors.New("empty queue")
	}

	t, q = q[0], q[1:]

	return t, q, nil
}
