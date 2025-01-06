package optimizer

import (
	"fmt"
	"sudonters/zootler/magicbeanvm/ast"
	"sudonters/zootler/magicbeanvm/symbols"
)

func InvokeBareFuncs(symbols *symbols.Table, funcs *ast.PartialFunctionTable) ast.Rewriter {
	promote := invokebarefuncs{symbols, funcs}
	return ast.Rewriter{
		Invoke:     ast.DontRewrite[ast.Invoke](),
		Identifier: promote.Identifier,
	}
}

type invokebarefuncs struct {
	symbols *symbols.Table
	funcs   *ast.PartialFunctionTable
}

func (this invokebarefuncs) Identifier(node ast.Identifier, _ ast.Rewriting) (ast.Node, error) {
	switch node.Symbol.Kind {
	case symbols.BUILT_IN:
		return ast.Invoke{Target: node, Args: nil}, nil
	case symbols.FUNCTION, symbols.COMPILED_FUNC:
		fn, exists := this.funcs.Get(node.Symbol.Name)
		if !exists {
			return nil, fmt.Errorf("fn %q was declared but not available in table", node.Symbol.Name)
		}

		if len(fn.Params) == 0 {
			return ast.Invoke{Target: node, Args: nil}, nil
		}
		return nil, fmt.Errorf("expected 0-arg function, but %q has %d args", node.Symbol.Name, len(fn.Params))
	default:
		return node, nil
	}
}
