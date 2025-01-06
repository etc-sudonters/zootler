package symbols

import (
	"fmt"
)

func NewTable() Table {
	return Table{
		names: make(map[string]int),
		syms:  nil,
	}
}

type Table struct {
	names   map[string]int
	syms    []Sym
	aliased int
}

func (tbl *Table) RawAll(f func(*Sym) bool) {
	for i := range tbl.syms {
		if !f(&tbl.syms[i]) {
			return
		}
	}

}

func (tbl *Table) All(f func(*Sym) bool) {
	for i := range tbl.syms {
		if tbl.syms[i].Index != Index(i) {
			continue
		}
		if !f(&tbl.syms[i]) {
			return
		}
	}
}

func (tbl *Table) DeclareMany(typ Kind, names []string) {
	for i := range names {
		tbl.Declare(names[i], typ)
	}
}

func (tbl *Table) Declare(name string, typ Kind) *Sym {
	if sym := tbl.byname(name); sym != nil {
		sym.SetKind(typ)
		return sym
	}

	var sym Sym
	sym.Name = name
	sym.Index = Index(len(tbl.syms))
	sym.Kind = typ
	tbl.syms = append(tbl.syms, sym)
	tbl.names[name] = int(sym.Index)
	return &tbl.syms[int(sym.Index)]
}

func (tbl *Table) LookUpByName(name string) *Sym {
	return tbl.byname(name)
}

func (tbl *Table) LookUpByIndex(index Index) *Sym {
	return tbl.byindex(int(index))
}

func (tbl *Table) byindex(idx int) *Sym {
	symbol := &tbl.syms[idx]
	if int(symbol.Index) != idx {
		symbol = &tbl.syms[int(symbol.Index)]
	}
	return symbol
}

func (tbl *Table) byname(name string) *Sym {
	if idx, exists := tbl.names[name]; exists {
		return tbl.byindex(idx)
	}
	return nil
}

func (tbl *Table) Alias(symbol *Sym, alias string) {
	aliasing := tbl.Declare(alias, ALIAS)
	aliasing.Index = symbol.Index
	tbl.aliased += 1
}

func (tbl *Table) Size() int {
	return len(tbl.syms) - tbl.aliased
}

func (tbl *Table) RawSize() int {
	return len(tbl.syms)
}

func (tbl *Table) AliasCount() int {
	return tbl.aliased
}

type Sym struct {
	Name  string
	Index Index
	Kind  Kind
}

func (s *Sym) String() string {
	return fmt.Sprintf("{Name: %q, Index: %d, Kind: %q}", s.Name, s.Index, s.Kind)
}

func (s *Sym) SetKind(t Kind) {
	switch {
	case t == UNKNOWN || t == ALIAS:
		return
		// event is more specific than token
	case t == EVENT && s.Kind == TOKEN:
		s.Kind = EVENT
	case s.Kind == EVENT && t == TOKEN:
		return
		// aliased over
	case t == COMPILED_FUNC && s.Kind == TOKEN:
		s.Kind = COMPILED_FUNC
		// function is least specific
	case t == FUNCTION && (s.Kind == BUILT_IN || s.Kind == COMP_FUNC || s.Kind == COMPILED_FUNC):
		return
	case s.Kind == UNKNOWN:
		s.Kind = t
		// rudimentary type checking
	case s.Kind != t:
		panic(fmt.Errorf("$%04X %q redeclared with different kind: %q -> %q", s.Index, s.Name, s.Kind, t))
	}
}

func (s *Sym) Eq(o *Sym) bool {
	this, other := s.Index, o.Index
	return this == other
}

type Index int
type Kind string

const (
	_             Kind = ""
	BUILT_IN           = "BUILT_IN"
	COMP_FUNC          = "COMP_FUNC"
	FUNCTION           = "FUNCTION"
	GLOBAL             = "GLOBAL"
	LOCATION           = "LOCATION"
	EVENT              = "EVENT"
	SETTING            = "SETTING"
	TOKEN              = "TOKEN"
	UNKNOWN            = "UNKNOWN"
	COMPILED_FUNC      = "COMPILED_FUNC"
	TRANSIT            = "TRANSIT"
	ALIAS              = "ALIAS"
	LOCAL              = "LOCAL"
)
