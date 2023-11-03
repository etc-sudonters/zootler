package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"regexp"
	"runtime/debug"

	"sudonters/zootler/cmd/zootler/tui"
	"sudonters/zootler/internal/entity"
	"sudonters/zootler/pkg/filler"
	"sudonters/zootler/pkg/logic"
	"sudonters/zootler/pkg/world"
	"sudonters/zootler/pkg/worldloader"

	"github.com/etc-sudonters/substrate/dontio"
	"github.com/etc-sudonters/substrate/stageleft"
)

type cliOptions struct {
	logicDir   string `short:"-l" description:"Path to logic files" required:"t"`
	dataDir    string `short:"-d" description:"Path to data files" required:"t"`
	visualizer bool   `short:"-v" description:"Open visualizer" required:"f"`
}

func (opts *cliOptions) init() {
	flag.StringVar(&opts.logicDir, "l", "", "Directory where logic files are located")
	flag.StringVar(&opts.dataDir, "d", "", "Directory where data files are stored")
	flag.BoolVar(&opts.visualizer, "v", false, "Open visualizer")
	flag.Parse()
}

func (c cliOptions) validate() error {
	if c.logicDir == "" {
		return missingRequired("-l")
	}

	if c.dataDir == "" {
		return missingRequired("-d")
	}

	return nil
}

func main() {
	var opts cliOptions
	var exit stageleft.ExitCode = stageleft.ExitSuccess
	stdio := dontio.Std{
		In:  os.Stdin,
		Out: os.Stdout,
		Err: os.Stdout,
	}
	defer func() {
		os.Exit(int(exit))
	}()
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				fmt.Fprintf(stdio.Err, "%s\n", err)
			}
			_, _ = stdio.Err.Write(debug.Stack())
			if exit != stageleft.ExitSuccess {
				exit = stageleft.AsExitCode(r, stageleft.ExitCode(126))
			}
		}
	}()

	ctx := context.Background()
	ctx = dontio.AddStdToContext(ctx, &stdio)

	(&opts).init()

	if cliErr := opts.validate(); cliErr != nil {
		fmt.Fprintf(stdio.Err, "%s\n", cliErr.Error())
		exit = stageleft.ExitCode(2)
		return
	}

	b := world.NewBuilder()
	loader := worldloader.FileSystemLoader{
		LogicDirectory: opts.logicDir,
		DataDirectory:  opts.dataDir,
	}

	if err := loader.LoadInto(ctx, b); err != nil {
		exit = stageleft.ExitCode(98)
		panic(err)
	}

	w := b.Build()

	if opts.visualizer {
		v := tui.Tui(w)
		if err := v.Run(ctx); err != nil {
			panic(err)
		}
		return
	}

	assumed := &filler.AssumedFill{
		Locations: []entity.Selector{entity.With[logic.Song]{}},
		Items:     []entity.Selector{entity.With[logic.Song]{}},
	}
	if err := assumed.Fill(ctx, w, filler.ConstGoal(true)); err != nil {
		exit = stageleft.ExitCodeFromErr(err, stageleft.ExitCode(2))
		fmt.Fprintf(stdio.Err, "Error during placement: %s\n", err.Error())
		return
	}

	if err := showTokenPlacements(ctx, w, entity.With[logic.Song]{}); err != nil {
		exit = stageleft.ExitCodeFromErr(err, stageleft.ExitCode(2))
		fmt.Fprintf(stdio.Err, "Error during placement review: %s\n", err.Error())
		return
	}
}

type missingRequired string // option name

func (arg missingRequired) Error() string {
	return fmt.Sprintf("%s is required", string(arg))
}

func compressWhiteSpace(r string) string {
	r = trailWhiteSpace.ReplaceAllLiteralString(r, "")
	r = leadWhiteSpace.ReplaceAllLiteralString(r, "")
	return compressWhiteSpaceRe.ReplaceAllLiteralString(r, " ")
}

var compressWhiteSpaceRe *regexp.Regexp = regexp.MustCompile(`\s+`)
var leadWhiteSpace *regexp.Regexp = regexp.MustCompile(`^\s+`)
var trailWhiteSpace *regexp.Regexp = regexp.MustCompile(`\s+$`)

func showTokenPlacements(ctx context.Context, w world.World, qs ...entity.Selector) error {
	filt := make([]entity.Selector, len(qs)+1)
	filt[0] = entity.With[logic.Inhabits]{}
	copy(filt[1:], qs)

	placed, err := w.Entities.Query(filt)
	if err != nil {
		return fmt.Errorf("while querying placements: %w", err)
	}
	stdio, _ := dontio.StdFromContext(ctx)

	for _, tok := range placed {
		var itemName world.Name
		var placementName world.Name
		var placement logic.Inhabits

		err = tok.Get(&itemName)
		if err != nil {
			return err
		}
		err = tok.Get(&placement)
		if err != nil {
			return err
		}
		w.Entities.Get(entity.Model(placement), []interface{}{&placementName})
		if placementName == "" {
			return fmt.Errorf("%v did not have an attached name", placement)
		}

		fmt.Fprintf(stdio.Out, "%s placed at %s", itemName, placementName)
	}

	return nil
}
