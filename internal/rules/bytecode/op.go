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
	OP_ROTATE

	OP_JMP_TRUE
	OP_JMP_FALSE

	OP_DEBUG = math.MaxUint8
)

var operandCount map[Bytecode]int

func OperandCount(op Bytecode) int {
	return operandCount[op]
}

func init() {
	operandCount = make(map[Bytecode]int)
	operandCount[OP_JMP_FALSE] = 3
	operandCount[OP_JMP_TRUE] = 3
	operandCount[OP_CONST] = 1
}
