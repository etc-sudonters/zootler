package shiro

import (
	"sudonters/zootler/icearrow/compiler"
	"sudonters/zootler/icearrow/zasm"
	"sudonters/zootler/internal/app"
	"sudonters/zootler/internal/settings"

	"github.com/etc-sudonters/substrate/dontio"
)

type WorldCompiler struct{}

func (wc *WorldCompiler) Setup(z *app.Zootlr) error {
	res := app.GetResource[zasm.Assembly](z)
	assembly := &res.Res
	symbols := compiler.CreateSymbolTable(assembly.Data())
	comp := compiler.RuleCompiler{
		Symbols: &symbols,
	}
	for unit := range assembly.Units {
		ct, unassembleErr := compiler.Unassemble(unit)
		if unassembleErr != nil {
			panic(unassembleErr)
		}
		ct = compiler.LastMileOptimizations(ct, &symbols)
		tape := comp.Compile(ct)
		dontio.WriteLineOut(z.Ctx(), "%q\n%s", unit.Name, compiler.ReadTape(tape))
	}

	return nil
}

type ZootrSettingResolver struct {
	ST       *compiler.SymbolTable
	Settings settings.ZootrSettings
}

// NOTE What's passed won't always be a dungeon, at least in one case
// we're passed a quasi location that we have to resolve specially
func (zsr *ZootrSettingResolver) DungeonHasShortcuts(string) bool {
	return false
}

func (zsr *ZootrSettingResolver) IsTrialSkipped(string) bool {
	return true
}

func (zsr *ZootrSettingResolver) LoadSetting(string) bool {
	return false
}

func (zsr *ZootrSettingResolver) LoadSetting2(string, string) bool {
	return false
}

// For TOD tracking, we either don't care or need to invoke a builtin
// that will trace TOD access.
func (zsr *ZootrSettingResolver) AtTimeOfDay(string) compiler.CompileTree {
	// TODO handle detecting and outputting tracing cases
	return compiler.Immediate{
		Value: true,
		Kind:  compiler.CT_IMMED_TRUE,
	}
}

func (zsr *ZootrSettingResolver) CheckStartCond(string) bool {
	return false
}

func (zsr *ZootrSettingResolver) CompareToSetting(string, compiler.CompileTree) bool {
	return false
}
