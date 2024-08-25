
package visitor

import (
	"fmt"
	"sudonters/zootler/internal/rules/parser"

	"github.com/etc-sudonters/substrate/slipup"

	"github.com/etc-sudonters/substrate/stageleft"
)

type Transformer interface {
	TransformBinOp(*parser.BinOp) (parser.Expression, error)
	TransformBoolOp(*parser.BoolOp) (parser.Expression, error)
	TransformCall(*parser.Call) (parser.Expression, error)
	TransformIdentifier(*parser.Identifier) (parser.Expression, error)
	TransformSubscript(*parser.Subscript) (parser.Expression, error)
	TransformTuple(*parser.Tuple) (parser.Expression, error)
	TransformUnary(*parser.UnaryOp) (parser.Expression, error)
	TransformLiteral(*parser.Literal) (parser.Expression, error)
}

func Transform(t Transformer, node parser.Expression) (parser.Expression, error) {
	switch node := node.(type) {
	case *parser.BinOp:
		return t.TransformBinOp(node)
	case *parser.BoolOp:
		return t.TransformBoolOp(node)
	case *parser.Call:
		return t.TransformCall(node)
	case *parser.Identifier:
		return t.TransformIdentifier(node)
	case *parser.Literal:
		return t.TransformLiteral(node)
	case *parser.Subscript:
		return t.TransformSubscript(node)
	case *parser.Tuple:
		return t.TransformTuple(node)
	case *parser.UnaryOp:
		return t.TransformUnary(node)
	default:
		panic(stageleft.AttachExitCode(
			fmt.Errorf("unknown node type %T", node),
			stageleft.ExitCode(91),
		))
	}
}

func TransformBinOp(t Transformer, op *parser.BinOp) (*parser.BinOp, error) {
	left, leftErr := Transform(t, op.Left)
	if leftErr != nil {
		return nil, slipup.Describef(leftErr, "while transforming left hand side %+v", op)
	}
	right, rightErr := Transform(t, op.Right)
	if rightErr != nil {
		return nil, slipup.Describef(rightErr, "while transforming right hand side %+v", op)
	}
    op = &parser.BinOp{
        Left: left,
        Op: op.Op,
        Right: right,
    }
	return op, nil
}

func TransformBoolOp(t Transformer, op *parser.BoolOp) (*parser.BoolOp, error) {
	left, leftErr := Transform(t, op.Left)
	if leftErr != nil {
		return nil, slipup.Describef(leftErr, "while transforming left hand side %+v", op)
	}
	right, rightErr := Transform(t, op.Right)
	if rightErr != nil {
		return nil, slipup.Describef(rightErr, "while transforming right hand side %+v", op)
	}
    op = &parser.BoolOp{
        Left: left,
        Op: op.Op,
        Right: right,
    }
	return op, nil
}

func TransformCall(t Transformer, call *parser.Call) (*parser.Call, error) {
	callee, calleeErr := Transform(t, call.Callee)
	if calleeErr != nil {
		return nil, slipup.Describef(calleeErr, "while parsing callee %+v", callee)
	}

    args := call.Args
    call = &parser.Call{
        Callee: callee,
        Args: make([]parser.Expression, len(args), len(args)),
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

func TransformIdentifier(t Transformer, ident *parser.Identifier) (*parser.Identifier, error) {
	return ident, nil
}

func TransformLiteral(t Transformer, literal *parser.Literal) (*parser.Literal, error) {
	return literal, nil
}

func TransformSubscript(t Transformer, subscript *parser.Subscript) (*parser.Subscript, error) {
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
        Index: index,
    }
	return subscript, nil
}

func TransformTuple(t Transformer, tuple *parser.Tuple) (*parser.Tuple, error) {
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

func TransformUnary(t Transformer, unary *parser.UnaryOp) (*parser.UnaryOp, error) {
	operand, operandErr := Transform(t, unary.Target)
	if operandErr != nil {
		return nil, slipup.Describef(operandErr, "while transforming operand %+v", unary)
	}
    unary = &parser.UnaryOp{
        Target: operand,
        Op: unary.Op,
    }
	return unary, nil
}
