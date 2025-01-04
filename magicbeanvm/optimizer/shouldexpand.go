package optimizer

import (
	"sudonters/zootler/magicbeanvm/ast"
	"sudonters/zootler/magicbeanvm/symbols"
)

func ShouldExpand(tbl *symbols.Table, node ast.Node) bool {
	marker := marker{tbl, false}
	v := ast.Visitor{
		Invoke:     ast.DontVisit[ast.Invoke](),
		Identifier: (&marker).Identifier,
	}
	v.Visit(node)
	return marker.expand
}

type marker struct {
	*symbols.Table
	expand bool
}

func (this *marker) Identifier(node ast.Identifier, _ ast.Visiting) error {
	symbol := this.LookUpByIndex(node.AsIndex())
	switch symbol.Type {
	case symbols.SETTING, symbols.TOKEN, symbols.FUNCTION:
		this.expand = true
	}
	return nil
}
