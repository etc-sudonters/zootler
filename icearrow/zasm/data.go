package zasm

import (
	"sudonters/zootler/internal"
	"sudonters/zootler/internal/intern"
)

type PackedValue float64

type Packable interface {
	~uint32 | ~uint16 | ~uint8 |
		~int32 | ~int16 | ~int8 |
		~float64 | ~float32
}

type Data struct {
	Strs   []uint8
	Consts []PackedValue
	Names  []string
}

type DataBuilder struct {
	Strs   intern.StrHeaper
	Consts intern.HashIntern[PackedValue]
	Names  intern.HashInternF[string, string]
}

func NewDataBuilder() DataBuilder {
	var db DataBuilder
	db.Strs = intern.NewStrHeaper()
	db.Consts = intern.NewInterner[PackedValue]()
	db.Names = intern.NewInternerF(
		func(s string) string { return string(internal.Normalize(s)) },
	)
	return db
}

func Pack[P Packable](p P) PackedValue {
	return PackedValue(float64(p))
}
