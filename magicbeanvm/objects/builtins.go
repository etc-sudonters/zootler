package objects

import "fmt"

func BuiltInFunctionNames() []string {
	names := make([]string, 0, len(builtInIndex))
	for name := range builtInIndex {
		names = append(names, name)
	}
	return names
}

func GetBuiltInIndex(name string) (Index, bool) {
	def, exists := builtInIndex[name]
	return def.index, exists
}

type builtindef struct {
	name   string
	index  Index
	params int
}

var builtInIndex = map[string]builtindef{
	// these functions have dedicated op codes and are almost always invoked
	// via those codes rather than being loaded into the stack like other calls
	"has":       {"has", Index(0), 2},
	"has_anyof": {"has_anyof", Index(1), -1},
	"has_every": {"has_every", Index(2), -1},
	"is_adult":  {"is_adult", Index(3), 0},
	"is_child":  {"is_child", Index(4), 0},

	// TODO evaluate these at compile time b/c they depend on settings (shuffle
	// notes, shuffle entrances) and can be turned into constants or replaced
	// w/ a new builtin
	"has_all_notes_for_song": {"has_all_notes_for_song", Index(5), 1},
	"at_dampe_time":          {"at_dampe_time", Index(6), 0},
	"at_day":                 {"at_day", Index(7), 0},
	"at_night":               {"at_night", Index(8), 0},

	// these are all invoked infrequently, with only has_bottle being noticable
	// but not often enough to warrant a dedicated op code
	"has_bottle":          {"has_bottle", Index(9), 0},
	"has_dungeon_rewards": {"has_dungeon_rewards", Index(10), 1},
	"has_hearts":          {"has_hearts", Index(11), 1},
	"has_medallions":      {"has_medallions", Index(12), 1},
	"has_stones":          {"has_stones", Index(13), 1},
	"is_starting_age":     {"is_starting_age", Index(14), 0},
}

func NewBuiltins(tbl BuiltInFunctionTable) BuiltInFunctions {
	var funcs BuiltInFunctions
	if tbl == nil {
		panic("nil function table")
	}

	funcs.BuiltInFunctionTable = tbl
	funcs.objs = make([]BuiltInFunc, len(builtInIndex))
	for _, def := range builtInIndex {
		var fn Callable
		switch def.name {
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
		default:
			panic(fmt.Errorf("unknown built in name %q", def.name))
		}

		funcs.objs[int(def.index)] = BuiltInFunc{
			Func:   fn,
			Name:   def.name,
			Params: def.params,
		}
	}

	return funcs
}

type BuiltInFunctions struct {
	BuiltInFunctionTable
	objs []BuiltInFunc
}

func (this *BuiltInFunctions) Get(idx Index) *BuiltInFunc {
	fn := &this.objs[int(idx)]
	return fn
}

type BuiltInFunctionTable interface {
	AtDampeTime([]Object) (Object, error)
	AtDay([]Object) (Object, error)
	AtNight([]Object) (Object, error)
	Has([]Object) (Object, error)
	HasAllNotesForSong([]Object) (Object, error)
	HasAnyOf([]Object) (Object, error)
	HasBottle([]Object) (Object, error)
	HasDungeonRewards([]Object) (Object, error)
	HasEvery([]Object) (Object, error)
	HasHearts([]Object) (Object, error)
	HasMedallions([]Object) (Object, error)
	HasStones([]Object) (Object, error)
	IsAdult([]Object) (Object, error)
	IsChild([]Object) (Object, error)
	IsStartingAge([]Object) (Object, error)
}
