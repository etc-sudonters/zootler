package slipup

import (
	"fmt"
)

func Trace(e error, s string) error {
	return fmt.Errorf("%s: %w", s, e)
}

func TraceMsg(e error, tpl string, v ...any) error {
	return Trace(e, fmt.Sprintf(tpl, v...))
}

func Create(tpl string, v ...any) error {
	return fmt.Errorf(tpl, v...)
}
