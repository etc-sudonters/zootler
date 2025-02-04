package objects

import (
	"fmt"
	"math"
)

type Object bits

func (this Object) Is(mask bits) bool {
	return bits(this)&mask == mask
}

func (this Object) Truthy() bool {
	if this == PackedFalse || this == Null {
		return false
	}
	return true
}

func (this Object) Type() string {
	field := bits(this)
	if !math.IsNaN(math.Float64frombits(uint64(field))) {
		return "F64"
	}
	if field&Ptr32 == Ptr32 {
		return "Ptr32"
	}
	if field&Str32 == Str32 {
		return "Str32"
	}
	if field&Array == Array {
		return "Array"
	}
	if field&I32 == I32 {
		return "I32"
	}
	if field&U32 == U32 {
		return "U32"
	}
	if field&Bool == Bool {
		return "Bool"
	}
	if field == bits(Null) {
		return "<untyped null>"
	}
	panic(fmt.Errorf("unknown object pattern 0x%08X", field))
}

func IsPtrWithTag(obj Object, tag uint8) bool {
	if !obj.Is(Ptr32) {
		return false
	}

	bits := bits(obj)
	return tag == bits.GetU8()
}

func PackF64(v float64) Object {
	if math.IsNaN(v) {
		panic("cannot pack NAN")
	}

	return Object(math.Float64bits(v))
}

func UnpackF64(p Object) float64 {
	f64 := math.Float64frombits(uint64(p))
	if math.IsNaN(f64) {
		panic(fmt.Errorf("%x: not a float", p))
	}
	return f64
}

func PackPtr32(tag uint8, addr uint32) Object {
	bits := Ptr32
	bits.PutU8(tag)
	bits.PutU32(addr)
	return Object(bits)
}

func UnpackPtr32(p Object) (uint8, uint32) {
	bits := mustMatch(Ptr32, p)
	return bits.GetU8(), bits.GetU32()
}

func PackStr32(len uint8, offset uint32) Object {
	bits := Str32
	bits.PutU8(len)
	bits.PutU32(offset)
	return Object(bits)
}

func UnpackStr32(p Object) (uint8, uint32) {
	bits := mustMatch(Str32, p)
	return bits.GetU8(), bits.GetU32()
}

func PackArray(array [6]byte) Object {
	var bytes bytes
	copy(bytes[0:6], array[:])
	return Object(Array | bytes.Bits())
}

func UnpackArray(p Object) [6]byte {
	bits := mustMatch(Array, p)
	bytes := bits.Bytes()
	array := [6]byte{}
	copy(array[:], bytes[0:6])
	return array
}

func PackBool(b bool) Object {
	if b {
		return PackedTrue
	}
	return PackedFalse
}

func UnpackBool(p Object) bool {
	mustMatch(Bool, p)
	return p == PackedTrue
}

func PackI32(i int32) Object {
	bits := I32
	bits.PutU32(uint32(i))
	return Object(bits)
}

func UnpackI32(p Object) int32 {
	bits := mustMatch(I32, p)
	return int32(bits.GetU32())
}

func PackU32(u uint32) Object {
	bits := U32
	bits.PutU32(u)
	return Object(bits)
}

func UnpackU32(p Object) uint32 {
	bits := mustMatch(U32, p)
	return bits.GetU32()
}

func mustMatch(mask bits, p Object) bits {
	bits := bits(p)
	if mask&bits != mask {
		panic("did not match mask")
	}
	return bits
}

type bits uint64
type bytes [8]uint8

var _ encoder = (*bits)(nil)
var _ encoder = bytes{}

type encoder interface {
	PutU8(uint8)
	GetU8() uint8
	PutU32(uint32)
	GetU32() uint32
}

func (this bytes) PutU8(u8 uint8) {
	this[0] = u8
}

func (this bytes) GetU8() uint8 {
	return this[0]
}

func (this bytes) PutU32(u32 uint32) {
	this[2] = byte(u32)
	this[3] = byte(u32 >> 8)
	this[4] = byte(u32 >> 16)
	this[5] = byte(u32 >> 24)
}

func (this bytes) GetU32() (u32 uint32) {
	u32 = uint32(this[2])
	u32 |= uint32(this[3]) << 8
	u32 |= uint32(this[4]) << 16
	u32 |= uint32(this[5]) << 24
	return
}

func (this *bits) PutU8(u8 uint8) {
	*this |= bits(u8 << 1)
}

func (this bits) GetU8() uint8 {
	return uint8(this >> 1)
}

func (this *bits) PutU32(u32 uint32) {
	*this |= bits(u32) << 9
}

func (this bits) GetU32() uint32 {
	return uint32(this >> 9)
}

func (this bits) Bytes() (bytes bytes) {
	bytes[0] = uint8(this)
	bytes[1] = uint8(this >> 8)
	bytes[2] = uint8(this >> 16)
	bytes[3] = uint8(this >> 24)
	bytes[4] = uint8(this >> 32)
	bytes[5] = uint8(this >> 40)
	bytes[6] = uint8(this >> 48)
	bytes[7] = uint8(this >> 56)
	return
}

func (this bytes) Bits() (bits_ bits) {
	bits_ = bits(this[0])
	bits_ = bits(this[1]) << 8
	bits_ = bits(this[2]) << 16
	bits_ = bits(this[3]) << 24
	bits_ = bits(this[4]) << 32
	bits_ = bits(this[5]) << 40
	bits_ = bits(this[6]) << 48
	bits_ = bits(this[7]) << 56
	return
}

const (
	LOW48    bits = (1 << 48) - 1
	MASKBITS bits = (1 << 63) | (1 << 49) | (1 << 48)

	ptrmask  bits = 1<<63 | 1<<49 | 0<<48
	strmask  bits = 1<<63 | 0<<49 | 1<<48
	array    bits = 1<<63 | 0<<49 | 0<<48
	i32mask  bits = 0<<63 | 1<<49 | 1<<48
	u32mask  bits = 0<<63 | 1<<49 | 0<<48
	boolmask bits = 0<<63 | 0<<49 | 1<<48
	QNAN     bits = 0x7FF8000000000001

	Ptr32 bits = QNAN | ptrmask
	Str32 bits = QNAN | strmask
	Array bits = QNAN | array
	I32   bits = QNAN | i32mask
	U32   bits = QNAN | u32mask
	Bool  bits = QNAN | boolmask
	F64   bits = ^QNAN

	PackedTrue  = Object(QNAN | boolmask | bits(1<<1))
	PackedFalse = Object(QNAN | boolmask | bits(1<<2))
	Null        = Object(0)
)
