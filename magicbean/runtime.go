package magicbean

import "sudonters/zootler/mido/objects"

func ConstBool(b bool) objects.BuiltInFunction {
	obj := objects.PackedTrue
	if !b {
		obj = objects.PackedFalse
	}
	return func(*objects.Table, []objects.Object) (objects.Object, error) {
		return obj, nil
	}
}
