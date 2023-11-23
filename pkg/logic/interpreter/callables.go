package interpreter

import (
	"fmt"

	"sudonters/zootler/pkg/rules/ast"
)

var _ Callable = (*Fn)(nil)
var _ Callable = (*PartiallyEvaluatedFn)(nil)
var _ Callable = (*BuiltIn)(nil)

type Callable interface {
	Value
	Arity() int
	Call(t Interpreter, args []Value) Value
}

type BuiltInCallable interface {
	Call(t Interpreter, args []Value) Value
}

type BuiltInFn func(Interpreter, []Value) Value

func (b BuiltInFn) Call(t Interpreter, args []Value) Value {
	return b(t, args)
}

type Fn struct {
	Params []string
	Body   ast.Expression
	Name   *ast.Identifier
}

func (f Fn) Type() Type { return CALL_TYPE }

func (f Fn) Eq(v Value) bool {
	if v == nil {
		return false
	}

	switch v := v.(type) {
	case Fn:
		return f.Name.Value == v.Name.Value

	default:
		return false
	}
}

func (f Fn) String() string {
	return fmt.Sprintf("%s", f.Name.Value)
}

func (f Fn) Arity() int {
	return len(f.Params)
}

func (f Fn) Call(t Interpreter, args []Value) Value {
	env := t.globals.Enclosed()
	for i := range args {
		env.Set(f.Params[i], args[i])
	}

	return t.Evaluate(f.Body, env)
}

type PartiallyEvaluatedFn struct {
	Body ast.Expression
	Env  Environment
	Name string
}

func (f PartiallyEvaluatedFn) String() string {
	return f.Name
}

func (f PartiallyEvaluatedFn) Type() Type { return CALL_TYPE }

func (f PartiallyEvaluatedFn) Eq(Value) bool { return false }

func (f PartiallyEvaluatedFn) Arity() int {
	return 0
}

func (f PartiallyEvaluatedFn) Call(t Interpreter, _ []Value) Value {
	return t.Evaluate(f.Body, f.Env)
}

type BuiltIn struct {
	N    int
	F    BuiltInCallable
	Name string
}

func (b BuiltIn) Arity() int {
	return b.N
}

func (b BuiltIn) Call(t Interpreter, args []Value) Value {
	return b.F.Call(t, args)
}

func (b BuiltIn) Eq(v Value) bool {
	if v == nil {
		return false
	}

	switch v := v.(type) {
	case BuiltIn:
		return v.Name == b.Name && v.N == b.N

	default:
		return false
	}
}

func (b BuiltIn) Type() Type {
	return CALL_TYPE
}

func (b BuiltIn) String() string {
	return b.Name
}

func FunctionDecl(decl, rule ast.Expression, env Environment) {
	switch decl.Type() {
	case ast.ExprIdentifier:
		decl = &ast.Call{
			Callee: decl,
			Args:   nil,
		}
		fallthrough
	case ast.ExprCall:
		call := decl.(*ast.Call)
		fn := Fn{
			Params: []string{},
			Body:   rule,
			Name:   call.Callee.(*ast.Identifier),
		}

		for i := range call.Args {
			name := call.Args[i].(*ast.Identifier)
			fn.Params = append(fn.Params, name.Value)
		}

		env.Set(fn.Name.Value, fn)
	default:
		panic(parseError("function decl must be identifier or call got %T", decl))
	}
}
