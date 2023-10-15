package bag

import "math/rand"

func Shuffle[T any, E ~[]T](elms E) {
	rand.Shuffle(len(elms), func(i, j int) {
		elms[i], elms[j] = elms[j], elms[i]
	})
}
