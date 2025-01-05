package optimizer

import (
	"fmt"
	"sudonters/zootler/magicbeanvm/ast"
	"sudonters/zootler/magicbeanvm/symbols"
)

func EnsureFuncs(symbols *symbols.Table, funcs *ast.FunctionTable) ast.Rewriter {
	promote := ensurefuncs{symbols, funcs}
	return ast.Rewriter{
		Invoke:     ast.DontRewrite[ast.Invoke](),
		Identifier: promote.Identifier,
	}
}

type ensurefuncs struct {
	symbols *symbols.Table
	funcs   *ast.FunctionTable
}

func (this ensurefuncs) Identifier(node ast.Identifier, _ ast.Rewriting) (ast.Node, error) {
	symbol := this.symbols.LookUpByIndex(node.AsIndex())

	switch symbol.Kind {
	case symbols.BUILT_IN:
		return ast.Invoke{Target: node, Args: nil}, nil
	case symbols.FUNCTION, symbols.COMPILED_FUNC:
		fn, exists := this.funcs.Get(node)
		if !exists {
			return nil, fmt.Errorf("fn %q was declared but not available in table", symbol.Name)
		}

		if len(fn.Params) == 0 {
			return ast.Invoke{Target: node, Args: nil}, nil
		}
		return nil, fmt.Errorf("expected 0-arg function, but %q has %d args", symbol.Name, len(fn.Params))
	default:
		return node, nil
	}
}
