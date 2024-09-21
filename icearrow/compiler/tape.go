package compiler

import (
	"fmt"
	"strings"

	"github.com/etc-sudonters/substrate/slipup"
)

func ReadTape(tape Tape) string {
	var str strings.Builder
	tt := tapeTranslator{&str}
	idx, end := 0, len(tape.ops)

	for idx < end {
		op := tape.ops[idx]
		tt.writeOp(op)
		switch op {
		case IA_RETURN:
			tt.padToGutter(0)
			tt.WriteString("RETURN")
			break
		case IA_HAS_QTY:
			tt.writeInlineByte(tape.ops[idx+1])
			tt.writeInlineByte(tape.ops[idx+2])
			tt.writeInlineByte(tape.ops[idx+3])
			idx += 3
			tt.padToGutter(3)
			tt.WriteString("IA_HAS")
			break
		case IA_HAS_ALL:
			tt.padToGutter(0)
			tt.WriteString("IA_HAS_ALL")
			break
		case IA_HAS_ANY:
			tt.padToGutter(0)
			tt.WriteString("IA_HAS_ANY")
			break
		case IA_LOAD_SYMBOL:
			tt.writeInlineByte(tape.ops[idx+1])
			tt.writeInlineByte(tape.ops[idx+2])
			tt.padToGutter(2)
			tt.WriteString("LOAD_SYMBOL")
			idx += 2
			break
		case IA_LOAD_CONST:
			tt.writeInlineByte(tape.ops[idx+1])
			tt.writeInlineByte(tape.ops[idx+2])
			tt.padToGutter(2)
			tt.WriteString("LOAD_CONST")
			idx += 2
			break
		case IA_LOAD_IMMED:
			tt.writeInlineByte(tape.ops[idx+1])
			tt.padToGutter(1)
			tt.WriteString("LOAD_U8")
			idx++
			break
		case IA_LOAD_IMMED2:
			tt.writeInlineByte(tape.ops[idx+1])
			tt.writeInlineByte(tape.ops[idx+2])
			tt.padToGutter(1)
			tt.WriteString("LOAD_U8")
			idx += 2
			break
		case IA_LOAD_FALSE:
			tt.padToGutter(0)
			tt.WriteString("LOAD_FALSE")
			break
		case IA_LOAD_TRUE:
			tt.padToGutter(0)
			tt.WriteString("LOAD_TRUE")
			break
		case IA_REDUCE_ALL:
			tt.padToGutter(0)
			tt.WriteString("REDUCE_ALL")
			break
		case IA_REDUCE_ANY:
			tt.padToGutter(0)
			tt.WriteString("REDUCE_ANY")
			break
		case IA_CALL_0:
			tt.padToGutter(0)
			tt.WriteString("CALL_0")
			break
		case IA_CALL_1:
			tt.padToGutter(0)
			tt.WriteString("CALL_1")
			break
		case IA_CALL_2:
			tt.padToGutter(0)
			tt.WriteString("CALL_2")
			break
		case TEMP_IA_LOAD_STR:
			tt.writeInlineByte(tape.ops[idx+1])
			tt.writeInlineByte(tape.ops[idx+2])
			tt.padToGutter(2)
			tt.WriteString("LOAD_STRING")
			idx += 2
			break
		default:
			panic(slipup.Createf("unknown op 0x%02X", op))
		}
		tt.WriteRune('\n')
		idx++
	}

	return str.String()
}

type tapeTranslator struct {
	*strings.Builder
}

func (tt tapeTranslator) writeOp(op uint8) {
	fmt.Fprintf(tt, "0x%02X | ", op)
}

func (tt tapeTranslator) writeInlineByte(val uint8) {
	fmt.Fprintf(tt, "0x%02X ", val)
}

func (tt tapeTranslator) padToGutter(inline int) {
	tt.WriteString(strings.Repeat(" ", 15-(5*inline)))
	tt.WriteRune('|')
	tt.WriteRune(' ')
}

type Tape struct {
	ops []uint8
}

func (t *Tape) write(op IceArrowOp, u8s ...uint8) {
	tt := t.ops
	tt = append(tt, uint8(op))
	tt = append(tt, u8s...)
	t.ops = tt
}

func (t *Tape) writeReturn() {
	t.write(IA_RETURN)
}

func (t *Tape) writeLoadConst(handle uint16) {
	t.write(IA_LOAD_CONST, encodeU16(handle)...)
}

func (t *Tape) writeLoadSymbol(handle uint16) {
	t.write(IA_LOAD_SYMBOL, encodeU16(handle)...)
}

func (t *Tape) writeLoadString(handle uint16) {
	t.write(TEMP_IA_LOAD_STR, encodeU16(handle)...)
}

func (t *Tape) writeLoadTrue() {
	t.write(IA_LOAD_TRUE)
}

func (t *Tape) writeLoadFalse() {
	t.write(IA_LOAD_FALSE)
}

func (t *Tape) writeLoadImmediateU8(u8 uint8) {
	t.write(IA_LOAD_IMMED, u8)
}

func (t *Tape) writeLoadImmediateU16(u16 uint16) {
	t.write(IA_LOAD_IMMED, encodeU16(u16)...)
}

func (t *Tape) writeReduceAll(qty uint8) {
	t.writeLoadImmediateU8(qty)
	t.write(IA_REDUCE_ALL)
}

func (t *Tape) writeReduceAny(qty uint8) {
	t.writeLoadImmediateU8(qty)
	t.write(IA_REDUCE_ANY)
}

func (t *Tape) writeJump() {
	t.writeLoadImmediateU16(0x00)
	t.write(IA_JUMP_UNCOND)
}

func (t *Tape) writeJumpTrue() {
	t.writeLoadImmediateU16(0x00)
	t.write(IA_JUMP_TRUE)
}

func (t *Tape) writeJumpFalse() {
	t.writeLoadImmediateU16(0x00)
	t.write(IA_JUMP_FALSE)
}

func (t *Tape) writeCall(handle uint16, args uint8) {
	t.writeLoadSymbol(handle)
	switch args {
	case 0:
		t.write(IA_CALL_0)
	case 1:
		t.write(IA_CALL_1)
	case 2:
		t.write(IA_CALL_2)
	default:
		panic("unsupported argument count")
	}
}

func encodeU16(u16 uint16) []uint8 {
	var enc []uint8 = []uint8{0, 0}
	enc[0] = uint8((0x00FF & u16))
	enc[1] = uint8((0xFF00 & u16) >> 8)
	return enc
}
