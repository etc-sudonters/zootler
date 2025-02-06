package optimizer

import (
	"sudonters/zootler/mido/ast"
	"sudonters/zootler/mido/symbols"
)

func FoldConstants(tbl *symbols.Table) ast.Rewriter {
	fold := fold{tbl}
	return ast.Rewriter{
		Compare: fold.Compare,
		Invert:  fold.Invert,
	}
}

type fold struct {
	*symbols.Table
}

func (this fold) Invert(node ast.Invert, rewrite ast.Rewriting) (ast.Node, error) {
	inner, err := rewrite(node.Inner)
	if err != nil {
		return nil, err
	}

	switch inner := inner.(type) {
	case ast.Boolean:
		return ast.Boolean(!inner), nil
	}

	return ast.Invert{Inner: inner}, nil
}

func (this fold) Compare(node ast.Compare, rewrite ast.Rewriting) (ast.Node, error) {
	if node.LHS.Kind() != node.RHS.Kind() {
		return rewrite.Compare(node)
	}

	switch node.LHS.Kind() {
	case ast.KindNumber:
		return ast.Boolean(this.compare_num(node)), nil
	case ast.KindString:
		return ast.Boolean(this.compare_str(node)), nil
	case ast.KindBool:
		return ast.Boolean(this.compare_bool(node)), nil
	case ast.KindIdentifier:
		return ast.Boolean(this.compare_ptr(node)), nil
	default:
		return rewrite.Compare(node)
	}
}

func (this fold) compare_str(node ast.Compare) bool {
	lhs := node.LHS.(ast.String)
	rhs := node.RHS.(ast.String)
	switch node.Op {
	case ast.CompareEq:
		return lhs == rhs
	case ast.CompareNq:
		return lhs != rhs
	case ast.CompareLt:
		return lhs < rhs
	default:
		panic("unsupported cmp op")
	}
}

func (this fold) compare_num(node ast.Compare) bool {
	lhs := node.LHS.(ast.Number)
	rhs := node.RHS.(ast.Number)
	switch node.Op {
	case ast.CompareEq:
		return lhs == rhs
	case ast.CompareNq:
		return lhs != rhs
	case ast.CompareLt:
		return lhs < rhs
	default:
		panic("unsupported cmp op")
	}
}

func (this fold) compare_bool(node ast.Compare) bool {
	lhs := node.LHS.(ast.Boolean)
	rhs := node.RHS.(ast.Boolean)
	switch node.Op {
	case ast.CompareEq:
		return lhs == rhs
	case ast.CompareNq:
		return lhs != rhs
	case ast.CompareLt:
		panic("booleans do not support ordering")
	default:
		panic("unsupported cmp op")
	}
}

func (this fold) compare_ptr(node ast.Compare) bool {
	lhs := this.LookUpByIndex(node.LHS.(ast.Identifier).AsIndex())
	rhs := this.LookUpByIndex(node.RHS.(ast.Identifier).AsIndex())
	switch node.Op {
	case ast.CompareEq:
		return lhs == rhs
	case ast.CompareNq:
		return lhs != rhs
	case ast.CompareLt:
		panic("ptrs do not support ordering")
	default:
		panic("unsupported cmp op")
	}

}
