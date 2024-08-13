package runtime

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
		case OP_NOP, OP_RETURN, DEBUG_STACK_OP,
			OP_EQ, OP_NEQ, OP_LT,
			OP_DUP, OP_ROTATE2,
			OP_AND, OP_OR, OP_SET_RETURN:
			writer.WriteOp(pc, bc)
			pc++
			break
		case OP_LOAD_CONST:
			idx := c.Ops[pc+1]
			v := c.Constants[int(idx)]
			writer.WriteConst(pc, idx, v)
			pc += 2
			break
		case OP_LOAD_IDENT:
			idx := c.Ops[pc+1]
			v := c.Names[int(idx)]
			writer.WriteLoad(pc, idx, v)
			pc += 2
			break
		case OP_JUMP_FALSE, OP_JUMP, OP_JUMP_TRUE:
			writer.WriteJump(pc, bc, c.Ops[pc+1], c.Ops[pc+2])
			pc += 3
		case OP_CALL0, OP_CALL1, OP_CALL2:
			idx := c.Ops[pc+1]
			v := c.Names[int(idx)]
			writer.WriteCall(pc, bc, idx, v)
			pc += 2
			break
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

func (b *bytecodewriter) WriteOp1(pos int, op Bytecode, _1 uint8) {
	b.WriteOp(pos, op)
	b.writeu8(_1)
}

func (b *bytecodewriter) WriteOp2(pos int, op Bytecode, _1, _2 uint8) {
	b.WriteOp1(pos, op, _1)
	b.writeu8(_2)
}

func (b *bytecodewriter) WriteConst(pos int, idx uint8, v Value) {
	b.WriteOp1(pos, OP_LOAD_CONST, idx)
	b.writeLiteral(v)
}

func (b *bytecodewriter) WriteCall(pos int, op Bytecode, idx uint8, name string) {
	b.writeIdentLoad(pos, op, idx, name)
}

func (b *bytecodewriter) WriteLoad(pos int, idx uint8, name string) {
	b.writeIdentLoad(pos, OP_LOAD_IDENT, idx, name)
}

func (b *bytecodewriter) writeIdentLoad(pos int, op Bytecode, idx uint8, name string) {
	b.WriteOp1(pos, op, idx)
	b.writeLiteral(name)
}

func (b *bytecodewriter) WriteJump(pos int, jmp Bytecode, lower, upper uint8) {
	dest := DecodeU16[uint16](lower, upper)
	b.WriteOp(pos, jmp)
	b.writeu8(lower)
	b.writeu8(upper)
	b.writeOperandLine()
	b.writeu16(dest)
	if jmp == OP_JUMP_FALSE || jmp == OP_JUMP_TRUE { // OP_JUMP doesn't branch
		b.writeOperandLine()
		b.writeu16(uint16(pos + OperandCount(jmp)))
	}
}

func (b *bytecodewriter) writeu8(u uint8) {
	fmt.Fprintf(b.b, " 0x%02X", u)
}

func (b *bytecodewriter) writeu16(u uint16) {
	fmt.Fprintf(b.b, " 0x%04X", u)
}

func (b *bytecodewriter) writeOperandLine() {
	fmt.Fprint(b.b, "\n                              ")
}

func (b *bytecodewriter) writeLiteral(v any) {
	fmt.Fprintf(b.b, " %+v", v)
}

func (op Bytecode) String() string {
	switch op {
	case OP_NOP:
		return "OP_NOP"
	case OP_RETURN:
		return "OP_RETURN"
	case OP_LOAD_CONST:
		return "OP_LOAD_CONST"
	case OP_LOAD_IDENT:
		return "OP_LOAD_IDENT"
	case OP_SET_RETURN:
		return "OP_SET_RETURN"
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
	case OP_JUMP_TRUE:
		return "OP_JUMP_TRUE"
	case OP_POP:
		return "OP_POP"
	case DEBUG_STACK_OP:
		return "OP_DEBUG_STACK"
	default:
		panic(slipup.Create("unknown bytecode op '%04X'", uint8(op)))
	}
}
