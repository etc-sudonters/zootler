package shiro

import (
	"sudonters/zootler/icearrow/compiler"
	"sudonters/zootler/icearrow/zasm"
	"sudonters/zootler/internal/app"
	"sudonters/zootler/internal/settings"

	"github.com/etc-sudonters/substrate/slipup"
)

func ReadDataIntoSymbolTable(data *zasm.Data) compiler.SymbolTable {
	return compiler.CreateSymbolTable(data)
}

func constTrue(compiler.Invocation, *compiler.Symbol, *compiler.SymbolTable) (compiler.CompileTree, error) {
	return compiler.Immediate{Value: true, Kind: compiler.CT_IMMED_TRUE}, nil

}

func createCompilerPlugin(st *compiler.SymbolTable, settings *intrinsics) compiler.Intrinsics {
	intrinsics := compiler.NewIntrinsics()

	nameMap := []struct {
		Name string
		Func settingResolver
	}{
		{"compareeqsetting", settings.CompareEq},
		{"compareltsetting", settings.CompareLt},
		{"comparenqsetting", settings.CompareNq},
		{"regionhasshortcuts", settings.HasShortcuts},
		{"hasdungeonshortcuts", settings.HasShortcuts},
		{"inverthasdungeonshortcuts", settings.InvertHasShortcuts},
		{"invertloadsetting", settings.LoadSetting},
		{"istrialskipped", settings.IsTrialSkipped},
		{"istrickenabled", settings.IsTrickEnabled},
		{"loadsetting", settings.LoadSetting},
	}

	for _, pair := range nameMap {
		sym := st.Named(pair.Name)
		if sym == nil {
			panic(slipup.Createf("could not find symbol for %q", pair.Name))
		}
		intrinsics.Add(sym, intoIntrinsicFunc(pair.Func))
	}

	for _, always := range []string{"hasallnotesforsong", "atday", "atnight", "atdampetime"} {
		intrinsics.Add(st.Named(always), constTrue)
	}

	return intrinsics
}

type WorldCompiler struct {
}

func (wc *WorldCompiler) Setup(z *app.Zootlr) error {
	prog := app.GetResource[zasm.Assembly](z)
	settings := app.GetResource[settings.ZootrSettings](z)

	assembly := &prog.Res
	compiled := CompiledWorldRules{
		Rules:   make(map[string]compiler.Tape, assembly.NumberOfUnits()),
		Symbols: ReadDataIntoSymbolTable(assembly.Data()),
	}
	settingIntrinsics := intoIntrinsics(&settings.Res)
	intrinsics := createCompilerPlugin(&compiled.Symbols, &settingIntrinsics)
	comp := compiler.RuleCompiler{
		Symbols: &compiled.Symbols,
	}
	lastMile := compiler.LastMileOptimizations(&compiled.Symbols, &intrinsics)

	for unit := range assembly.Units {
		ct, unassembleErr := compiler.Unassemble(unit, &compiled.Symbols)
		if unassembleErr != nil {
			panic(unassembleErr)
		}
		ct = lastMile(ct)
		tapeWriter := comp.Compile(ct)
		compiled.Rules[unit.Name] = tapeWriter.Tape
	}
	app.AddResource(z, compiled)
	return nil
}

type CompiledWorldRules struct {
	Rules   map[string]compiler.Tape
	Symbols compiler.SymbolTable
}
