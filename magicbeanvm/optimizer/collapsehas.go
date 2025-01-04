package optimizer

import (
	"slices"
	"sudonters/zootler/magicbeanvm/ast"
	"sudonters/zootler/magicbeanvm/symbols"
)

func CollapseHas(tbl *symbols.Table) ast.Rewriter {
	has, hasEvery, hasAny := tbl.Declare("has", symbols.FUNCTION), tbl.Declare("has_every", symbols.FUNCTION), tbl.Declare("has_anyof", symbols.FUNCTION)
	collapse := collapse{tbl, has, hasEvery, hasAny}
	return ast.Rewriter{
		Every: collapse.Every,
		AnyOf: collapse.AnyOf,
	}
}

type collapse struct {
	tbl                     *symbols.Table
	has, hasEvery, hasAnyOf *symbols.Sym
}

func (this collapse) Every(node ast.Every, rewrite ast.Rewriting) (ast.Node, error) {
	nodes, err := this.collapse(this.hasEvery, node, rewrite)
	return ast.Every(nodes).Reduce(), err
}

func (this collapse) AnyOf(node ast.AnyOf, rewrite ast.Rewriting) (ast.Node, error) {
	nodes, err := this.collapse(this.hasAnyOf, node, rewrite)
	return ast.AnyOf(nodes).Reduce(), err
}

type hasInvoke struct {
	what ast.Node
	qty  ast.Number
}

func (this collapse) isHasInvoke(node ast.Node) (hasInvoke, bool) {
	var invoke hasInvoke
	switch node := node.(type) {
	case ast.Invoke:
		sym := ast.LookUpNodeInTable(this.tbl, node.Target)
		if sym != nil && sym.Eq(this.has) && len(node.Args) == 2 && node.Args[1].Kind() == ast.KindNumber {
			invoke.what = node.Args[0]
			invoke.qty = node.Args[1].(ast.Number)
			return invoke, true
		}
	}

	return invoke, false
}

func (this collapse) isHasManyInvoke(which *symbols.Sym, node ast.Node) ([]ast.Node, bool) {
	switch node := node.(type) {
	case ast.Invoke:
		sym := ast.LookUpNodeInTable(this.tbl, node.Target)
		if sym != nil && sym.Eq(which) {
			return node.Args, true
		}
	}
	return nil, false
}

func (this collapse) collapse(as *symbols.Sym, nodes []ast.Node, rewrite ast.Rewriting) ([]ast.Node, error) {
	var hasMany []ast.Node
	var collected []ast.Node

	for _, node := range nodes {
		if invoke, isHasInvoke := this.isHasInvoke(node); isHasInvoke && invoke.qty == 1 {
			hasMany = append(hasMany, invoke.what)
		} else if many, isHasMany := this.isHasManyInvoke(as, node); isHasMany {
			hasMany = slices.Concat(hasMany, many)
		} else {
			var err error
			node, err = rewrite(node)
			if err != nil {
				return nil, err
			}
			collected = append(collected, node)
		}
	}

	switch len(hasMany) {
	case 0:
	case 1:
		collected = append(collected, ast.Invoke{
			Target: ast.IdentifierFrom(this.has),
			Args:   []ast.Node{hasMany[0], ast.Number(1)},
		})
	default:
		collected = append(collected, ast.Invoke{
			Target: ast.IdentifierFrom(as),
			Args:   hasMany,
		})
	}

	return collected, nil
}
