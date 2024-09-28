package nan

import (
	"fmt"
	"math"
)

const (
	qnan     uint64 = 0x00 | 0x7ff8000000000001
	zbool           = qnan | 0x0000610000000001
	ztrue           = qnan | 0x0000611000000001
	zfalse          = qnan | 0x0000612000000001
	zunum           = qnan | 0x0000620000000001
	zsnum           = qnan | 0x0000630000000001
	zstr            = qnan | 0x0000640000000001
	zfunc           = qnan | 0x0000650000000001
	ztoken          = qnan | 0x0000660000000001
	ztrick          = qnan | 0x0000670000000001
	zsetting        = qnan | 0x0000680000000001
	zvar            = qnan | 0x0000690000000001
)

type PackedType uint8
type PtrType uint8

const (
	_ PackedType = iota
	PT_F64
	PT_BOOL
	PT_UNUM
	PT_SNUM
	PT_STRING
	PT_FUNC
	PT_TOKEN
	PT_TRICK
	PT_SETTING
	PT_VAR

	_ PtrType = iota
	PTR_FUNC
	PTR_SETTING
	PTR_STR
	PTR_TOKEN
	PTR_TRICK
	PTR_VAR
)

type Packed float64

func (pv Packed) Type() PackedType {
	if pv == pv {
		return PT_F64
	}

	u := pv.Bits()
	if u&zbool == zbool {
		return PT_BOOL
	} else if u&zunum == zunum {
		return PT_UNUM
	} else if u&zsnum == zsnum {
		return PT_SNUM
	} else if u&zstr == zstr {
		return PT_STRING
	} else if u&zfunc == zfunc {
		return PT_FUNC
	} else if u&ztoken == ztoken {
		return PT_TOKEN
	} else if u&ztrick == ztrick {
		return PT_TRICK
	} else if u&zsetting == zsetting {
		return PT_SETTING
	} else if u&zvar == zvar {
		return PT_VAR
	} else {
		panic("unknown packed type")
	}
}

func (pv Packed) Bits() uint64 {
	return unpack(pv)
}

func (pv Packed) String() string {
	if pv == pv {
		return fmt.Sprintf("%f", float64(pv))
	}

	u := unpack(pv)
	if u&ztoken == ztoken {
		return fmt.Sprintf("TOK: %d", (^ztoken&u)>>1)
	} else if u&zunum == zunum {
		return fmt.Sprintf("0x%04X", (^zunum&u)>>1)
	} else if u&zsnum == zsnum {
		return fmt.Sprintf("%d", int32((^zsnum&u)>>1))
	} else if u&ztrue == ztrue {
		return "true"
	} else if u&zfalse == zfalse {
		return "false"
	} else {
		return "???wtf"
	}
}

func (pv Packed) Equals(p Packed) bool {
	if pv == p {
		return true
	}

	return unpack(pv) == unpack(p)
}

func (pv Packed) Int() (int, bool) {
	u, isInt := unpackWithMask(pv, zsnum)
	if !isInt {
		return 0, false
	}

	val := (^zsnum & u) >> 1
	return int(val), true
}

func (pv Packed) Uint() (uint32, bool) {
	u, isInt := unpackWithMask(pv, zunum)
	if !isInt {
		return 0, false
	}

	val := (^zsnum & u) >> 1
	return uint32(val), true
}

func (pv Packed) Token() (uint32, bool) {
	u, isTok := unpackWithMask(pv, ztoken)
	if !isTok {
		return 0, false
	}

	return uint32((^ztoken & u) >> 1), true
}

func (pv Packed) Bool() (bool, bool) {
	u, isBool := unpackWithMask(pv, zbool)
	return u&ztrue == ztrue, isBool
}

type PackableUint interface {
	~uint16 | ~uint8
}

type PackableInt interface {
	~int16 | ~int8
}

func PackUint[P PackableUint](p P) Packed {
	return pack(zunum | uint64(p)<<1)
}

func PackInt[P PackableInt](p P) Packed {
	return PackUint(uint16(int16(p)))
}

func PackFloat64(p float64) Packed {
	return Packed(p)
}

func PackBool(b bool) Packed {
	if b {
		return pack(ztrue)
	}
	return pack(zfalse)
}

func PackToken(ptr uint32) Packed {
	return pack(ztoken | uint64(ptr)<<1)
}

func PackPtr(typ PtrType, ptr uint16) Packed {
	ptr64 := uint64(ptr) << 1

	switch typ {
	case PTR_FUNC:
		return pack(zfunc | ptr64)
	case PTR_SETTING:
		return pack(zsetting | ptr64)
	case PTR_STR:
		return pack(zstr | ptr64)
	case PTR_TOKEN:
		return pack(ztoken | ptr64)
	case PTR_TRICK:
		return pack(ztrick | ptr64)
	case PTR_VAR:
		return pack(zvar | ptr64)
	default:
		panic("unknown ptr type")
	}
}

func pack(u uint64) Packed {
	return Packed(math.Float64frombits(u))
}

func unpack(pv Packed) uint64 {
	return math.Float64bits(float64(pv))
}

func unpackWithMask(pv Packed, mask uint64) (uint64, bool) {
	u := unpack(pv)
	return u, u&mask == mask
}
