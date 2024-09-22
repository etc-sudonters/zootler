package zasm

import (
	"sudonters/zootler/icearrow/nan"
	"sudonters/zootler/internal"
	"sudonters/zootler/internal/intern"
)

type Data struct {
	Strs   []string
	Consts []nan.PackedValue
	Names  []string
}

type DataBuilder struct {
	Strs   intern.HashIntern[string]
	Consts intern.HashIntern[nan.PackedValue]
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
	db.Consts = intern.NewInterner[nan.PackedValue]()
	db.Names = intern.NewInternerF(
		func(s string) string { return string(internal.Normalize(s)) },
	)
	return db
}
