package shiro

import (
	"runtime"
	"sudonters/zootler/icearrow/compiler"
	"sudonters/zootler/icearrow/zasm"
	"sudonters/zootler/internal/app"
	"sudonters/zootler/internal/settings"

	"github.com/etc-sudonters/substrate/dontio"
	"github.com/etc-sudonters/substrate/slipup"
)

func ReadDataIntoSymbolTable(data *zasm.Data) compiler.SymbolTable {
	return compiler.CreateSymbolTable(data)
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

	return intrinsics
}

type WorldCompiler struct {
}

func (wc *WorldCompiler) Setup(z *app.Zootlr) error {
	prog := app.GetResource[zasm.Assembly](z)
	settings := app.GetResource[settings.ZootrSettings](z)

	assembly := &prog.Res
	symbols := ReadDataIntoSymbolTable(assembly.Data())
	settingIntrinsics := intoIntrinsics(&settings.Res)
	intrinsics := createCompilerPlugin(&symbols, &settingIntrinsics)
	comp := compiler.RuleCompiler{
		Symbols: &symbols,
	}
	lastMile := compiler.LastMileOptimizations(&symbols, &intrinsics)

	for unit := range assembly.Units {
		dontio.WriteLineOut(z.Ctx(), unit.Name)
		ct, unassembleErr := compiler.Unassemble(unit, &symbols)
		if unassembleErr != nil {
			panic(unassembleErr)
		}
		ct = lastMile(ct)
		tape := comp.Compile(ct)
		runtime.KeepAlive(tape)
	}

	return nil
}
