package packed

import "math"

const (
	quiet  uint64 = 0x7FFC000000000000
	zbool  uint64 = 0x7FFCB00000000000
	ztrue  uint64 = 0x7FFCBA0000000000
	zfalse uint64 = 0x7FFCBB0000000000
	ptr32  uint64 = 0x7FFCC00000000000
	str32  uint64 = 0x7FFCD00000000000
	u32    uint64 = 0x7FFCE00000000000
	i32    uint64 = 0x7FFCF00000000000

	Boolean = Tag(zbool)
	Ptr32   = Tag(ptr32)
	Str32   = Tag(str32)
	U32     = Tag(u32)
	I32     = Tag(i32)
	F64     = Tag(^quiet)
)

var (
	packedTrue  = pack(ztrue)
	packedFalse = pack(zfalse)
)

type Value uint64
type Tag uint64

func (this Value) Is(tag Tag) bool {
	bits := uint64(this)
	t := uint64(tag)
	switch tag {
	case F64:
		return bits&quiet != quiet
	default:
		return bits&t == t
	}
}

func PackFloat(v float64) Value {
	return Value(v)
}

func pack(u64 uint64) Value {
	return Value(quiet | u64)
}

func UnpackFloat(v Value) float64 {
	return math.Float64frombits(uint64(v))
}

func PackI32(v int32) Value {
	return pack(i32 | uint64(uint32(v)))
}

func UnpackI32(v Value) int32 {
	bits := uint64(v)
	mask := quiet | i32
	bits = bits & (^mask)
	return int32(uint32(bits))
}

func PackU32(v uint32) Value {
	return pack(u32 | uint64(v))
}

func UnpackU32(v Value) uint32 {
	bits := uint64(v)
	mask := quiet | u32
	return uint32(bits & (^mask))
}

func PackStr32(len uint8, offset uint32) Value {
	return pack(str32 | uint64(len)<<32 | uint64(offset))
}

func UnpackStr32(v Value) (uint8, uint32) {
	bits := uint64(v)
	mask := quiet | str32
	bits = bits & (^mask)

	len := uint8(bits >> 32)
	off := uint32(bits)
	return len, off
}

func PackPtr32(ptr uint32) Value {
	return pack(ptr32 | uint64(ptr))
}

func UnpackPtr32(v Value) uint32 {
	bits := uint64(v)
	mask := quiet | ptr32
	return uint32(bits & (^mask))
}

func PackTaggedPtr32(tag uint8, ptr uint32) Value {
	return pack(ptr32 | uint64(tag)<<32 | uint64(ptr))
}

func UnpackTaggedPtr32(v Value) (uint8, uint32) {
	bits := uint64(v)
	mask := quiet | ptr32
	bits = bits & (^mask)

	tag := uint8(bits >> 32)
	ptr := uint32(bits)
	return tag, ptr
}

func UnpackBool(v Value) bool {
	return v == packedTrue
}

func PackTrue() Value {
	return packedTrue
}

func PackFalse() Value {
	return packedFalse
}
