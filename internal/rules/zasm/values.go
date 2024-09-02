package zasm

import (
	"math"
)

type Packable interface {
	float64 | float32 |
		uint32 | uint16 | uint8 |
		int32 | int16 | int8
}

type Value float64 // nanpacked

func Pack[P Packable](v P) Value {
	return Value(math.NaN())
}
