package compiler

import (
	"slices"
	"sudonters/zootler/magicbeanvm/code"
)

type Tape struct {
	tape   code.Instructions
	offset int
}

func (bc *Tape) Read() code.Instructions {
	return bc.tape[:]
}

func (bc *Tape) Write(tape code.Instructions) int {
	offset := bc.offset
	bc.tape = slices.Concat(bc.tape, tape)
	bc.offset += len(tape)
	return offset
}
