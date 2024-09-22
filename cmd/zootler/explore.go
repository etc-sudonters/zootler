package main

import (
	"sudonters/zootler/carpenters/shiro"
	"sudonters/zootler/icearrow/compiler"
	"sudonters/zootler/icearrow/runtime"
	"sudonters/zootler/internal/app"
	"sudonters/zootler/internal/entities"
	"sudonters/zootler/internal/world"

	"github.com/etc-sudonters/substrate/dontio"
	"github.com/etc-sudonters/substrate/skelly/graph"
	"github.com/etc-sudonters/substrate/slipup"
)

func ExploreBasicGraph(z *app.Zootlr) error {
	root := app.GetResource[world.Root](z).Res
	prog := app.GetResource[shiro.CompiledWorldRules](z)
	world, err := BuildWorldGraph(z)
	if err != nil {
		return err
	}
	exploration := world.BFS(root)
	vmState := FakeVMState{}
	vm := runtime.VM{}
	symbols := &prog.Res.Symbols

	for traversal := range exploration.Walk {
		dontio.WriteLineOut(z.Ctx(), "executing %q", traversal.EdgeAttrs.Name())
		vm.Execute(&traversal.EdgeAttrs.Tape, vmState, symbols)
		exploration.Accept(traversal.Destination)
	}

	return nil
}

type CompiledEdge struct {
	entities.Edge
	Tape compiler.Tape
}

func BuildWorldGraph(z *app.Zootlr) (world.Graph[entities.Location, CompiledEdge], error) {
	worldgraph := app.GetResource[graph.Builder](z)
	world := world.NewGraph[entities.Location, CompiledEdge](worldgraph.Res.G)
	locations := app.GetResource[entities.Locations](z)
	edges := app.GetResource[entities.Edges](z)
	rules := app.GetResource[shiro.CompiledWorldRules](z)

	for loc := range locations.Res.All {
		world.SetNodeAttributes(graph.Node(loc.Id()), loc)
	}

	for edge := range edges.Res.All {
		origin, wasOrigin := edge.Retrieve("originId").(graph.Origination)
		dest, wasDest := edge.Retrieve("destId").(graph.Destination)
		if !wasOrigin || !wasDest {
			return world, slipup.Createf(
				"%q is incomplete: {Origin: %v, Dest: %v}",
				edge.Name(), origin, dest,
			)
		}

		tape, hadTape := rules.Res.Rules[string(edge.Name())]
		if !hadTape {
			return world, slipup.Createf("%q does not have a compiled rule", edge.Name())
		}
		compiledEdge := CompiledEdge{
			Edge: edge,
			Tape: tape,
		}

		world.SetEdgeAttributes(origin, dest, compiledEdge)
	}

	return world, nil
}

type FakeVMState struct{}

func (_ FakeVMState) HasQty(uint32, uint8) bool { return true }
func (_ FakeVMState) HasAny(...uint32) bool     { return true }
func (_ FakeVMState) HasAll(...uint32) bool     { return true }
func (_ FakeVMState) HasBottle() bool           { return true }
func (_ FakeVMState) IsAdult() bool             { return true }
func (_ FakeVMState) IsChild() bool             { return true }
func (_ FakeVMState) AtTod(uint8) bool          { return true }
