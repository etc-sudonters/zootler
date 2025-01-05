package optimizer

import (
	"sudonters/zootler/magicbeanvm/ast"
	"sudonters/zootler/magicbeanvm/symbols"
)

func FoldConstants(tbl *symbols.Table) ast.Rewriter {
	fold := fold{tbl}
	return ast.Rewriter{
		AnyOf:   fold.AnyOf,
		Compare: fold.Compare,
		Every:   fold.Every,
		Invert:  fold.Invert,
	}
}

type fold struct {
	*symbols.Table
}

func (f fold) Invert(node ast.Invert, rewrite ast.Rewriting) (ast.Node, error) {
	inner, err := rewrite(node.Inner)
	if err != nil {
		return nil, err
	}

	switch inner := inner.(type) {
	case ast.Bool:
		return ast.Bool(!inner), nil
	default:
		return ast.Invert{inner}, nil
	}

}

func (f fold) AnyOf(node ast.AnyOf, rewrite ast.Rewriting) (ast.Node, error) {
	rewritten, err := rewrite.All(node)
	if err != nil {
		return nil, err
	}

	return ast.AnyOf(rewritten).Flatten().Reduce(), nil
}

func (f fold) Every(node ast.Every, rewrite ast.Rewriting) (ast.Node, error) {
	rewritten, err := rewrite.All(node)
	if err != nil {
		return nil, err
	}

	return ast.Every(rewritten).Flatten().Reduce(), nil
}

func (f fold) Compare(node ast.Compare, rewrite ast.Rewriting) (ast.Node, error) {
	if (node.Op != ast.CompareEq && node.Op != ast.CompareNq) || node.LHS.Kind() != ast.KindIdentifier || node.RHS.Kind() != ast.KindIdentifier {
		return node, nil
	}

	lhs := f.LookUpByIndex(node.LHS.(ast.Identifier).AsIndex())
	rhs := f.LookUpByIndex(node.RHS.(ast.Identifier).AsIndex())

	switch node.Op {
	case ast.CompareEq:
		return ast.Bool(lhs.Eq(rhs)), nil
	case ast.CompareNq:
		return !ast.Bool(lhs.Eq(rhs)), nil
	default:
		panic("unreachable")
	}
}
