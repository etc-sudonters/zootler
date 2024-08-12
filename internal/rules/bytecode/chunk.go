package bytecode

import (
	"fmt"
	"math"
	"strings"
	"sudonters/zootler/internal/slipup"
)

const maxConsts = math.MaxUint8

type PC int
type ConstIdx int
type JumpOffset PC
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

type Chunk struct {
	Ops       Ops
	Constants Values
}

func (c Chunk) Disassemble(tag string) string {
	return disassemble(&c, tag)
}

func (c Chunk) ReadU16(pc PC) uint16 {
	return DecodeU16(c.Ops[pc], c.Ops[pc+1])
}

func (c Chunk) Len() int {
	return len(c.Ops)
}

type ChunkBuilder struct {
	Chunk
}

func (c *ChunkBuilder) LessThan() PC {
	return c.write(OP_LT)
}

func (c *ChunkBuilder) Equal() PC {
	return c.write(OP_EQ)
}

func (c *ChunkBuilder) NotEqual() PC {
	return c.write(OP_NEQ)
}

func (c *ChunkBuilder) JumpFalse() (PC, JumpOffset) {
	return c.writeJump(OP_JUMP_FALSE)
}

func (c *ChunkBuilder) UnconditionalJump() (PC, JumpOffset) {
	return c.writeJump(OP_JUMP)
}

func (c *ChunkBuilder) Pop() PC {
	return c.write(OP_POP)
}

func (c *ChunkBuilder) Dup() PC {
	return c.write(OP_DUP)
}

func (c *ChunkBuilder) Rotate() PC {
	return c.write(OP_ROTATE2)
}

func (c *ChunkBuilder) And() PC {
	return c.write(OP_AND)
}

func (c *ChunkBuilder) Or() PC {
	return c.write(OP_OR)
}

func (c *ChunkBuilder) PatchJump(jmp JumpOffset, target PC) {
	lower, upper := EncodeU16(uint16(target) - uint16(jmp))
	c.Ops[jmp+1] = lower
	c.Ops[jmp+2] = upper
}

func (c *ChunkBuilder) PushConst(v Value) (PC, ConstIdx) {
	where := len(c.Constants)
	if where >= maxConsts {
		panic(slipup.Create("too many constants; maximum allowed %d", maxConsts))
	}
	c.Constants = append(c.Constants, v)
	pc := c.write(OP_CONST, uint8(where))
	return pc, ConstIdx(where)
}

func (c *ChunkBuilder) Return() PC {
	return c.write(OP_RETURN)
}

func (c *ChunkBuilder) DumpStack() PC {
	return c.write(OP_DEBUG_STACK)
}

func (c *ChunkBuilder) write(o Bytecode, operands ...uint8) PC {
	c.Ops = append(c.Ops, uint8(o))
	c.Ops = append(c.Ops, operands...)
	return c.pc()
}

func (c *ChunkBuilder) pc() PC {
	return PC(len(c.Ops) - 1)
}

func (c *ChunkBuilder) writeJump(op Bytecode) (PC, JumpOffset) {
	c.write(op, 0xFF, 0xFF)
	return c.pc(), JumpOffset(c.pc() - 2)
}

func EncodeU16(u16 uint16) (lower uint8, upper uint8) {
	lower, upper = uint8(u16&0x00FF), uint8((u16>>4)&0x00FF)
	return
}

func DecodeU16(lower, upper uint8) uint16 {
	return uint16((upper << 4) | lower)
}
