package interpreter

import (
	"errors"
	"fmt"
	"sudonters/zootler/internal/entity"
	"sudonters/zootler/pkg/rules/ast"
)

type Type int

const (
	_ Type = iota
	BOOL_TYPE
	NUM_TYPE
	STR_TYPE
	CALL_TYPE
	ENT_TYPE
)

type Value interface {
	Type() Type
	Eq(Value) bool
}

type Boolean struct {
	Value bool
}

func (b Boolean) Type() Type { return BOOL_TYPE }

func (b Boolean) Eq(v Value) bool {
	if v == nil {
		return false
	}

	switch v := v.(type) {
	case Boolean:
		return v == b
	default:
		return false
	}
}

type Number struct {
	Value float64
}

func (n Number) Type() Type { return NUM_TYPE }
func (n Number) Eq(v Value) bool {
	if v == nil {
		return false
	}

	switch v := v.(type) {
	case Number:
		return v == n
	default:
		return false
	}
}

type String struct {
	Value string
}

func (s String) Type() Type { return STR_TYPE }
func (s String) Eq(v Value) bool {
	if v == nil {
		return false
	}

	switch v := v.(type) {
	case String:
		return v == s
	default:
		return false
	}
}

type Entity struct {
	Value entity.Model
}

func (e Entity) Type() Type {
	return ENT_TYPE
}

func (e Entity) Eq(o Value) bool {
	if o == nil {
		return false
	}

	switch o := o.(type) {
	case Entity:
		return e.Value == o.Value
	default:
		return false
	}
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

func (f Fn) Arity() int {
	return len(f.Params)
}

func (f Fn) Call(t interpreter, args []Value) Value {
	env := t.globals.Enclosed()
	for i := range args {
		env.Set(f.Params[i], args[i])
	}

	return t.Evaluate(f.Body, env)
}

func (t interpreter) IsTruthy(v Value) bool {
	switch v := v.(type) {
	case Boolean:
		return bool(v.Value)
	case Number:
		return v.Value != 0
	case String:
		return v.Value != ""
	case Entity:
		if v.Value == entity.INVALID_ENTITY {
			// this shouldn't ever happen but we need to catch it
			panic(entity.ErrInvalidEntity)
		}
		return true
	case Callable:
		if v.Arity() == 0 {
			return t.IsTruthy(v.Call(t, nil))
		}
		panic(fmt.Errorf("%d arity func must be evaluated!", v.Arity()))
	default:
		panic(errors.New("unknown truthiness kind"))
	}
}

type Callable interface {
	Arity() int
	Call(t interpreter, args []Value) Value
}

type BuiltIn struct {
	N int
	F func(t interpreter, args []Value) Value
}

func (b BuiltIn) Arity() int {
	return b.N
}

func (b BuiltIn) Call(t interpreter, args []Value) Value {
	return b.F(t, args)
}

func Box(v any) Value {
	switch v := v.(type) {
	case bool:
		return Boolean{Value: v}
	case float64:
		return Number{Value: v}
	case int:
		return Number{Value: float64(v)}
	case string:
		return String{Value: v}
	case entity.Model:
		return Entity{Value: v}
	default:
		panic(fmt.Errorf("cannot box %T", v))
	}
}

func Unbox(v Value) any {
	switch v := v.(type) {
	case Number:
		return v.Value
	case Boolean:
		return v.Value
	case String:
		return v.Value
	default:
		panic(fmt.Errorf("cannot unbox %T", v))
	}
}

func CanLiteralfy(v Value) bool {
	return v.Type() == NUM_TYPE || v.Type() == BOOL_TYPE || v.Type() == STR_TYPE
}

func Literalify(v any) ast.Expression {
	if v == nil {
		panic(errors.New("nil value"))
	}
	switch v := v.(type) {
	case float64:
		return &ast.Literal{Value: v, Kind: ast.LiteralNum}
	case Number:
		return Literalify(v.Value)
	case bool:
		return &ast.Literal{Value: v, Kind: ast.LiteralBool}
	case Boolean:
		return Literalify(v.Value)
	case string:
		return &ast.Literal{Value: v, Kind: ast.LiteralStr}
	case String:
		return Literalify(v.Value)
	default:
		panic(fmt.Errorf("cannot literalfy %T", v))
	}
}

func ReifyLiteral(expr *ast.Literal) Value {
	switch expr.Kind {
	case ast.LiteralNum:
		return Number{Value: expr.Value.(float64)}
	case ast.LiteralBool:
		return Boolean{Value: expr.Value.(bool)}
	case ast.LiteralStr:
		return String{Value: expr.Value.(string)}
	default:
		panic(fmt.Errorf("cannot Deliteralfy %s", expr.Kind))
	}
}

type PartiallyEvaluatedFn struct {
	Body ast.Expression
	Env  Environment
	Name string
}

func (f PartiallyEvaluatedFn) Type() Type { return CALL_TYPE }

func (f PartiallyEvaluatedFn) Eq(Value) bool { return false }

func (f PartiallyEvaluatedFn) Arity() int {
	return 0
}

func (f PartiallyEvaluatedFn) Call(t interpreter, _ []Value) Value {
	return t.Evaluate(f.Body, f.Env)
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
