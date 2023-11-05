package ast

import (
	"fmt"
	"github.com/etc-sudonters/substrate/stageleft"
)

type Visitor interface {
	VisitAttrAccess(*AttrAccess) error
	VisitBinOp(*BinOp) error
	VisitBoolOp(*BoolOp) error
	VisitBoolean(*Boolean) error
	VisitCall(*Call) error
	VisitIdentifier(*Identifier) error
	VisitNumber(*Number) error
	VisitString(*String) error
	VisitSubscript(*Subscript) error
	VisitTuple(*Tuple) error
	VisitUnary(*UnaryOp) error
}

func Visit(v Visitor, node Expression) error {
	switch node := node.(type) {
	case *AttrAccess:
		return v.VisitAttrAccess(node)
	case *BinOp:
		return v.VisitBinOp(node)
	case *BoolOp:
		return v.VisitBoolOp(node)
	case *Boolean:
		return v.VisitBoolean(node)
	case *Call:
		return v.VisitCall(node)
	case *Identifier:
		return v.VisitIdentifier(node)
	case *Number:
		return v.VisitNumber(node)
	case *String:
		return v.VisitString(node)
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
