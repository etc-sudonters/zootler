package ast

import (
	"fmt"
	"sudonters/zootler/internal/ruleparser"
	"sudonters/zootler/magicbeanvm/symbols"

	"github.com/etc-sudonters/substrate/peruse"
)

func BuildCompilingFunctionTable(symbolTable *symbols.Table, grammar peruse.Grammar[ruleparser.Tree], decls map[string]string) (CompilingFunctions, error) {
	funcTable := make(functbl, len(decls))
	bodies := make(map[symbols.Index]string, len(decls))

	for header, body := range decls {
		var decl CompilingFunction
		head, declErr := Parse(header, symbolTable, grammar)
		if declErr != nil {
			panic(declErr)
		}

		switch head := head.(type) {
		case Invoke:
			decl.Symbol = LookUpNodeInTable(symbolTable, head.Target)
			if decl.Symbol.Kind == symbols.BUILT_IN {
				continue
			}
			if decl.Symbol == nil {
				panic(fmt.Errorf("did not find entry for %#v", head))
			}
			decl.Params = make([]Identifier, len(head.Args))
			for i := range decl.Params {
				param := head.Args[i].(Identifier)
				param.Symbol.SetKind(symbols.LOCAL)
				decl.Params[i] = param
			}
		case Identifier:
			decl.Symbol = symbolTable.LookUpByIndex(head.AsIndex())
			if decl.Symbol.Kind == symbols.BUILT_IN {
				continue
			}
			decl.Params = nil
		default:
			panic(fmt.Errorf("expected identifier or invoke style declaration, got %#v", head))
		}

		decl.Symbol.SetKind(symbols.COMPILED_FUNC)
		id := IdentifierFrom(decl.Symbol)
		bodies[id.AsIndex()] = body
		funcTable[id.AsIndex()] = decl
	}

	for id, decl := range funcTable {
		body := bodies[id]
		nodes, bodyErr := Parse(body, symbolTable, grammar)
		if bodyErr != nil {
			panic(bodyErr)
		}
		decl.Body = nodes
		funcTable[id] = decl
	}

	return CompilingFunctions{funcTable}, nil
}

type CompilingFunction struct {
	Symbol *symbols.Sym
	Params []Identifier
	Body   Node
}

type CompilingFunctions struct {
	tbl functbl
}

func (ft *CompilingFunctions) Get(which Identifier) (CompilingFunction, bool) {
	decl, exists := ft.tbl[which.AsIndex()]
	return decl, exists
}

type functbl map[symbols.Index]CompilingFunction
