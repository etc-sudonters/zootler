package runtime

import "math"

var mask uint64 = 0

type PackedValue float64

func (p PackedValue) AsBool() (bool, error)       { panic("not implemented") }
func (p PackedValue) AsFloat64() (float64, error) { panic("not implemented") }
func (p PackedValue) AsToken() (uint64, error)    { panic("not implemented") }

func IsRealNan(f float64) bool {
	panic("not implemented")
}

func PackBool(b bool) PackedValue {
	panic("not implemented")
}

func PackFloat(f float64) PackedValue {
	return PackedValue(f)
}

func PackSigil[T ~uint64](t T) PackedValue {
	panic("not implemented")
}

func PackShortStr(s string) (PackedValue, error) {
	panic("not implemented")
}

func bits(p PackedValue) uint64 {
	return math.Float64bits(float64(p))
}

func pack(u uint64) PackedValue {
	return PackedValue(math.Float64frombits(u))
}
