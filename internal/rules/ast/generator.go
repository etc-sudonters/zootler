package ast

import (
	"errors"
	"github.com/etc-sudonters/substrate/slipup"
	"sudonters/zootler/internal/rules/parser"
)

type Ast struct{}

func (a *Ast) Lower(expr parser.Expression) (Node, error) {
	return parser.Transform(a, expr)
}

func (a *Ast) TransformBinOp(pt *parser.BinOp) (Node, error) {
	if pt.Op == parser.BinOpContains {
		return a.rewriteToCall("load_setting_2", pt.Right, pt.Left) // setting in block
	}

	lhs, lhsErr := parser.Transform(a, pt.Left)
	rhs, rhsErr := parser.Transform(a, pt.Right)

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
	lhs, lhsErr := parser.Transform(a, pt.Left)
	rhs, rhsErr := parser.Transform(a, pt.Right)

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
	target, err := parser.Transform(a, pt.Target)
	if err != nil {
		return nil, err
	}

	if target.Type() == AST_NODE_LITERAL {
		target := target.(*Literal)
		if target.Kind == AST_LIT_BOOL {
			target.Value = !(target.Value.(bool))
			return target, nil
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
