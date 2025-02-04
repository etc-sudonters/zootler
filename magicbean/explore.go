package magicbean

import (
	"fmt"
	"strings"
	"sudonters/zootler/internal"
	"sudonters/zootler/internal/skelly/bitset32"
	"sudonters/zootler/internal/skelly/graph32"
	"sudonters/zootler/mido"
	"sudonters/zootler/mido/code"
	"sudonters/zootler/mido/compiler"
	"sudonters/zootler/mido/objects"
	"sudonters/zootler/zecs"
)

type ExplorableEdge struct {
	Kind   EdgeKind
	Entity zecs.Entity
	Rule   RuleCompiled
	Name   Name
}

type ExplorableWorld struct {
	Graph graph32.Directed
	Edges map[Connection]ExplorableEdge
}

func (this ExplorableWorld) Edge(from, to graph32.Node) (ExplorableEdge, bool) {
	edge, exists := this.Edges[Connection{zecs.Entity(from), zecs.Entity(to)}]
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

func (this *Exploration) CanTransit(world *ExplorableWorld, from, to graph32.Node) bool {
	edge, exists := world.Edge(from, to)
	if !exists {
		panic(fmt.Errorf("no edge registered between %d %d", from, to))
	}
	fmt.Printf("exploring %q\n", edge.Name)
	var dis strings.Builder
	this.VM.Dis(&dis, compiler.Bytecode(edge.Rule))
	fmt.Println(dis.String())
	result := this.evaluateRule(compiler.Bytecode(edge.Rule))
	fmt.Printf("\tcrossed? %t\n", result)
	return result
}

type Results struct {
	Workset bitset32.Bitset
	Reached bitset32.Bitset
}

func (this *ExplorableWorld) ExploreAvailableEdges(xplr Exploration) Results {
	var results Results

	for current := range nodeiter(xplr.Workset).UntilEmpty {
		neighbors, err := this.Graph.Successors(current)
		internal.PanicOnError(err)
		neighbors = neighbors.Difference(*xplr.Visited)
		for neighbor := range nodeiter(&neighbors).All {
			if xplr.CanTransit(this, current, neighbor) {
				bitset32.Unset(&neighbors, neighbor)
				bitset32.Set(xplr.Workset, neighbor)
				bitset32.Set(&results.Reached, neighbor)
				bitset32.Set(xplr.Visited, neighbor)
			}
		}

		if !neighbors.IsEmpty() {
			bitset32.Set(&results.Workset, current)
		}
	}

	return results
}

func nodeiter(set *bitset32.Bitset) bitset32.IterOf[graph32.Node] {
	return bitset32.IterT[graph32.Node](set)
}
