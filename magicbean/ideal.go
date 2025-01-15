package magicbean

import (
	"sudonters/zootler/internal/query"
	"sudonters/zootler/internal/skelly/bitset"
	"sudonters/zootler/mido/compiler"
	"sudonters/zootler/mido/objects"
	"sudonters/zootler/mido/vm"

	"github.com/etc-sudonters/substrate/skelly/graph"
)

type logic map[graph.Edge]compiler.Bytecode

func (this logic) rule(o, d graph.Node) (compiler.Bytecode, bool) {
	code, exists := this[graph.Edge{
		O: graph.Origination(o),
		D: graph.Destination(d),
	}]

	return code, exists
}

func createmap(data query.Engine) logic {
	q := data.CreateQuery()
	q.Load(query.MustAsColumnId[graph.Edge](data))
	q.Load(query.MustAsColumnId[compiler.Bytecode](data))

	results, err := data.Retrieve(q)
	if err != nil {
		panic(err)
	}
	logic := make(logic, results.Len())
	for _, tuple := range results.All {
		between := tuple.Values[0].(graph.Edge)
		code := tuple.Values[1].(compiler.Bytecode)
		logic[between] = code
	}
	return logic
}

func visitall() {
	var mido vm.VM
	var physical graph.Directed
	var workset bitset.Bitset32
	var visited bitset.Bitset32
	var reached bitset.Bitset32

	logic := createmap(nil)

	for visiting := range nodebiter(&workset).All {
		neighbors := successors(visiting, &physical)
		for neighbor := range nodebiter(&neighbors).All {
			if bitset.IsSet32(&visited, neighbor) {
				bitset.Unset32(&neighbors, neighbor)
			}

			code, exists := logic.rule(visiting, neighbor)
			if !exists {
				panic("no rule found between declared edge")
			}
			obj, execErr := mido.Execute(code)
			if execErr != nil {
				panic(execErr)
			}

			canTransit, isBool := obj.(objects.Boolean)
			if !isBool {
				panic("rule produced non-boolean result")
			}

			if canTransit {
				bitset.Unset32(&neighbors, neighbor)
				bitset.Set32(&workset, neighbor)
				bitset.Set32(&visited, neighbor)
				bitset.Set32(&reached, neighbor)
			}
		}

		if neighbors.IsEmpty() {
			bitset.Unset32(&workset, visiting)
		}
	}
}
