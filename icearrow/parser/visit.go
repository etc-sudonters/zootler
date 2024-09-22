package parser

import (
	"errors"
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

type Transformer[T any] interface {
	TransformBinOp(*BinOp) (T, error)
	TransformBoolOp(*BoolOp) (T, error)
	TransformCall(*Call) (T, error)
	TransformIdentifier(*Identifier) (T, error)
	TransformSubscript(*Subscript) (T, error)
	TransformTuple(*Tuple) (T, error)
	TransformUnary(*UnaryOp) (T, error)
	TransformLiteral(*Literal) (T, error)
}

func TransformAll[T any](trans Transformer[T], nodes []Expression) ([]T, error) {
	transed := make([]T, len(nodes))
	var err error

	for i, node := range nodes {
		t, tErr := Transform(trans, node)
		transed[i] = t
		err = errors.Join(err, tErr)
	}

	return transed, err
}

func Transform[T any](trans Transformer[T], node Expression) (T, error) {
	switch node := node.(type) {
	case *BinOp:
		return trans.TransformBinOp(node)
	case *BoolOp:
		return trans.TransformBoolOp(node)
	case *Call:
		return trans.TransformCall(node)
	case *Identifier:
		return trans.TransformIdentifier(node)
	case *Literal:
		return trans.TransformLiteral(node)
	case *Subscript:
		return trans.TransformSubscript(node)
	case *Tuple:
		return trans.TransformTuple(node)
	case *UnaryOp:
		return trans.TransformUnary(node)
	default:
		panic(stageleft.AttachExitCode(
			fmt.Errorf("unknown node type %T", node),
			stageleft.ExitCode(91),
		))
	}
}

func TransformBinOpAst(t Transformer[Expression], op *BinOp) (*BinOp, error) {
	left, leftErr := Transform(t, op.Left)
	if leftErr != nil {
		return nil, slipup.Describef(leftErr, "while transforming left hand side %+v", op)
	}
	right, rightErr := Transform(t, op.Right)
	if rightErr != nil {
		return nil, slipup.Describef(rightErr, "while transforming right hand side %+v", op)
	}
	op = &BinOp{
		Left:  left,
		Op:    op.Op,
		Right: right,
	}
	return op, nil
}

func TransformBoolOpAst(t Transformer[Expression], op *BoolOp) (*BoolOp, error) {
	left, leftErr := Transform(t, op.Left)
	if leftErr != nil {
		return nil, slipup.Describef(leftErr, "while transforming left hand side %+v", op)
	}
	right, rightErr := Transform(t, op.Right)
	if rightErr != nil {
		return nil, slipup.Describef(rightErr, "while transforming right hand side %+v", op)
	}
	op = &BoolOp{
		Left:  left,
		Op:    op.Op,
		Right: right,
	}
	return op, nil
}

func TransformCallAst(t Transformer[Expression], call *Call) (*Call, error) {
	callee, calleeErr := Transform(t, call.Callee)
	if calleeErr != nil {
		return nil, slipup.Describef(calleeErr, "while parsing callee %+v", callee)
	}

	args := call.Args
	call = &Call{
		Callee: callee,
		Args:   make([]Expression, len(args), len(args)),
	}

	for i, a := range args {
		arg, argErr := Transform(t, a)
		if argErr != nil {
			return nil, slipup.Describef(argErr, "while transforming parameter %d %+v", i, call)
		}
		call.Args[i] = arg
	}
	return call, nil
}

func TransformIdentifierAst(t Transformer[Expression], ident *Identifier) (*Identifier, error) {
	return ident, nil
}

func TransformLiteral(t Transformer[Expression], literal *Literal) (*Literal, error) {
	return literal, nil
}

func TransformSubscriptAst(t Transformer[Expression], subscript *Subscript) (*Subscript, error) {
	target, targetErr := Transform(t, subscript.Target)
	if targetErr != nil {
		return nil, slipup.Describef(targetErr, "while transforming target %+v", subscript)
	}
	index, indexErr := Transform(t, subscript.Index)
	if indexErr != nil {
		return nil, slipup.Describef(indexErr, "while transforming index %+v", subscript)
	}
	subscript = &Subscript{
		Target: target,
		Index:  index,
	}
	return subscript, nil
}

func TransformTupleAst(t Transformer[Expression], tuple *Tuple) (*Tuple, error) {
	elems := tuple.Elems
	tuple = &Tuple{
		Elems: make([]Expression, len(elems), len(elems)),
	}
	for i, elm := range tuple.Elems {
		elem, elemErr := Transform(t, elm)
		if elemErr != nil {
			return nil, slipup.Describef(elemErr, "while transforming element %d %+v", i, tuple)
		}
		tuple.Elems[i] = elem
	}
	return tuple, nil
}

func TransformUnaryAst(t Transformer[Expression], unary *UnaryOp) (*UnaryOp, error) {
	operand, operandErr := Transform(t, unary.Target)
	if operandErr != nil {
		return nil, slipup.Describef(operandErr, "while transforming operand %+v", unary)
	}
	unary = &UnaryOp{
		Target: operand,
		Op:     unary.Op,
	}
	return unary, nil
}
