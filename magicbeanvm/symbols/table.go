package symbols

import (
	"fmt"
	"sudonters/zootler/magicbeanvm/nan"
)

func NewTable() Table {
	return Table{
		names: make(map[string]int),
		syms:  nil,
	}
}

type Table struct {
	names map[string]int
	syms  []Sym
}

func (tbl *Table) DeclareMany(typ Type, names []string) {
	for i := range names {
		tbl.Declare(names[i], typ)
	}
}

func (tbl *Table) Declare(name string, typ Type) *Sym {
	if idx, exists := tbl.names[name]; exists {
		sym := &tbl.syms[idx]
		sym.SetType(typ)
		return sym
	}

	var sym Sym
	sym.Name = name
	sym.Index = Index(len(tbl.syms))
	sym.Type = typ
	tbl.syms = append(tbl.syms, sym)
	tbl.names[name] = int(sym.Index)
	return &tbl.syms[int(sym.Index)]
}

func (tbl *Table) LookUpByName(name string) *Sym {
	if idx, exists := tbl.names[name]; exists {
		return &tbl.syms[idx]
	}
	return nil
}

func (tbl *Table) LookUpByIndex(index Index) *Sym {
	return &tbl.syms[int(index)]
}

func (tbl *Table) Alias(symbol *Sym, alias string) {
	aliasing := tbl.Declare(alias, symbol.Type)
	aliasing.Value = nan.PackPtr(uint32(symbol.Index))
}

type Sym struct {
	Name  string
	Index Index
	Type  Type
	Value nan.PackedValue
}

func (s *Sym) SetType(t Type) {
	switch {
	case t == UNKNOWN:
		return
	case t == FUNCTION && s.Type == TOKEN:
		s.Type = FUNCTION
	case t == FUNCTION && (s.Type == BUILT_IN || s.Type == COMP_TIME):
		return
	case s.Type == UNKNOWN:
		s.Type = t
	case s.Type != t:
		panic(fmt.Errorf("$%04X %q redeclared with different type: %q -> %q", s.Index, s.Name, s.Type, t))
	}
}

func (s *Sym) Eq(o *Sym) bool {
	this, other := s.Index, o.Index

	if thisPtr, isPtrThis := s.Value.Pointer(); isPtrThis {
		this = Index(thisPtr)
	}

	if otherPtr, isPtrOther := o.Value.Pointer(); isPtrOther {
		other = Index(otherPtr)
	}

	return this == other
}

type Index int
type Type string

const (
	_         Type = ""
	BUILT_IN       = "BUILT_IN"
	COMP_TIME      = "COMP_TIME"
	FUNCTION       = "FUNCTION"
	GLOBAL         = "GLOBAL"
	LOCATION       = "LOCATION"
	SETTING        = "SETTING"
	TOKEN          = "TOKEN"
	UNKNOWN        = "UNKNOWN"
)
