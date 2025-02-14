package main

import (
	"context"
	"math/rand/v2"
	"os"
	"path/filepath"
	"sudonters/libzootr/internal"
	"sudonters/libzootr/internal/skelly/bitset32"
	"sudonters/libzootr/magicbean/tracking"
	"sudonters/libzootr/settings"

	"github.com/etc-sudonters/substrate/dontio"
	"github.com/etc-sudonters/substrate/rng"
	"github.com/etc-sudonters/substrate/stageleft"

	"runtime/pprof"
	"sudonters/libzootr/cmd/zoodle/bootstrap"
	"sudonters/libzootr/magicbean"
	"sudonters/libzootr/mido"
	"sudonters/libzootr/mido/objects"
)

func runMain(ctx context.Context, opts cliOptions) stageleft.ExitCode {
	stopProfiling := profileto(opts.profile)
	defer stopProfiling()

	paths := bootstrap.LoadPaths{
		Tokens:     filepath.Join(opts.dataDir, "items.json"),
		Placements: filepath.Join(opts.dataDir, "locations.json"),
		Scripts:    filepath.Join(opts.logicDir, "..", "helpers.json"),
		Relations:  opts.logicDir,
		Spoiler:    opts.spoiler,
	}

	theseSettings := settings.Default()
	FinalizeSettings(&theseSettings)
	generation := setup(ctx, paths, &theseSettings)
	generation.Settings = theseSettings
	CollectStartingItems(&generation)

	visited := bitset32.Bitset{}
	workset := generation.World.Graph.Roots()
	xplr := magicbean.Exploration{
		Visited: &visited,
		Workset: &workset,
	}
	results := explore(ctx, &xplr, &generation, AgeAdult)
	std, err := dontio.StdFromContext(ctx)
	internal.PanicOnError(err)
	std.WriteLineOut("Visited %d", visited.Len())
	std.WriteLineOut("Reached %d", results.Reached.Len())
	std.WriteLineOut("Pending %d", results.Pending.Len())
	return stageleft.ExitCode(0)
}

func FinalizeSettings(these *settings.Zootr) {
	these.Seed = 0x76E76E14E9691280
	these.Shuffling.OcarinaNotes = true
	these.Spawns.StartingAge = settings.StartAgeAdult
	these.Locations.OpenDoorOfTime = true
}

func setup(ctx context.Context, paths bootstrap.LoadPaths, settings *settings.Zootr) (generation magicbean.Generation) {
	ocm := bootstrap.Phase1_InitializeStorage(nil)
	trackSet := tracking.NewTrackingSet(&ocm)

	phase2Error := bootstrap.Phase2_ImportFromFiles(ctx, settings, &ocm, &trackSet, paths)
	internal.PanicOnError(phase2Error)

	compileEnv := bootstrap.Phase3_ConfigureCompiler(&ocm, settings)

	codegen := mido.Compiler(&compileEnv)

	bootstrap.PanicWhenErr(bootstrap.Phase4_Compile(
		&ocm, &codegen,
	))

	world := bootstrap.Phase5_CreateWorld(&ocm, settings, objects.TableFrom(compileEnv.Objects))

	generation.Ocm = ocm
	generation.World = world
	generation.Objects = objects.TableFrom(compileEnv.Objects)
	generation.Inventory = magicbean.EmptyInventory()
	generation.Rng = *rand.New(rng.NewXoshiro256PPFromU64(settings.Seed))

	return generation
}

func profileto(path string) func() {
	if path == "" {
		return func() {}
	}

	f, err := os.Create(path)
	bootstrap.PanicWhenErr(err)
	pprof.StartCPUProfile(f)
	return func() { pprof.StopCPUProfile() }
}
