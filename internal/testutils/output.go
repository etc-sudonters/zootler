package testutils

import "testing"

func Dump(t *testing.T, v any) {
	t.Logf("%+v\n", v)
}
