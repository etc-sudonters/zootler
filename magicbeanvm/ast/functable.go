package ast

import (
	"fmt"
	"sudonters/zootler/internal/ruleparser"
	"sudonters/zootler/magicbeanvm/symbols"

	"github.com/etc-sudonters/substrate/peruse"
)

func BuildFunctionTable(tbl *symbols.Table, grammar peruse.Grammar[ruleparser.Tree], decls map[string]string) (FunctionTable, error) {
	funcs := make(FunctionTable, len(decls))
	bodies := make(map[Identifier]string, len(decls))

	for header, body := range decls {
		var decl FunctionDecl
		head, declErr := Parse(header, tbl, grammar)
		if declErr != nil {
			panic(declErr)
		}

		switch head := head.(type) {
		case Invoke:
			decl.Symbol = LookUpNodeInTable(tbl, head.Target)
			if decl.Symbol.Type == symbols.BUILT_IN {
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
			decl.Symbol = tbl.LookUpByIndex(head.AsIndex())
			if decl.Symbol.Type == symbols.BUILT_IN {
				continue
			}
			decl.Params = nil
		default:
			panic(fmt.Errorf("expected identifier or invoke style declaration, got %#v", head))
		}

		decl.Symbol.SetType(symbols.FUNCTION)
		id := IdentifierFrom(decl.Symbol)
		bodies[id] = body
		funcs[id] = decl
	}

	for id, decl := range funcs {
		body := bodies[id]
		nodes, bodyErr := Parse(body, tbl, grammar)
		if bodyErr != nil {
			panic(bodyErr)
		}
		decl.Body = nodes
		funcs[id] = decl
	}

	return funcs, nil
}

type FunctionDecl struct {
	Symbol *symbols.Sym
	Params []Identifier
	Body   Node
}

type FunctionTable map[Identifier]FunctionDecl
