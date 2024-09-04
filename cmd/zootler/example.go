package main

/*

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"regexp"
	"sudonters/zootler/internal"
	"sudonters/zootler/internal/app"
	"sudonters/zootler/internal/rules/ast"
	"sudonters/zootler/internal/rules/parser"
	"sudonters/zootler/internal/rules/visitor"
	"sudonters/zootler/internal/rules/zasm"

	"github.com/etc-sudonters/substrate/dontio"
	"github.com/etc-sudonters/substrate/slipup"
)



type MacroExpander struct {
	Macros        ast.MacroDecls
	CurrentOrigin string
	expansions    map[string]int
}

func (x *MacroExpander) Expand(node *ast.Call, scope zasm.AssemblerScope) (zasm.Instructions, error) {
	decl, _ := x.Macros.Get(node.Callee)
	if !decl.IsBuiltIn() {
		return x.expandScript(node, decl, scope)
	}
	return x.expandBuiltIn(node, decl, scope)
}

func (x *MacroExpander) expandBuiltIn(call *ast.Call, _ ast.MacroDecl, scope zasm.AssemblerScope) (zasm.Instructions, error) {
	if call.Callee != "at" && call.Callee != "here" {
		return zasm.Tape().WriteLoadBool(true).Instructions(), nil
	}

	if x.expansions == nil {
		x.expansions = map[string]int{}
	}

	var target string
	if call.Callee == "at" {
		target = ast.MustAssertAs[*ast.Identifier](call.Args[0]).Name
	} else {
		if x.CurrentOrigin == "" {
			panic(slipup.Createf("here invoked w/o origin set on expander: %s", call))
		}
		target = x.CurrentOrigin
	}

	name := fmt.Sprintf("<CANREACH>%s[%d]<%s>", target, x.expansions[target], x.CurrentOrigin)
	x.expansions[target] += 1
	ident := scope.Assembler.Data.Registers.Intern(name)
	return zasm.Tape().WriteLoadIdent(uint32(ident)).WriteLoadU24(1).WriteOp(zasm.OP_CHK_QTY).Instructions(), nil
}

func (x *MacroExpander) expandScript(call *ast.Call, decl ast.MacroDecl, scope zasm.AssemblerScope) (zasm.Instructions, error) {
	argErr := scope.InParentScope(func() error {
		for idx, param := range decl.Params {
			scratch := zasm.Scratch()
			if err := scope.Assembler.AssembleInto(&scratch, call.Args[idx]); err != nil {
				return slipup.Describef(err, "macro arg %d", idx)
			}
			scope.Replacements[param] = scratch.Instr
		}
		return nil
	})
	if argErr != nil {
		return nil, argErr
	}

	replacement := zasm.Scratch()
	if err := scope.Assembler.AssembleInto(&replacement, decl.Body); err != nil {
		return nil, err
	}
	return replacement.Instr, nil
}


func Setup(z *app.Zootlr) error {
	ctx := z.Ctx()
	macros, bodies := load_json_macros(filepath.Join("", "..", "helpers.json"))
	sni := SpecialNameIdentifier{macros, make(map[string]struct{})}
	xpander := MacroExpander{
		Macros: macros,
	}
	assembly := zasm.NewAssembly()
	assembler := zasm.Assembler{
		Data:   zasm.NewDataBuilder(),
		Macros: &xpander,
		Functions: map[string]int{
			"has":                    2,
			"has_all_notes_for_song": 1,
			"has_hearts":             1,
			"has_medallions":         1,
			"has_stones":             1,
			"load_setting":           1,
			"load_setting_2":         2,
			"load_trick":             1,
			"region_has_shortcuts":   1,
		},
	}

	assembler.DebugOutput = func(tpl string, v ...any) {
		dontio.WriteLineOut(ctx, tpl, v...)
	}

	macros.DeclareBuiltIn("at", []string{"target", "rule"})
	macros.DeclareBuiltIn("at_dampe_time", nil)
	macros.DeclareBuiltIn("at_day", nil)
	macros.DeclareBuiltIn("at_night", nil)
	macros.DeclareBuiltIn("had_night_start", nil)
	macros.DeclareBuiltIn("here", []string{"rule"})

	sni.DeclareSetting("shuffle_individual_ocarina_notes")

	astGen := ast.NewGenerator(&sni)
	initialize_script_macros(astGen, macros, bodies)

	filepath.Walk("", func(path string, info fs.FileInfo, err error) error {
		dontio.WriteLineOut(ctx, "at logic dir file %s", path)
		ext := filepath.Ext(path)
		if ext != ".json" {
			dontio.WriteLineOut(ctx, "skipping non json file: %s", path)
			return nil
		}

		locs, _ := internal.ReadJsonFileAs[[]fileloc](path)

		for edge := range all_rules_from(locs) {
			astGen.ClearStacks()
			xpander.CurrentOrigin = edge.Origin
			dontio.WriteLineOut(ctx, "%s -> %s: %s", edge.Origin, edge.Dest, edge.Rule)
			pt, parseErr := parser.Parse(edge.Rule)
			if parseErr != nil {
				return slipup.Describe(parseErr, "while parsing")
			}

			genAst, astErr := visitor.Transform(astGen, pt)
			if astErr != nil {
				return slipup.Describe(astErr, "pt -> ast on rule")
			}

			repr := ast.AstRender{}
			ast.Visit(&repr, genAst)
			dontio.WriteLineOut(ctx, repr.String())

			asmBlock, blockErr := assembly.Block(fmt.Sprintf("EDGE<%s>%s->%s", edge.Kind, edge.Origin, edge.Dest))
			if blockErr != nil {
				panic(blockErr)
			}
			if err := assembler.AssembleInto(&asmBlock, genAst); err != nil {
				panic(err)
			}

			dis := zasm.ZasmDisassembler{}
			dontio.WriteLineOut(ctx, dis.Disassemble(asmBlock.Instr))
		}

		return nil
	})

	assembly.Load(assembler.Data)
	return nil
}
*/
