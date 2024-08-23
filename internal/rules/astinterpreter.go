package rules

import (
	"errors"
	"strings"
	"sudonters/zootler/internal/rules/parser"

	"github.com/etc-sudonters/substrate/slipup"
)

type AstEnvironment map[string]parser.Literal
type FunctionTracker map[string]parser.FunctionDecl
type DelayedRules map[string][]DelayedRule
type DelayedRule struct {
	Target, Name string
	Rule         parser.Expression
}

type AstInterpreter struct {
	Current     string
	Environment AstEnvironment
	Functions   FunctionTracker
	Delayed     DelayedRules
}

func (a *AstInterpreter) TransformBinOp(op *parser.BinOp) (parser.Expression, error) {
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

        settingName := flattenNestedSetting(ident.Value, str)
        return a.resolveSettingValue(settingName)
	}

	lhs, lhsErr := parser.Transform(a, op.Left)
	rhs, rhsErr := parser.Transform(a, op.Right)

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

func (a *AstInterpreter) TransformBoolOp(op *parser.BoolOp) (parser.Expression, error) {
	lhs, lhsErr := parser.Transform(a, op.Left)
	rhs, rhsErr := parser.Transform(a, op.Right)

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

func (a *AstInterpreter) TransformCall(call *parser.Call) (parser.Expression, error) {
	//SAFETY: don't call parser.Transform to visit child nodes here
	id, invalidFnIdentifier := parser.AssertAs[*parser.Identifier](call.Callee)
	if invalidFnIdentifier != nil {
		return nil, slipup.Createf("invalid function identifier: %+v", call)
	}

	if decl, decld := a.Functions[id.Value]; decld {
		return a.tryInlineCall(call, decl)
	}

	switch id.Value {
	case "at":
		target, targetErr := parser.AssertAs[*parser.Literal](call.Args[0])
		if targetErr != nil || target.Kind != parser.LiteralStr {
			return nil, slipup.Createf("expected location target name as first argument: %+v", call)
		}
		str, _ := target.AsString()
		return a.handleMacro(str, call.Args[1])
	case "here":
		return a.handleMacro(a.Current, call.Args[0])
	default:
		return nil, slipup.Createf("could not transform call %+v", call)
	}
}

func (a *AstInterpreter) TransformIdentifier(id *parser.Identifier) (parser.Expression, error) {
	if value, envDecld := a.Environment[id.Value]; envDecld {
		return &value, nil
	}

	if fn, isDecld := a.Functions[id.Value]; isDecld {
		if len(fn.Parameters) != 0 {
			return nil, slipup.Createf("expected function with 0 arguments but found %+v", fn)
		}

		return parser.MakeCall(id, nil), nil
	}

	if tok, tokErr := a.resolveTokenLiteral(id.Value); tokErr == nil {
		return tok, nil
	}

	if setting, settingErr := a.resolveSettingValue(id.Value); settingErr == nil {
		return setting, nil
	}

	return nil, slipup.Createf("could not resolve %+v to any value", id)
}

func (a *AstInterpreter) TransformSubscript(lookup *parser.Subscript) (parser.Expression, error) {
	target, targetErr := parser.AssertAs[*parser.Identifier](lookup.Target)
	index, indexErr := parser.AssertAs[*parser.Identifier](lookup.Index)

	if targetErr != nil || indexErr != nil {
		return nil, slipup.Createf("cannot inline subscript: %+v", lookup)
	}

	settingName := flattenNestedSetting(target.Value, index.Value)
	return a.resolveSettingValue(settingName)
}

func (a *AstInterpreter) TransformTuple(tup *parser.Tuple) (parser.Expression, error) {
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
			setting, settingErr := a.resolveSettingValue(i.Value)
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

	return a.makeHasCall(ident, amount)
}

func (a *AstInterpreter) TransformUnary(unary *parser.UnaryOp) (parser.Expression, error) {
	if unary.Op != parser.UnaryNot {
		panic("unknown op: %+v")
	}

	value, transformError := parser.Transform(a, unary.Target)
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

func (a *AstInterpreter) TransformLiteral(lit *parser.Literal) (parser.Expression, error) {
	str, isStr := lit.AsString()
	if !isStr {
		return lit, nil
	}

	tok, err := a.resolveTokenLiteral(str)
	if err == nil {
		return tok, nil
	}

	return lit, nil
}

func (a *AstInterpreter) resolveSettingValue(setting string) (*parser.Literal, error) {
	return nil, nil
}

func (a *AstInterpreter) resolveTokenLiteral(given string) (*parser.Literal, error) {
	return nil, nil
}

func (a *AstInterpreter) makeHasCall(token string, amount float64) (*parser.Call, error) {
	return nil, nil
}

func (a *AstInterpreter) handleMacro(target string, rule parser.Expression) (*parser.Call, error) {
	return nil, nil
}

func (a *AstInterpreter) tryInlineCall(call *parser.Call, decl parser.FunctionDecl) (parser.Expression, error) {
	return call, nil
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

func flattenNestedSetting(tiers ...string) string {
	return strings.Join(tiers, "__")
}


