package vm

import "sudonters/zootler/magicbeanvm/objects"

func BuiltInFunctionNames() []string {
	return builtInNames[:]
}

func GlobalNames() []string {
	return globalNames[:]
}

var globalNames = []string{
	"age",
	"Forest",
	"Fire",
	"Water",
	"Shadow",
	"Spirit",
	"Light",
}

var builtInNames = []string{
	"at_dampe_time",
	"at_day",
	"at_night",
	"has",
	"has_all_notes_for_song",
	"has_anyof",
	"has_bottle",
	"has_dungeon_rewards",
	"has_every",
	"has_hearts",
	"has_medallions",
	"has_stones",
	"is_adult",
	"is_child",
	"is_starting_age",
	"region_has_shortcuts",
}

type BuiltInFunctionTable interface {
	AtDampeTime(...objects.Object) (objects.Object, error)
	AtDay(...objects.Object) (objects.Object, error)
	AtNight(...objects.Object) (objects.Object, error)
	Has(...objects.Object) (objects.Object, error)
	HasAllNotesForSong(...objects.Object) (objects.Object, error)
	HasAnyOf(...objects.Object) (objects.Object, error)
	HasBottle(...objects.Object) (objects.Object, error)
	HasDungeonRewards(...objects.Object) (objects.Object, error)
	HasEvery(...objects.Object) (objects.Object, error)
	HasHearts(...objects.Object) (objects.Object, error)
	HasMedallions(...objects.Object) (objects.Object, error)
	HasStones(...objects.Object) (objects.Object, error)
	IsAdult(...objects.Object) (objects.Object, error)
	IsChild(...objects.Object) (objects.Object, error)
	IsStartingAge(...objects.Object) (objects.Object, error)
	RegionHasShortcuts(...objects.Object) (objects.Object, error)
}

type Runner struct {
	builtins BuiltInFunctionTable
}

func New(builtins BuiltInFunctionTable) Runner {
	var r Runner
	r.builtins = builtins
	return r
}
