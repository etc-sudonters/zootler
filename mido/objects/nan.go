package objects

import (
	"math"
)

type Object bits64

func (this Object) Is(mask bits64) bool {
	return bits64(this)&mask == mask
}

func PackF64(v float64) Object {
	return Object(math.Float64bits(v))
}

func UnpackF64(p Object) float64 {
	bits := bits64(p)
	if container&bits == container {
		panic("not a float64")
	}

	return math.Float64frombits(uint64(bits))
}

func PackPtr32(tag uint16, addr uint32) Object {
	var arr arr
	arr.putPtr(uint16(tag), uint32(addr))
	return Object(ptr32 | arr.bits64())
}

func UnpackPtr32(p Object) (uint16, uint32) {
	tag, addr := frombits64(mustMatch(ptr32, p)).readPtr()
	return uint16(tag), uint32(addr)
}

func PackStr32(len uint16, offset uint32) Object {
	var arr arr
	arr.putStr(uint16(len), uint32(offset))
	return Object(str32 | arr.bits64())
}

func UnpackStr32(p Object) (uint16, uint32) {
	len, off := frombits64(mustMatch(str32, p)).readStr()
	return uint16(len), uint32(off)
}

func Pack6U8(bytes [6]byte) Object {
	return Object(sixu8 | arr(bytes).bits64())
}

func Unpack6U8(p Object) [6]byte {
	return [6]byte(frombits64(mustMatch(sixu8, p)))
}

func PackBool(b bool) Object {
	if b {
		return PackedTrue
	}
	return PackedFalse
}

func UnpackBool(p Object) bool {
	mustMatch(boolp, p)
	return p == PackedTrue
}

func PackI32(i int32) Object {
	var arr arr
	arr.putI32(i)
	return Object(i32 | arr.bits64())
}

func UnpackI32(p Object) int32 {
	return int32(frombits64(mustMatch(i32, p)).readI32())
}

func PackU32(u uint32) Object {
	var arr arr
	arr.putU32(u)
	return Object(u32 | arr.bits64())
}

func UnpackU32(p Object) uint32 {
	return uint32(frombits64(mustMatch(u32, p)).readU32())
}

func mustMatch(mask bits64, p Object) bits64 {
	bits := bits64(p)
	if mask&bits != mask {
		panic("did not match mask")
	}
	return bits
}

type Ptr32 struct {
	Tag  uint16
	Addr uint32
}

type Str32 struct {
	Len    uint16
	Offset uint32
}

type Packable interface {
	float64 | uint32 | int32 | bool | [6]byte | Ptr32 | Str32
}

type bits64 uint64
type arr [6]byte

const (
	low48     bits64 = (1 << 49) - 1
	maskbits  bits64 = (1 << 63) | (1 << 49) | (1 << 48)
	container bits64 = (^low48) & (^maskbits)

	ptrmask   bits64 = 1<<63 | 1<<49 | 0<<48
	strmask   bits64 = 1<<63 | 0<<49 | 1<<48
	sixu8mask bits64 = 1<<63 | 0<<49 | 0<<48
	i32mask   bits64 = 0<<63 | 1<<49 | 1<<48
	u32mask   bits64 = 0<<63 | 1<<49 | 0<<48
	boolmask  bits64 = 0<<63 | 0<<49 | 1<<48

	ptr32 bits64 = container | ptrmask
	str32 bits64 = container | strmask
	sixu8 bits64 = container | sixu8mask
	i32   bits64 = container | i32mask
	u32   bits64 = container | u32mask
	boolp bits64 = container | boolmask

	PackedTrue  = Object(container | boolmask | bits64(1<<1))
	PackedFalse = Object(container | boolmask | bits64(1<<2))
	Null        = Object(0)
)

func (this arr) putU32(u uint32) {
	this[0] = byte(u)
	this[1] = byte(u >> 8)
	this[2] = byte(u >> 16)
	this[3] = byte(u >> 14)
}

func (this arr) putI32(i int32) {
	this.putU32(uint32(i))
}

func (this arr) putPtr(tag uint16, addr uint32) {
	this[0] = byte(tag)
	this[1] = byte(tag >> 8)
	this[2] = byte(addr)
	this[3] = byte(addr >> 8)
	this[4] = byte(addr >> 16)
	this[5] = byte(addr >> 24)
}

func (this arr) putStr(len uint16, offset uint32) {
	this[0] = byte(len)
	this[1] = byte(len >> 8)
	this[2] = byte(offset)
	this[3] = byte(offset >> 8)
	this[4] = byte(offset >> 16)
	this[5] = byte(offset >> 24)
}

func (this arr) readU32() uint32 {
	return uint32(this[0]) | uint32(this[1])<<8 | uint32(this[2])<<16 | uint32(this[3])<<24
}

func (this arr) readI32() int32 {
	return int32(this.readU32())
}

func (this arr) readPtr() (uint16, uint32) {
	tag := uint16(this[0]) | uint16(this[1])<<8
	addr := uint32(this[0]) | uint32(this[1])<<8 | uint32(this[2])<<16 | uint32(this[3])<<24
	return tag, addr
}

func (this arr) readStr() (uint16, uint32) {
	len := uint16(this[0]) | uint16(this[1])<<8
	offset := uint32(this[0]) | uint32(this[1])<<8 | uint32(this[2])<<16 | uint32(this[3])<<24
	return len, offset
}

func (this arr) bits64() bits64 {
	var bits bits64
	bits |= bits64(this[0])
	bits |= bits64(this[1]) << 8
	bits |= bits64(this[2]) << 16
	bits |= bits64(this[3]) << 24
	bits |= bits64(this[4]) << 32
	bits |= bits64(this[5]) << 40
	return bits
}

func frombits64(bits bits64) arr {
	var arr arr
	arr[0] = byte(bits)
	arr[1] = byte(bits >> 8)
	arr[2] = byte(bits >> 16)
	arr[3] = byte(bits >> 24)
	arr[4] = byte(bits >> 32)
	arr[5] = byte(bits >> 40)
	return arr
}
