package main

import (
	"iter"
	"os"
	"path"
	"strings"
	"sudonters/zootler/internal"
	"sudonters/zootler/internal/app"
	"sudonters/zootler/internal/components"
	"sudonters/zootler/internal/query"
	"sudonters/zootler/internal/rules/parser"
	"sudonters/zootler/internal/rules/preprocessor"
	"sudonters/zootler/internal/table"

	"github.com/etc-sudonters/substrate/skelly/graph"
	"github.com/etc-sudonters/substrate/slipup"
)

type Loader struct {
	Path, Helpers string
}

type loadingstate struct {
	pre       preprocessor.P
	graph     graph.Builder
	locations locations
}

type locations struct {
	eng  query.Engine
	locs map[internal.NormalizedStr]*locbuilder
}

type locbuilder struct {
	parent *locations
	id     table.RowId
	name   components.Name
	rules  map[string]struct {
		raw        string
		components []table.Value
	}
}

type edge struct {
	Components []table.Value
	Rule       string
	Name       components.Name
}

type filelocation struct {
	BossRoom      bool              `json:"is_boss_room"`
	DoesTimePass  bool              `json:"time_passes"`
	Dungeon       string            `json:"dungeon"`
	HintRegion    string            `json:"hint"`     // hint zone
	HintRegionAlt string            `json:"alt_hint"` // more specific location -- restricted to gannon's chambers?
	Name          string            `json:"region_name"`
	Savewarp      string            `json:"savewarp"` // effectively another exit
	Scene         string            `json:"scene"`    // similar to hint?
	Events        map[string]string `json:"events"`
	Exits         map[string]string `json:"exits"`
	Locations     map[string]string `json:"locations"`
}

func (l *Loader) Setup(z *app.Zootlr) error {
	state, stateErr := l.createState(z, l.Helpers)
	if stateErr != nil {
		return slipup.Describe(stateErr, "while building state")
	}
	for loaded := range state.LoadDirectory(l.Path) {
		location, err := state.locations.Build(loaded)
		if err != nil {
			panic(err)
		}

		for edge := range location.Edges {
			firstPassAst, firstPassErr := state.parse(edge.Rule)
			if firstPassErr != nil {
				panic(firstPassErr)
			}
			secondPassAst, secondPassErr := state.optimize(firstPassAst)
			if secondPassErr != nil {
				panic(secondPassErr)
			}
			compiled, compileErr := state.compile(secondPassAst)
			if compileErr != nil {
				panic(compileErr)
			}

			location.LinkVia(edge, compiled)
		}
	}

	return nil
}

func (l *Loader) createState(z *app.Zootlr, helpersPath string) (*loadingstate, error) {
	return nil, nil
}

func (s *loadingstate) LoadDirectory(logicDir string) iter.Seq[filelocation] {
	directory, dirErr := os.ReadDir(logicDir)
	if dirErr != nil {
		panic(dirErr)
	}

	readLocations := func(entry os.DirEntry) []filelocation {
		locations, err := internal.ReadJsonFileAs[[]filelocation](path.Join(logicDir, entry.Name()))
		if err != nil {
			panic(slipup.Describef(err, "while reading %s/%s", logicDir, entry.Name()))
		}

		return locations
	}

	return func(yield func(filelocation) bool) {
		for _, entry := range directory {
			if !internal.IsFile(entry) {
				continue
			}

			if strings.Contains(entry.Name(), "mq") {
				continue
			}

			for _, loc := range readLocations(entry) {
				if !yield(loc) {
					return
				}
			}
		}
	}
}

func (s *loadingstate) parse(rule string) (parser.Expression, error) {
	return nil, nil
}

func (s *loadingstate) optimize(ast parser.Expression) (parser.Expression, error) {
	return nil, nil
}

func (s *loadingstate) compile(ast parser.Expression) (parser.Expression, error) {
	return nil, nil
}

func (l *locations) Build() {}
