package code

import (
	"fmt"
	"io"
	"strings"
)

type dis struct {
	io.Writer
}

func (this dis) WriteOp(offset int, def Defintion) {
	fmt.Fprintf(this, "0x%02X | ", offset)
	this.writeu8(uint8(def.Op))
	fmt.Fprintf(this, " | %-16s", def.Name)
}

func (this dis) CopyU16(tape []byte) {
	fmt.Fprintf(this, "0x%02X", ReadU16(tape))
}

func (this dis) CopyU8(tape []byte) {
	this.writeu8(tape[0])
}

func (this dis) Clear() {
	fmt.Fprintln(this)
}

func (this dis) writeu8(u8 uint8) {
	fmt.Fprintf(this, "0x%02X", u8)

}

func Disassemble(tape Instructions) string {
	var s strings.Builder
	var offset int
	dis := dis{&s}

	length := len(tape)

	for offset < length {
		def := definitions[Op(tape[offset])]
		dis.WriteOp(offset, def)
		offset += 1
		for _, width := range def.Operands {
			fmt.Fprint(dis, " | ")
			switch width {
			case 1:
				dis.CopyU8(tape[offset:])
			case 2:
				dis.CopyU16(tape[offset:])
			}
			offset += width
		}
		dis.Clear()
	}

	return s.String()
}
