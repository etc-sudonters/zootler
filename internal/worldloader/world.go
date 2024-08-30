package worldloader

import (
	"iter"
	"sudonters/zootler/internal/app"
)

type World struct {
	HelperPath, WorldPath string
}

func (w World) Setup(z *app.Zootlr) error {
	ctx, eng := z.Ctx(), z.Engine()

	logic := new(LogicLoader)
	if err := logic.Init(eng); err != nil {
		return nil
	}

}

func (w World) allLocations() iter.Seq[LogicLocation] {

}
