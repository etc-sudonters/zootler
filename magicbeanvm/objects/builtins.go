package objects

import "fmt"

func BuiltInFunctionNames() []string {
	names := make([]string, 0, len(builtInIndex))
	for name := range builtInIndex {
		names = append(names, name)
	}
	return names
}

func GetBuiltInIndex(name string) (index Index, exists bool) {
	index, exists = builtInIndex[name]
	return
}

var builtInIndex = map[string]Index{
	"at_dampe_time":          Index(0),
	"at_day":                 Index(1),
	"at_night":               Index(2),
	"has":                    Index(3),
	"has_all_notes_for_song": Index(4),
	"has_anyof":              Index(5),
	"has_bottle":             Index(6),
	"has_dungeon_rewards":    Index(7),
	"has_every":              Index(8),
	"has_hearts":             Index(9),
	"has_medallions":         Index(10),
	"has_stones":             Index(11),
	"is_starting_age":        Index(12),
	"region_has_shortcuts":   Index(13),
	"is_adult":               Index(14),
	"is_child":               Index(15),
}

func NewBuiltins(tbl BuiltInFunctionTable) BuiltIns {
	var funcs BuiltIns
	if tbl == nil {
		panic("nil function table")
	}

	funcs.tbl = tbl
	funcs.objs = make([]BuiltInFunc, len(builtInIndex))
	for name, index := range builtInIndex {
		var fn Callable
		switch name {
		case "at_dampe_time":
			fn = tbl.AtDampeTime
		case "at_day":
			fn = tbl.AtDay
		case "at_night":
			fn = tbl.AtNight
		case "has":
			fn = tbl.Has
		case "has_all_notes_for_song":
			fn = tbl.HasAllNotesForSong
		case "has_anyof":
			fn = tbl.HasAnyOf
		case "has_bottle":
			fn = tbl.HasBottle
		case "has_dungeon_rewards":
			fn = tbl.HasDungeonRewards
		case "has_every":
			fn = tbl.HasEvery
		case "has_hearts":
			fn = tbl.HasHearts
		case "has_medallions":
			fn = tbl.HasMedallions
		case "has_stones":
			fn = tbl.HasStones
		case "is_adult":
			fn = tbl.IsAdult
		case "is_child":
			fn = tbl.IsChild
		case "is_starting_age":
			fn = tbl.IsStartingAge
		case "region_has_shortcuts":
			fn = tbl.RegionHasShortcuts
		default:
			panic(fmt.Errorf("unknown built in name %q", name))
		}

		funcs.objs[int(index)] = BuiltInFunc{
			Name: name,
			Func: fn,
		}
	}

	return funcs
}

type BuiltIns struct {
	tbl  BuiltInFunctionTable
	objs []BuiltInFunc
}

func (this *BuiltIns) Get(idx Index) *BuiltInFunc {
	fn := &this.objs[int(idx)]
	return fn
}

type BuiltInFunctionTable interface {
	AtDampeTime(...Object) (Object, error)
	AtDay(...Object) (Object, error)
	AtNight(...Object) (Object, error)
	Has(...Object) (Object, error)
	HasAllNotesForSong(...Object) (Object, error)
	HasAnyOf(...Object) (Object, error)
	HasBottle(...Object) (Object, error)
	HasDungeonRewards(...Object) (Object, error)
	HasEvery(...Object) (Object, error)
	HasHearts(...Object) (Object, error)
	HasMedallions(...Object) (Object, error)
	HasStones(...Object) (Object, error)
	IsAdult(...Object) (Object, error)
	IsChild(...Object) (Object, error)
	IsStartingAge(...Object) (Object, error)
	RegionHasShortcuts(...Object) (Object, error)
}
