package objects

import (
	"fmt"
	"sudonters/zootler/mido/nan"
)

type BuiltInFn func([]Object) (Object, error)
type BuiltInFunctions []BuiltInFunction
type BuiltInFunctionDefs []BuiltInFunctionDef

type Kind string
type Boolean bool
type BuiltInFunction struct {
	Name   string
	Params int
	Fn     BuiltInFn
}

type BuiltInFunctionDef struct {
	Name   string
	Params int
}

type Number float64
type Ptr nan.PackedValue
type String string

type OpaquePointer uint16
type PtrTag uint8

const (
	_          PtrTag = 0
	PtrToken          = 0x0A
	PtrSetting        = 0x0B
)

func (this Ptr) String() string {
	ptr, _ := nan.PackedValue(this).Pointer()
	return fmt.Sprintf("0x%04X", ptr)
}

func Pointer(ptr OpaquePointer, tag PtrTag) Ptr {
	u32 := uint32(tag) << 16
	u32 = u32 | uint32(ptr)
	return Ptr(nan.PackPtr(u32))
}

func UnpackPointer(opaque Ptr) (OpaquePointer, PtrTag) {
	u32, isPtr := nan.PackedValue(opaque).Pointer()
	if !isPtr {
		panic("derefencing non-pointer")
	}

	tag := uint8(u32 >> 16)
	ptr := uint16(u32)
	return OpaquePointer(ptr), PtrTag(tag)
}

const (
	_        Kind = ""
	BOOLEAN       = "BOOLEAN"
	BUILT_IN      = "BUILT_IN"
	NUMBER        = "NUMBER"
	POINTER       = "POINTER"
	STRING        = "STRING"
)

type Object interface {
	Kind() Kind
}

func (this Boolean) Kind() Kind {
	return BOOLEAN
}

func (this String) Kind() Kind {
	return STRING
}

func (this Number) Kind() Kind {
	return NUMBER
}

func (this *BuiltInFunction) Kind() Kind {
	return BUILT_IN
}

func (this Ptr) Kind() Kind {
	return POINTER
}
