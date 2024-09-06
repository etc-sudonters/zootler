package mutoh

import (
	"sudonters/zootler/carpenters/ichiro"
	"sudonters/zootler/carpenters/jiro"
	"sudonters/zootler/carpenters/saburo"
	"sudonters/zootler/carpenters/shiro"
	"sudonters/zootler/internal/app"
)

type Bootstrapper struct{}

func (b *Bootstrapper) Setup(_ *app.Zootlr) error {
	ichiro.LoadDataTables()
	jiro.LoadWorldGraph()
	saburo.CompileAllRules()
	shiro.CreateWorldTemplate()
	return nil
}
