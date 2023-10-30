package reitertools

func ToSlice[E any](it Iterator[E]) []E {
	var slice []E

	for it.MoveNext() {
		slice = append(slice, it.Current())
	}

	return slice
}

func Next[E any](it Iterator[E]) (E, bool) {
	if !it.MoveNext() {
		var e E
		return e, false
	}
	return it.Current(), true
}
