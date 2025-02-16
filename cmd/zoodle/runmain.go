package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sudonters/libzootr/components"
	"sudonters/libzootr/internal"
	"sudonters/libzootr/magicbean/tracking"
	"sudonters/libzootr/settings"
	"sudonters/libzootr/zecs"
	"text/tabwriter"

	"github.com/etc-sudonters/substrate/skelly/bitset32"
	"github.com/etc-sudonters/substrate/slipup"

	"github.com/etc-sudonters/substrate/dontio"
	"github.com/etc-sudonters/substrate/stageleft"

	"runtime/pprof"
	"sudonters/libzootr/boot"
	"sudonters/libzootr/magicbean"
)

func runMain(ctx context.Context, std *dontio.Std, opts *cliOptions) stageleft.ExitCode {
	stopProfiling := profileto(opts.profile)
	defer stopProfiling()

	paths := boot.LoadPaths{
		Tokens:     filepath.Join(opts.dataDir, "items.json"),
		Placements: filepath.Join(opts.dataDir, "locations.json"),
		Scripts:    filepath.Join(opts.logicDir, "..", "helpers.json"),
		Relations:  opts.logicDir,
		Spoiler:    opts.spoiler,
	}

	theseSettings := settings.Default()
	FinalizeSettings(&theseSettings)

	generation, err := boot.Default(ctx, paths, &theseSettings)
	if err != nil {
		fmt.Fprintf(std.Err, "failed to boot zootr engine: %w", err)
		return 5
	}
	magicbean.CollectStartingItems(&generation)

	tokenNames, err := tracking.NameTableFrom(&generation.Ocm, zecs.With[components.TokenMarker])
	internal.PanicOnError(err)

	internal.PanicOnError(PrintInventory(std.Out, generation.Inventory, tokenNames))

	adult := Search{
		Workset:    generation.World.Graph.Roots(),
		Generation: &generation,
		Age:        AgeAdult,
	}
	child := Search{
		Workset:    generation.World.Graph.Roots(),
		Generation: &generation,
		Age:        AgeChild,
	}
	var i int
	for {
		i++
		std.WriteLineOut("")
		std.WriteLineOut("Sphere %d", i)
		adultReached := adult.Visit(ctx)
		childReached := child.Visit(ctx)
		reached := adultReached.Union(childReached)
		if reached.Len() == 0 {
			std.WriteLineOut("did not reach more nodes in either age")
			break
		}
		internal.PanicOnError(magicbean.CollectTokensFrom(
			&generation.Ocm,
			reached,
			generation.Inventory,
		))

		std.WriteLineOut("Inventory after collection")
		internal.PanicOnError(PrintInventory(std.Out, generation.Inventory, tokenNames))

		tabber := tabwriter.NewWriter(std.Out, 8, 2, 1, ' ', 0)
		fmt.Fprintf(tabber, "\tAdult\tChild\n")
		fmt.Fprintf(tabber, "Visited\t%3d\t%3d\n", adult.Visited.Len(), child.Visited.Len())
		fmt.Fprintf(tabber, "Reached\t%3d\t%3d\n", adultReached.Len(), childReached.Len())
		fmt.Fprintf(tabber, "Pending\t%3d\t%3d\n", adult.Workset.Len(), child.Workset.Len())
		tabber.Flush()
	}
	return stageleft.ExitCode(0)
}

type Search struct {
	Visited    bitset32.Bitset
	Workset    bitset32.Bitset
	Generation *magicbean.Generation
	Age        Age
}

func (this *Search) Visit(ctx context.Context) (reached bitset32.Bitset) {
	xplr := magicbean.Exploration{
		Visited: &this.Visited,
		Workset: &this.Workset,
	}

	results := explore(ctx, &xplr, this.Generation, this.Age)
	this.Workset = results.Pending
	return results.Reached
}

func FinalizeSettings(these *settings.Model) {
	these.Seed = 0x76E76E14E9691280
	these.Logic.Shuffling.Flags |= settings.ShuffleOcarinaNotes
	these.Logic.Spawns.StartAge = settings.StartAgeChild
	these.Logic.Connections.Flags |= settings.ConnectionOpenDoorOfTime
}

func profileto(path string) func() {
	if path == "" {
		return func() {}
	}

	f, err := os.Create(path)
	slipup.NeedsErrorHandling(err)
	pprof.StartCPUProfile(f)
	return func() { pprof.StopCPUProfile() }
}

func PrintInventory(w io.Writer, inventory magicbean.Inventory, names tracking.NameTable) error {
	tabber := tabwriter.NewWriter(w, 8, 8, 2, ' ', 0)
	if _, err := fmt.Fprintln(tabber, "Inventory:"); err != nil {
		return err
	}
	for ent, qty := range inventory {
		name, hasName := names[ent]
		if !hasName {
			continue
		}

		if _, err := fmt.Fprintf(tabber, "%s\t%d\n", name, qty); err != nil {
			return err
		}
	}
	return tabber.Flush()
}
