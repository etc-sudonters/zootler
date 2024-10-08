package zasm

import (
	"slices"
	"sudonters/zootler/icearrow/nan"
	"sudonters/zootler/internal/intern"

	"github.com/etc-sudonters/substrate/slipup"
)

type Instruction uint32

func (i Instruction) Bytes() [4]uint8 {
	var b [4]uint8
	b[OPER_IDX] = uint8((i & ONLY_OPER) >> OPER_SHIFT)
	b[ARG1_IDX] = uint8((i & ONLY_ARG1) >> ARG1_SHIFT)
	b[ARG2_IDX] = uint8((i & ONLY_ARG2) >> ARG2_SHIFT)
	b[ARG3_IDX] = uint8((i & ONLY_ARG3) >> ARG3_SHIFT)
	return b
}

func (i Instruction) Op() Op {
	return Op((i & ONLY_OPER) >> OPER_SHIFT)
}

func (i Instruction) OpAndPayload() (Op, [3]uint8) {
	bytes := i.Bytes()
	return Op(bytes[0]), [3]uint8(bytes[1:4])
}

func (i Instruction) MatchOp(op Op) bool {
	return Op((i&ONLY_OPER)>>OPER_SHIFT)&op == op
}

type Op uint8

const (
	OP_LOAD_CONST Op = 0x22 // 24bit index to const table, push into stack
	OP_LOAD_STR      = 0x23 // 24bit index to str table
	OP_LOAD_BOOL     = 0x24 // arg1 is 1 for true, 2 for false
	OP_LOAD_TOK      = 0x25 // 24bit index to names table
	OP_LOAD_VAR      = 0x26 // 24bit index to names table
	OP_LOAD_SYM      = 0x27 // 24bit index to names table
	OP_LOAD_TRK      = 0x28
	OP_LOAD_SET      = 0x29
	OP_BOOL_AND      = 0x41 // pop 2 from stack, push AND result to stack
	OP_BOOL_OR       = 0x42 // pop 2 from stack, push OR result to stack
	OP_CALL_0        = 0x51 // 24bit index to name table to use as func name
	OP_CALL_1        = 0x52 // 24bit index... pop 1 from stack as arg
	OP_CALL_2        = 0x53 // 24bit index... pop 2 from stack as args

	OPER_SHIFT uint = 24
	ARG1_SHIFT      = 16
	ARG2_SHIFT      = 8
	ARG3_SHIFT      = 0

	OPER_IDX int = 0
	ARG1_IDX     = 1
	ARG2_IDX     = 2
	ARG3_IDX     = 3

	ONLY_OPER Instruction = 0xFF000000
	ONLY_ARGS             = ^ONLY_OPER
	ONLY_ARG1             = 0x00FF0000
	ONLY_ARG2             = 0x0000FF00
	ONLY_ARG3             = 0x000000FF

	U24_MASK = 0x00FFFFFF
)

func EncodeOp(o Op) Instruction {
	return Instruction(o) << OPER_SHIFT
}

func Encode(o Op, data [3]uint8) Instruction {
	var i Instruction
	i |= Instruction(o) << OPER_SHIFT
	i |= Instruction(data[0]) << ARG1_SHIFT
	i |= Instruction(data[1]) << ARG2_SHIFT
	i |= Instruction(data[2]) << ARG3_SHIFT
	return i
}

func EncodeU16(u uint16) [3]uint8 {
	var payload [3]uint8
	payload[0] = uint8(0x00FF & u)
	payload[1] = uint8((0xFF00 & u) >> 8)
	payload[2] = 0
	return payload
}

func DecodeBool(payload [3]uint8) bool {
	return payload[0] == 1
}

func DecodeU16(payload [3]uint8) uint16 {
	return (uint16(payload[1]) << 8) | uint16(payload[0])
}

func DecodeU24(payload [3]uint8) uint32 {
	var u uint32
	u |= uint32(payload[0]) << ARG1_SHIFT
	u |= uint32(payload[1]) << ARG2_SHIFT
	u |= uint32(payload[2]) << ARG3_SHIFT
	return u
}

func EncodeOpAndU24(o Op, u uint32) Instruction {
	return EncodeOp(o) | (Instruction(u) & ONLY_ARGS)
}

func AssertU24[U ~uint32](u U) uint32 {
	if U(U24_MASK)&u != u {
		panic(slipup.Createf("u24 overflow: %d", u))
	}
	return uint32(u)
}

type Instructions []Instruction

func (i Instructions) Eq(o Instructions) bool {
	if len(i) != len(o) {
		return false
	}

	for idx := range len(i) {
		if i[idx] != o[idx] {
			return false
		}
	}

	return true
}

type InstructionWriter struct {
	I Instructions
}

func IW() *InstructionWriter {
	return &InstructionWriter{}
}

func IntoIW(i Instructions) *InstructionWriter {
	iw := IW()
	iw.I = i
	return iw
}

func (iw *InstructionWriter) Write(i Instruction) *InstructionWriter {
	iw.I = append(iw.I, i)
	return iw
}

func (iw *InstructionWriter) Union(i Instructions) *InstructionWriter {
	iw.I = slices.Concat(iw.I, i)
	return iw
}

func (iw *InstructionWriter) WriteOp(o Op) *InstructionWriter {
	return iw.Write(EncodeOp(o))
}

func (iw *InstructionWriter) WriteLoadConst(h intern.Handle[nan.PackedValue]) *InstructionWriter {
	return iw.Write(EncodeOpAndU24(OP_LOAD_CONST, AssertU24(h)))
}

func (iw *InstructionWriter) WriteLoadIdent(kind Op, h intern.Handle[string]) *InstructionWriter {
	return iw.Write(EncodeOpAndU24(kind, AssertU24(h)))
}

func (iw *InstructionWriter) WriteLoadStr(s intern.Handle[string]) *InstructionWriter {
	return iw.Write(EncodeOpAndU24(OP_LOAD_STR, AssertU24(s)))
}

func (iw *InstructionWriter) WriteLoadBool(b bool) *InstructionWriter {
	var v uint8 = 2 // falsey
	if b {
		v = 1
	}

	return iw.Write(Encode(OP_LOAD_BOOL, [3]uint8{v, 0, 0}))
}

func (iw *InstructionWriter) WriteAnd() *InstructionWriter {
	return iw.WriteOp(OP_BOOL_AND)
}

func (iw *InstructionWriter) WriteOr() *InstructionWriter {
	return iw.WriteOp(OP_BOOL_OR)
}

func (iw *InstructionWriter) WriteCall(name intern.Handle[string], arity int) *InstructionWriter {
	var op Op
	switch arity {
	case 0:
		op = OP_CALL_0
		break
	case 1:
		op = OP_CALL_1
		break
	case 2:
		op = OP_CALL_2
		break
	default:
		panic(slipup.Createf("unsupported arity %d", arity))
	}

	return iw.Write(EncodeOpAndU24(op, AssertU24(name)))
}
