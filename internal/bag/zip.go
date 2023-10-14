package bag

type ziptwo[A any, B any] struct {
	a   []A
	b   []B
	cur int
	k   int
}

type zippedtwo[A any, B any] struct {
	A A
	B B
}

func ZipTwo[A any, B any](a []A, b []B) *ziptwo[A, B] {
	return &ziptwo[A, B]{
		a:   a,
		b:   b,
		cur: -1,
		k:   Min(len(a), len(b)),
	}
}

func (z *ziptwo[A, B]) Next() *zippedtwo[A, B] {
	if z.cur+1 >= z.k {
		return nil
	}
	z.cur += 1
	return &zippedtwo[A, B]{A: z.a[z.cur], B: z.b[z.cur]}
}
