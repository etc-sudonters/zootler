package code

import (
	"encoding/binary"
	"fmt"
)

var definitions = map[Op]Defintion{
	BEAN_NOP:          {"NOP", BEAN_NOP, nil},
	BEAN_ERR:          {"ERR", BEAN_ERR, nil},
	BEAN_PUSH_T:       {"PUSH_T", BEAN_PUSH_T, nil},
	BEAN_PUSH_F:       {"PUSH_F", BEAN_PUSH_F, nil},
	BEAN_PUSH_CONST:   {"PUSH_CONST", BEAN_PUSH_CONST, []int{2}},
	BEAN_PUSH_PTR:     {"PUSH_PTR", BEAN_PUSH_PTR, []int{2}},
	BEAN_PUSH_BUILTIN: {"PUSH_BUILTIN", BEAN_PUSH_BUILTIN, []int{2}},
	BEAN_PUSH_OPP:     {"PUSH_OPP", BEAN_PUSH_OPP, nil},
	BEAN_NEED_ALL:     {"NEED_ALL", BEAN_NEED_ALL, []int{2}},
	BEAN_NEED_ANY:     {"NEED_ANY", BEAN_NEED_ANY, []int{2}},
	BEAN_CHK_QTY:      {"CHK_QTY", BEAN_CHK_QTY, []int{2, 1}},
	BEAN_CHK_ALL:      {"CHK_ALL", BEAN_CHK_ALL, []int{2}},
	BEAN_CHK_ANY:      {"CHK_ANY", BEAN_CHK_ANY, []int{2}},
	BEAN_IS_CHILD:     {"IS_CHILD", BEAN_IS_CHILD, nil},
	BEAN_IS_ADULT:     {"IS_ADULT", BEAN_IS_ADULT, nil},
	BEAN_CALL:         {"CALL", BEAN_CALL, []int{2}},
	BEAN_CMP_EQ:       {"CMP_EQ", BEAN_CMP_EQ, nil},
	BEAN_CMP_NQ:       {"CMP_NQ", BEAN_CMP_NQ, nil},
	BEAN_CMP_LT:       {"CMP_LT", BEAN_CMP_LT, nil},
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
	BEAN_NOP          Op = 0x00
	BEAN_ERR             = 0xFF
	BEAN_PUSH_T          = 0x21
	BEAN_PUSH_F          = 0x22
	BEAN_PUSH_CONST      = 0x23
	BEAN_PUSH_PTR        = 0x25
	BEAN_PUSH_BUILTIN    = 0x26
	BEAN_PUSH_OPP        = 0x27
	BEAN_NEED_ALL        = 0x31
	BEAN_NEED_ANY        = 0x32
	BEAN_CHK_QTY         = 0x41
	BEAN_CHK_ALL         = 0x42
	BEAN_CHK_ANY         = 0x43
	BEAN_IS_CHILD        = 0x44
	BEAN_IS_ADULT        = 0x45
	BEAN_CALL            = 0x51
	BEAN_CMP_EQ          = 0x61
	BEAN_CMP_NQ          = 0x62
	BEAN_CMP_LT          = 0x63
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
