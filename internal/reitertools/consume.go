package reitertools

func ToSlice[E any](it Iterator[E]) []E {
	var slice []E

	for it.MoveNext() {
		slice = append(slice, it.Current())
	}

	return slice
}
