package parser

import (
	"fmt"
	"github.com/etc-sudonters/substrate/slipup"

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
			stageleft.ExitCode(90),
		))
	}
}

type Transformer interface {
	TransformBinOp(*BinOp) (Expression, error)
	TransformBoolOp(*BoolOp) (Expression, error)
	TransformCall(*Call) (Expression, error)
	TransformIdentifier(*Identifier) (Expression, error)
	TransformSubscript(*Subscript) (Expression, error)
	TransformTuple(*Tuple) (Expression, error)
	TransformUnary(*UnaryOp) (Expression, error)
	TransformLiteral(*Literal) (Expression, error)
}

func Transform(t Transformer, node Expression) (Expression, error) {
	switch node := node.(type) {
	case *BinOp:
		return t.TransformBinOp(node)
	case *BoolOp:
		return t.TransformBoolOp(node)
	case *Call:
		return t.TransformCall(node)
	case *Identifier:
		return t.TransformIdentifier(node)
	case *Literal:
		return t.TransformLiteral(node)
	case *Subscript:
		return t.TransformSubscript(node)
	case *Tuple:
		return t.TransformTuple(node)
	case *UnaryOp:
		return t.TransformUnary(node)
	default:
		panic(stageleft.AttachExitCode(
			fmt.Errorf("unknown node type %T", node),
			stageleft.ExitCode(91),
		))
	}
}

func TransformBinOp(t Transformer, op *BinOp) (*BinOp, error) {
	left, leftErr := Transform(t, op.Left)
	if leftErr != nil {
		return nil, slipup.Describef(leftErr, "while transforming left hand side %+v", op)
	}
	right, rightErr := Transform(t, op.Right)
	if rightErr != nil {
		return nil, slipup.Describef(rightErr, "while transforming right hand side %+v", op)
	}
	op.Left = left
	op.Right = right
	return op, nil
}

func TransformBoolOp(t Transformer, op *BoolOp) (*BoolOp, error) {
	left, leftErr := Transform(t, op.Left)
	if leftErr != nil {
		return nil, slipup.Describef(leftErr, "while transforming left hand side %+v", op)
	}
	right, rightErr := Transform(t, op.Right)
	if rightErr != nil {
		return nil, slipup.Describef(rightErr, "while transforming right hand side %+v", op)
	}
	op.Left = left
	op.Right = right
	return op, nil
}

func TransformCall(t Transformer, call *Call) (*Call, error) {
	callee, calleeErr := Transform(t, call.Callee)
	if calleeErr != nil {
		return nil, slipup.Describef(calleeErr, "while parsing callee %+v", callee)
	}

	for i, a := range call.Args {
		arg, argErr := Transform(t, a)
		if argErr != nil {
			return nil, slipup.Describef(argErr, "while transforming parameter %d %+v", i, call)
		}
		call.Args[i] = arg
	}
	call.Callee = callee
	return call, nil
}

func TransformIdentifier(t Transformer, ident *Identifier) (*Identifier, error) {
	return ident, nil
}

func TransformLiteral(t Transformer, literal *Literal) (*Literal, error) {
	return literal, nil
}

func TransformSubscript(t Transformer, subscript *Subscript) (*Subscript, error) {
	target, targetErr := Transform(t, subscript.Target)
	if targetErr != nil {
		return nil, slipup.Describef(targetErr, "while transforming target %+v", subscript)
	}
	index, indexErr := Transform(t, subscript.Index)
	if indexErr != nil {
		return nil, slipup.Describef(indexErr, "while transforming index %+v", subscript)
	}
	subscript.Target = target
	subscript.Index = index
	return subscript, nil
}

func TransformTuple(t Transformer, tuple *Tuple) (*Tuple, error) {
	for i, elm := range tuple.Elems {
		elem, elemErr := Transform(t, elm)
		if elemErr != nil {
			return nil, slipup.Describef(elemErr, "while transforming element %d %+v", i, tuple)
		}
		tuple.Elems[i] = elem
	}
	return tuple, nil
}

func TransformUnary(t Transformer, unary *UnaryOp) (*UnaryOp, error) {
	operand, operandErr := Transform(t, unary.Target)
	if operandErr != nil {
		return nil, slipup.Describef(operandErr, "while transforming operand %+v", unary)
	}
	unary.Target = operand
	return unary, nil
}
