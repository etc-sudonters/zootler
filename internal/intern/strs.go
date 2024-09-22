package intern

import "github.com/etc-sudonters/substrate/slipup"

func NewStrPiler() StrPiler {
	var h StrPiler
	h.pile = make(StrPile, 0, 2048)
	h.strs = make(map[string]Str, 32)
	return h
}

type StrPiler struct {
	pile StrPile
	strs map[string]Str
}

func (h *StrPiler) Pile() StrPile {
	return h.pile
}

func (h *StrPiler) Intern(interning string) Str {
	if s, interned := h.strs[interning]; interned {
		return s
	}

	bytes := []uint8(interning)
	str := encodeStrFromIntPair(len(h.pile), len(bytes))
	h.pile = append(h.pile, bytes...)
	return str
}

func encodeStrFromIntPair(offset, length int) Str {
	var s Str
	off16 := uint16(offset)
	len8 := uint8(length)

	if int(off16) != offset {
		panic(slipup.Createf("string pile would become too large: %d", offset+length))
	}

	if int(len8) != length {
		panic(slipup.Createf("string too long: %d", length))
	}

	s[0] = uint8(off16 & 0x00FF)
	s[1] = uint8((off16 & 0xFF00) >> 8)
	s[2] = len8
	return s
}

type StrPile []uint8

func (pile StrPile) Retrieve(str Str) string {
	bytes := pile[str.Offset() : str.Offset()+str.Len()]
	s := string(bytes)
	return s
}

type Str [3]uint8

func (s Str) Bytes() [3]uint8 {
	return [3]uint8(s)
}

func (s Str) Offset() int {
	var i int

	i = i | (int(s[1]) << 8)
	i = i | int(s[0])
	return i
}

func (s Str) Len() int {
	return int(s[2])
}
