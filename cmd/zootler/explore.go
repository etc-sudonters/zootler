package main

import (
	"fmt"
	"sudonters/zootler/internal/skelly/bitset32"
	"sudonters/zootler/magicbean"
	"sudonters/zootler/mido"
	"sudonters/zootler/mido/objects"
	"sudonters/zootler/zecs"
)

type Age bool

const AgeAdult Age = true
const AgeChild Age = false

func explore(artifacts *Artifacts, age Age) {
	q := artifacts.Ocm.Query()
	q.Build(zecs.With[magicbean.Region], zecs.Load[magicbean.Name])
	roots := bitset32.Bitset{}

	for ent, tup := range q.Rows {
		name := tup.Values[0].(magicbean.Name)
		if name == "Root" {
			bitset32.Set(&roots, ent)
			break
		}
	}

	qty := magicbean.QtyBuiltInFunctions{Pocket: magicbean.EmptyPockets()}
	isAdult := magicbean.ConstBool(age == AgeAdult)
	isChild := magicbean.ConstBool(age == AgeChild)
	isStartingAge := magicbean.ConstBool(false)
	hasBottle := magicbean.ConstBool(false)
	hasDungeonRewards := magicbean.ConstBool(false)
	hasHearts := magicbean.ConstBool(false)
	hasMeds := magicbean.ConstBool(false)
	hasStones := magicbean.ConstBool(false)

	vm := mido.VM{
		Objects: &artifacts.Objects,
		Funcs: objects.BuiltInFunctions{
			qty.Has,
			qty.HasAnyOf,
			qty.HasEvery,
			isAdult,
			isChild,
			hasBottle,
			hasDungeonRewards,
			hasHearts,
			hasMeds,
			hasStones,
			isStartingAge,
		},
	}

	workset := bitset32.Copy(roots)

	for range 10 {
		if workset.Len() == 0 {
			fmt.Printf("empty exploration set, exiting explore loop")
			break
		}

		result := artifacts.World.ExploreAvailableEdges(magicbean.Exploration{
			Workset: &workset,
			Visited: &bitset32.Bitset{},
			VM:      vm,
		})

		if result.Reached.Len() == 0 {
			fmt.Printf("no progress made in exploration, exiting")
			break
		}

		workset = result.Workset
	}
}
