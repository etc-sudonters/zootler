package zasm

import (
	"fmt"
	"strings"

	"github.com/etc-sudonters/substrate/slipup"
)

type ZasmDisassembler struct{}
type diswriter struct {
	strings.Builder
}

func (_ ZasmDisassembler) Disassemble(instrs Instructions) string {
	var dw diswriter
	for _, zasm := range instrs {
		fmt.Fprintf(&dw, "  0x%08X |\t", uint32(zasm))
		dw.write(dis(zasm.Bytes()))
		dw.WriteRune('\n')
	}

	return dw.String()
}

type dis [4]uint8
