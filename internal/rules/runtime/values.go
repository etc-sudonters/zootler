package runtime

import (
	"errors"
	"fmt"
	"sudonters/zootler/internal/slipup"
)

var ErrUnsupportedType error = errors.New("unsupported type")

func ValueFrom(t interface{}) (Value, error) {
	switch v := t.(type) {
	case bool:
		return ValueFromBool(v), nil
	case int:
		return ValueFromInt(v), nil
	case float64:
		return ValueFromFloat(v), nil
	case string:
		return ValueFromStr(v), nil
	default:
		return NullValue(), ErrUnsupportedType
	}
}

func ValueOrPanic(t interface{}) Value {
	v, err := ValueFrom(t)
	if err != nil {
		panic(err)
	}

	return v
}

func NullValue() Value {
	return Value{
		kind: VAL_NULL,
	}
}

func ValueFromBool(v bool) Value {
	return Value{
		kind: VAL_BOOL,
		v:    v,
	}
}

func ValueFromInt(v int) Value {
	return Value{
		kind: VAL_INT,
		v:    v,
	}
}

func ValueFromFloat(v float64) Value {
	return Value{
		kind: VAL_FLOAT,
		v:    v,
	}
}

func ValueFromStr(v string) Value {
	return Value{
		kind: VAL_STR,
		v:    v,
	}
}

type ValueKind uint8

func (v ValueKind) String() string {
	switch v {
	case VAL_NULL:
		return "null"
	case VAL_INT:
		return "int"
	case VAL_FLOAT:
		return "float"
	case VAL_BOOL:
		return "bool"
	case VAL_STR:
		return "str"
	default:
		panic(fmt.Errorf("unknown value kind: %02X", uint8(v)))
	}
}

const (
	VAL_NULL ValueKind = iota
	VAL_INT
	VAL_FLOAT
	VAL_BOOL
	VAL_STR
)

type Value struct {
	kind ValueKind
	v    any
}
type Values []Value

func (v Value) Eq(o Value) bool {
	if v.kind == VAL_INT && o.kind == VAL_FLOAT {
		return float64(v.v.(int)) == o.v.(float64)
	}

	if v.kind == VAL_FLOAT && o.kind == VAL_INT {
		return v.v.(float64) == float64(o.v.(int))
	}

	if v.kind != o.kind {
		return false
	}

	switch v.kind {
	case VAL_NULL:
		panic("null dereference")
	case VAL_INT:
		return v.v.(int) == o.v.(int)
	case VAL_FLOAT:
		return v.v.(float64) == o.v.(float64)
	case VAL_BOOL:
		return v.v.(bool) == o.v.(bool)
	case VAL_STR:
		return v.v.(string) == o.v.(string)
	default:
		return false
	}
}

func (v Value) Lt(o Value) bool {
	switch v.kind {
	case VAL_NULL:
		panic("null dereference")
	case VAL_INT:
		if o.kind == VAL_FLOAT {
			return float64(v.v.(int)) < o.v.(float64)
		}
		return v.v.(int) < o.v.(int)
	case VAL_FLOAT:
		if o.kind == VAL_INT {
			return v.v.(float64) == float64(o.v.(int))
		}
		return v.v.(float64) < o.v.(float64)
	}

	panic(fmt.Errorf("unorderable types: '%s' and '%s'", v.kind, o.kind))
}

func (v Value) Truthy() bool {
	switch v.kind {
	case VAL_NULL:
		return false
	case VAL_FLOAT:
		return v.v.(float64) != 0
	case VAL_INT:
		return v.v.(int) != 0
	case VAL_BOOL:
		return v.v.(bool)
	}

	return true
}

func (v Value) Unwrap() interface{} {
	return v.v
}

func (v Value) AsInt() (int, error) {
	switch i := v.v.(type) {
	case int:
		return i, nil
	case float64:
		return int(i), nil
	default:
		return 0, slipup.Describef(ErrUnsupportedType, "cannot convert '%+v' to int", v)
	}
}

func (v Value) AsFloat() (float64, error) {
	switch i := v.v.(type) {
	case int:
		return float64(i), nil
	case float64:
		return i, nil
	default:
		return 0, slipup.Describef(ErrUnsupportedType, "cannot convert '%+v' to int", v)
	}
}

func (v Value) AsBool() (bool, error) {
	switch i := v.v.(type) {
	case bool:
		return i, nil
	default:
		return false, slipup.Describef(ErrUnsupportedType, "cannot convert '%+v' to bool", v)
	}
}

func (v Value) AsStr() (string, error) {
	switch i := v.v.(type) {
	case string:
		return i, nil
	default:
		return "", slipup.Describef(ErrUnsupportedType, "cannot convert '%+v' to string", v)
	}
}
