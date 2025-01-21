package magicbean

import (
	"sudonters/zootler/mido/objects"
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
	CheckTodAccess    objects.BuiltInFunction `zootler:"check_tod_access,params=0"`
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
		{Name: "check_tod_access", Params: 0},
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
