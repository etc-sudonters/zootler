package nan

import (
	"fmt"
	"math"
)

const (
	qnan   uint64 = 0x00000007ff8000000000001
	zconst        = qnan | 0x0000670000000000
	zstr          = qnan | 0x0000680000000000
	zptr          = qnan | 0x0000690000000000
	unum          = qnan | 0x00006a0000000000
	snum          = qnan | 0x00006b0000000000
	zbool         = qnan | 0x00006c0000000000
	ztrue         = qnan | 0x00006c1000000000
	zfalse        = qnan | 0x00006c2000000000
)

type PackedValue float64
type Pointer uint32
type ConstantPointer uint32
type StringPointer uint32

func (pv PackedValue) String() string {
	if pv == pv {
		return fmt.Sprintf("%f", float64(pv))
	}

	u := unpack(pv)

	switch {
	case u&zconst == zconst:
		return fmt.Sprintf("CONST: %d", (^zptr&u)>>1)
	case u&zstr == zstr:
		return fmt.Sprintf("STR: %d", (^zptr&u)>>1)
	case u&zptr == zptr:
		return fmt.Sprintf("PTR: %d", (^zptr&u)>>1)
	case u&unum == unum:
		return fmt.Sprintf("0x%04X", (^unum&u)>>1)
	case u&snum == snum:
		return fmt.Sprintf("%d", int32((^snum&u)>>1))
	case u&ztrue == ztrue:
		return "true"
	case u&zfalse == zfalse:
		return "false"
	default:
		panic("not implemented")
	}
}

func unpack(pv PackedValue) uint64 {
	return math.Float64bits(float64(pv))
}

func unpackWithMask(pv PackedValue, mask uint64) (uint64, bool) {
	u := unpack(pv)
	return u, u&mask == mask
}

func (pv PackedValue) Equals(p PackedValue) bool {
	if pv == p {
		return true
	}

	return unpack(pv) == unpack(p)
}

func (pv PackedValue) Int() (int, bool) {
	u, isInt := unpackWithMask(pv, snum)
	if !isInt {
		return 0, false
	}

	val := (^snum & u) >> 1
	return int(val), true
}

func (pv PackedValue) Uint() (uint32, bool) {
	u, isInt := unpackWithMask(pv, unum)
	if !isInt {
		return 0, false
	}

	val := (^snum & u) >> 1
	return uint32(val), true
}

func (pv PackedValue) Pointer() (Pointer, bool) {
	return unpackPtr[Pointer](pv, zptr)
}

func (pv PackedValue) ConstantPointer() (ConstantPointer, bool) {
	return unpackPtr[ConstantPointer](pv, zconst)
}

func (pv PackedValue) StringPointer() (StringPointer, bool) {
	return unpackPtr[StringPointer](pv, zconst)
}

func (pv PackedValue) Bool() (bool, bool) {
	u, isBool := unpackWithMask(pv, zbool)
	return u&ztrue == ztrue, isBool
}

func unpackPtr[P ~uint32](pv PackedValue, mask uint64) (P, bool) {
	u, isPtr := unpackWithMask(pv, mask)
	if !isPtr {
		return 0, false
	}

	return P((^mask & u) >> 1), true
}

type PackableUint interface {
	~uint32 | ~uint16 | ~uint8
}

type PackableInt interface {
	~int32 | ~int16 | ~int8
}

func PackUint[P PackableUint](p P) PackedValue {
	return pack(unum | uint64(p)<<1)
}

func PackInt[P PackableInt](p P) PackedValue {
	return PackUint(uint32(int32(p)))
}

func PackFloat64(p float64) PackedValue {
	return PackedValue(p)
}

func PackBool(b bool) PackedValue {
	if b {
		return pack(ztrue)
	}
	return pack(zfalse)
}

func PackPtr(ptr uint32) PackedValue {
	return pack(zptr | uint64(ptr)<<1)
}

func PackStr(strPtr uint32) PackedValue {
	return pack(zstr | uint64(strPtr)<<1)
}

func PackConst(conPtr uint32) PackedValue {
	return pack(zconst | uint64(conPtr)<<1)
}

func pack(u uint64) PackedValue {
	return PackedValue(math.Float64frombits(u))
}
