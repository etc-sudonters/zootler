package optimizer

import (
	"strings"
	"sudonters/zootler/magicbeanvm/ast"
	"sudonters/zootler/magicbeanvm/symbols"
)

func PromoteTokens(tbl *symbols.Table) ast.Rewriter {
	promote := promotetokens{tbl}
	return ast.Rewriter{
		Invoke:     ast.DontRewrite[ast.Invoke](),
		Identifier: promote.Identifier,
		String:     promote.String,
	}
}

type promotetokens struct {
	*symbols.Table
}

func (this promotetokens) Identifier(node ast.Identifier, _ ast.Rewriting) (ast.Node, error) {
	symbol := this.LookUpByIndex(node.AsIndex())

	switch symbol.Kind {
	case symbols.TOKEN:
		return this.has(node), nil
	case symbols.SETTING:
		return this.loadSetting(symbol.Name), nil
	default:
		if strings.Contains(symbol.Name, "_") {
			rawSymbol := this.LookUpByName(strings.ReplaceAll(symbol.Name, "_", " "))
			if rawSymbol != nil {
				return this.has(ast.IdentifierFrom(rawSymbol)), nil
			}
		}
	}
	return node, nil
}

func (this promotetokens) String(node ast.String, _ ast.Rewriting) (ast.Node, error) {
	symbol := this.LookUpByName(string(node))

	switch {
	case symbol == nil:
		return node, nil
	case symbol.Kind == symbols.TOKEN:
		return this.has(ast.IdentifierFrom(symbol)), nil
	case symbol.Kind == symbols.SETTING:
		return this.loadSetting(symbol.Name), nil
	default:
		if strings.Contains(symbol.Name, "_") {
			rawSymbol := this.LookUpByName(strings.ReplaceAll(symbol.Name, "_", " "))
			if rawSymbol != nil {
				return this.has(ast.IdentifierFrom(rawSymbol)), nil
			}
		}
	}
	return node, nil
}

func (this promotetokens) has(what ast.Node) ast.Invoke {
	return ast.Invoke{
		Target: ast.IdentifierFrom(this.Declare("has", symbols.FUNCTION)),
		Args:   []ast.Node{what, ast.Number(1)},
	}
}

func (this promotetokens) loadSetting(which string) ast.Invoke {
	return ast.Invoke{
		Target: ast.IdentifierFrom(this.Declare("load_setting", symbols.FUNCTION)),
		Args:   []ast.Node{ast.String(which)},
	}
}
