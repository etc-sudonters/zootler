package parser

import (
	"fmt"

	"github.com/etc-sudonters/substrate/stageleft"
)

type Visitor interface {
	VisitBinOp(*BinOp) error
	VisitBoolOp(*BoolOp) error
	VisitCall(*Call) error
	VisitIdentifier(*Identifier) error
	VisitSubscript(*Subscript) error
	VisitTuple(*Tuple) error
	VisitUnary(*UnaryOp) error
	VisitLiteral(*Literal) error
}

func Visit(v Visitor, node Expression) error {
	switch node := node.(type) {
	case *BinOp:
		return v.VisitBinOp(node)
	case *BoolOp:
		return v.VisitBoolOp(node)
	case *Call:
		return v.VisitCall(node)
	case *Identifier:
		return v.VisitIdentifier(node)
	case *Literal:
		return v.VisitLiteral(node)
	case *Subscript:
		return v.VisitSubscript(node)
	case *Tuple:
		return v.VisitTuple(node)
	case *UnaryOp:
		return v.VisitUnary(node)
	default:
		panic(stageleft.AttachExitCode(
			fmt.Errorf("unknown node type %T", node),
			stageleft.ExitCode(86),
		))
	}
}
