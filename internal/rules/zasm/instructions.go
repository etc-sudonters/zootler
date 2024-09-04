package zasm

import (
	"fmt"
	"slices"
	"strings"
	"sudonters/zootler/internal/intern"

	"github.com/etc-sudonters/substrate/slipup"
)

// Three Addr style quad more wasted space per operation but avoids variable
// length operations
type Instruction uint32

func (i Instruction) Bytes() [4]uint8 {
	var bytes [4]uint8
	bytes[0] = uint8((i & INSTR_OPER) >> OPER_SHIFT)
	bytes[1] = uint8((i & INSTR_ARG1) >> ARG1_SHIFT)
	bytes[2] = uint8((i & INSTR_ARG2) >> ARG2_SHIFT)
	bytes[3] = uint8((i & INSTR_ARG3) >> ARG3_SHIFT)
	return bytes
}

func (i Instruction) String() string {
	return fmt.Sprintf("0x%X", uint32(i))
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
	return Instruction(o) << OPER_SHIFT
}

func Encode(o Op, data [3]uint8) Instruction {
	var i Instruction
	i = i | Instruction(o)<<OPER_SHIFT
	i = i | Instruction(data[0])<<ARG1_SHIFT
	i = i | Instruction(data[1])<<ARG2_SHIFT
	i = i | Instruction(data[2])<<ARG3_SHIFT
	return i
}

func DecodeU24(data [3]uint8) uint32 {
	var u uint32
	u |= uint32(data[0]) << ARG1_SHIFT
	u |= uint32(data[1]) << ARG2_SHIFT
	u |= uint32(data[2]) << ARG3_SHIFT
	return u
}

func Encode24Bit(o Op, raw uint32) Instruction {
	data := Instruction(raw)
	if !IsU24(data) {
		panic("24bit overflow")
	}

	return (Instruction(o) << OPER_SHIFT) | data
}

type Instructions []Instruction

func (is Instructions) String() string {
	var sb strings.Builder

	sb.WriteRune('{')
	for idx, inst := range is {
		if idx != 0 {
			sb.WriteRune(' ')
		}
		fmt.Fprintf(&sb, "%s", inst)
	}
	sb.WriteRune('}')

	return sb.String()
}

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

type InstructionWriter struct {
	i Instructions
}

func Tape() *InstructionWriter {
	var iw InstructionWriter
	iw.i = make(Instructions, 0, 16)
	return &iw
}

func (iw *InstructionWriter) Instructions() Instructions {
	return iw.i
}

func (iw *InstructionWriter) Len() int {
	return len(iw.i)
}

func (iw *InstructionWriter) Concat(o ...Instructions) *InstructionWriter {
	all := make([]Instructions, len(o)+1)
	all[0] = iw.i
	copy(all[1:len(o)], o)
	iw.i = slices.Concat(all...)
	return iw
}

func (iw *InstructionWriter) Write(i Instruction) *InstructionWriter {
	iw.i = append(iw.i, i)
	return iw
}

func (iw *InstructionWriter) WriteOp(o Op) *InstructionWriter {
	return iw.Write(EncodeOp(o))
}

func (iw *InstructionWriter) WriteFull(o Op, data [3]uint8) *InstructionWriter {
	return iw.Write(Encode(o, data))
}

func (iw *InstructionWriter) WriteU24(o Op, u uint32) *InstructionWriter {
	return iw.Write(Encode24Bit(o, u))
}

func (iw *InstructionWriter) WriteLoadConst(constId uint32) *InstructionWriter {
	return iw.WriteU24(OP_LOAD_CONST, constId)
}

func (iw *InstructionWriter) WriteLoadIdent(ident uint32) *InstructionWriter {
	return iw.WriteU24(OP_LOAD_IDENT, ident)
}

func (iw *InstructionWriter) WriteLoadStr(str intern.Str) *InstructionWriter {
	return iw.WriteFull(OP_LOAD_STR, str.Bytes())
}

func (iw *InstructionWriter) WriteLoadU24(u uint32) *InstructionWriter {
	return iw.WriteU24(OP_LOAD_U24, u)
}

func (iw *InstructionWriter) WriteLoadBool(b bool) *InstructionWriter {
	var v uint8 = 2
	if b {
		v = 1
	}
	return iw.WriteFull(OP_LOAD_BOOL, [3]uint8{v, 0, 0})
}

func (iw *InstructionWriter) WriteCall(arity int) (*InstructionWriter, error) {
	switch arity {
	case 0:
		iw.WriteOp(OP_CALL_0)
		break
	case 1:
		iw.WriteOp(OP_CALL_1)
		break
	case 2:
		iw.WriteOp(OP_CALL_2)
		break
	default:
		return iw, slipup.Createf("maximum 2 arguments supported, got %d", arity)
	}
	return iw, nil
}

type Op uint8

const (
	OPER_SHIFT = 24
	ARG1_SHIFT = 16
	ARG2_SHIFT = 8
	ARG3_SHIFT = 0

	// masks
	U24_MASK        uint32      = 0x00FFFFFF
	INSTR_OPER      Instruction = 0xFF000000
	INSTR_ARG1                  = 0x00FF0000
	INSTR_ARG2                  = 0x0000FF00
	INSTR_ARG3                  = 0x000000FF
	INSTR_MASK_OPER             = ^INSTR_OPER

	// common things
	INSTR_ANY_LOAD        Instruction = 0x20000000
	INSTR_LOAD_BOOL                   = 0x25000000
	INSTR_LOAD_BOOL_FALSE             = (OP_LOAD_BOOL << OPER_SHIFT) | 2<<ARG1_SHIFT
	INSTR_LOAD_BOOL_TRUE              = (OP_LOAD_BOOL << OPER_SHIFT) | 1<<ARG1_SHIFT

	// operations
	OP_LOAD_CONST Op = 0x21 // next 24 bits is const array index
	OP_LOAD_IDENT    = 0x22 // next 24 bits is the name array index
	OP_LOAD_STR      = 0x23 // next 24 bits that form a str pointer
	OP_LOAD_U24      = 0x24 // push next 24 bits as nanpacked int to stack
	OP_LOAD_BOOL     = 0x25 // arg 1 is 1 for true and 2 for false

	OP_CMP_EQ Op = 0x31
	OP_CMP_NQ    = 0x32
	OP_CMP_LT    = 0x33

	OP_BOOL_AND    Op = 0x41
	OP_BOOL_OR        = 0x42
	OP_BOOL_NEGATE    = 0x43

	OP_CALL_0    Op = 0x51
	OP_CALL_1       = 0x52
	OP_CALL_2       = 0x53
	OP_CHK_SET_1    = 0x54
	OP_CHK_SET_2    = 0x55
	OP_CHK_TRK      = 0x56
	OP_CHK_QTY      = 0x57

	OP_JMP_TRUE   Op = 0x61
	OP_JMP_FALSE     = 0x62
	OP_JMP_UNCOND    = 0x63
)
