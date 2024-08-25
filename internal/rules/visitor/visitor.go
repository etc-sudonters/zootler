package visitor

import (
	"fmt"
	"sudonters/zootler/internal/rules/parser"
	"github.com/etc-sudonters/substrate/stageleft"
)

type Visitor interface {
	VisitBinOp(*parser.BinOp) error
	VisitBoolOp(*parser.BoolOp) error
	VisitCall(*parser.Call) error
	VisitIdentifier(*parser.Identifier) error
	VisitSubscript(*parser.Subscript) error
	VisitTuple(*parser.Tuple) error
	VisitUnary(*parser.UnaryOp) error
	VisitLiteral(*parser.Literal) error
}

func Visit(v Visitor, node parser.Expression) error {
	switch node := node.(type) {
	case *parser.BinOp:
		return v.VisitBinOp(node)
	case *parser.BoolOp:
		return v.VisitBoolOp(node)
	case *parser.Call:
		return v.VisitCall(node)
	case *parser.Identifier:
		return v.VisitIdentifier(node)
	case *parser.Literal:
		return v.VisitLiteral(node)
	case *parser.Subscript:
		return v.VisitSubscript(node)
	case *parser.Tuple:
		return v.VisitTuple(node)
	case *parser.UnaryOp:
		return v.VisitUnary(node)
	default:
		panic(stageleft.AttachExitCode(
			fmt.Errorf("unknown node type %T", node),
			stageleft.ExitCode(90),
		))
	}
}
