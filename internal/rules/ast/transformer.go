package ast

import "github.com/etc-sudonters/substrate/slipup"

type Transformer[T any] interface {
	Comparison(ast *Comparison) (T, error)
	BooleanOp(ast *BooleanOp) (T, error)
	Call(ast *Call) (T, error)
	Identifier(ast *Identifier) (T, error)
	Literal(ast *Literal) (T, error)
	Empty(ast *Empty) (T, error)
}

func Transform[T any](trans Transformer[T], ast Node) (T, error) {
	switch ast := ast.(type) {
	case *Comparison:
		return trans.Comparison(ast)
	case *BooleanOp:
		return trans.BooleanOp(ast)
	case *Call:
		return trans.Call(ast)
	case *Identifier:
		return trans.Identifier(ast)
	case *Literal:
		return trans.Literal(ast)
	case *Empty:
		return trans.Empty(ast)
	default:
		panic("aaahh!!!")
	}
}

func AssertIs[T Node](node Node) (T, bool) {
	t, ok := node.(T)
	return t, ok
}

func MustAssertAs[T Node](node Node) T {
	t, ok := AssertIs[T](node)
	if !ok {
		panic(slipup.Createf("could not assert node as %T: %+v", t, node))
	}

	return t
}

func Map(node Node, f func(Node) Node) Node {
	node, _ = Transform(mapast{f}, node)
	return node
}

type mapast struct {
	f func(Node) Node
}

func (c mapast) Identifier(ast *Identifier) (Node, error) {
	return c.f(&Identifier{
		Name: ast.Name,
		Kind: ast.Kind,
	}), nil
}
func (c mapast) Literal(ast *Literal) (Node, error) {
	return c.f(&Literal{Value: ast.Value, Kind: ast.Kind}), nil
}

func (c mapast) Comparison(ast *Comparison) (Node, error) {
	lhs, _ := Transform(c, ast.LHS)
	rhs, _ := Transform(c, ast.RHS)

	return c.f(&Comparison{
		LHS: lhs, RHS: rhs, Op: ast.Op,
	}), nil
}

func (c mapast) BooleanOp(ast *BooleanOp) (Node, error) {
	lhs, _ := Transform(c, ast.LHS)
	rhs, _ := Transform(c, ast.RHS)

	return c.f(&BooleanOp{
		LHS: lhs, RHS: rhs, Op: ast.Op,
	}), nil
}

func (c mapast) Empty(ast *Empty) (Node, error) {
	return c.f(&Empty{}), nil
}

func (c mapast) Call(ast *Call) (Node, error) {
	n := &Call{
		Callee: ast.Callee,
		Macro:  ast.Macro,
		Args:   make([]Node, len(ast.Args)),
	}

	for i := range len(ast.Args) {
		n.Args[i], _ = Transform(c, ast.Args[i])
	}

	return c.f(n), nil
}
