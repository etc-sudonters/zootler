package interpreter

import (
	"errors"
	"fmt"
	"sudonters/zootler/pkg/rules/ast"
)

var _ Evaluation[Value] = interpreter{}

var UnknownIdentifierErr = errors.New("unknown identifier")

func New(globals Environment) interpreter {
	return interpreter{globals}
}

type interpreter struct {
	globals Environment
}

func (t interpreter) Evaluate(ex ast.Expression, env Environment) Value {
	return Evaluate(t, ex, env)
}

func (t interpreter) EvalAttrAccess(access *ast.AttrAccess, env Environment) Value {
	panic("not implemented") // TODO: Implement
}

func (t interpreter) EvalBinOp(op *ast.BinOp, env Environment) Value {
	left := t.Evaluate(op.Left, env)
	right := t.Evaluate(op.Right, env)

	switch op.Op {
	case ast.BinOpEq:
		return Boolean(left.Eq(right))
	case ast.BinOpNotEq:
		return Boolean(!left.Eq(right))
	case ast.BinOpLt:
		if left.Type() == right.Type() && left.Type() == NUM_TYPE {
			l := left.(Number)
			r := right.(Number)
			return Boolean(l.Value < r.Value)
		}
		panic(fmt.Errorf("only numbers can be compared not %T and %T", left, right))
	}
	panic("unreachable")
}

func (t interpreter) EvalBoolOp(op *ast.BoolOp, env Environment) Value {
	left := t.Evaluate(op.Left, env)

	if op.Op == ast.BoolOpOr {
		if IsTruthy(left) {
			return left
		}
	} else {
		if !IsTruthy(left) {
			return left
		}
	}

	return t.Evaluate(op.Right, env)
}

func (t interpreter) EvalBoolean(bool *ast.Boolean) Value {
	return Boolean(bool.Value)
}

func (t interpreter) EvalCall(call *ast.Call, env Environment) Value {
	callee := t.Evaluate(call.Callee, env)
	fn, ok := callee.(Callable)
	if !ok {
		panic(fmt.Errorf("%v is not callable", callee))
	}

	if fn.Arity() != len(call.Args) {
		panic(fmt.Errorf(
			"Expected %d arguments but got %d",
			fn.Arity(),
			len(call.Args),
		))
	}

	args := make([]Value, len(call.Args))
	for i := range args {
		args[i] = t.Evaluate(call.Args[i], env)
	}

	return fn.Call(t, args)
}

func (t interpreter) EvalIdentifier(ident *ast.Identifier, env Environment) Value {
	v, ok := env.Get(ident.Value)
	if !ok {
		panic(fmt.Errorf("%w: %q", UnknownIdentifierErr, ident.Value))
	}

	return v
}

func (t interpreter) EvalNumber(num *ast.Number) Value {
	return Number{Value: num.Value}
}

func (t interpreter) EvalString(str *ast.String) Value {
	return String{Value: str.Value}
}

func (t interpreter) EvalSubscript(subscript *ast.Subscript, env Environment) Value {
	panic("not implemented") // TODO: Implement
}

func (t interpreter) EvalTuple(tup *ast.Tuple, env Environment) Value {
	panic("not implemented") // TODO: Implement
}

func (t interpreter) EvalUnary(unary *ast.UnaryOp, env Environment) Value {
	panic("not implemented") // TODO: Implement
}
