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
)

type Value interface {
	Type() Type
	Eq(Value) bool
	String() string
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

func (b Boolean) String() string {
	return fmt.Sprintf("%t", b.Value)
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

func (n Number) String() string {
	return fmt.Sprintf("%f", n.Value)
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

func (s String) String() string {
	return fmt.Sprintf("%q", s.Value)
}

func (t Interpreter) IsTruthy(v Value) bool {
	switch v := v.(type) {
	case Boolean:
		return bool(v.Value)
	case Number:
		return v.Value != 0
	case String:
		return v.Value != ""
	case Callable:
		if v.Arity() == 0 {
			return t.IsTruthy(v.Call(t, nil))
		}
		panic(fmt.Errorf("%d arity func must be evaluated!", v.Arity()))
	default:
		panic(errors.New("unknown truthiness kind"))
	}
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
