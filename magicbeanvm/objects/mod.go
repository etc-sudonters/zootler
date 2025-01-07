package objects

import "sudonters/zootler/magicbeanvm/nan"

type Callable = func([]Object) (Object, error)

type Kind string
type Boolean bool
type BuiltInFunc struct {
	Name   string
	Func   Callable
	Params int
}
type Number float64
type Ptr nan.PackedValue
type String string

type PtrTag uint8

const (
	_          PtrTag = 0
	PtrToken          = 0x0A
	PtrSetting        = 0x0B
)

func Pointer(ptr uint16, tag PtrTag) Ptr {
	u32 := uint32(tag)
	u32 = u32 | uint32(ptr<<8)
	return Ptr(nan.PackPtr(u32))
}

func UnpackPointer(ptr Ptr) (uint16, PtrTag) {
	u32, isPtr := nan.PackedValue(ptr).Pointer()
	if !isPtr {
		panic("derefencing non-pointer")
	}

	tag := uint8(u32)
	index := uint16(u32 >> 8)
	return index, PtrTag(tag)
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

func (this *BuiltInFunc) Kind() Kind {
	return BUILT_IN
}

func (this Ptr) Kind() Kind {
	return POINTER
}
