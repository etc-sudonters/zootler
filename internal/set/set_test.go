package set

import (
	"testing"
)

func TestDifference(t *testing.T) {
	fives := FromSlice([]int{5, 10, 15, 20, 25, 30})
	threes := FromSlice([]int{3, 6, 9, 12, 15, 18, 21, 24, 27, 30})

	actual := Difference(fives, threes)

	if !Equal(actual, FromSlice([]int{5, 10, 20, 25})) {
		t.Fail()
		t.Logf("unexpected set difference %+v", actual)
	}
}

func TestIntersection(t *testing.T) {
	fives := FromSlice([]int{5, 10, 15, 20, 25, 30})
	threes := FromSlice([]int{3, 6, 9, 12, 15, 18, 21, 24, 27, 30})

	actual := Intersection(fives, threes)

	if !Equal(actual, FromSlice([]int{15, 30})) {
		t.Fail()
		t.Logf("unexpected set difference %+v", actual)
	}
}

func TestUnion(t *testing.T) {
	fives := FromSlice([]int{5, 10, 15, 20, 25, 30})
	threes := FromSlice([]int{3, 6, 9, 12, 15, 18, 21, 24, 27, 30})

	actual := Union(fives, threes)

	if !Equal(actual, FromSlice([]int{3, 5, 6, 9, 10, 12, 15, 18, 20, 21, 24, 25, 27, 30})) {
		t.Fail()
		t.Logf("unexpected set difference %+v", actual)
	}
}
