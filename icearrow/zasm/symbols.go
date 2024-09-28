package zasm

import (
	"slices"
	"sudonters/zootler/icearrow/ast"
	"sudonters/zootler/icearrow/nan"

	"github.com/etc-sudonters/substrate/slipup"
)

func NewTable() SymbolTable {
	var st SymbolTable
	st.named = make(map[string]int)
	st.interned = make(map[string]int)
	st.constants = make(map[uint64]int)
	st.symbols = []Symbol{{Idx: SymbolIdx(0x00), Type: SYM_NULL}}
	st.strings = []uint8{0x0}
	st.interned[""] = 0
	return st
}

type SymbolIdx uint16
type SymbolType uint8
type Symbol struct {
	Idx   SymbolIdx
	Name  SymbolIdx
	Type  SymbolType
	Value uint64
}

type SymbolTable struct {
	strings   []uint8
	symbols   []Symbol
	named     map[string]int
	interned  map[string]int
	constants map[uint64]int
}

func (st *SymbolTable) DeclareConst(value nan.Packed) *Symbol {
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

func (st *SymbolTable) InternString(value string) *Symbol {
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

func (st *SymbolTable) DeclareName(name string, typ SymbolType) *Symbol {
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

func (st *SymbolTable) mint() *Symbol {
	var sym Symbol
	idx := len(st.symbols)
	symdex := SymbolIdx(idx)
	if idx != int(symdex) {
		panic("symbol table grew too large")
	}
	sym.Idx = symdex
	st.symbols = append(st.symbols, sym)
	return &st.symbols[idx]
}

const (
	SYM_NULL SymbolType = iota
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

func symTypeFromPacked(value nan.Packed) SymbolType {
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

func symTypeFromAst(ident *ast.Identifier) SymbolType {
	switch ident.Kind {
	case ast.AST_IDENT_TOKEN, ast.AST_IDENT_EVENT:
		return SYM_TOK
	case ast.AST_IDENT_VAR:
		return SYM_VAR
	case ast.AST_IDENT_SETTING:
		return SYM_SET
	case ast.AST_IDENT_TRICK:
		return SYM_TRK
	case ast.AST_IDENT_BUILTIN:
		return SYM_FUNC
	case ast.AST_IDENT_SYMBOL:
		return SYM_NAME
	default:
		panic(slipup.Createf("unsupported symbol type: %s", ident.Name))
	}
}
