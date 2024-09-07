package jiro

import (
	"io/fs"
	"path/filepath"
	"sudonters/zootler/internal"
	"sudonters/zootler/internal/app"
	"sudonters/zootler/internal/entities"

	"github.com/etc-sudonters/substrate/skelly/graph"
	"github.com/etc-sudonters/substrate/slipup"
)

type WorldGraph struct {
	LogicDir string
}

func (wg WorldGraph) Setup(z *app.Zootlr) error {
	var state loadstate
	state.locs = app.GetResource[entities.Locations](z).Res
	state.edge = app.GetResource[entities.Edges](z).Res
	state.grph = graph.Builder{graph.New()}

	if err := wg.loaddir(&state); err != nil {
		return slipup.Describef(err, "while loading dir '%s'", wg.LogicDir)
	}

	z.AddResource(state.grph)

	return nil
}

func (wg WorldGraph) loaddir(state *loadstate) error {
	return filepath.WalkDir(wg.LogicDir, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return slipup.Describe(err, "logic directory walk called with err")
		}

		info, err := entry.Info()
		if err != nil || info.Mode() != (^fs.ModeType)&info.Mode() {
			// either we couldn't get the info, which doesn't bode well
			// or it's some kind of not file thing which we also don't want
			return nil
		}

		if ext := filepath.Ext(path); ext != ".json" {
			return nil
		}

		nodes, readErr := internal.ReadJsonFileAs[[]worldnode](path)
		if readErr != nil {
			return slipup.Describef(readErr, "while reading file '%s'", path)
		}

		for _, node := range nodes {
			if err := state.load(node); err != nil {
				return slipup.Describef(err, "while handling node '%s' in file '%s'", node.RegionName, path)
			}
		}

		return nil
	})
}
