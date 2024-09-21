package main

import (
	"sudonters/zootler/internal/app"
	"sudonters/zootler/internal/entities"
	"sudonters/zootler/internal/world"

	"github.com/etc-sudonters/substrate/dontio"
	"github.com/etc-sudonters/substrate/skelly/bitset"
	"github.com/etc-sudonters/substrate/skelly/graph"
	"github.com/etc-sudonters/substrate/skelly/stack"
	"github.com/etc-sudonters/substrate/slipup"
)

func ExploreBasicGraph(z *app.Zootlr) error {
	worldGraph := app.GetResource[graph.Builder](z).Res.G
	locations := app.GetResource[entities.Locations](z)
	root := app.GetResource[world.Root](z).Res

	graphLocations := map[graph.Node]entities.Location{}
	for loc := range locations.Res.All {
		graphLocations[graph.Node(loc.Id())] = loc
	}
	exploration := stack.Make[graph.Node](0, 32)
	exploration.Push(graph.Node(root))
	visited := bitset.New(bitset.Buckets(uint64(worldGraph.NodeCount())))

	for exploration.Len() > 0 {
		current, _ := exploration.Pop()
		if !bitset.Set(&visited, current) {
			continue
		}
		currentLoc, exists := graphLocations[current]
		if !exists {
			panic(slipup.Createf("did not find origination 0x%04X in location map", current))
		}

		if currentLoc.Name() == "Farores Wind Warp" {
			continue
		}

		successors, _ := worldGraph.Successors(current)
		if successors.Len() == 0 {
			continue
		}

		dontio.WriteLineOut(z.Ctx(), "From %q can traverse to", currentLoc.Name())
		for dest := range bitset.Iter64T[graph.Node](successors).All {
			destLoc, exists := graphLocations[dest]
			if !exists {
				panic(slipup.Createf("did not find destination 0x%04X in location map", current))
			}
			dontio.WriteLineOut(z.Ctx(), "- %q", destLoc.Name())
			exploration.Push(dest)
		}
		dontio.WriteLineOut(z.Ctx(), "")
	}

	return nil
}
