package magicbean

import (
	"fmt"
	"sudonters/zootler/mido/objects"
	"sudonters/zootler/zecs"
)

func ConstBool(b bool) objects.BuiltInFunction {
	obj := objects.PackedTrue
	if !b {
		obj = objects.PackedFalse
	}
	return func(*objects.Table, []objects.Object) (objects.Object, error) {
		return obj, nil
	}
}

type BuiltIns struct {
	CheckTodAccess    objects.BuiltInFunction `zootler:"check_tod_access,params=1"`
	Has               objects.BuiltInFunction `zootler:"has,params=2"`
	HasAnyOf          objects.BuiltInFunction `zootler:"has_anyof,params=-1"`
	HasBottle         objects.BuiltInFunction `zootler:"has_bottle,params=0"`
	HasDungeonRewards objects.BuiltInFunction `zootler:"has_dungeon_rewards,params=1"`
	HasEvery          objects.BuiltInFunction `zootler:"has_every,params=-1"`
	HasHearts         objects.BuiltInFunction `zootler:"has_hearts,params=1"`
	HasMedallions     objects.BuiltInFunction `zootler:"has_medallions,params=1"`
	HasNotesForSong   objects.BuiltInFunction `zootler:"has_notes_for_song,params=1"`
	HasStones         objects.BuiltInFunction `zootler:"has_stones,params=1"`
	IsAdult           objects.BuiltInFunction `zootler:"is_adult,params=0"`
	IsChild           objects.BuiltInFunction `zootler:"is_child,params=0"`
	IsStartingAge     objects.BuiltInFunction `zootler:"is_starting_age,params=0"`
}

func (this BuiltIns) Table() objects.BuiltInFunctions {
	return objects.BuiltInFunctions{
		this.CheckTodAccess,
		this.Has,
		this.HasAnyOf,
		this.HasBottle,
		this.HasDungeonRewards,
		this.HasEvery,
		this.HasHearts,
		this.HasMedallions,
		this.HasNotesForSong,
		this.HasStones,
		this.IsAdult,
		this.IsChild,
		this.IsStartingAge,
	}
}

func CreateBuiltInDefs() []objects.BuiltInFunctionDef {
	return []objects.BuiltInFunctionDef{
		{Name: "check_tod_access", Params: 1},
		{Name: "has", Params: 2},
		{Name: "has_anyof", Params: -1},
		{Name: "has_bottle", Params: 0},
		{Name: "has_dungeon_rewards", Params: 1},
		{Name: "has_every", Params: -1},
		{Name: "has_hearts", Params: 1},
		{Name: "has_medallions", Params: 1},
		{Name: "has_notes_for_song", Params: 1},
		{Name: "has_stones", Params: 1},
		{Name: "is_adult", Params: 0},
		{Name: "is_child", Params: 0},
		{Name: "is_starting_age", Params: 0},
	}
}

type ShuffleFlags uint64

const (
	SHUFFLE_OCARINA_NOTES = 1
)

func CreateBuiltInHasFuncs(builtins *BuiltIns, pocket *Pocket, flags ShuffleFlags) {
	builtins.Has = func(tbl *objects.Table, args []objects.Object) (objects.Object, error) {
		if len(args) != 2 {
			return objects.Null, fmt.Errorf("has expects 2 arguments, got %d", len(args))
		}

		target := args[0]
		qty := args[1]

		if !target.Is(objects.Ptr32) {
			return objects.Null, fmt.Errorf("has expects token ptr as first argument, got: %q", target.Type())
		}

		tag, item := objects.UnpackPtr32(target)
		if tag != objects.PtrToken {
			return objects.Null, fmt.Errorf("has expects token ptr (%X) as first argument, tag: %X", objects.PtrToken, tag)
		}

		var n uint64

		switch {
		case qty.Is(objects.F64):
			n = uint64(objects.UnpackF64(qty))
		case qty.Is(objects.U32):
			n = uint64(objects.UnpackU32(qty))
		case qty.Is(objects.I32):
			n = uint64(int64(objects.UnpackI32(qty)))
		default:
			return objects.Null, fmt.Errorf("has expects number as second argument, got %q", qty.Type())
		}

		result := pocket.Has(zecs.Entity(item), uint64(n))
		return objects.PackBool(result), nil
	}

	builtins.HasAnyOf = func(_ *objects.Table, args []objects.Object) (objects.Object, error) {
		items := make([]zecs.Entity, len(args))
		for i, arg := range args {
			if !objects.IsPtrWithTag(arg, objects.PtrToken) {
				return objects.Null, fmt.Errorf("has_anyof expects all args to be ptrs, %d was not", i+1)
			}

			_, item := objects.UnpackPtr32(arg)
			items[i] = zecs.Entity(item)
		}

		result := pocket.HasAny(items)
		return objects.PackBool(result), nil
	}

	builtins.HasEvery = func(_ *objects.Table, args []objects.Object) (objects.Object, error) {
		items := make([]zecs.Entity, len(args))
		for i, arg := range args {
			if !objects.IsPtrWithTag(arg, objects.PtrToken) {
				return objects.Null, fmt.Errorf("has_anyof expects all args to be ptrs, %d was not", i+1)
			}

			_, item := objects.UnpackPtr32(arg)
			items[i] = zecs.Entity(item)
		}

		result := pocket.HasEvery(items)
		return objects.PackBool(result), nil
	}

	builtins.HasBottle = func(_ *objects.Table, _ []objects.Object) (objects.Object, error) {
		return objects.PackBool(pocket.HasBottle()), nil
	}

	builtins.HasDungeonRewards = func(_ *objects.Table, args []objects.Object) (objects.Object, error) {
		qty := args[0]
		var n uint64
		switch {
		case qty.Is(objects.F64):
			n = uint64(objects.UnpackF64(qty))
		case qty.Is(objects.U32):
			n = uint64(objects.UnpackU32(qty))
		case qty.Is(objects.I32):
			n = uint64(int64(objects.UnpackI32(qty)))
		default:
			return objects.Null, fmt.Errorf("has_dungeon_rewards expects number as first argument")
		}

		return objects.PackBool(pocket.HasDungeonRewards(n)), nil
	}

	builtins.HasHearts = func(_ *objects.Table, args []objects.Object) (objects.Object, error) {
		qty := args[0]
		var n uint64
		switch {
		case qty.Is(objects.F64):
			n = uint64(objects.UnpackF64(qty))
		case qty.Is(objects.U32):
			n = uint64(objects.UnpackU32(qty))
		case qty.Is(objects.I32):
			n = uint64(int64(objects.UnpackI32(qty)))
		default:
			return objects.Null, fmt.Errorf("has_hearts expects number as first argument")
		}

		return objects.PackBool(pocket.HasHearts(n)), nil
	}

	builtins.HasStones = func(_ *objects.Table, args []objects.Object) (objects.Object, error) {
		qty := args[0]
		var n uint64
		switch {
		case qty.Is(objects.F64):
			n = uint64(objects.UnpackF64(qty))
		case qty.Is(objects.U32):
			n = uint64(objects.UnpackU32(qty))
		case qty.Is(objects.I32):
			n = uint64(int64(objects.UnpackI32(qty)))
		default:
			return objects.Null, fmt.Errorf("has_stones expects number as first argument")
		}

		return objects.PackBool(pocket.HasStones(n)), nil
	}

	builtins.HasMedallions = func(_ *objects.Table, args []objects.Object) (objects.Object, error) {
		qty := args[0]
		var n uint64
		switch {
		case qty.Is(objects.F64):
			n = uint64(objects.UnpackF64(qty))
		case qty.Is(objects.U32):
			n = uint64(objects.UnpackU32(qty))
		case qty.Is(objects.I32):
			n = uint64(int64(objects.UnpackI32(qty)))
		default:
			return objects.Null, fmt.Errorf("has_medallions expects number as first argument")
		}

		return objects.PackBool(pocket.HasMedallions(n)), nil
	}

	if flags&SHUFFLE_OCARINA_NOTES == SHUFFLE_OCARINA_NOTES {
		builtins.HasNotesForSong = func(_ *objects.Table, args []objects.Object) (objects.Object, error) {
			if !objects.IsPtrWithTag(args[0], objects.PtrToken) {
				return objects.Null, fmt.Errorf("has_notes_for_song expects song ptr as argument")
			}

			_, song := objects.UnpackPtr32(args[0])
			return objects.PackBool(pocket.HasAllNotes(zecs.Entity(song))), nil
		}
	} else {
		builtins.HasNotesForSong = ConstBool(true)
	}
}
