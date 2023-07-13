package bag

// Items in this package have no obvious home but are useful across a variety
// of domains

import (
	"fmt"
	"reflect"

	"golang.org/x/exp/constraints"
)

// returns a if a < b otherwise b
func Min[A constraints.Ordered](a, b A) A {
	if a < b {
		return a
	}
	return b
}

// determines if E is present in T
func Contains[E comparable, T ~[]E](needle E, haystack T) bool {
	for i := range haystack {
		if needle == haystack[i] {
			return true
		}
	}

	return false
}

// returns the name of the type represented, if it is a pointer & is prefixed
func NiceTypeName(t reflect.Type) string {
	if t.Kind() != reflect.Pointer {
		return t.Name()
	}

	t = t.Elem()
	return fmt.Sprintf("&%s", t.Name())
}
