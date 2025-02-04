package main

import (
	"context"
	"fmt"
	"sudonters/zootler/cmd/zootler/bootstrap"
	"sudonters/zootler/cmd/zootler/z16"
	"sudonters/zootler/internal/settings"
	"sudonters/zootler/internal/skelly/bitset32"
	"sudonters/zootler/magicbean"
	"sudonters/zootler/mido"
	"sudonters/zootler/mido/objects"
	"sudonters/zootler/zecs"

	"github.com/etc-sudonters/substrate/dontio"
	"golang.org/x/text/collate"
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

func explore(ctx context.Context, artifacts *Artifacts, age Age, these *settings.Zootr) {
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

	std, noStd := dontio.StdFromContext(ctx)
	if noStd != nil {
		panic("no std found in context")
	}

	vm := mido.VM{
		Objects: &artifacts.Objects,
		Funcs:   funcs.Table(),
		Std:     std,
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
	type collecting struct {
		entity zecs.Entity
		qty    float64
	}
	var starting []collecting

	collect := func(token z16.Token, qty float64) {
		starting = append(starting, collecting{token.Entity(), qty})
	}

	collectOneEach := func(token ...z16.Token) {
		new := make([]collecting, len(starting)+len(token))
		copy(new[len(token):], starting)
		for i, t := range token {
			new[i] = collecting{t.Entity(), 1}
		}

		starting = new
	}

	tokens := z16.NewTokens(&artifacts.Ocm)

	if these.Locations.OpenDoorOfTime {
		collect(tokens.MustGet("Time Travel"), 1)
	}

	collectOneEach(
		tokens.MustGet("Ocarina"),
		tokens.MustGet("Deku Stick (1)"),
		tokens.MustGet("Deku Shield"),
	)

	for _, collect := range starting {
		artifacts.Inventory.Collect(collect.entity, collate.qty)
	}
}
