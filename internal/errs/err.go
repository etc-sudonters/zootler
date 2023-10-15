package errs

import "errors"

var NotImplErr = errors.New("not implemented")

func PanicNotImpled() error {
	panic(NotImplErr)
}

func NotImpled[T any]() (T, error) {
	var t T
	return t, NotImplErr
}
