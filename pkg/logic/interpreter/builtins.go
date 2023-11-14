package interpreter

import (
	"sudonters/zootler/internal/entity"
	"sudonters/zootler/pkg/logic"
)

var (
	AtDay   BuiltInFn = atTod
	AtNigt  BuiltInFn = atTod
	AtDampe BuiltInFn = atTod
)

func atTod(_ Interpreter, _ []Value) Value {
	return Box(true)
}

// State.py
// ("item name", qty) tuples and "raw_item_name" w/ implicit qty = 1, having more is fine
type Zoot_HasQuantityOf struct {
	Entities entity.Queryable
	Selector logic.TypedStringSelector
}

func (z Zoot_HasQuantityOf) Call(t Interpreter, args []Value) Value {
	kind := args[0].(Token)
	qty := int(args[1].(Number).Value) // always safe

	ents, err := z.Entities.Query([]entity.Selector{
		entity.With[logic.Collected]{},
		z.Selector.With(kind.Literal),
	})

	if err != nil {
		panic(err)
	}

	return Box(qty <= len(ents))
}

type Zoot_HasAnyOf struct{}
type Zoot_HasAllOf struct{}
type Zoot_CountOf struct{}
type Zoot_HeartCount struct{}
type Zoot_HasHearts struct{}
type Zoot_HasMedallions struct{}
type Zoot_HasStones struct{}
type Zoot_HasDungeonRewards struct{}
type Zoot_HasOcarinaButtons struct{}
type Zoot_HasItemGoal struct{}
type Zoot_ItemCount struct{}
type Zoot_ItemNameCount struct{}
type Zoot_HasBottle struct{}
type Zoot_HasFullItemGoal struct{}
type Zoot_HasAllItemGoals struct{}
type Zoot_HadNightStart struct{}
type Zoot_CanLiveDmg struct{}
type Zoot_GuaranteeHint struct{}
type Zoot_RegionHasShortcuts struct{}
type Zoot_HasNotesForSong struct{}
