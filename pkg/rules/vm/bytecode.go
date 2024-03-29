package vm

import (
	"fmt"
	"io"
	"strings"
)

type Op uint8
type Ops []Op

const (
	OP_RETURN Op = iota
	OP_COUNT
	OP_CONSTANT // constant id
	OP_SELECT   // column id
	OP_WITH     // column id
	OP_WITHOUT  // column id
)

func (o Op) String() string {
	switch o {
	case OP_RETURN:
		return "OP_RETURN"
	case OP_COUNT:
		return "OP_COUNT"
	case OP_CONSTANT:
		return "OP_CONSTANT"
	case OP_SELECT:
		return "OP_SELECT"
	case OP_WITH:
		return "OP_WITH"
	case OP_WITHOUT:
		return "OP_WITHOUT"
	default:
		return fmt.Sprintf("UNKNOWN_OP %03X", uint8(o))
	}
}

type ValKind uint8

const (
	ValNum ValKind = iota
	ValBool
	ValStr
)

func (v ValKind) String() string {
	switch v {
	case ValNum:
		return "Num "
	case ValBool:
		return "Bool"
	case ValStr:
		return "Str "
	default:
		return "????"
	}
}

type Value struct {
	Kind  ValKind
	Value any
}

type Constants []Value

type Chunk struct {
	Name      string
	Code      Ops
	Constants Constants
}

func (c *Chunk) AddConstant(v Value) Op {
	c.Constants = append(c.Constants, v)
	return Op(len(c.Constants) - 1)
}

func (c Chunk) String() string {
	var s strings.Builder
	fmt.Fprintf(&s, "=== %s ===\n", c.Name)
	fmt.Fprintf(&s, "%4s\t%2s %-16s ARGUMENTS\n\n", "PC", "OP", "NAME")
	ops := c.Code
	for i := 0; i < len(ops); {
		i += c.writeOp(i, &s)
	}
	return s.String()
}

func (c Chunk) writeOp(offset int, w io.Writer) int {
	o := c.Code[offset]
	fmt.Fprintf(w, "%04d\t%02X ", offset, uint8(o))
	switch o {
	case OP_RETURN:
		return c.writeZeroArgOp(w, offset)
	case OP_CONSTANT, OP_SELECT, OP_WITH, OP_WITHOUT:
		return c.writeConstantOp(w, offset)
	default:
		fmt.Fprintf(w, "UNKNOWN: %s\n", o)
		return 1
	}
}

func (c Chunk) writeZeroArgOp(w io.Writer, offset int) int {
	fmt.Fprintf(w, "%s\n", c.Code[offset])
	return 1
}

func (c Chunk) writeConstantOp(w io.Writer, offset int) int {
	idx := c.Code[offset+1]
	constant := c.Constants[idx]
	fmt.Fprintf(w, "%-16s %03d %4s %v\n", c.Code[offset], idx, constant.Kind, constant.Value)
	return 2
}
