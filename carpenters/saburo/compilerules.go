package saburo

import (
	"sudonters/zootler/icearrow/ast"
	parsing "sudonters/zootler/icearrow/parser"
	"sudonters/zootler/icearrow/zasm"
	"sudonters/zootler/internal/app"
	"sudonters/zootler/internal/entities"

	"github.com/etc-sudonters/substrate/dontio"
)

func paniconerr(e error) {
	if e != nil {
		panic(e)
	}
}

type RuleCompilation struct {
	ScriptPath string
}

func (rc RuleCompilation) Setup(z *app.Zootlr) error {
	macros := parsing.DefaultCoven()
	rc.loadMacros(parsing.InitiateCoven(&macros))
	parser := rc.createParser(macros)
	return rc.assembleAllRules(z, parser)
}

func (rc RuleCompilation) assembleAllRules(z *app.Zootlr, rp parsing.RulesParser) error {
	assembler := rc.createAssembler()
	edges := app.GetResource[entities.Edges](z)

	for edge := range edges.Res.All {
		pt, ptErr := rp.ParseString(string(edge.GetRawRule()))
		paniconerr(ptErr)
		ast, astErr := parsing.Transform(&ast.Ast{}, pt)
		paniconerr(astErr)
		asm, asmErr := assembler.Assemble(ast)
		paniconerr(asmErr)

		dis := zasm.Disassemble(asm.I)
		dontio.WriteLineOut(z.Ctx(), string(edge.Name()))
		dontio.WriteLineOut(z.Ctx(), string(edge.GetRawRule()))
		dontio.WriteLineOut(z.Ctx(), dis)
	}

	return nil
}

func (rc RuleCompilation) createParser(mc parsing.MacroCoven) parsing.RulesParser {
	return parsing.NewRulesParser(mc)
}

func (rc RuleCompilation) createAssembler() zasm.Assembler {
	return zasm.Assembler{
		Data: zasm.NewDataBuilder(),
	}
}

func (rc RuleCompilation) loadMacros(mb parsing.MacroBuilder) error {
	return LoadScriptedMacros(mb, rc.ScriptPath)
}
