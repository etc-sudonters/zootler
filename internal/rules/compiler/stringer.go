package compiler

import (
	"fmt"
	"strings"
	"sudonters/zootler/internal/slipup"
)

func disassemble(c *Chunk, tag string) string {
	b := new(strings.Builder)
	fmt.Fprintf(b, "== %s ==\n", tag)
	fmt.Fprintf(b, "0x00PC   0xOP OP_NAME     OPERANDS\n")
	fmt.Fprintf(b, "----------------------------------\n")

	pc := 0
	size := len(c.Ops)

	W := &bytecodeStringer{b: b, pc: &pc}

	for pc < size {
		op := c.Ops[pc]
		pc++
		switch bc := Bytecode(op); bc {
		case OP_NOP, OP_RETURN:
			W.WriteOp(bc)
			break

		case OP_CONST:
			idx := int(c.Ops[pc])
			pc++
			v := c.Constants[idx]
			W.WriteConst(idx, v)
			break
		default:
			panic(slipup.Create("unknown opcode: 0x%02X @ 0x%04X", op, pc))
		}
		b.WriteRune('\n')
	}

	return b.String()

}

type bytecodeStringer struct {
	b  *strings.Builder
	pc *int
}

func (b *bytecodeStringer) WriteOp(op Bytecode) {
	fmt.Fprintf(b.b, "0x%04X   0x%02X %-16s", *b.pc, uint8(op), op)
}

func (b *bytecodeStringer) WriteConst(idx int, v Value) {
	b.WriteOp(OP_CONST)
	fmt.Fprintf(b.b, " %02d %v", idx, v)
}

func (op Bytecode) String() string {
	switch op {
	case OP_NOP:
		return "OP_NOP"
	case OP_RETURN:
		return "OP_RETURN"
	case OP_CONST:
		return "OP_CONST"
	default:
		panic(slipup.Create("unknown bytecode op '%04X'", uint8(op)))
	}
}
