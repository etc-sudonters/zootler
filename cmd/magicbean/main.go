package main

import (
	"os"
	"runtime/pprof"
	"sudonters/zootler/cmd/magicbean/bootstrap"
	"sudonters/zootler/internal/query"
	"sudonters/zootler/internal/settings"
	"sudonters/zootler/internal/skelly/bitset32"
	"sudonters/zootler/magicbean"
	"sudonters/zootler/mido"
	"sudonters/zootler/mido/objects"
	"sudonters/zootler/zecs"
)

func consttrue(*objects.Table, []objects.Object) (objects.Object, error) {
	return objects.PackedTrue, nil
}

func main() {
	paths := bootstrap.LoadPaths{
		Tokens:     ".data/data/items.json",
		Placements: ".data/data/locations.json",
		Scripts:    ".data/logic/helpers.json",
		Relations:  ".data/logic/glitchless/",
	}
	theseSettings := settings.Default()

	stopProfiling := profileto("zootr.prof")
	defer stopProfiling()

	artifacts := setup(paths, &theseSettings)
	eng := artifacts.Ocm.Engine()
	tbl, err := query.ExtractTable(eng)
	bootstrap.PanicWhenErr(err)
	displaycolstats(tbl)

	{
		q := artifacts.Ocm.Query()
		q.Build(zecs.With[magicbean.Region], zecs.Load[magicbean.Name])
		roots := bitset32.Bitset{}

		for ent, tup := range q.Rows {
			name := tup.Values[0].(magicbean.Name)
			if name == "Root" {
				bitset32.Set(&roots, ent)
				break
			}
		}

		for range 10 {
			artifacts.World.ExploreAvailableEdges(magicbean.Exploration{
				Workset: bitset32.Copy(roots),
				Visited: bitset32.Bitset{},
				VM: mido.VM{
					Objects: &artifacts.Objects,
					Funcs: objects.BuiltInFunctions{
						consttrue,
						consttrue,
						consttrue,
						consttrue,
						consttrue,
						consttrue,
						consttrue,
						consttrue,
						consttrue,
						consttrue,
						consttrue,
					},
				},
			})
		}
	}
}

type Artifacts struct {
	Ocm     zecs.Ocm
	World   magicbean.ExplorableWorld
	Objects objects.Table
}

func setup(paths bootstrap.LoadPaths, settings *settings.Zootr) (artifacts Artifacts) {
	ocm := bootstrap.Phase1_InitializeStorage(nil)
	bootstrap.PanicWhenErr(bootstrap.Phase2_ImportFromFiles(&ocm, paths))

	compileEnv := bootstrap.Phase3_ConfigureCompiler(&ocm, settings)

	codegen := mido.Compiler(&compileEnv)

	bootstrap.PanicWhenErr(bootstrap.Phase4_Compile(
		&ocm, &codegen,
	))

	world := bootstrap.Phase5_CreateWorld(&ocm, settings, objects.TableFrom(compileEnv.Objects))

	artifacts.Ocm = ocm
	artifacts.World = world
	artifacts.Objects = objects.TableFrom(compileEnv.Objects)

	return artifacts
}

func noop() {}

func profileto(path string) func() {
	if path == "" {
		return noop
	}

	f, err := os.Create(path)
	bootstrap.PanicWhenErr(err)
	pprof.StartCPUProfile(f)
	return func() { pprof.StopCPUProfile() }
}
