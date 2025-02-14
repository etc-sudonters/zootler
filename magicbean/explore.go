package magicbean

import (
	"context"
	"fmt"
	"sudonters/libzootr/components"
	"sudonters/libzootr/internal"
	"sudonters/libzootr/internal/skelly/bitset32"
	"sudonters/libzootr/internal/skelly/graph32"
	"sudonters/libzootr/mido"
	"sudonters/libzootr/mido/code"
	"sudonters/libzootr/mido/compiler"
	"sudonters/libzootr/mido/objects"
	"sudonters/libzootr/zecs"
)

type ExplorableEdge struct {
	Kind   components.EdgeKind
	Entity zecs.Entity
	Rule   components.RuleCompiled
	Src    components.RuleSource
	Name   components.Name
}

type ExplorableWorld struct {
	Graph graph32.Directed
	Edges map[components.Connection]ExplorableEdge
}

func (this ExplorableWorld) Edge(from, to graph32.Node) (ExplorableEdge, bool) {
	conn := components.Connection{From: zecs.Entity(from), To: zecs.Entity(to)}
	edge, exists := this.Edges[conn]
	return edge, exists
}

type Exploration struct {
	VM      mido.VM
	Visited *bitset32.Bitset
	Workset *bitset32.Bitset
	Objects *objects.Table
}

func (this *Exploration) evaluateRule(bytecode compiler.Bytecode) bool {
	if len(bytecode.Tape) == 1 {
		switch bytecode.Tape[0] {
		case code.PUSH_T:
			return true
		case code.PUSH_F:
			return false
		}
	}

	answer, vmErr := this.VM.Execute(bytecode)
	if vmErr != nil {
		fmt.Println(vmErr)
		answer = objects.PackedFalse
	}

	return this.VM.Truthy(answer)
}

func (this *Exploration) CanTransit(ctx context.Context, world *ExplorableWorld, from, to graph32.Node) bool {
	edge, exists := world.Edge(from, to)
	if !exists {
		panic(fmt.Errorf("no edge registered between %d %d", from, to))
	}
	bytecode := compiler.Bytecode(edge.Rule)
	result := this.evaluateRule(bytecode)
	return result
}

type ExplorationResults struct {
	Pending bitset32.Bitset
	Reached bitset32.Bitset
}

func (this *ExplorableWorld) ExploreAvailableEdges(ctx context.Context, xplr *Exploration) ExplorationResults {
	var results ExplorationResults
	for current := range nodeiter(xplr.Workset).UntilEmpty {
		neighbors, err := this.Graph.Successors(current)
		internal.PanicOnError(err)
		neighbors = neighbors.Difference(*xplr.Visited)

		for neighbor := range nodeiter(&neighbors).All {
			if xplr.CanTransit(ctx, this, current, neighbor) {
				bitset32.Unset(&neighbors, neighbor)
				bitset32.Set(xplr.Workset, neighbor)
				bitset32.Set(&results.Reached, neighbor)
				bitset32.Set(xplr.Visited, neighbor)
			}
		}

		if !neighbors.IsEmpty() {
			bitset32.Set(&results.Pending, current)
		}
	}

	return results
}

func nodeiter(set *bitset32.Bitset) bitset32.IterOf[graph32.Node] {
	return bitset32.IterT[graph32.Node](set)
}
