package optimizer

import (
	"sudonters/libzootr/mido/ast"
	"sudonters/libzootr/mido/symbols"
)

type CanUse struct {
	symbols *symbols.Table
	cache   map[symbols.Index]ast.Node
}
