package testutils

import (
	"testing"

	"github.com/etc-sudonters/zootler/internal/bag"
)

func ArrEq[U comparable, S ~[]U, T ~[]U](expected S, actual T, t *testing.T) {
	ArrEqF(expected, actual, func(exp, act U) bool { return exp == act }, t)
}

func ArrEqF[U any, S ~[]U, T ~[]U](expected S, actual T, eqBy func(U, U) bool, t *testing.T) {
	if len(expected) != len(actual) {
		t.Fail()
		t.Logf("mismatched lengths\nexpected\t%d\nactual\t\t%d", len(expected), len(actual))
	}

	min := bag.Min(len(expected), len(actual))

	for i := 0; i < min; i++ {
		if !eqBy(expected[i], actual[i]) {
			t.Fail()
			t.Logf("mismatched at index %d\nexpected\t%+v\nactual\t\t%+v", i, expected[i], actual[i])
		}
	}

	if len(expected) > min {
		t.Logf("trailing values in expected: %+v", expected[min:])
	}

	if len(actual) > min {
		t.Logf("trailing values in actual: %+v", actual[min:])
	}

}
