package zasm

import (
	"sudonters/zootler/internal"
	"sudonters/zootler/internal/intern"
)

type PackedValue float64

func (pv PackedValue) Equals(p PackedValue) bool {
	return p == pv
}

type Packable interface {
	~uint32 | ~uint16 | ~uint8 |
		~int32 | ~int16 | ~int8 |
		~float64 | ~float32
}

type Data struct {
	Strs   []string
	Consts []PackedValue
	Names  []string
}

type DataBuilder struct {
	Strs   intern.HashIntern[string]
	Consts intern.HashIntern[PackedValue]
	Names  intern.HashInternF[string, string]
}

func CreateDataTables(db DataBuilder) Data {
	var d Data
	d.Strs = db.Strs.IntoTable()
	d.Consts = db.Consts.IntoTable()
	d.Names = db.Names.IntoTable()
	return d

}

func NewDataBuilder() DataBuilder {
	var db DataBuilder
	db.Strs = intern.NewInterner[string]()
	db.Consts = intern.NewInterner[PackedValue]()
	db.Names = intern.NewInternerF(
		func(s string) string { return string(internal.Normalize(s)) },
	)
	return db
}

func Pack[P Packable](p P) PackedValue {
	return PackedValue(float64(p))
}
