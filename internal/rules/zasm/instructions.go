package zasm

import (
	"slices"
	"sudonters/zootler/internal/intern"

	"github.com/etc-sudonters/substrate/slipup"
)

// Three Addr style quad more wasted space per operation but avoids variable
// length operations
type Instruction uint32

func (i Instruction) Bytes() [4]uint8 {
	var bytes [4]uint8
	bytes[0] = uint8((i & INSTR_OPER) >> 24)
	bytes[1] = uint8((i & INSTR_ARG1) >> 16)
	bytes[2] = uint8((i & INSTR_ARG2) >> 8)
	bytes[3] = uint8(i & INSTR_ARG3)
	return bytes
}

func IsU8[U ~uint32](u U) bool {
	return U(uint8(u)) == u
}

func IsU16[U ~uint32](u U) bool {
	return U(uint16(u)) == u
}

func IsU24[U ~uint32](u U) bool {
	return u&U(U24_MASK) == u
}

func EncodeOp(o Op) Instruction {
	return Instruction(o) << 24
}

func Encode(o Op, data [3]uint8) Instruction {
	var i Instruction
	i = i | Instruction(0)<<24
	i = i | Instruction(data[0])<<16
	i = i | Instruction(data[1])<<8
	i = i | Instruction(data[2])
	return i
}

func Encode24Bit(o Op, raw uint32) Instruction {
	data := Instruction(raw)
	if !IsU24(data) {
		panic("24bit overflow")
	}

	return (Instruction(o) << 24) | data
}

type Instructions []Instruction

func (is Instructions) MatchOne(i Instruction) bool {
	return len(is) == 1 && i&is[0] == i
}

func (is Instructions) Match(mask Instructions) bool {
	if len(is) != len(mask) {
		return false
	}

	for idx, instr := range is {
		match := mask[idx]
		if match != match&instr {
			return false
		}
	}

	return true
}

func (is Instructions) Concat(o ...Instructions) Instructions {
	all := make([]Instructions, len(o)+1)
	all[0] = is
	copy(all[1:len(o)], o)
	return slices.Concat(all...)
}

func (is Instructions) Write(i Instruction) Instructions {
	return append(is, i)
}

func (is Instructions) WriteOp(o Op) Instructions {
	return is.Write(Instruction(o) << 24)
}

func (is Instructions) WriteFull(o Op, data [3]uint8) Instructions {
	return is.Write(Encode(o, data))
}

func (is Instructions) WriteU24(o Op, u uint32) Instructions {
	return is.Write(Encode24Bit(o, u))
}

func (is Instructions) WriteLoadConst(constId uint32) Instructions {
	return is.WriteU24(OP_LOAD_CONST, constId)
}

func (is Instructions) WriteLoadIdent(ident uint32) Instructions {
	return is.WriteU24(OP_LOAD_IDENT, ident)
}

func (is Instructions) WriteLoadStr(str intern.Str) Instructions {
	return is.WriteFull(OP_LOAD_STR, str.Bytes())
}

func (is Instructions) WriteLoadU24(u uint32) Instructions {
	return is.WriteU24(OP_LOAD_U24, u)
}

func (is Instructions) WriteLoadBool(b bool) Instructions {
	if b {
		return is.WriteU24(OP_LOAD_BOOL, 1)
	}
	return is.WriteU24(OP_LOAD_BOOL, 0)
}

func (is Instructions) WriteCall(arity int) (Instructions, error) {
	switch arity {
	case 0:
		return is.WriteOp(OP_CALL_0), nil
	case 1:
		return is.WriteOp(OP_CALL_1), nil
	case 2:
		return is.WriteOp(OP_CALL_2), nil
	default:
		return nil, slipup.Createf("maximum 2 arguments supported, got %d", arity)
	}
}

type Op uint8

const (
	U24_MASK uint32 = 0x00FFFFFF

	// masks
	INSTR_OPER      Instruction = 0xFF000000
	INSTR_ARG1                  = 0x00FF0000
	INSTR_ARG2                  = 0x0000FF00
	INSTR_ARG3                  = 0x000000FF
	INSTR_MASK_OPER             = ^INSTR_OPER

	// common things
	INSTR_ANY_LOAD        Instruction = 0x20000000
	INSTR_LOAD_BOOL                   = 0x25000000
	INSTR_LOAD_BOOL_FALSE             = 0x25000002
	INSTR_LOAD_BOOL_TRUE              = 0x25000001

	// operations
	OP_LOAD_CONST Op = 0x21 // next 24 bits is const array index
	OP_LOAD_IDENT    = 0x22 // next 24 bits is the name array index
	OP_LOAD_STR      = 0x23 // next 24 bits that form a str pointer
	OP_LOAD_U24      = 0x24 // push next 24 bits as nanpacked int to stack
	OP_LOAD_BOOL     = 0x25 // next 24 bit int is 1 for true, 2 for false

	OP_CMP_EQ Op = 0x31
	OP_CMP_NQ    = 0x32
	OP_CMP_LT    = 0x33

	OP_BOOL_AND    Op = 0x41
	OP_BOOL_OR        = 0x42
	OP_BOOL_NEGATE    = 0x43

	OP_CALL_0 Op = 0x51
	OP_CALL_1    = 0x52
	OP_CALL_2    = 0x53
)
