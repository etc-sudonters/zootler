package ast

import (
	"sudonters/zootler/internal/rules/parser"

	"github.com/etc-sudonters/substrate/slipup"
)

// At this level macros carry semantic meaning only. We identify all
// invocations of macros via MacroDecls. Parsing AST requires this table to
// carry at least all declarations. Builtin macros are not required to become
// initialized and all uninitialized macros are considered built in.
type MacroDecls struct {
	decls map[string]MacroDecl
}

func NewMacros(capacity int) MacroDecls {
	return MacroDecls{
		decls: make(map[string]MacroDecl, capacity),
	}
}

type MacroDecl struct {
	Name   string
	Params []string
	Body   Node
}

func (md MacroDecl) IsBuiltIn() bool {
	_, ok := md.Body.(*Empty)
	return ok
}

func (md MacroDecls) DeclareBuiltIn(name string, params []string) {
	decl := MacroDecl{
		Name:   name,
		Params: params,
		Body:   &Empty{},
	}
	md.decls[name] = decl
}

func (md MacroDecls) Declare(name string, params []string) MacroDecl {
	decl := MacroDecl{
		Name:   name,
		Params: params,
	}
	md.decls[name] = decl
	return decl
}

func (md MacroDecls) Get(name string) (MacroDecl, bool) {
	macro, exists := md.decls[name]
	return macro, exists
}

func (md MacroDecls) Exists(name string) bool {
	_, ok := md.decls[name]
	return ok
}

func (md MacroDecls) Initialize(name string, body Node) error {
	macro, exists := md.decls[name]
	if !exists {
		return slipup.Createf("macro not declared: '%s'", name)
	}

	if macro.Body != nil {
		return slipup.Createf("macro already Initialized: '%s'", name)
	}

	macro.Body = body
	md.decls[name] = macro
	return nil
}

func DeclareFromParseTree(decls MacroDecls, decl *parser.Call) MacroDecl {
	var params []string
	name := parser.MustAssertAs[*parser.Identifier](decl.Callee).Value

	if len(decl.Args) != 0 {
		params = make([]string, len(decl.Args))
		for i, arg := range decl.Args {
			params[i] = parser.MustAssertAs[*parser.Identifier](arg).Value
		}
	}

	return decls.Declare(name, params)
}
