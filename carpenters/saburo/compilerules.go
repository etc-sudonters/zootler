package saburo

import (
	"slices"
	"strings"
	"sudonters/zootler/icearrow/ast"
	parsing "sudonters/zootler/icearrow/parser"
	"sudonters/zootler/icearrow/zasm"
	"sudonters/zootler/internal/app"
	"sudonters/zootler/internal/entities"

	"github.com/etc-sudonters/substrate/dontio"
	"github.com/etc-sudonters/substrate/slipup"
)

type RuleCompilation struct {
	ScriptPath string
}

func (rc RuleCompilation) Setup(z *app.Zootlr) error {
	macros := parsing.DefaultCoven()
	rc.loadMacros(parsing.InitiateCoven(&macros))
	parser := parsing.NewParserStack(macros, parsing.MACROS_DISABLE)
	return rc.assembleAllRules(z, parser)
}

func (rc RuleCompilation) assembleAllRules(z *app.Zootlr, rp *parsing.ParserStack) error {
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

	for _, edge = range collected {
		pt, ptErr := rp.ParseString(string(edge.GetRawRule()))
		if ptErr != nil {
			return whileHandlingRule(ptErr, "parsing")

		}
		ast, astErr := parsing.Transform(&ast.Ast{}, pt)
		if astErr != nil {
			return whileHandlingRule(astErr, "lowering tree")

		}
		asm, asmErr := assembler.Assemble(ast)
		if asmErr != nil {
			return whileHandlingRule(asmErr, "assembling")
		}

		dis := zasm.Disassemble(asm.I)
		dontio.WriteLineOut(z.Ctx(), string(edge.Name()))
		dontio.WriteLineOut(z.Ctx(), string(edge.GetRawRule()))
		dontio.WriteLineOut(z.Ctx(), dis)
	}

	return nil
}

func (rc RuleCompilation) createAssembler() zasm.Assembler {
	return zasm.Assembler{
		Data: zasm.NewDataBuilder(),
	}
}

func (rc RuleCompilation) loadMacros(mb parsing.MacroBuilder) error {
	return LoadScriptedMacros(mb, rc.ScriptPath)
}
