package bytecode

import "math"

type Bytecode uint8

const (
	OP_NOP Bytecode = iota
	OP_RETURN
	OP_CONST

	OP_EQ
	OP_NEQ
	OP_LT

	OP_CONTAIN
	OP_AND
	OP_OR
	OP_NOT

	OP_DUP
	OP_ROTATE2

	OP_JUMP
	OP_JUMP_FALSE
	OP_POP

	OP_DEBUG_STACK = math.MaxUint8
)

var operandCount map[Bytecode]int

func OperandCount(op Bytecode) int {
	return operandCount[op]
}

func init() {
	operandCount = make(map[Bytecode]int)
	operandCount[OP_JUMP_FALSE] = 3
	operandCount[OP_JUMP] = 3
	operandCount[OP_CONST] = 1
}
