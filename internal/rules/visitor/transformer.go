package visitor

import (
	"errors"
	"fmt"
	"sudonters/zootler/internal/rules/parser"

	"github.com/etc-sudonters/substrate/slipup"

	"github.com/etc-sudonters/substrate/stageleft"
)

type Transformer[T any] interface {
	TransformBinOp(*parser.BinOp) (T, error)
	TransformBoolOp(*parser.BoolOp) (T, error)
	TransformCall(*parser.Call) (T, error)
	TransformIdentifier(*parser.Identifier) (T, error)
	TransformSubscript(*parser.Subscript) (T, error)
	TransformTuple(*parser.Tuple) (T, error)
	TransformUnary(*parser.UnaryOp) (T, error)
	TransformLiteral(*parser.Literal) (T, error)
}

func TransformAll[T any](trans Transformer[T], nodes []parser.Expression) ([]T, error) {
	transed := make([]T, len(nodes))
	var err error

	for i, node := range nodes {
		t, tErr := Transform(trans, node)
		transed[i] = t
		err = errors.Join(err, tErr)
	}

	return transed, err
}

func Transform[T any](trans Transformer[T], node parser.Expression) (T, error) {
	switch node := node.(type) {
	case *parser.BinOp:
		return trans.TransformBinOp(node)
	case *parser.BoolOp:
		return trans.TransformBoolOp(node)
	case *parser.Call:
		return trans.TransformCall(node)
	case *parser.Identifier:
		return trans.TransformIdentifier(node)
	case *parser.Literal:
		return trans.TransformLiteral(node)
	case *parser.Subscript:
		return trans.TransformSubscript(node)
	case *parser.Tuple:
		return trans.TransformTuple(node)
	case *parser.UnaryOp:
		return trans.TransformUnary(node)
	default:
		panic(stageleft.AttachExitCode(
			fmt.Errorf("unknown node type %T", node),
			stageleft.ExitCode(91),
		))
	}
}

func TransformBinOpAst(t Transformer[parser.Expression], op *parser.BinOp) (*parser.BinOp, error) {
	left, leftErr := Transform(t, op.Left)
	if leftErr != nil {
		return nil, slipup.Describef(leftErr, "while transforming left hand side %+v", op)
	}
	right, rightErr := Transform(t, op.Right)
	if rightErr != nil {
		return nil, slipup.Describef(rightErr, "while transforming right hand side %+v", op)
	}
	op = &parser.BinOp{
		Left:  left,
		Op:    op.Op,
		Right: right,
	}
	return op, nil
}

func TransformBoolOpAst(t Transformer[parser.Expression], op *parser.BoolOp) (*parser.BoolOp, error) {
	left, leftErr := Transform(t, op.Left)
	if leftErr != nil {
		return nil, slipup.Describef(leftErr, "while transforming left hand side %+v", op)
	}
	right, rightErr := Transform(t, op.Right)
	if rightErr != nil {
		return nil, slipup.Describef(rightErr, "while transforming right hand side %+v", op)
	}
	op = &parser.BoolOp{
		Left:  left,
		Op:    op.Op,
		Right: right,
	}
	return op, nil
}

func TransformCallAst(t Transformer[parser.Expression], call *parser.Call) (*parser.Call, error) {
	callee, calleeErr := Transform(t, call.Callee)
	if calleeErr != nil {
		return nil, slipup.Describef(calleeErr, "while parsing callee %+v", callee)
	}

	args := call.Args
	call = &parser.Call{
		Callee: callee,
		Args:   make([]parser.Expression, len(args), len(args)),
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

func TransformIdentifierAst(t Transformer[parser.Expression], ident *parser.Identifier) (*parser.Identifier, error) {
	return ident, nil
}

func TransformLiteral(t Transformer[parser.Expression], literal *parser.Literal) (*parser.Literal, error) {
	return literal, nil
}

func TransformSubscriptAst(t Transformer[parser.Expression], subscript *parser.Subscript) (*parser.Subscript, error) {
	target, targetErr := Transform(t, subscript.Target)
	if targetErr != nil {
		return nil, slipup.Describef(targetErr, "while transforming target %+v", subscript)
	}
	index, indexErr := Transform(t, subscript.Index)
	if indexErr != nil {
		return nil, slipup.Describef(indexErr, "while transforming index %+v", subscript)
	}
	subscript = &parser.Subscript{
		Target: target,
		Index:  index,
	}
	return subscript, nil
}

func TransformTupleAst(t Transformer[parser.Expression], tuple *parser.Tuple) (*parser.Tuple, error) {
	elems := tuple.Elems
	tuple = &parser.Tuple{
		Elems: make([]parser.Expression, len(elems), len(elems)),
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

func TransformUnaryAst(t Transformer[parser.Expression], unary *parser.UnaryOp) (*parser.UnaryOp, error) {
	operand, operandErr := Transform(t, unary.Target)
	if operandErr != nil {
		return nil, slipup.Describef(operandErr, "while transforming operand %+v", unary)
	}
	unary = &parser.UnaryOp{
		Target: operand,
		Op:     unary.Op,
	}
	return unary, nil
}
