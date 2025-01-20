package main

import (
	"os"
	"runtime/pprof"
	"sudonters/zootler/cmd/magicbean/bootstrap"
	"sudonters/zootler/internal/query"
	"sudonters/zootler/internal/settings"
	"sudonters/zootler/magicbean"
	"sudonters/zootler/mido"
	"sudonters/zootler/mido/objects"
	"sudonters/zootler/zecs"
)

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

	ocm, world := setup(paths, &theseSettings)
	_, _ = ocm, world
	eng := ocm.Engine()
	tbl, err := query.ExtractTable(eng)
	bootstrap.PanicWhenErr(err)
	displaycolstats(tbl)
}

func setup(paths bootstrap.LoadPaths, settings *settings.Zootr) (zecs.Ocm, magicbean.ExplorableWorld) {
	ocm := bootstrap.Phase1_InitializeStorage(nil)
	bootstrap.PanicWhenErr(bootstrap.Phase2_ImportFromFiles(&ocm, paths))

	compileEnv := bootstrap.Phase3_ConfigureCompiler(&ocm, settings)

	codegen := mido.Compiler(&compileEnv)

	bootstrap.PanicWhenErr(bootstrap.Phase4_Compile(
		&ocm, &codegen,
	))

	world := bootstrap.Phase5_CreateWorld(&ocm, settings, objects.TableFrom(compileEnv.Objects))

	return ocm, world
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
