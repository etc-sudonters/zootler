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

type Visitor interface {
	Comparison(ast *Comparison) error
	BooleanOp(ast *BooleanOp) error
	Call(ast *Call) error
	Identifier(ast *Identifier) error
	Literal(ast *Literal) error
	Empty(ast *Empty) error
}

func Visit(guest Visitor, ast Node) error {
	switch ast := ast.(type) {
	case *Comparison:
		return guest.Comparison(ast)
	case *BooleanOp:
		return guest.BooleanOp(ast)
	case *Call:
		return guest.Call(ast)
	case *Identifier:
		return guest.Identifier(ast)
	case *Literal:
		return guest.Literal(ast)
	case *Empty:
		return guest.Empty(ast)
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
