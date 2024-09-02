package intern

func NewStrHeaper() StrHeaper {
	var h StrHeaper
	h.heap = make(StrHeap, 0, 2048)
	h.strs = make(map[string]Str, 32)
	return h
}

type StrHeaper struct {
	heap StrHeap
	strs map[string]Str
}

func (h *StrHeaper) Heap() StrHeap {
	return h.heap
}

func (h *StrHeaper) Intern(interning string) Str {
	if s, interned := h.strs[interning]; interned {
		return s
	}

	bytes := []uint8(interning)
	str := encodeStrFromIntPair(len(h.heap), len(bytes))
	h.heap = append(h.heap, bytes...)
	return str
}

func encodeStrFromIntPair(offset, length int) Str {
	var s Str
	off16 := uint16(offset)
	len8 := uint8(length)

	if int(off16) != offset || int(len8) != length {
		panic("string heap too big")
	}

	s[0] = uint8(off16 & 0x00FF)
	s[1] = uint8((off16 & 0xFF00) >> 8)
	s[2] = len8
	return s
}

type StrHeap []uint8

func (heap StrHeap) Retrieve(str Str) string {
	bytes := heap[str.Offset() : str.Offset()+str.Len()]
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
