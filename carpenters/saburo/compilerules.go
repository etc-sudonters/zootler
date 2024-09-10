package saburo

import (
	"slices"
	"strings"
	"sudonters/zootler/icearrow/ast"
	"sudonters/zootler/icearrow/debug"
	"sudonters/zootler/icearrow/macros"
	parsing "sudonters/zootler/icearrow/parser"
	"sudonters/zootler/icearrow/zasm"
	"sudonters/zootler/internal"
	"sudonters/zootler/internal/app"
	"sudonters/zootler/internal/entities"

	"github.com/etc-sudonters/substrate/dontio"
	"github.com/etc-sudonters/substrate/peruse"
	"github.com/etc-sudonters/substrate/slipup"
)

type RuleCompilation struct {
	ScriptPath string
}

func (rc RuleCompilation) Setup(z *app.Zootlr) error {
	assembler := rc.createAssembler()
	edges := app.GetResource[entities.Edges](z)

	collected := slices.Collect(edges.Res.All)
	slices.SortFunc(collected, func(a, b entities.Edge) int {
		return strings.Compare(string(a.Name()), string(b.Name()))
	})

	var edge entities.Edge

	whileHandlingRule := func(err error, action string) error {
		return slipup.Describef(err, "while %s rule %q", action, edge.GetRawRule())
	}
	xpndr := macros.NewMacroExpansions()
	loadAllMacros(&xpndr, rc.ScriptPath)
	grammar := parsing.NewRulesGrammar()

	for _, edge = range collected {
		dontio.WriteLineOut(z.Ctx(), string(edge.Name()))
		dontio.WriteLineOut(z.Ctx(), string(edge.GetRawRule()))
		tokens := macros.ExpandWith(&xpndr, parsing.NewRulesLexer(string(edge.GetRawRule())))
		parser := peruse.NewParser(&grammar, tokens)
		pt, ptErr := parser.ParseAt(parsing.LOWEST)
		if ptErr != nil {
			return whileHandlingRule(ptErr, "parsing")
		}
		ast, astErr := parsing.Transform(&ast.Ast{}, pt)
		if astErr != nil {
			return whileHandlingRule(astErr, "lowering tree")
		}
		dontio.WriteLineOut(z.Ctx(), debug.AstSexpr(ast))
		asm, asmErr := assembler.Assemble(ast)
		if asmErr != nil {
			return whileHandlingRule(asmErr, "assembling")
		}
		dis := zasm.Disassemble(asm.I)
		dontio.WriteLineOut(z.Ctx(), dis)
	}

	return nil
}

func (rc RuleCompilation) createAssembler() zasm.Assembler {
	return zasm.Assembler{
		Data: zasm.NewDataBuilder(),
	}
}

func macro(rule string) []peruse.Token {
	return slices.Collect(peruse.AllTokens(parsing.NewRulesLexer(rule)))
}

func loadAllMacros(xpndrs *macros.Expansions, path string) error {
	all, fileErr := internal.ReadJsonFileStringMap(path)
	if fileErr != nil {
		return fileErr
	}

	for decl, body := range all {
		if err := macros.CreateScriptedMacro(xpndrs, decl, body); err != nil {
			return slipup.Describef(err, "while loading macro %s", decl)
		}
	}
	return nil
}
