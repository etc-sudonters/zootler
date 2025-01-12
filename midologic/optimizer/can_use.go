package optimizer

import (
	"sudonters/zootler/midologic/ast"
	"sudonters/zootler/midologic/symbols"
)

type CanUse struct {
	symbols *symbols.Table
	cache   map[symbols.Index]ast.Node
}
