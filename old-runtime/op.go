package runtime

import (
	"fmt"
	"math"
	"strings"
)

type Bytecode uint8
type Operand uint8

type Ops []uint8

func (o Ops) String() string {
	var b strings.Builder
	columns := 8
	col := 0

	for i := range o {
		fmt.Fprintf(&b, "0x%02X", o[i])
		col++
		if col%columns == 0 {
			b.WriteRune('\n')
		} else {
			b.WriteRune(' ')
		}
	}

	return b.String()
}

const (
	OP_NOP Bytecode = iota
	OP_RETURN
	OP_SET_RETURN
	OP_LOAD_CONST
	OP_LOAD_IDENT
	OP_EQ
	OP_NEQ
	OP_LT
	OP_AND
	OP_OR
	OP_NOT
	OP_JUMP
	OP_JUMP_FALSE
	OP_JUMP_TRUE
	OP_POP
	OP_CALL0
	OP_CALL1
	OP_CALL2
	// breaks naming scheme on purpose
	DEBUG_STACK_OP = math.MaxUint8
)

var operandCount map[Bytecode]int

func OperandCount(op Bytecode) int {
	return operandCount[op]
}

func init() {
	operandCount = make(map[Bytecode]int)
	operandCount[OP_CALL0] = 1 // func idx
	operandCount[OP_CALL1] = 1
	operandCount[OP_CALL2] = 1
	operandCount[OP_LOAD_CONST] = 1 // const idx
	operandCount[OP_JUMP] = 2       //  PC 1: lower 8bits 2: upper 8bits
	operandCount[OP_JUMP_FALSE] = 2
	operandCount[OP_JUMP_TRUE] = 2
}
