package ast

import (
	"errors"
	"sudonters/zootler/internal/rules/parser"
	"sudonters/zootler/internal/rules/visitor"

	"github.com/etc-sudonters/substrate/skelly/stack"
	"github.com/etc-sudonters/substrate/slipup"
)

type IdentityChecker interface {
	Special(string, bool) Node
	Macro(string) bool
}

func NewGenerator(idents IdentityChecker) *AstGenerator {
	return &AstGenerator{
		calls:  stack.S[struct{}]{},
		idents: idents,
	}
}

type AstGenerator struct {
	calls  stack.S[struct{}]
	idents IdentityChecker
}

func (a *AstGenerator) ClearStacks() {
	a.calls = stack.S[struct{}]{}
}

func (a *AstGenerator) TransformBinOp(pt *parser.BinOp) (Node, error) {
	if pt.Op == parser.BinOpContains {
		return a.rewriteToCall("load_setting", pt.Left, pt.Right)
	}

	lhs, lhsErr := visitor.Transform(a, pt.Left)
	rhs, rhsErr := visitor.Transform(a, pt.Right)

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

func (a *AstGenerator) TransformBoolOp(pt *parser.BoolOp) (Node, error) {
	lhs, lhsErr := visitor.Transform(a, pt.Left)
	rhs, rhsErr := visitor.Transform(a, pt.Right)

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

func (a *AstGenerator) TransformTuple(pt *parser.Tuple) (Node, error) {
	if len(pt.Elems) != 2 {
		return nil, slipup.Createf("invalid tuple construction: %+v", pt)
	}

	return a.rewriteToCall("has", pt.Elems[0], pt.Elems[1])
}

func (a *AstGenerator) TransformCall(pt *parser.Call) (Node, error) {
	callee := parser.MustAssertAs[*parser.Identifier](pt.Callee)
	return a.rewriteToCall(callee.Value, pt.Args...)
}

func (a *AstGenerator) TransformSubscript(pt *parser.Subscript) (Node, error) {
	return a.rewriteToCall("load_setting_2",
		parser.MustAssertAs[*parser.Identifier](pt.Target),
		parser.MustAssertAs[*parser.Identifier](pt.Index),
	)
}

func (a *AstGenerator) TransformUnary(pt *parser.UnaryOp) (Node, error) {
	target, err := visitor.Transform(a, pt.Target)
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

func (a *AstGenerator) TransformLiteral(pt *parser.Literal) (Node, error) {
	if pt.Kind == parser.LiteralStr {
		str := pt.Value.(string)
		if replacement := a.idents.Special(str, a.visitingCall()); replacement != nil {
			return replacement, nil
		}
	}

	lit := &Literal{
		Value: pt.Value,
	}

	var exists bool
	lit.Kind, exists = ptToAstLitKind[pt.Kind]
	if !exists {
		panic(slipup.Createf("unknown literal kind %+v", pt))
	}

	return lit, nil
}

var ptToAstLitKind = map[parser.LiteralKind]AstLiteralKind{
	parser.LiteralBool: AST_LIT_BOOL,
	parser.LiteralNum:  AST_LIT_NUM,
	parser.LiteralStr:  AST_LIT_STR,
}

func (a *AstGenerator) replaceCompositeSetting(block, name string) (Node, error) {
	return a.rewriteToCall("load_setting_2", parser.Identify(block), parser.Identify(name))
}

func (a *AstGenerator) TransformIdentifier(pt *parser.Identifier) (Node, error) {
	if replacement := a.idents.Special(pt.Value, a.visitingCall()); replacement != nil {
		return replacement, nil
	}

	ident := &Identifier{
		Name: pt.Value,
		Kind: AST_IDENT_UNK,
	}

	return ident, nil
}

func (a *AstGenerator) rewriteToCall(name string, inputs ...parser.Expression) (*Call, error) {
	stopCall := a.startCallVisit()
	defer stopCall()
	args, argErr := visitor.TransformAll(a, inputs)
	if argErr != nil {
		return nil, argErr
	}

	call := &Call{
		Callee: name,
		Args:   args,
		Macro:  a.idents.Macro(name),
	}

	return call, nil
}

func (a *AstGenerator) startCallVisit() func() {
	a.calls.Push(struct{}{})
	return func() { a.calls.Pop() }
}

func (a *AstGenerator) visitingCall() bool {
	return a.calls.Len() != 0
}
