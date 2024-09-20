package compiler

import (
	"sudonters/zootler/icearrow/zasm"
	"sudonters/zootler/internal"
)

func CreateSymbolTable(data *zasm.Data) SymbolTable {
	var st SymbolTable
	st.byname = make(map[string]int)
	st.strings = make([]String, len(data.Strs))
	st.symbols = make([]Symbol, len(data.Names))
	st.consts = make([]Const, len(data.Consts))

	for i := range data.Strs {
		st.strings[i] = String{
			Id:    uint32(i),
			Value: data.Strs[i],
		}
	}

	for i := range data.Consts {
		st.consts[i] = Const{
			Id:    uint32(i),
			Value: data.Consts[i],
		}
	}

	for i := range data.Names {
		st.symbols[i] = Symbol{
			Id:   uint32(i),
			Name: data.Names[i],
		}
		st.byname[data.Names[i]] = i
	}

	st.Declare("has_all", SYM_KIND_CALLABLE)
	st.Declare("has_any", SYM_KIND_CALLABLE)

	return st
}

type SymbolTable struct {
	byname  map[string]int
	symbols []Symbol
	consts  []Const
	strings []String
}

func (st *SymbolTable) ConstOf(pv zasm.PackedValue) *Const {
	for i := range st.consts {
		sym := &st.consts[i]
		if sym.Value.Equals(pv) {
			return sym
		}
	}
	return nil

}

func (st *SymbolTable) Named(name string) *Symbol {
	normaled := string(internal.Normalize(name))
	if idx, found := st.byname[normaled]; found {
		return &st.symbols[idx]
	}

	for i := range st.symbols {
		if st.symbols[i].Name == name {
			st.byname[normaled] = i
			return &st.symbols[i]
		}
	}

	return nil
}

func (st *SymbolTable) Const(id uint32) *Const {
	return &st.consts[id-1]
}

func (st *SymbolTable) String(id uint32) *String {
	return &st.strings[id-1]
}

func (st *SymbolTable) Symbol(id uint32) *Symbol {
	return &st.symbols[id-1]
}

func (st *SymbolTable) Declare(name string, kind SymbolKind) *Symbol {
	id := uint32(len(st.symbols) + 1)
	st.symbols = append(st.symbols, Symbol{
		Id:   id,
		Name: name,
		Kind: kind,
	})
	st.byname[string(internal.Normalize(name))] = int(id - 1)
	return &st.symbols[id-1]
}

type SymbolKind uint8

const (
	_ SymbolKind = iota
	SYM_KIND_CALLABLE
	SYM_KIND_TOKEN
	SYM_KIND_VAR
	SYM_KIND_SYMBOL
	SYM_KIND_TRICK
	SYM_KIND_SETTING
)

type Symbol struct {
	Id   uint32
	Kind SymbolKind
	Name string
}

func (s *Symbol) Set(kind SymbolKind) {
	if s.Kind != 0 && s.Kind != kind {
		panic("symbol kind already set")
	}

	s.Kind = kind
}

type Const struct {
	Id    uint32
	Value zasm.PackedValue
}

type String struct {
	Id    uint32
	Value string
}
