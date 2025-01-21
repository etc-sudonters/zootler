package objects

import (
	"fmt"
	"math"
)

type Object bits

func (this Object) Is(mask bits) bool {
	return bits(this)&mask == mask
}

func IsPtrWithTag(obj Object, tag uint16) bool {
	if !obj.Is(Ptr32) {
		return false
	}

	bits := bits(obj)
	return tag == bits.GetU16()
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

func PackPtr32(tag uint16, addr uint32) Object {
	bits := Ptr32
	bits.PutU16(tag)
	bits.PutU32(addr)
	return Object(bits)
}

func UnpackPtr32(p Object) (uint16, uint32) {
	bits := mustMatch(Ptr32, p)
	return bits.GetU16(), bits.GetU32()
}

func PackStr32(len uint16, offset uint32) Object {
	bits := Str32
	bits.PutU16(len)
	bits.PutU32(offset)
	return Object(bits)
}

func UnpackStr32(p Object) (uint16, uint32) {
	bits := mustMatch(Str32, p)
	return bits.GetU16(), bits.GetU32()
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
	PutU16(uint16)
	GetU16() uint16
	PutU32(uint32)
	GetU32() uint32
}

func (this bytes) PutU16(u16 uint16) {
	this[0] = byte(u16)
	this[1] = byte(u16 >> 8)
}

func (this bytes) GetU16() uint16 {
	return uint16(this[0]) | uint16(this[1])<<8
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

func (this *bits) PutU16(u16 uint16) {
	*this |= bits(u16)
}

func (this bits) GetU16() uint16 {
	return uint16(this)
}

func (this *bits) PutU32(u32 uint32) {
	*this |= bits(u32) << 16
}

func (this bits) GetU32() uint32 {
	return uint32(this >> 16)
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
	low48     bits = (1 << 48) - 1
	maskbits  bits = (1 << 63) | (1 << 49) | (1 << 48)
	container bits = (^low48) & (^maskbits)

	ptrmask  bits = 1<<63 | 1<<49 | 0<<48
	strmask  bits = 1<<63 | 0<<49 | 1<<48
	array    bits = 1<<63 | 0<<49 | 0<<48
	i32mask  bits = 0<<63 | 1<<49 | 1<<48
	u32mask  bits = 0<<63 | 1<<49 | 0<<48
	boolmask bits = 0<<63 | 0<<49 | 1<<48
	qnan     bits = 0x7FF8000000000001

	Ptr32 bits = container | ptrmask
	Str32 bits = container | strmask
	Array bits = container | array
	I32   bits = container | i32mask
	U32   bits = container | u32mask
	Bool  bits = container | boolmask
	F64   bits = ^qnan

	PackedTrue  = Object(container | boolmask | bits(1<<1))
	PackedFalse = Object(container | boolmask | bits(1<<2))
	Null        = Object(0)
)
