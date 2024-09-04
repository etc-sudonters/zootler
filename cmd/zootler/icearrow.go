package main

import (
	"iter"
	"sudonters/zootler/icearrow/zasm"
	"sudonters/zootler/internal/app"
	"sudonters/zootler/internal/rules/ast"

	"github.com/etc-sudonters/substrate/dontio"
)

type IceArrowRuntime struct{}

type IceArrowEdge struct {
	AstEdge
	Instructions zasm.Instructions
}

func (iar *IceArrowRuntime) Setup(z *app.Zootlr) error {
	macros := app.GetResource[ast.MacroDecls](z)
	ass := zasm.Assembler{
		Data:   zasm.NewDataBuilder(),
		Macros: &xpander{macros: macros.Res},
	}

	for edge := range iar.assembleAllEdge(&ass, app.GetResource[AstAllRuleEdges](z).Res.Iter) {
		dis := zasm.Disassemble(edge.Instructions)
		dontio.WriteLineOut(z.Ctx(), "%s -> %s", edge.Origin, edge.Dest)
		dontio.WriteLineOut(z.Ctx(), edge.Rule)
		dontio.WriteLineOut(z.Ctx(), dis)
	}

	return nil
}

func (iar *IceArrowRuntime) assembleAllEdge(assembler *zasm.Assembler, edges iter.Seq[AstEdge]) iter.Seq[IceArrowEdge] {
	return func(yield func(IceArrowEdge) bool) {
		for edge := range edges {
			asm, err := assembler.Assemble(edge.Ast)
			if err != nil {
				panic(err)
			}
			if !yield(IceArrowEdge{edge, asm.I}) {
				return
			}
		}
	}
}

type xpander struct {
	macros  ast.MacroDecls
	current string
}

func (m *xpander) ExpandMacro(assembler *zasm.Assembler, call *ast.Call) zasm.Instructions {
	decl, exists := m.macros.Get(call.Callee)
	if exists && !decl.IsBuiltIn() {
		return m.expandScriptMacro(assembler, call, decl)
	}

	return nil
}

func (m *xpander) expandScriptMacro(assembler *zasm.Assembler, call *ast.Call, decl ast.MacroDecl) zasm.Instructions {
	renames := make(map[string]string, 2)
	for i := range len(call.Args) {
		renames[decl.Params[i]] = ast.MustAssertAs[*ast.Identifier](call.Args[i]).Name
	}

	body := ast.Map(decl.Body, func(node ast.Node) ast.Node {
		id, ok := node.(*ast.Identifier)
		if ok {
			if rename, has := renames[id.Name]; has {
				id.Name = rename
			}
		}
		return node
	})

	assembly, errs := assembler.Assemble(body)
	if errs != nil {
		panic(errs)
	}
	return assembly.I
}
