package slipup

import (
	"errors"
	"fmt"
)

var ErrNotImplemented = errors.New("not implemented")

func Describe(e error, s string) error {
	return fmt.Errorf("%s: %w", s, e)
}

func Describef(e error, tpl string, v ...any) error {
	return Describe(e, fmt.Sprintf(tpl, v...))
}

func Createf(tpl string, v ...any) error {
	return fmt.Errorf(tpl, v...)
}

func NotImplemented(name string) error {
	return Describe(ErrNotImplemented, name)
}
