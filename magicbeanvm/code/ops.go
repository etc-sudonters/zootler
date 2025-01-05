package code

import (
	"encoding/binary"
	"fmt"
)

var definitions = map[Op]Defintion{
	BEAN_NOP:        {"BEAN_NOP", BEAN_NOP, nil},
	BEAN_ERR:        {"BEAN_ERR", BEAN_ERR, nil},
	BEAN_PUSH_T:     {"BEAN_PUSH_T", BEAN_PUSH_T, nil},
	BEAN_PUSH_F:     {"BEAN_PUSH_F", BEAN_PUSH_F, nil},
	BEAN_PUSH_CONST: {"BEAN_PUSH_CONST", BEAN_PUSH_CONST, []int{2}},
	BEAN_PUSH_PTR:   {"BEAN_PUSH_PTR", BEAN_PUSH_PTR, []int{2}},
	BEAN_PUSH_FUNC:  {"BEAN_PUSH_FUNC", BEAN_PUSH_FUNC, []int{2}},
	BEAN_NEED_ALL:   {"BEAN_NEED_ALL", BEAN_NEED_ALL, []int{2}},
	BEAN_NEED_ANY:   {"BEAN_NEED_ANY", BEAN_NEED_ANY, []int{2}},
	BEAN_CHK_QTY:    {"BEAN_CHK_QTY", BEAN_CHK_QTY, []int{2, 1}},
	BEAN_CHK_ALL:    {"BEAN_CHK_ALL", BEAN_CHK_ALL, []int{2}},
	BEAN_CHK_ANY:    {"BEAN_CHK_ANY", BEAN_CHK_ANY, []int{2}},
	BEAN_CALL:       {"BEAN_CALL", BEAN_CALL, []int{2}},
}

func Make(op Op, operands ...int) Instructions {
	def, ok := definitions[op]
	if !ok {
		return nil
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
	BEAN_NOP        Op = 0x00
	BEAN_ERR           = 0xFF
	BEAN_PUSH_T        = 0x21
	BEAN_PUSH_F        = 0x22
	BEAN_PUSH_CONST    = 0x23
	BEAN_PUSH_PTR      = 0x25
	BEAN_PUSH_FUNC     = 0x26
	BEAN_NEED_ALL      = 0x31
	BEAN_NEED_ANY      = 0x32
	BEAN_CHK_QTY       = 0x41
	BEAN_CHK_ALL       = 0x42
	BEAN_CHK_ANY       = 0x43
	BEAN_IS_CHILD      = 0x44
	BEAN_IS_ADULT      = 0x45
	BEAN_CALL          = 0x51
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
