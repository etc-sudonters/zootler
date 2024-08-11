package bytecode

import "fmt"

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
	default:
		panic(fmt.Errorf("unknown value kind: %02X", uint8(v)))
	}
}

const (
	VAL_NULL ValueKind = iota
	VAL_INT
	VAL_FLOAT
	VAL_BOOL
)

type Value struct {
	kind ValueKind
	v    any
}
type Values []Value

func (v Value) Eq(o Value) bool {
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
	default:
		return false
	}
}

func (v Value) Lt(o Value) bool {
	if v.kind != o.kind {
		panic(fmt.Errorf("unorderable types: '%s' and '%s'", v.kind, o.kind))
	}

	switch v.kind {
	case VAL_NULL:
		panic("null dereference")
	case VAL_INT:
		return v.v.(int) < o.v.(int)
	case VAL_FLOAT:
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
