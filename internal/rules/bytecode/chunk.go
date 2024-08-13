package bytecode

import (
	"math"
)

const maxConsts = math.MaxUint8
const maxNames = math.MaxUint8

type PC uint16
type ConstIdx int
type NameIdx int
type Names []string

type Chunk struct {
	Ops       Ops
	Constants Values
	Names     Names
}

func (c Chunk) Disassemble(tag string) string {
	return disassemble(&c, tag)
}

func (c Chunk) ReadU16(pc PC) uint16 {
	return DecodeU16[uint16](c.Ops[pc], c.Ops[pc+1])
}

func (c Chunk) Len() int {
	return len(c.Ops)
}

type ChunkBuilder struct {
	Chunk
}

func (c *ChunkBuilder) LoadConst(v Value) (PC, ConstIdx) {
	id := uint8(len(c.Constants))
	c.Constants = append(c.Constants, v)
	return c.write(OP_LOAD_CONST, id), ConstIdx(id)
}

func (c *ChunkBuilder) LoadIdentifier(v string) (PC, NameIdx) {
	id := uint8(len(c.Names))
	c.Names = append(c.Names, v)
	return c.write(OP_LOAD_IDENT, id), NameIdx(id)

}

func (c *ChunkBuilder) Equal() PC {
	return c.write(OP_EQ)
}

func (c *ChunkBuilder) JumpIfTrue() PC {
	return c.writeJump(OP_JUMP_TRUE)
}

func (c *ChunkBuilder) JumpIfFalse() PC {
	return c.writeJump(OP_JUMP_FALSE)
}

func (c *ChunkBuilder) UnconditionalJump() PC {
	return c.writeJump(OP_JUMP)
}

func (c *ChunkBuilder) DumpStack() PC {
	return c.write(DEBUG_STACK_OP)
}

func (c *ChunkBuilder) SetReturn() PC {
	return c.write(OP_SET_RETURN)
}

func (c *ChunkBuilder) Return() PC {
	return c.write(OP_RETURN)
}

func (c *ChunkBuilder) PatchJump(jump, target PC) {
	lower, upper := EncodeU16(target)
	c.Ops[jump+1] = lower
	c.Ops[jump+2] = upper
}

func (c *ChunkBuilder) write(op Bytecode, operands ...uint8) PC {
	pc := len(c.Ops)
	c.Ops = append(c.Ops, uint8(op))
	c.Ops = append(c.Ops, operands...)
	return PC(pc)
}

func (c *ChunkBuilder) writeJump(op Bytecode) PC {
	return c.write(op, 0x00, 0x00)
}

func EncodeU16[U ~uint16](u16 U) (lower uint8, upper uint8) {
	lower, upper = uint8(u16&0x00FF), uint8((u16&0xFF00)>>4)
	return
}

func DecodeU16[U ~uint16](lower, upper uint8) U {
	return U((upper << 4) | lower)
}
