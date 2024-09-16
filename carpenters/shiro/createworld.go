package shiro

import (
	"sudonters/zootler/icearrow/compiler"
	"sudonters/zootler/internal/app"
	"sudonters/zootler/internal/settings"
)

type WorldCompiler struct{}

func (wc *WorldCompiler) Setup(z *app.Zootlr) error {
	return nil
}

type ZootrSettingResolver struct {
	ST       *compiler.SymbolTable
	Settings settings.ZootrSettings
}

func (res *ZootrSettingResolver) Resolve(string) (sym compiler.Symbol) {
	return
}

func (res *ZootrSettingResolver) ResolveNested(string, string) (sym compiler.Symbol) {
	return
}

func (res *ZootrSettingResolver) ResolveTrick(string) (sym compiler.Symbol) {
	return
}
