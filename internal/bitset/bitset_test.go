package bitset

import "testing"

func TestUnionBitSets(t *testing.T) {
	a := New(1)
	b := New(1)
	a.Set(1)
	b.Set(6)

	expected := New(1)
	expected.Set(1)
	expected.Set(6)

	if !expected.Eq(a.Union(b)) {
		t.Fatal("expected sets A and B to union")
	}
}

func TestBitSetIntersection(t *testing.T) {
	a := New(3)
	b := New(3)

	SetMany(&a, 1, 10, 16, 9, 99)
	SetMany(&b, 99, 16, 76, 108)

	expected := New(3)
	SetMany(&expected, 16, 99)

	if !expected.Eq(a.Intersect(b)) {
		t.Fatal("expected sets A and B to intersect")
	}
}

func TestSparseIntersection(t *testing.T) {
	fizz := New(10)
	buzz := New(10)
	fizzbuzz := New(10)

	for i := int64(0); i < 10*bs64Size; i++ {
		if i%3 == 0 {
			(&fizz).Set(i)
		}

		if i%5 == 0 {
			(&buzz).Set(i)
		}

		if fizz.Test(i) && buzz.Test(i) {
			(&fizzbuzz).Set(i)
		}
	}

	if !fizzbuzz.Eq(fizz.Intersect(buzz)) {
		t.Fatalf("expected big fizzbuzz data!")
	}
}
