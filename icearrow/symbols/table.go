package symbols

import (
	"slices"
	"sudonters/zootler/icearrow/nan"
)

func NewTable() Table {
	var st Table
	st.named = make(map[string]int)
	st.interned = make(map[string]int)
	st.constants = make(map[uint64]int)
	st.symbols = []Entry{{Idx: Idx(0x00), Type: SYM_NULL}}
	st.strings = []uint8{0x0}
	st.interned[""] = 0
	return st
}

type Idx uint16
type Type uint8
type Entry struct {
	Idx   Idx
	Name  Idx
	Type  Type
	Value uint64
}

type Table struct {
	symbols   []Entry
	strings   []uint8
	constants map[uint64]int

	named    map[string]int
	interned map[string]int
}

func (st *Table) ByIdx(i Idx) *Entry {
	return &st.symbols[int(i)]
}

func (st *Table) DeclareConst(value nan.Packed) *Entry {
	bits := value.Bits()
	if sym, exists := st.constants[bits]; exists {
		return &st.symbols[int(sym)]
	}

	sym := st.mint()
	sym.Value = bits
	sym.Type = symTypeFromPacked(value)
	st.constants[bits] = int(sym.Idx)
	return sym
}

func RepackSymbol(s *Entry) nan.Packed {
	return nan.PackBits(s.Value)
}

func (st *Table) InternString(value string) *Entry {
	if sym, exists := st.interned[value]; exists {
		return &st.symbols[sym]
	}

	bytes := append([]uint8(value), 0x0)
	length := uint64(len(bytes))
	startOf := uint64(len(st.strings))
	st.strings = slices.Concat(st.strings, bytes)

	sym := st.mint()
	sym.Type = SYM_STR
	sym.Value = length<<32 | startOf
	return sym
}

func (st *Table) DivulgeString(sym *Entry) string {
	if sym.Type != SYM_STR {
		panic("not a string")
	}

	start := (sym.Value << 32) >> 32
	length := sym.Value >> 32
	bytes := st.strings[start : start+length]
	return string(bytes)
}

func (st *Table) DeclareName(name string, typ Type) *Entry {
	if typ == SYM_NULL {
		panic("cannot declare null symbol")
	}

	if name == "" {
		panic("cannot declare empty name")
	}

	if sym, exists := st.named[name]; exists {
		return &st.symbols[sym]
	}

	str := st.InternString(name)
	sym := st.mint()
	sym.Type = typ
	sym.Name = str.Idx
	return sym
}

func (st *Table) mint() *Entry {
	var sym Entry
	idx := len(st.symbols)
	symdex := Idx(idx)
	if idx != int(symdex) {
		panic("symbol table grew too large")
	}
	sym.Idx = symdex
	st.symbols = append(st.symbols, sym)
	return &st.symbols[idx]
}

const (
	SYM_NULL Type = iota
	SYM_BOOL
	SYM_STR
	SYM_UNUM
	SYM_SNUM
	SYM_F64
	SYM_FUNC
	SYM_TOK
	SYM_VAR
	SYM_SET
	SYM_TRK
	SYM_NAME
)

func symTypeFromPacked(value nan.Packed) Type {
	switch value.Type() {
	case nan.PT_F64:
		return SYM_F64
	case nan.PT_BOOL:
		return SYM_BOOL
	case nan.PT_UNUM:
		return SYM_UNUM
	case nan.PT_SNUM:
		return SYM_SNUM
	case nan.PT_STRING:
		return SYM_STR
	case nan.PT_FUNC:
		return SYM_FUNC
	case nan.PT_TOKEN:
		return SYM_TOK
	case nan.PT_TRICK:
		return SYM_TRK
	case nan.PT_SETTING:
		return SYM_SET
	case nan.PT_VAR:
		return SYM_VAR
	default:
		panic("unknown packed value type")
	}
}
