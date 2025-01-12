package optimizer

import (
	"sudonters/zootler/midologic/ast"
	"sudonters/zootler/midologic/symbols"
)

type CompilerFunction func([]ast.Node, ast.Rewriting) (ast.Node, error)

type CompilerFunctionTable map[string]CompilerFunction

func NewCompilerFuncs(symbolTable *symbols.Table, funcs CompilerFunctionTable) ast.Rewriter {
	compileFuncs := CompilerFunctions{funcs, symbolTable}
	for funcName := range funcs {
		symbolTable.Declare(funcName, symbols.COMPILER_FUNCTION)
	}
	return ast.Rewriter{Invoke: compileFuncs.Invoke}
}

type CompilerFunctions struct {
	funcs   CompilerFunctionTable
	symbols *symbols.Table
}

func (this CompilerFunctions) Invoke(node ast.Invoke, rewrite ast.Rewriting) (ast.Node, error) {
	symbol := ast.LookUpNodeInTable(this.symbols, node.Target)
	if symbol == nil {
		return node, nil
	}

	fn, exists := this.funcs[symbol.Name]
	if !exists {
		return node, nil
	}
	return fn(node.Args, rewrite)
}
