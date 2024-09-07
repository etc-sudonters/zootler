package saburo

import (
	"iter"
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
	return rc.assembleAllRules(z)
}

func (rc RuleCompilation) assembleAllRules(z *app.Zootlr) error {
	parser := parsing.NewBetterRulesParser()
	assembler := rc.createAssembler()

	for edge := range rc.edges(z) {
		pt, ptErr := parser.ParseString(string(edge.GetRawRule()))
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

func (rc RuleCompilation) createParser() parsing.RulesParser {
	return parsing.NewBetterRulesParser()
}

func (rc RuleCompilation) createAssembler() zasm.Assembler {
	return zasm.Assembler{
		Data: zasm.NewDataBuilder(),
	}
}

func (rc RuleCompilation) edges(z *app.Zootlr) iter.Seq[entities.Edge] {
	edges := app.GetResource[entities.Edges](z)
	return func(yield func(entities.Edge) bool) {
		for edge := range edges.Res.All {
			if !yield(edge) {
				return
			}
		}
	}
}
