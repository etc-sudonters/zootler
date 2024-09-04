package ast

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
