package preprocessor

import (
	"errors"
	"sudonters/zootler/internal/rules/parser"
	"sudonters/zootler/internal/rules/visitor"
	"unicode"

	"github.com/etc-sudonters/substrate/slipup"
)

type state struct {
	location string
	parent   *P
}

type P struct {
	// at rules sink
	Delayed   DelayedRules
	Functions FunctionTable
	Env       ValuesTable
}

func (parent *P) Process(origin string, rule parser.Expression) (parser.Expression, error) {
	return visitor.Transform(state{origin, parent}, rule)
}

func (s state) TransformBinOp(op *parser.BinOp) (parser.Expression, error) {
	if areSameIdentifier(op.Left, op.Right) {
		return parser.BoolLiteral(op.Op == parser.BinOpEq), nil
	}

	if areDifferentIdentifier(op.Left, op.Right) {
		return parser.BoolLiteral(op.Op == parser.BinOpNotEq), nil
	}

	if op.Op == parser.BinOpContains && op.Left.Type() == parser.ExprLiteral && op.Right.Type() == parser.ExprIdentifier {
		str, isStr := assertAsStr(op.Left)
		if !isStr {
			return nil, slipup.Createf("expected construction of 'str in identifier', received: %+v", op)
		}

		ident := op.Right.(*parser.Identifier)

		return s.parent.Env.MustResolveNested(ident.Value, str)
	}

	lhs, lhsErr := visitor.Transform(s, op.Left)
	rhs, rhsErr := visitor.Transform(s, op.Right)

	if joined := errors.Join(lhsErr, rhsErr); joined != nil {
		return nil, slipup.Describef(joined, "while transforming boolop %+v", op)
	}

	lhl, lhlErr := parser.AssertAs[*parser.Literal](lhs)
	rhl, rhlErr := parser.AssertAs[*parser.Literal](rhs)

	if lhlErr != nil || rhlErr != nil {
		return &parser.BinOp{
			Left:  lhs,
			Op:    op.Op,
			Right: rhs,
		}, nil
	}

	switch op.Op {
	case parser.BinOpEq:
		return parser.BoolLiteral(lhl.Value == rhl.Value), nil
	case parser.BinOpNotEq:
		return parser.BoolLiteral(lhl.Value != rhl.Value), nil
	case parser.BinOpLt:
		lhn, isLhn := lhl.AsNumber()
		rhn, isRhn := rhl.AsNumber()
		if !isLhn || !isRhn {
			return nil, slipup.Createf("unorderable values: %+v", op)
		}
		return parser.BoolLiteral(lhn < rhn), nil
	default:
		return nil, slipup.Createf("unknown binop: %+v", op)
	}
}

func (s state) TransformBoolOp(op *parser.BoolOp) (parser.Expression, error) {
	lhs, lhsErr := visitor.Transform(s, op.Left)
	rhs, rhsErr := visitor.Transform(s, op.Right)

	if joined := errors.Join(lhsErr, rhsErr); joined != nil {
		return nil, slipup.Describef(joined, "while transforming boolop %+v", op)
	}

	lhb, isLhb := assertAsBool(lhs)
	rhb, isRhb := assertAsBool(rhs)

	if !isLhb && !isRhb {
		return &parser.BoolOp{
			Left:  lhs,
			Op:    op.Op,
			Right: rhs,
		}, nil
	}

	if isLhb && isRhb {
		if op.Op == parser.BoolOpAnd {
			return parser.BoolLiteral(lhb && rhb), nil
		}

		return parser.BoolLiteral(lhb || rhb), nil
	}

	if isLhb {
		if op.Op == parser.BoolOpAnd && !lhb {
			return parser.BoolLiteral(false), nil
		}
		if op.Op == parser.BoolOpOr && lhb {
			return parser.BoolLiteral(true), nil
		}
	}

	if isRhb {
		if op.Op == parser.BoolOpAnd && !rhb {
			return parser.BoolLiteral(false), nil
		}
		if op.Op == parser.BoolOpOr && rhb {
			return parser.BoolLiteral(true), nil
		}
	}

	return &parser.BoolOp{
		Left:  lhs,
		Op:    op.Op,
		Right: rhs,
	}, nil
}

func (s state) TransformCall(call *parser.Call) (parser.Expression, error) {
	//SAFETY: don't call visitor.Transform(s, ...) to visit child nodes here
	id, invalidFnIdentifier := parser.AssertAs[*parser.Identifier](call.Callee)
	if invalidFnIdentifier != nil {
		return nil, slipup.Createf("invalid function identifier: %+v", call)
	}

	if decl, _ := s.parent.Functions.Retrieve(id.Value); decl != nil {
		return s.tryInlineCall(call, decl)
	}

	switch id.Value {
	case "at":
		target, targetErr := parser.AssertAs[*parser.Literal](call.Args[0])
		if targetErr != nil || target.Kind != parser.LiteralStr {
			return nil, slipup.Createf("expected location target name as first argument: %+v", call)
		}
		targetName, _ := target.AsString()
		return s.handleMacro(targetName, call.Args[1])
	case "here":
		return s.handleMacro(s.location, call.Args[0])
	default:
		return nil, slipup.Createf("'%s' is not a known function", call.Callee)
	}
}

func (s state) TransformIdentifier(id *parser.Identifier) (parser.Expression, error) {
	if value, decld := s.parent.Env.Resolve(id.Value); decld {
		return value, nil
	}

	if fn, _ := s.parent.Functions.Retrieve(id.Value); fn != nil {
		if len(fn.Parameters) != 0 {
			return nil, slipup.Createf("expected function with 0 arguments but found %+v", fn)
		}

		return parser.MakeCall(id, nil), nil
	}

	if tok, tokErr := s.parent.Env.ResolveAsToken(id.Value); tokErr == nil {
		return tok, nil
	}

	if setting, exists := s.parent.Env.Resolve(id.Value); exists {
		return setting, nil
	}

	return nil, slipup.Createf("could not resolve %+v to any value", id)
}

func (s state) TransformSubscript(lookup *parser.Subscript) (parser.Expression, error) {
	target, targetErr := parser.AssertAs[*parser.Identifier](lookup.Target)
	index, indexErr := parser.AssertAs[*parser.Identifier](lookup.Index)

	if targetErr != nil || indexErr != nil {
		return nil, slipup.Createf("cannot inline subscript: %+v", lookup)
	}

	return s.parent.Env.MustResolveNested(target.Value, index.Value)
}

func (s state) TransformTuple(tup *parser.Tuple) (parser.Expression, error) {
	if len(tup.Elems) != 2 {
		return tup, slipup.Createf("expected exactly 2 elements -- identifier and amount -- for tuple expression\nrecieved: %+v", tup)
	}

	ident, identErr := parser.Unify(
		tup.Elems[0],
		func(i *parser.Identifier) (string, error) { return i.Value, nil },
		func(l *parser.Literal) (string, error) {
			if str, ok := l.AsString(); ok {
				return str, nil
			}
			return "", slipup.Createf("expected string literal, got %+v", tup.Elems[0])
		})

	amount, amountErr := parser.Unify(
		tup.Elems[1],
		func(i *parser.Identifier) (float64, error) {
			setting, settingErr := s.parent.Env.MustResolve(i.Value)
			if settingErr != nil {
				return -1, slipup.Describef(settingErr, "while resolving identifier %+v", i)
			}

			if number, isNumber := setting.AsNumber(); isNumber {
				return number, nil
			}

			return -1, slipup.Createf("expected to resolve %+v as number but resolved %+v", i, setting)
		},
		func(l *parser.Literal) (float64, error) {
			if number, isNumber := l.AsNumber(); isNumber {
				return number, nil
			}

			return -1, slipup.Createf("expected number literal but recieved %+v", l)
		})

	if joined := errors.Join(identErr, amountErr); joined != nil {
		return nil, joined
	}

	return s.makeHasCall(ident, amount)
}

func (s state) TransformUnary(unary *parser.UnaryOp) (parser.Expression, error) {
	if unary.Op != parser.UnaryNot {
		panic("unknown op: %+v")
	}

	value, transformError := visitor.Transform(s, unary.Target)
	if transformError != nil {
		return nil, slipup.Describef(transformError, "while transforming %+v", unary.Target)
	}

	switch value := value.(type) {
	case *parser.Literal:
		b, isBool := value.AsBool()
		if !isBool {
			return nil, slipup.Createf("cannot negate non boolean/call %+v", value)
		}
		return parser.BoolLiteral(!b), nil
	case *parser.Call:
		return &parser.UnaryOp{
			Op:     unary.Op,
			Target: value,
		}, nil
	default:
		return nil, slipup.Createf("cannot negate non boolean/call %+v", value)
	}

}

func (s state) TransformLiteral(lit *parser.Literal) (parser.Expression, error) {
	str, isStr := lit.AsString()
	if !isStr {
		return lit, nil
	}

	tok, err := s.parent.Env.ResolveAsToken(str)
	if err == nil {
		return tok, nil
	}

	return lit, nil
}

func (s state) resolveTokenLiteral(given string) (*parser.Literal, error) {
	first := []rune(given)[0]
	if !unicode.IsUpper(first) {
		return nil, slipup.Createf("token literals must begin with upper case letter: %s", given)
	}

	value, exists := s.parent.Env.Resolve(given)
	if !exists {
		return nil, slipup.Createf("token %s not found", given)
	}

	if value.Kind != parser.LiteralToken {
		return nil, slipup.Createf("expected token but resolved: %+v", value)
	}

	return value, nil
}

func (s state) makeHasCall(token string, amount float64) (*parser.Call, error) {
	return parser.MakeCallSplat(parser.Identify("has"), parser.StringLiteral(token), parser.NumberLiteral(amount)), nil
}

func (s state) handleMacro(target string, rule parser.Expression) (*parser.Call, error) {
	tokenname := s.parent.Delayed.Add(target, rule)
	return parser.MakeCallSplat(
		parser.Identify("has"),
		parser.StringLiteral(tokenname),
		parser.NumberLiteral(1),
	), nil
}

func (s state) tryInlineCall(call *parser.Call, decl *parser.FunctionDecl) (parser.Expression, error) {
	// built in, impossible to inline
	if decl.Body == nil {
		return call, nil
	}

	if len(call.Args) != len(decl.Parameters) {
		return nil, slipup.Createf("%s expected %d arguments but received %d", decl.Identifier, len(decl.Parameters), len(call.Args))
	}

	replacements := make(map[string]parser.Expression, len(call.Args))

	for i := range len(call.Args) {
		replacements[decl.Parameters[i]] = call.Args[i]
	}

	renamed, renameErr := visitor.Transform(rewritefuncparams(replacements), decl.Body)
	if renameErr != nil {
		return nil, slipup.Describef(renameErr, "while inlining %s", call.Callee)
	}

	return visitor.Transform(s, renamed)
}

func assertAsBool(ast parser.Expression) (bool, bool) {
	lit, litErr := parser.AssertAs[*parser.Literal](ast)
	if litErr != nil {
		return false, false
	}

	return lit.AsBool()
}

func assertAsStr(ast parser.Expression) (string, bool) {
	lit, litErr := parser.AssertAs[*parser.Literal](ast)
	if litErr != nil {
		return "", false
	}

	return lit.AsString()
}

func areSameIdentifier(lhs, rhs parser.Expression) bool {
	lhi, lhiErr := parser.AssertAs[*parser.Identifier](lhs)
	rhi, rhiErr := parser.AssertAs[*parser.Identifier](rhs)

	if lhiErr != nil || rhiErr != nil {
		return false
	}

	return lhi.Value == rhi.Value
}

func areDifferentIdentifier(lhs, rhs parser.Expression) bool {
	lhi, lhiErr := parser.AssertAs[*parser.Identifier](lhs)
	rhi, rhiErr := parser.AssertAs[*parser.Identifier](rhs)

	if lhiErr != nil || rhiErr != nil {
		return false
	}

	return lhi.Value != rhi.Value
}
