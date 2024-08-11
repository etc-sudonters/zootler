package compiler

import (
	"math"
	"sudonters/zootler/internal/slipup"
)

type Value float64
type Values []Value

const maxConsts = math.MaxUint8

type Bytecode uint8

const (
	OP_NOP Bytecode = iota
	OP_RETURN
	OP_CONST
)

type Chunk struct {
	Ops       []uint8
	Constants Values
}

type ChunkBuilder struct {
	Chunk
}

func (c *ChunkBuilder) WriteOp(o Bytecode, operands ...uint8) {
	c.Ops = append(c.Ops, uint8(o))
	c.Ops = append(c.Ops, operands...)
}

func (c *ChunkBuilder) PushConst(v Value) int {
	where := len(c.Constants)
	if where >= maxConsts {
		panic(slipup.Create("too many constants; maximum allowed %d", maxConsts))
	}
	c.Constants = append(c.Constants, v)
	c.WriteOp(OP_CONST, uint8(where))
	return where
}

func (c Chunk) Disassemble(tag string) string {
	return disassemble(&c, tag)
}
