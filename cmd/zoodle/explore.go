package main

import (
	"context"
	"io"
	"sudonters/libzootr/magicbean"
	"sudonters/libzootr/mido"

	"github.com/etc-sudonters/substrate/dontio"
)

type Age bool

const AgeAdult Age = true
const AgeChild Age = false

func explore(ctx context.Context, xplr *magicbean.Exploration, generation *magicbean.Generation, age Age) magicbean.ExplorationResults {
	pockets := magicbean.NewPockets(generation.Inventory, &generation.Ocm)

	funcs := magicbean.BuiltIns{}
	magicbean.CreateBuiltInHasFuncs(&funcs, &pockets, generation.Settings.Logic.Shuffling.Flags)
	funcs.CheckTodAccess = magicbean.ConstBool(true)
	funcs.IsAdult = magicbean.ConstBool(age == AgeAdult)
	funcs.IsChild = magicbean.ConstBool(age == AgeChild)
	funcs.IsStartingAge = magicbean.ConstBool(age == Age(generation.Settings.Logic.Spawns.StartAge))

	vm := mido.VM{
		Objects: &generation.Objects,
		Funcs:   funcs.Table(),
		Std: &dontio.Std{
			Out: io.Discard,
			Err: io.Discard,
			In:  eof{},
		},
		ChkQty: funcs.Has,
	}

	xplr.VM = vm
	xplr.Objects = &generation.Objects

	return generation.World.ExploreAvailableEdges(ctx, xplr)
}

type eof struct{}

func (_ eof) Read([]byte) (int, error) {
	return 0, io.EOF
}
