package code

import (
	"encoding/binary"
	"fmt"
)

var definitions = map[Op]Defintion{
	NOP:          {"NOP", NOP, nil},
	ERR:          {"ERR", ERR, nil},
	PUSH_T:       {"PUSH_T", PUSH_T, nil},
	PUSH_F:       {"PUSH_F", PUSH_F, nil},
	PUSH_CONST:   {"PUSH_CONST", PUSH_CONST, []int{2}},
	PUSH_PTR:     {"PUSH_PTR", PUSH_PTR, []int{2}},
	PUSH_BUILTIN: {"PUSH_BUILTIN", PUSH_BUILTIN, []int{2}},
	INVERT:       {"INVERT", INVERT, nil},
	NEED_ALL:     {"NEED_ALL", NEED_ALL, []int{2}},
	NEED_ANY:     {"NEED_ANY", NEED_ANY, []int{2}},
	CHK_QTY:      {"CHK_QTY", CHK_QTY, []int{2, 1}},
	CHK_ALL:      {"CHK_ALL", CHK_ALL, []int{2}},
	CHK_ANY:      {"CHK_ANY", CHK_ANY, []int{2}},
	IS_CHILD:     {"IS_CHILD", IS_CHILD, nil},
	IS_ADULT:     {"IS_ADULT", IS_ADULT, nil},
	INVOKE:       {"INVOKE", INVOKE, []int{2}},
	CMP_EQ:       {"CMP_EQ", CMP_EQ, nil},
	CMP_NQ:       {"CMP_NQ", CMP_NQ, nil},
	CMP_LT:       {"CMP_LT", CMP_LT, nil},
}

func Make(op Op, operands ...int) Instructions {
	def, ok := definitions[op]
	if !ok {
		return nil
	}

	if len(operands) != len(def.Operands) {
		panic(fmt.Errorf(
			"0x%02X expects %d operands, received %d",
			op, len(def.Operands), len(operands),
		))
	}

	instructionLen := 1
	for _, width := range def.Operands {
		instructionLen += width
	}

	tape := make(Instructions, instructionLen)
	tape[0] = byte(op)
	offset := 1
	for idx, operand := range operands {
		width := def.Operands[idx]
		switch width {
		case 1:
			tape[offset] = byte(operand)
		case 2:
			binary.LittleEndian.PutUint16(tape[offset:], uint16(operand))
		default:
			panic(fmt.Errorf("unsupport operand length: %d", width))
		}
		offset += width
	}

	return tape
}

const (
	NOP          Op = 0x00
	ERR             = 0xFF
	PUSH_T          = 0x21
	PUSH_F          = 0x22
	PUSH_CONST      = 0x23
	PUSH_PTR        = 0x25
	PUSH_BUILTIN    = 0x26
	INVERT          = 0x27
	NEED_ALL        = 0x31
	NEED_ANY        = 0x32
	CHK_QTY         = 0x41
	CHK_ALL         = 0x42
	CHK_ANY         = 0x43
	IS_CHILD        = 0x44
	IS_ADULT        = 0x45
	INVOKE          = 0x51
	CMP_EQ          = 0x61
	CMP_NQ          = 0x62
	CMP_LT          = 0x63
)

type Instructions []byte
type Op uint8

func LookUp(op Op) (Defintion, error) {
	var err error
	def, ok := definitions[op]
	if !ok {
		err = fmt.Errorf("unknown op: 0x%02X", op)
	}
	return def, err
}

type Defintion struct {
	Name     string
	Op       Op
	Operands []int
}

func ReadU16(program []byte) uint16 {
	return binary.LittleEndian.Uint16(program)
}

func ReadU8(program []byte) uint8 {
	return program[0]
}
