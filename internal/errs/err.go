package errs

import "errors"

var NotImplErr = errors.New("not implemented")

func PanicNotImpled() error {
	panic(NotImplErr)
}
