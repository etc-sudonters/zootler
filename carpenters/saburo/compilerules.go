package saburo

import (
	"sudonters/zootler/icearrow/ast"
	parsing "sudonters/zootler/icearrow/parser"
	"sudonters/zootler/icearrow/zasm"
	"sudonters/zootler/internal/app"
	"sudonters/zootler/internal/entities"

	"github.com/etc-sudonters/substrate/dontio"
)

type RuleCompilation struct {
	ScriptPath string
}

func (rc RuleCompilation) Setup(z *app.Zootlr) error {
	macros := parsing.DefaultCoven()
	rc.loadMacros(parsing.InitiateCoven(&macros))
	parser := parsing.NewParserStack(macros)
	return rc.assembleAllRules(z, parser)
}

func (rc RuleCompilation) assembleAllRules(z *app.Zootlr, rp *parsing.ParserStack) error {
	assembler := rc.createAssembler()
	edges := app.GetResource[entities.Edges](z)

	for edge := range edges.Res.All {
		pt, ptErr := rp.ParseString(string(edge.GetRawRule()))
		if ptErr != nil {
			return ptErr

		}
		ast, astErr := parsing.Transform(&ast.Ast{}, pt)
		if astErr != nil {
			return astErr

		}
		asm, asmErr := assembler.Assemble(ast)
		if asmErr != nil {
			return asmErr
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
