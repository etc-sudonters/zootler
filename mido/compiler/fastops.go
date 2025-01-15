package compiler

import (
	"sudonters/zootler/mido/ast"
	"sudonters/zootler/mido/code"
	"sudonters/zootler/mido/objects"
	"sudonters/zootler/mido/symbols"
)

type FastOpCompiler func(ast.Invoke, *symbols.Table, *objects.TableBuilder, ast.Visiting) (code.Instructions, error)

type FastOps map[string]FastOpCompiler

func FastHasOp(node ast.Invoke, symbolTable *symbols.Table, objTable *objects.TableBuilder, _ ast.Visiting) (code.Instructions, error) {
	what := ast.LookUpNodeInTable(symbolTable, node.Args[0])
	qty, isQty := node.Args[1].(ast.Number)
	if what != nil && isQty {
		ptr := objTable.GetPointerFor(what.Name)
		return code.Make(code.CHK_QTY, int(ptr), int(qty)), nil
	}
	return nil, nil
}
