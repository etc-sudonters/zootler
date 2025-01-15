package ast

import (
	"fmt"
	"sudonters/zootler/internal/ruleparser"
	"sudonters/zootler/mido/symbols"

	"github.com/etc-sudonters/substrate/peruse"
)

func BuildCompilingFunctionTable(symbolTable *symbols.Table, grammar peruse.Grammar[ruleparser.Tree], decls map[string]string) (PartialFunctionTable, error) {
	funcTable := PartialFunctionTable{make(map[string]PartialFunction, len(decls))}
	bodies := make(map[string]string, len(decls))

	for header, body := range decls {
		var decl PartialFunction
		head, declErr := Parse(header, symbolTable, grammar)
		if declErr != nil {
			panic(declErr)
		}

		switch head := head.(type) {
		case Invoke:
			decl.Symbol = LookUpNodeInTable(symbolTable, head.Target)
			if decl.Symbol.Kind == symbols.BUILT_IN_FUNCTION {
				continue
			}
			if decl.Symbol == nil {
				panic(fmt.Errorf("did not find entry for %#v", head))
			}
			decl.Params = make([]Identifier, len(head.Args))
			for i := range decl.Params {
				param := head.Args[i].(Identifier)
				symbol := symbolTable.LookUpByIndex(param.AsIndex())
				symbol.SetKind(symbols.LOCAL)
				decl.Params[i] = param
			}
		case Identifier:
			decl.Symbol = symbolTable.LookUpByIndex(head.AsIndex())
			if decl.Symbol.Kind == symbols.BUILT_IN_FUNCTION {
				continue
			}
			decl.Params = nil
		default:
			panic(fmt.Errorf("expected identifier or invoke style declaration, got %#v", head))
		}

		decl.Symbol.SetKind(symbols.SCRIPTED_FUNC)
		bodies[decl.Symbol.Name] = body
		funcTable.add(decl.Symbol, decl)
	}

	for name, compiling := range funcTable.tbl {
		body := bodies[name]
		nodes, bodyErr := Parse(body, symbolTable, grammar)
		if bodyErr != nil {
			panic(bodyErr)
		}
		compiling.Body = nodes
		funcTable.tbl[name] = compiling
	}

	return funcTable, nil
}

type PartialFunction struct {
	Symbol *symbols.Sym
	Params []Identifier
	Body   Node
}

type PartialFunctionTable struct {
	tbl map[string]PartialFunction
}

func (this *PartialFunctionTable) Get(name string) (PartialFunction, bool) {
	fn, exists := this.tbl[name]
	return fn, exists
}

func (this *PartialFunctionTable) add(sym *symbols.Sym, body PartialFunction) {
	if _, exists := this.tbl[sym.Name]; exists {
		return
	}

	this.tbl[sym.Name] = body
}
