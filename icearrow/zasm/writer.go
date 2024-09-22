package zasm

import (
	"encoding/binary"
	"io"
)

const (
	ARCH_RULE_START = 0xCEE5A1C0
	ARCH_ARCH_START = 0xD50AF10F
)

type countwriter struct {
	w io.Writer
	n int
}

func (c *countwriter) Write(u8 []byte) (int, error) {
	i, e := c.w.Write(u8)
	c.n += i
	return i, e
}

func WriteAssembly(w io.Writer, asm *Assembly) (n int, err error) {
	counting := countwriter{w, 0}
	defer func() { n = counting.n }()
	header := [20]uint32{
		// marker
		ARCH_ARCH_START,
		// logic files hash
		0, 0, 0, 0,
		0, 0, 0, 0,
		// rules + length
		0, 0,
		// consts + length
		0, 0,
		// names + length
		0, 0,
		// strings + length
		0, 0,
		// labels + length
		0, 0,
	}
	err = binary.Write(&counting, binary.LittleEndian, header)
	if err != nil {
		return
	}

	for _, unit := range asm.units {
		if err = WriteUnit(&counting, unit); err != nil {
			return
		}
	}

	err = EncodeDataTables(&counting, asm.data)
	return
}

func WriteUnit(w io.Writer, unit Unit) (err error) {
	markers := []uint32{
		ARCH_RULE_START, uint32(unit.Id), uint32(len(unit.I)),
	}
	if err = binary.Write(w, binary.LittleEndian, markers); err != nil {
		return
	}
	err = binary.Write(w, binary.LittleEndian, unit.I)
	return
}

func EncodeDataTables(w io.Writer, data Data) (err error) {
	return
}
