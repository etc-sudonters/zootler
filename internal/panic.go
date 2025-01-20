package internal

func PanicOnError(e error) {
	if e != nil {
		panic(e)
	}
}
