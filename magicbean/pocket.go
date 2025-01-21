package magicbean

import (
	"fmt"
	"sudonters/zootler/mido/objects"
	"sudonters/zootler/zecs"
)

func EmptyPockets() Pocket {
	return Pocket{make(map[zecs.Entity]float64, 32)}
}

type Pocket struct {
	tokens map[zecs.Entity]float64
}

func (this Pocket) Quantity(object objects.Object) float64 {
	if !objects.IsPtrWithTag(object, objects.PtrToken) {
		panic(fmt.Errorf("%x is not a token pointer", object))
	}

	_, entity := objects.UnpackPtr32(object)
	return this.tokens[zecs.Entity(entity)]
}

func (this Pocket) Collect(object objects.Object, qty float64) {
	if !objects.IsPtrWithTag(object, objects.PtrToken) {
		panic(fmt.Errorf("%x is not a token pointer", object))
	}

	_, ent := objects.UnpackPtr32(object)
	entity := zecs.Entity(ent)

	already := this.tokens[entity]
	this.tokens[entity] = already + qty
}

type QtyBuiltInFunctions struct {
	Pocket Pocket
}

func (this QtyBuiltInFunctions) Has(_ *objects.Table, args []objects.Object) (objects.Object, error) {
	if !args[0].Is(objects.F64) {
		return objects.PackedFalse, nil
	}

	f64 := objects.UnpackF64(args[1])
	enough := this.Pocket.Quantity(args[0]) >= f64

	switch enough {
	case true:
		return objects.PackedTrue, nil
	default:
		return objects.PackedFalse, nil
	}
}

func (this QtyBuiltInFunctions) HasEvery(_ *objects.Table, args []objects.Object) (objects.Object, error) {

	for _, ptr := range args {
		if this.Pocket.Quantity(ptr) < 1 {
			return objects.PackedFalse, nil
		}
	}

	return objects.PackedTrue, nil
}

func (this QtyBuiltInFunctions) HasAnyOf(_ *objects.Table, args []objects.Object) (objects.Object, error) {
	for _, ptr := range args {
		if this.Pocket.Quantity(ptr) >= 1 {
			return objects.PackedTrue, nil
		}
	}

	return objects.PackedFalse, nil
}
