package ast

import (
	"fmt"
	"sudonters/zootler/internal/ruleparser"
	"sudonters/zootler/magicbeanvm/symbols"

	"github.com/etc-sudonters/substrate/peruse"
)

func BuildFunctionTable(symbolTable *symbols.Table, grammar peruse.Grammar[ruleparser.Tree], decls map[string]string) (FunctionTable, error) {
	funcTable := make(functbl, len(decls))
	bodies := make(map[symbols.Index]string, len(decls))

	for header, body := range decls {
		var decl FunctionDecl
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
				decl.Params[i] = head.Args[i].(Identifier)
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

	return FunctionTable{funcTable}, nil
}

type FunctionDecl struct {
	Symbol *symbols.Sym
	Params []Identifier
	Body   Node
}

type FunctionTable struct {
	tbl functbl
}

func (ft *FunctionTable) Get(which Identifier) (FunctionDecl, bool) {
	decl, exists := ft.tbl[which.AsIndex()]
	return decl, exists
}

type functbl map[symbols.Index]FunctionDecl
