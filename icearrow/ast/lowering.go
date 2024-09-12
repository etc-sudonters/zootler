package ast

import (
	"errors"
	"sudonters/zootler/icearrow/parser"
	"sudonters/zootler/internal"

	"github.com/etc-sudonters/substrate/slipup"
)

type Ast struct{}

func (a *Ast) Lower(expr parser.Expression) (Node, error) {
	return parser.Transform(a, expr)
}

func (a *Ast) TransformBinOp(pt *parser.BinOp) (Node, error) {
	if pt.Op == parser.BinOpContains {
		return a.rewriteToCall("load_setting_2", pt.Right, pt.Left) // setting in block
	}

	newNode, didEliminate, eliminationErr := a.eliminateConstCompare(pt)
	if eliminationErr != nil {
		return nil, eliminationErr
	}
	if didEliminate {
		return newNode, nil
	}

	lhs, lhsErr := a.Lower(pt.Left)
	rhs, rhsErr := a.Lower(pt.Right)

	if joined := errors.Join(lhsErr, rhsErr); joined != nil {
		panic(joined)
	}

	switch pt.Op {
	case parser.BinOpEq:
		return &Comparison{LHS: lhs, RHS: rhs, Op: AST_CMP_EQ}, nil
	case parser.BinOpNotEq:
		return &Comparison{LHS: lhs, RHS: rhs, Op: AST_CMP_NQ}, nil
	case parser.BinOpLt:
		return &Comparison{LHS: lhs, RHS: rhs, Op: AST_CMP_LT}, nil
	}

	panic(slipup.Createf("unsupported binop: %+v", pt))
}

func (a *Ast) TransformBoolOp(pt *parser.BoolOp) (Node, error) {

	lhs, lhsErr := a.Lower(pt.Left)
	rhs, rhsErr := a.Lower(pt.Right)

	newNode, didEliminate, eliminationErr := a.eliminateConstBranch(lhs, rhs, pt.Op)
	if eliminationErr != nil {
		return nil, eliminationErr
	}
	if didEliminate {
		return newNode, nil
	}

	if joined := errors.Join(lhsErr, rhsErr); joined != nil {
		panic(joined)
	}

	switch pt.Op {
	case parser.BoolOpAnd:
		return &BooleanOp{LHS: lhs, RHS: rhs, Op: AST_BOOL_AND}, nil
	case parser.BoolOpOr:
		return &BooleanOp{LHS: lhs, RHS: rhs, Op: AST_BOOL_OR}, nil
	}

	panic(slipup.Createf("unsupported boolop: %+v", pt))
}

func (a *Ast) TransformCall(pt *parser.Call) (Node, error) {
	callee := parser.MustAssertAs[*parser.Identifier](pt.Callee)
	return a.rewriteToCall(callee.Value, pt.Args...)
}

func (a *Ast) TransformIdentifier(pt *parser.Identifier) (Node, error) {
	var ident Identifier
	ident.Name = pt.Value
	ident.Kind = AST_IDENT_UNK
	return &ident, nil
}

func (a *Ast) TransformSubscript(pt *parser.Subscript) (Node, error) {
	return a.rewriteToCall("load_setting_2", // block[setting]
		parser.MustAssertAs[*parser.Identifier](pt.Target),
		parser.MustAssertAs[*parser.Identifier](pt.Index),
	)
}

func (a *Ast) TransformTuple(pt *parser.Tuple) (Node, error) {
	if len(pt.Elems) != 2 {
		return nil, slipup.Createf("invalid tuple construction: %+v", pt)
	}

	return a.rewriteToCall("has", pt.Elems[0], pt.Elems[1])
}

func (a *Ast) TransformUnary(pt *parser.UnaryOp) (Node, error) {
	target, err := a.Lower(pt.Target)
	if err != nil {
		return nil, err
	}

	if target.Type() == AST_NODE_LITERAL {
		target := target.(*Literal)
		if target.Kind == AST_LIT_BOOL {
			return LiteralBool(!(target.Value.(bool))), nil
		}
	}

	return &BooleanOp{
		LHS: target,
		RHS: &Empty{},
		Op:  AST_BOOL_NEGATE,
	}, nil
}

func (a *Ast) TransformLiteral(pt *parser.Literal) (Node, error) {
	var lit Literal
	var exists bool
	lit.Kind, exists = ptToAstLitKind[pt.Kind]
	lit.Value = pt.Value
	if !exists {
		panic(slipup.Createf("unknown literal kind %+v", pt))
	}
	return &lit, nil
}

var ptToAstLitKind = map[parser.LiteralKind]AstLiteralKind{
	parser.LiteralBool: AST_LIT_BOOL,
	parser.LiteralNum:  AST_LIT_NUM,
	parser.LiteralStr:  AST_LIT_STR,
}

func (a *Ast) rewriteToCall(name string, inputs ...parser.Expression) (*Call, error) {
	var call Call
	args, argErr := parser.TransformAll(a, inputs)
	if argErr != nil {
		return nil, argErr
	}
	call.Callee = name
	call.Args = args
	return &call, nil
}

func (a *Ast) eliminateConstCompare(pt *parser.BinOp) (Node, bool, error) {
	if pt.Op != parser.BinOpEq && pt.Op != parser.BinOpNotEq {
		return nil, false, nil
	}
	// accomodate macro expansion that explodes things like Longshot into a has call
	lhs, lhsNotIdent := extractIdent(pt.Left)
	rhs, rhsNotIdent := extractIdent(pt.Right)

	if lhsNotIdent != nil || rhsNotIdent != nil {
		return nil, false, nil
	}

	sameIdent := internal.Normalize(lhs) == internal.Normalize(rhs)
	if !sameIdent {
		return nil, false, nil
	}

	node := Literal{
		Value: sameIdent && pt.Op == parser.BinOpEq,
		Kind:  AST_LIT_BOOL,
	}
	return &node, true, nil
}

func (a *Ast) eliminateConstBranch(lhs, rhs Node, op parser.BoolOpKind) (Node, bool, error) {
	getBoolLiteralFrom := func(node Node) (bool, bool) {
		lit, isNotLit := node.(*Literal)
		if !isNotLit {
			return false, false
		}

		b, isB := lit.Value.(bool)
		return b, isB
	}

	if leftHandLiteral, lhsIsBool := getBoolLiteralFrom(lhs); lhsIsBool {
		if op == parser.BoolOpAnd && leftHandLiteral {
			return rhs, true, nil
		}
		if op == parser.BoolOpAnd && !leftHandLiteral {
			return LiteralBool(false), true, nil
		}
		if op == parser.BoolOpOr && leftHandLiteral {
			return LiteralBool(true), true, nil
		}
		if op == parser.BoolOpOr && !leftHandLiteral {
			return rhs, true, nil
		}
	}

	if rightHandLiteral, rhsIsBool := getBoolLiteralFrom(rhs); rhsIsBool {
		if op == parser.BoolOpAnd && rightHandLiteral {
			return lhs, true, nil
		}
		if op == parser.BoolOpAnd && !rightHandLiteral {
			return LiteralBool(false), true, nil
		}
		if op == parser.BoolOpOr && rightHandLiteral {
			return LiteralBool(true), true, nil
		}
		if op == parser.BoolOpOr && !rightHandLiteral {
			return lhs, true, nil
		}
	}

	return nil, false, nil
}

func extractIdent(pt parser.Expression) (string, error) {
	switch pt := pt.(type) {
	case *parser.Identifier:
		return pt.Value, nil
	case *parser.Tuple:
		return extractIdent(pt.Elems[0])
	case *parser.Call:
		return extractIdent(pt.Callee)
	case *parser.Literal:
		if pt.Kind != parser.LiteralStr {
			return "", errors.New("not a string literal")
		}
		return pt.Value.(string), nil
	default:
		return "", errors.New("no identifier to extract")
	}
}
