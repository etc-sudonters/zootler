package bag

// Items in this package have no obvious home but are useful across a variety
// of domains

import (
	"fmt"
	"reflect"

	"golang.org/x/exp/constraints"
)

func Max[A constraints.Ordered](a, b A) A {
	if a > b {
		return a
	}
	return b
}

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

func Map[A any, AT ~[]A, B any, BT ~[]B](as AT, f func(A) B) BT {
	bs := make(BT, len(as))

	for i, a := range as {
		bs[i] = f(a)
	}

	return bs
}

func Filter[A any, AT ~[]A, F func(A) bool](as AT, f F) AT {
	var na AT

	for _, a := range as {
		if f(a) {
			na = append(na, a)
		}
	}

	return na
}

// returns the name of the type represented, if it is a pointer & is prefixed
func NiceTypeName(t reflect.Type) string {
	if t == nil {
		return "nil"
	}

	if t.Kind() != reflect.Pointer {
		return t.Name()
	}

	t = t.Elem()
	return fmt.Sprintf("&%s", t.Name())
}
