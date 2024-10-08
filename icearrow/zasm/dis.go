package zasm

import (
	"fmt"
	"strings"

	"github.com/etc-sudonters/substrate/slipup"
)

type diswriter struct {
	strings.Builder
}

func Disassemble(i Instructions) string {
	var dw diswriter
	for _, zasm := range i {
		fmt.Fprintf(&dw, "  0x%08X |\t", uint32(zasm))
		dw.write(zasm.Bytes())
		dw.WriteRune('\n')
	}

	return dw.String()

}

func (w *diswriter) write(d [4]uint8) {
	op := Op(d[0])
	payload := [3]uint8(d[1:4])

	switch op {
	case OP_LOAD_CONST:
		fmt.Fprintf(w, "ldc #consts[0x%08X]", DecodeU24(payload))
		return
	case OP_LOAD_STR:
		fmt.Fprintf(w, "lds #strs[0x%08X]", DecodeU24(payload))
		return
	case OP_LOAD_BOOL:
		fmt.Fprintf(w, "ldb %t", payload[0] == 1)
		return
	case OP_LOAD_SYM:
		fmt.Fprintf(w, "ldi #idents[0x%08X]", DecodeU24(payload))
		return
	case OP_LOAD_TOK:
		fmt.Fprintf(w, "ldt #idents[0x%08X]", DecodeU24(payload))
		return
	case OP_LOAD_VAR:
		fmt.Fprintf(w, "ldv #idents[0x%08X]", DecodeU24(payload))
		return
	case OP_LOAD_TRK:
		fmt.Fprintf(w, "ldr #idents[0x%08X]", DecodeU24(payload))
		return
	case OP_LOAD_SET:
		fmt.Fprintf(w, "lde #idents[0x%08X]", DecodeU24(payload))
		return
	case OP_BOOL_AND:
		fmt.Fprint(w, "and")
		return
	case OP_BOOL_OR:
		fmt.Fprint(w, "orr")
		return
	case OP_CALL_0:
		fmt.Fprintf(w, "cl0 #idents[0x%08X]", DecodeU24(payload))
		return
	case OP_CALL_1:
		fmt.Fprintf(w, "cl1 #idents[0x%08X]", DecodeU24(payload))
		return
	case OP_CALL_2:
		fmt.Fprintf(w, "cl2 #idents[0x%08X]", DecodeU24(payload))
		return
	case 0x0:
		fmt.Fprint(w, "greg")
		return
	default:
		panic(slipup.Createf("unknown op (0x%2X) in instr %s", d[0], d))
	}
}
