package compiler

type RuleCompiler struct{}

type Compiling struct {
	tape tape
}

type tape []uint8

func (t *tape) write(op IceArrowOp, u8s ...uint8) int {
	tt := *t
	l := len(tt)
	tt = append(tt, uint8(op))
	tt = append(tt, u8s...)
	*t = tt
	return l
}

func (t *tape) writeReturn() int {
	return t.write(IA_RETURN)
}

func (t *tape) writeLoadConst(handle uint16) int {
	return t.write(IA_LOAD_CONST, encodeU16(handle)...)
}

func (t *tape) writeLoadSymbol(handle uint16) int {
	return t.write(IA_LOAD_SYMBOL, encodeU16(handle)...)
}

func (t *tape) writeLoadTrue() int {
	return t.write(IA_LOAD_TRUE)
}

func (t *tape) writeLoadFalse() int {
	return t.write(IA_LOAD_FALSE)
}

func (t *tape) writeLoadImmediateU8(u8 uint8) int {
	return t.write(IA_LOAD_IMMED, u8)
}

func (t *tape) writeLoadImmediateU16(u16 uint16) int {
	return t.write(IA_LOAD_IMMED, encodeU16(u16)...)
}

func (t *tape) writeReduceAnd() int {
	return t.write(IA_REDUCE_AND)
}

func (t *tape) writeReduceOr() int {
	return t.write(IA_REDUCE_OR)
}

func (t *tape) writeReduceAll(qty uint8) int {
	start := t.writeLoadImmediateU8(qty)
	t.write(IA_REDUCE_ALL)
	return start
}

func (t *tape) writeReduceAny(qty uint8) int {
	start := t.writeLoadImmediateU8(qty)
	t.write(IA_REDUCE_ANY)
	return start
}

func (t *tape) writeJump() int {
	dest := t.writeLoadImmediateU16(0x00)
	t.write(IA_JUMP_UNCOND)
	return dest
}

func (t *tape) writeJumpTrue() int {
	dest := t.writeLoadImmediateU16(0x00)
	t.write(IA_JUMP_TRUE)
	return dest
}

func (t *tape) writeJumpFalse() int {
	dest := t.writeLoadImmediateU16(0x00)
	t.write(IA_JUMP_FALSE)
	return dest
}

func (t *tape) writeCall(handle uint16, args int) int {
	t.writeLoadSymbol(handle)
	switch args {
	case 0:
		return t.write(IA_CALL_0)
	case 1:
		return t.write(IA_CALL_1)
	case 2:
		return t.write(IA_CALL_2)
	default:
		panic("unreachable")
	}
}

func encodeU16(u16 uint16) []uint8 {
	var enc []uint8 = []uint8{0, 0}
	enc[0] = uint8((0x00FF & u16))
	enc[1] = uint8((0xFF00 & u16) >> 8)
	return enc
}
