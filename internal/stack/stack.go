package stack

import "errors"

type S[T any] []T

func From[E any, T ~[]E](src T) S[E] {
	l := len(src)
	dest := make(S[E], l)

	for i, e := range src {
		dest[l-1-i] = e
	}
	return dest
}

func (s S[T]) Push(t T) S[T] {
	return append([]T{t}, s...)
}

func (s S[T]) Pop() (T, []T, error) {
	var t T
	if len(s) == 0 {
		return t, nil, errors.New("empty stack")
	}

	t, s = s[0], s[1:]
	return t, s, nil
}
