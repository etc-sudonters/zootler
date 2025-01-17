package objects

import (
	"fmt"
	"math"
)

const (
	quiet  Object = 0x7FFC000000000000
	u8     Object = 0x7FFC800000000000
	zbool  Object = 0x7FFCB00000000000
	ztrue  Object = 0x7FFCBA0000000000
	zfalse Object = 0x7FFCBB0000000000
	ptr32  Object = 0x7FFCC00000000000
	str32  Object = 0x7FFCD00000000000

	U8      = nantag(u8)
	Boolean = nantag(zbool)
	Ptr32   = nantag(ptr32)
	Str32   = nantag(str32)
	F64     = nantag(^quiet)
	Func    = nantag(ptr32 | Object(PtrFunc)<<32)

	PtrToken   PtrTag = 0xE0
	PtrEvent   PtrTag = 0xE0
	PtrSetting PtrTag = 0xC0
	PtrFunc    PtrTag = 0xF0
	PtrLoc     PtrTag = 0xAC
	PtrEdge    PtrTag = 0x23

	True         = ztrue
	False        = zfalse
	Null  Object = 0
)

type Object uint64
type PtrTag uint8
type nantag uint64

func (this Object) Is(tag nantag) bool {
	t := Object(tag)
	switch tag {
	case F64:
		return this&quiet != quiet
	default:
		return this&t == t
	}
}

func PackU8(v uint8) Object {
	return pack(u8 | Object(v))
}

func UnpackU8(bits Object) uint8 {
	mask := (quiet | u8)
	return uint8(bits & (^mask))
}

func PackFloat64(v float64) Object {
	return Object(v)
}

func pack(u64 Object) Object {
	return quiet | u64
}

func UnpackFloat64(v Object) float64 {
	return math.Float64frombits(uint64(v))
}

func PackStr32(len int, offset int) Object {
	if len > 0xFF || offset > 0xFFFFFFFF {
		panic(fmt.Errorf("Str{len: %d, offset: %d} is too big to pack", len, offset))
	}

	return pack(str32 | Object(len)<<32 | Object(offset))
}

func UnpackStr32(bits Object) (int, int) {
	mask := quiet | str32
	bits = bits & (^mask)

	len := uint8(bits >> 32)
	off := uint32(bits)
	return int(len), int(off)
}

func PackPtr32(ptr uint32) Object {
	return pack(ptr32 | Object(ptr))
}

func UnpackPtr32(bits Object) uint32 {
	mask := quiet | ptr32
	return uint32(bits & (^mask))
}

func PackTaggedPtr32(tag PtrTag, ptr uint32) Object {
	return pack(ptr32 | Object(tag)<<32 | Object(ptr))
}

func UnpackTaggedPtr32(bits Object) (uint8, uint32) {
	mask := quiet | ptr32
	bits = bits & (^mask)

	tag := uint8(bits >> 32)
	ptr := uint32(bits)
	return tag, ptr
}

func UnpackBool(v Object) bool {
	return v == ztrue
}

func PackBool(b bool) Object {
	if b {
		return True
	}
	return False
}
