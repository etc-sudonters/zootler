package objects

import (
	"fmt"
	"slices"
	"sudonters/zootler/mido/symbols"
)

type Index uint16

func TableFrom(builder *Builder) Table {
	var tbl Table
	tbl.strings = make([]byte, len(builder.strings))
	tbl.values = make([]Object, len(builder.values))
	copy(tbl.strings, builder.strings)
	copy(tbl.values, builder.values)
	return tbl
}

type Table struct {
	strings []byte
	values  []Object
}

func (this Table) DecodeString(obj Object) string {
	if !obj.Is(Str32) {
		panic("non-string dereference")
	}

	len, offset := UnpackStr32(obj)
	return string(this.strings[offset : offset+len])
}

func (this Table) AtIndex(idx Index) Object {
	return this.values[idx]
}

type Builder struct {
	nums map[float64]Index
	strs map[string]Index
	ptrs map[symbols.Index]Index
	defs map[symbols.Index]BuiltInFunctionDef

	strings []byte
	values  []Object
}

func NewTableBuilder() Builder {
	return Builder{
		nums:    make(map[float64]Index),
		strs:    make(map[string]Index),
		ptrs:    make(map[symbols.Index]Index),
		defs:    make(map[symbols.Index]BuiltInFunctionDef),
		strings: make([]byte, 0),
		values:  make([]Object, 0),
	}
}

func (this *Builder) InternNumber(number float64) Index {
	if idx, exists := this.nums[number]; exists {
		return idx
	}

	idx := this.insert(PackFloat64(number))
	this.nums[number] = idx
	return idx
}

func (this *Builder) DefineFunction(symbol *symbols.Sym, ptr Object, def BuiltInFunctionDef) {
	this.AssociateSymbol(symbol, ptr)
	this.defs[symbol.Index] = def
}

func (this *Builder) FunctionDefinition(symbol *symbols.Sym) BuiltInFunctionDef {
	def, exists := this.defs[symbol.Index]
	if !exists {
		panic(fmt.Errorf("%q does not have a mapped definition", symbol.Name))
	}
	return def
}

func (this *Builder) AssociateSymbol(symbol *symbols.Sym, ptr Object) Index {
	if index, exists := this.ptrs[symbol.Index]; exists {
		already := this.values[index]
		if already != Object(ptr) {
			panic(fmt.Errorf("ptr for %#v moved: %x -> %x", symbol, already, ptr))
		}

		return index
	}

	idx := this.insert(Object(ptr))
	this.ptrs[symbol.Index] = idx
	return idx
}

func (this *Builder) PtrFor(symbol *symbols.Sym) Index {
	index, exists := this.ptrs[symbol.Index]
	if !exists {
		panic(fmt.Errorf("compile time dereference of null pointer %#v", symbol))
	}

	return index
}

func (this *Builder) InternStr(str string) Index {
	if index, exists := this.strs[str]; exists {
		return index
	}

	offset := len(this.strings)
	bytes := []byte(str)
	this.strings = slices.Concat(this.strings, bytes)
	idx := this.insert(PackStr32(len(bytes), offset))
	this.strs[str] = idx
	return idx
}

func (this *Builder) insert(v Object) Index {
	idx := Index(len(this.values))
	this.values = append(this.values, v)
	return idx
}
