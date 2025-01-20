package main

import (
	"sudonters/zootler/internal/skelly/bitset32"
	"sudonters/zootler/magicbean"
	"sudonters/zootler/mido"
	"sudonters/zootler/mido/objects"
	"sudonters/zootler/zecs"
)

func consttrue(*objects.Table, []objects.Object) (objects.Object, error) {
	return objects.PackedTrue, nil
}

func explore(artifacts *Artifacts) {
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

	vm := mido.VM{
		Objects: &artifacts.Objects,
		Funcs: objects.BuiltInFunctions{
			consttrue,
			consttrue,
			consttrue,
			consttrue,
			consttrue,
			consttrue,
			consttrue,
			consttrue,
			consttrue,
			consttrue,
			consttrue,
		},
	}

	for range 10 {
		artifacts.World.ExploreAvailableEdges(magicbean.Exploration{
			Workset: bitset32.Copy(roots),
			Visited: bitset32.Bitset{},
			VM:      vm,
		})
	}
}
