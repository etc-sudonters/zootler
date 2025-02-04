package main

import (
	"fmt"
	"sudonters/zootler/cmd/zootler/bootstrap"
	"sudonters/zootler/cmd/zootler/z16"
	"sudonters/zootler/internal/settings"
	"sudonters/zootler/internal/skelly/bitset32"
	"sudonters/zootler/magicbean"
	"sudonters/zootler/mido"
	"sudonters/zootler/mido/objects"
	"sudonters/zootler/zecs"
)

type Age bool

const AgeAdult Age = true
const AgeChild Age = false

func fromStartingAge(start settings.StartingAge) Age {
	switch start {
	case settings.StartAgeAdult:
		return AgeAdult
	case settings.StartAgeChild:
		return AgeChild
	default:
		panic("unknown starting age")

	}
}

func explore(artifacts *Artifacts, age Age, these *settings.Zootr) {
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

	pockets := magicbean.NewPockets(&artifacts.Inventory, &artifacts.Ocm)

	var shuffleFlags magicbean.ShuffleFlags
	if these.Shuffling.OcarinaNotes {
		shuffleFlags = shuffleFlags | magicbean.SHUFFLE_OCARINA_NOTES
	}

	funcs := magicbean.BuiltIns{}
	magicbean.CreateBuiltInHasFuncs(&funcs, &pockets, shuffleFlags)
	funcs.CheckTodAccess = magicbean.ConstBool(true)
	funcs.IsAdult = magicbean.ConstBool(age == AgeAdult)
	funcs.IsChild = magicbean.ConstBool(age == AgeChild)
	funcs.IsStartingAge = magicbean.ConstBool(age == fromStartingAge(these.Spawns.StartingAge))

	vm := mido.VM{
		Objects: &artifacts.Objects,
		Funcs:   funcs.Table(),
	}

	workset := bitset32.Copy(roots)

	artifacts.World.ExploreAvailableEdges(magicbean.Exploration{
		Workset: &workset,
		Visited: &bitset32.Bitset{},
		VM:      vm,
		Objects: &artifacts.Objects,
	})
}

func PtrsMatching(ocm *zecs.Ocm, query ...zecs.BuildQuery) []objects.Object {
	q := ocm.Query()
	q.Build(zecs.Load[magicbean.Ptr], zecs.With[magicbean.Token])
	rows, err := q.Execute()
	bootstrap.PanicWhenErr(err)
	ptrs := make([]objects.Object, 0, rows.Len())

	for _, tup := range rows.All {
		ptr := tup.Values[0].(magicbean.Ptr)
		ptrs = append(ptrs, objects.Object(ptr))
	}

	return ptrs
}

func CollectStartingItems(artifacts *Artifacts, these *settings.Zootr) {
	tokens := z16.NewTokens(&artifacts.Ocm)

	if these.Locations.OpenDoorOfTime {
		timeTravel := tokens.MustGet("Time Travel")
		artifacts.Inventory.CollectOne(timeTravel.Entity())
		fmt.Printf("Time Travel ID: %x\n", timeTravel.Entity())
	}
}
