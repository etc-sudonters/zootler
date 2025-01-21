package main

import (
	"sudonters/zootler/cmd/zootler/bootstrap"
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

	qty := magicbean.QtyBuiltInFunctions{
		Pocket:         magicbean.EmptyPockets(),
		OcarinaButtons: PtrsMatching(&artifacts.Ocm, zecs.With[magicbean.OcarinaButton]),
		HeartPieces:    HeartPieces(&artifacts.Ocm),
		Bottles:        PtrsMatching(&artifacts.Ocm, zecs.With[magicbean.Bottle]),
		Medallions:     PtrsMatching(&artifacts.Ocm, zecs.With[magicbean.Medallion]),
		DungeonRewards: PtrsMatching(&artifacts.Ocm, zecs.With[magicbean.DungeonReward]),
		Stones:         PtrsMatching(&artifacts.Ocm, zecs.With[magicbean.Stone]),
	}

	hasNotesForSong := magicbean.ConstBool(true)

	funcs := magicbean.BuiltIns{
		CheckTodAccess:    magicbean.ConstBool(true),
		Has:               qty.Has,
		HasAnyOf:          qty.HasAnyOf,
		HasBottle:         qty.HasBottle,
		HasDungeonRewards: qty.HasDungeonRewards,
		HasEvery:          qty.HasEvery,
		HasHearts:         qty.HasHearts,
		HasMedallions:     qty.HasMedallions,
		HasNotesForSong:   hasNotesForSong,
		HasStones:         qty.HasStones,
		IsAdult:           magicbean.ConstBool(age == AgeAdult),
		IsChild:           magicbean.ConstBool(age == AgeChild),
		IsStartingAge:     magicbean.ConstBool(age == fromStartingAge(these.Spawns.StartingAge)),
	}

	vm := mido.VM{
		Objects: &artifacts.Objects,
		Funcs:   funcs.Table(),
	}

	workset := bitset32.Copy(roots)

	artifacts.World.ExploreAvailableEdges(magicbean.Exploration{
		Workset: &workset,
		Visited: &bitset32.Bitset{},
		VM:      vm,
	})
}

func PtrsMatching(ocm *zecs.Ocm, query ...zecs.BuildQuery) []objects.Object {
	q := ocm.Query()
	q.Build(zecs.Load[magicbean.Ptr], zecs.With[magicbean.Bottle], zecs.With[magicbean.Token])
	rows, err := q.Execute()
	bootstrap.PanicWhenErr(err)
	ptrs := make([]objects.Object, 0, rows.Len())

	for _, tup := range rows.All {
		ptr := tup.Values[0].(magicbean.Ptr)
		ptrs = append(ptrs, objects.Object(ptr))
	}

	return ptrs
}

func HeartPieces(ocm *zecs.Ocm) map[objects.Object]magicbean.HeartPieceCount {
	q := ocm.Query()
	q.Build(zecs.Load[magicbean.Ptr], zecs.Load[magicbean.HeartPieceCount], zecs.With[magicbean.Token])
	rows, err := q.Execute()
	bootstrap.PanicWhenErr(err)
	ptrs := make(map[objects.Object]magicbean.HeartPieceCount, rows.Len())

	for _, tup := range rows.All {
		ptr := tup.Values[0].(magicbean.Ptr)
		hp := tup.Values[1].(magicbean.HeartPieceCount)
		ptrs[objects.Object(ptr)] = hp
	}

	return ptrs
}
