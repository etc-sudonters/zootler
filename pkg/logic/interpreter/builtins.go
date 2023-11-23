package interpreter

import (
	"sudonters/zootler/internal/entity"
	"sudonters/zootler/pkg/world/components"

	"github.com/etc-sudonters/substrate/mirrors"
	"github.com/etc-sudonters/substrate/skelly/hashset"
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
}

func (z Zoot_HasQuantityOf) Call(t Interpreter, args []Value) Value {
	token := args[0].(Token)
	qty := int(args[1].(Number).Value)

	filter := entity.FilterBuilder{}.
		With(mirrors.TypeOf[components.Collected]()).
		With(token.Component).
		Build()

	ents, err := z.Entities.Query(filter)

	if err != nil {
		panic(err)
	}

	return Box(qty <= len(ents))
}

type Zoot_HasMedallions struct {
	Has Zoot_HasQuantityOf
}

func (z Zoot_HasMedallions) Call(t Interpreter, args []Value) Value {
	return z.Has.Call(t, []Value{
		Token{
			Component: mirrors.TypeOf[components.Medallion](),
			Literal:   "",
		},
		args[0],
	})
}

type Zoot_RegionHasShortcuts struct {
	RegionalShortcuts hashset.Hash[string]
}

func (z Zoot_RegionHasShortcuts) Call(t Interpreter, args []Value) Value {
	region := args[0].(String)
	return Box(z.RegionalShortcuts.Exists(region.Value))
}

type Zoot_HasBottle struct {
	Has Zoot_HasQuantityOf
}

func (z Zoot_HasBottle) Call(t Interpreter, args []Value) Value {
	return z.Has.Call(
		t, []Value{
			Token{
				Component: mirrors.TypeOf[components.Bottle](),
			},
			Box(1),
		})
}

type Zoot_HasAnyOf struct{}
type Zoot_HasAllOf struct{}
type Zoot_CountOf struct{}
type Zoot_HeartCount struct{}
type Zoot_HasHearts struct{}
type Zoot_HasStones struct{}
type Zoot_HasDungeonRewards struct{}
type Zoot_HasOcarinaButtons struct{}
type Zoot_HasItemGoal struct{}
type Zoot_ItemCount struct{}
type Zoot_ItemNameCount struct{}
type Zoot_HasFullItemGoal struct{}
type Zoot_HasAllItemGoals struct{}
type Zoot_HadNightStart struct{}
type Zoot_CanLiveDmg struct{}
type Zoot_GuaranteeHint struct{}
type Zoot_HasNotesForSong struct{}
