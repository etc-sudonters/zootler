package optimizer

import (
	"sudonters/zootler/mido/ast"
	"sudonters/zootler/mido/symbols"
)

type CanUse struct {
	symbols *symbols.Table
	cache   map[symbols.Index]ast.Node
}
