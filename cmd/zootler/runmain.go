package main

import (
	"context"
	"os"
	"path/filepath"
	"sudonters/zootler/internal/settings"

	"github.com/etc-sudonters/substrate/stageleft"

	"runtime/pprof"
	"sudonters/zootler/cmd/zootler/bootstrap"
	"sudonters/zootler/internal/query"
	"sudonters/zootler/magicbean"
	"sudonters/zootler/mido"
	"sudonters/zootler/mido/objects"
	"sudonters/zootler/zecs"
)

func runMain(ctx context.Context, opts cliOptions) stageleft.ExitCode {
	stopProfiling := profileto(opts.profile)
	defer stopProfiling()

	paths := bootstrap.LoadPaths{
		Tokens:     filepath.Join(opts.dataDir, "items.json"),
		Placements: filepath.Join(opts.dataDir, "locations.json"),
		Scripts:    filepath.Join(opts.logicDir, "..", "helpers.json"),
		Relations:  opts.logicDir,
	}

	theseSettings := settings.Default()
	artifacts := setup(paths, &theseSettings)
	eng := artifacts.Ocm.Engine()
	tbl, err := query.ExtractTable(eng)
	bootstrap.PanicWhenErr(err)

	explore(&artifacts)
	displaycolstats(tbl)
	return stageleft.ExitCode(0)
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
