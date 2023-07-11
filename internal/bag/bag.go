package bag

import "golang.org/x/exp/constraints"

func Min[A constraints.Ordered](a, b A) A {
	if a < b {
		return a
	}
	return b
}
