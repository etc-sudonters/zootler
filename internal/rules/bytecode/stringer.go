package bytecode

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
	pos := -1

	writer := &bytecodewriter{b: b}

	for pc < size {
		if pos == pc {
			panic(fmt.Errorf(
				"'%s' did not increment program counter",
				Bytecode(c.Ops[pos]),
			))
		}
		pos = pc
		op := c.Ops[pc]
		switch bc := Bytecode(op); bc {
		case OP_NOP, OP_RETURN, OP_DEBUG_STACK,
			OP_EQ, OP_NEQ, OP_LT,
			OP_DUP, OP_ROTATE2,
			OP_AND, OP_OR:
			writer.WriteOp(pc, bc)
			pc++
			break
		case OP_CONST:
			idx := int(c.Ops[pc+1])
			v := c.Constants[idx]
			writer.WriteConst(pc, idx, v)
			pc += 2
			break
		case OP_JUMP_FALSE, OP_JUMP:
			writer.WriteJump(pc, bc, c.Ops[pc+1], c.Ops[pc+2])
			pc += 3
		default:
			panic(slipup.Create("unknown opcode: 0x%02X @ 0x%04X", op, pc))
		}
		b.WriteRune('\n')
	}

	return b.String()

}

type bytecodewriter struct {
	b *strings.Builder
}

func (b *bytecodewriter) WriteOp(pos int, op Bytecode) {
	fmt.Fprintf(b.b, "0x%04X   0x%02X %-16s", pos, uint8(op), op)
}

func (b *bytecodewriter) WriteConst(pos int, idx int, v Value) {
	b.WriteOp(pos, OP_CONST)
	fmt.Fprintf(b.b, " %02X %+v", idx, v)
}

func (b *bytecodewriter) WriteJump(pos int, jmp Bytecode, lower, upper uint8) {
	b.WriteOp(pos, jmp)
	b.writeInlineOperand(lower)
	offset := DecodeU16(lower, upper)
	dest := pos + int(offset)
	fmt.Fprintf(b.b, " 0x%04X", offset)
	fmt.Fprintf(b.b, " 0x%04X", dest)
	b.writeOperandLine(upper)
	if jmp == OP_JUMP_FALSE {
		fmt.Fprintf(b.b, " 0x%04X", pos+OperandCount(jmp))
	}
}

func (b *bytecodewriter) writeInlineOperand(v uint8) {
	fmt.Fprintf(b.b, " 0x%02X", v)
}

func (b *bytecodewriter) writeOperandLine(v uint8) {
	fmt.Fprintf(b.b, "\n                               0x%02X", v)
}

func (op Bytecode) String() string {
	switch op {
	case OP_NOP:
		return "OP_NOP"
	case OP_RETURN:
		return "OP_RETURN"
	case OP_CONST:
		return "OP_CONST"
	case OP_EQ:
		return "OP_EQ"
	case OP_NEQ:
		return "OP_NEQ"
	case OP_LT:
		return "OP_LT"
	case OP_CONTAIN:
		return "OP_CONTAIN"
	case OP_AND:
		return "OP_AND"
	case OP_OR:
		return "OP_OR"
	case OP_NOT:
		return "OP_NOT"
	case OP_DUP:
		return "OP_DUP"
	case OP_ROTATE2:
		return "OP_ROTATE2"
	case OP_JUMP:
		return "OP_JUMP"
	case OP_JUMP_FALSE:
		return "OP_JUMP_FALSE"
	case OP_POP:
		return "OP_POP"
	case OP_DEBUG_STACK:
		return "OP_DEBUG_STACK"
	default:
		panic(slipup.Create("unknown bytecode op '%04X'", uint8(op)))
	}
}
