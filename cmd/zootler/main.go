package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime/debug"

	"sudonters/zootler/internal/ioutil"
	"sudonters/zootler/internal/rules"
	"sudonters/zootler/pkg/entity"
	"sudonters/zootler/pkg/logic"
	"sudonters/zootler/pkg/world"
)

type cliOptions struct {
	logicDir string `short:"-l" description:"Path to logic files" required:"t"`
	dataDir  string `short:"-d" description:"Path to data files" required:"t"`
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
	var exit ioutil.ExitCode = ioutil.ExitSuccess
	stdio := ioutil.Std{
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
			if exit != ioutil.ExitSuccess {
				exit = ioutil.AsExitCode(r, ioutil.ExitPanic)
			}
		}
	}()

	ctx := context.Background()
	ctx = ioutil.AddStdToContext(ctx, &stdio)
	ctx = ioutil.AddExitCodeToContext(ctx, &exit)

	flag.StringVar(&opts.logicDir, "l", "", "Directory where logic files are located")
	flag.StringVar(&opts.dataDir, "d", "", "Directory where data files are stored")
	flag.Parse()

	if cliErr := opts.validate(); cliErr != nil {
		fmt.Fprintf(stdio.Err, "%s\n", cliErr.Error())
		exit = ioutil.ExitBadFlag
		return
	}

	b, err := buildWorldFromFiles(ctx, opts)
	if err != nil {
		exit = ioutil.GetExitCodeOr(err, ioutil.ExitCode(2))
		fmt.Fprintf(stdio.Err, "Error while parsing files: %s\n", err.Error())
		return
	}

	w := b.Build()

	assumed := &world.AssumedFill{
		Locations: []entity.Selector{entity.With[logic.Song]{}},
		Items:     []entity.Selector{entity.With[logic.Song]{}},
	}
	if err := assumed.Fill(ctx, w, world.ConstGoal(true)); err != nil {
		exit = ioutil.GetExitCodeOr(err, ioutil.ExitCode(2))
		fmt.Fprintf(stdio.Err, "Error during placement: %s\n", err.Error())
		return
	}

	if err := showTokenPlacements(ctx, w, entity.With[logic.Song]{}); err != nil {
		exit = ioutil.GetExitCodeOr(err, ioutil.ExitCode(2))
		fmt.Fprintf(stdio.Err, "Error during placement review: %s\n", err.Error())
		return
	}
}

type missingRequired string // option name

func (arg missingRequired) Error() string {
	return fmt.Sprintf("%s is required", string(arg))
}

func buildWorldFromFiles(ctx context.Context, opts cliOptions) (b *world.Builder, err error) {
	b = world.NewBuilder(1)
	locs, err := loadLocationFile(ctx, b, path.Join(opts.dataDir, "locations.json"))

	if err != nil {
		return
	}
	err = loadItemFile(ctx, b, path.Join(opts.dataDir, "items.json"))
	if err != nil {
		return
	}
	err = loadLogicFiles(ctx, b, locs, opts.logicDir)
	return
}

func loadLocationFile(ctx context.Context, b *world.Builder, path string) (map[string]entity.View, error) {
	locs, err := logic.ReadLocationFile(path)
	if err != nil {
		return nil, err
	}

	lookup := make(map[string]entity.View, len(locs))

	for _, loc := range locs {
		ent, err := b.Pool.Create(logic.Name(loc.Name))
		if err != nil {
			// bitpool ran out of IDs
			panic(ioutil.AttachExitCode(err, ioutil.ExitCode(100)))
		}
		for _, comp := range logic.GetAllLocationComponents(loc) {
			ent.Add(comp)
		}
		ent.Add(logic.Location{})
		b.AddNode(ent)
		lookup[loc.Name] = ent
	}

	return lookup, nil
}

func loadItemFile(ctx context.Context, b *world.Builder, path string) error {
	items, err := logic.ReadItemFile(path)
	if err != nil {
		return err
	}

	for _, item := range items {
		ent, err := b.Pool.Create(logic.Name(item.Name))
		if err != nil {
			panic(ioutil.AttachExitCode(err, ioutil.ExitCode(100)))
		}

		ent.Add(logic.Token{})
		ent.Add(item.Importance)
		for _, comp := range item.Components {
			ent.Add(comp)
		}
	}

	return nil
}

func loadLogicFiles(ctx context.Context, b *world.Builder, locCache map[string]entity.View, path string) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("failed to read logic from %s: %w", path, err)
	}

	for _, entry := range entries {
		if filepath.Ext(entry.Name()) != "json" {
			continue
		}

		fp := filepath.Join(path, entry.Name())
		locs, err := rules.ReadLogicFile(fp)
		if err != nil {
			return fmt.Errorf("failed to read logic file %s: %w", fp, err)
		}

		for _, loc := range locs {
			region := locCache[string(loc.Region)]
			for _, comp := range loc.Components() {
				region.Add(comp)
			}

			for evt, rule := range loc.Events {
				ent, err := b.Pool.Create(logic.Name(evt))
				if err != nil {
					panic(ioutil.AttachExitCode(err, ioutil.ExitCode(100)))
				}
				ent.Add(logic.Event{})
				ent.Add(logic.RawRule(rule))
			}

			for check, rule := range loc.Locations {
				ent, err := b.Pool.Create(logic.Name(check))
				if err != nil {
					panic(ioutil.AttachExitCode(err, ioutil.ExitParseFailure))
				}
				ent.Add(logic.RawRule(rule))
			}

			for exit, rule := range loc.Exits {
				raw := compressWhiteSpace(string(rule))
				name := fmt.Sprintf("%s -> %s", loc.Region, exit)
				exit := locCache[string(exit)]
				edge, err := b.AddEdge(region, exit)
				if err != nil {
					panic(ioutil.AttachExitCode(err, ioutil.ExitParseFailure))
				}
				edge.Add(logic.RawRule(raw))

				lex := rules.NewLexer(name, raw)
				parser := rules.NewParser(lex)
				rule, err := parser.ParseTotalRule()
				if err != nil {
					panic(ioutil.AttachExitCode(err, ioutil.ExitParseFailure))
				}
				edge.Add(logic.ParsedRule{E: rule.Rule})
			}
		}
	}

	return nil
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
	placed, err := w.Entities.Query(entity.With[logic.Inhabits]{}, qs...)
	if err != nil {
		return fmt.Errorf("while querying placements: %w", err)
	}
	stdio, _ := ioutil.StdFromContext(ctx)

	for _, tok := range placed {
		var itemName logic.Name
		var placementName logic.Name
		var placement logic.Inhabits

		err = tok.Get(&itemName)
		if err != nil {
			return err
		}
		err = tok.Get(&placement)
		if err != nil {
			return err
		}
		w.Entities.Get(entity.Model(placement), &placementName)
		if placementName == "" {
			return fmt.Errorf("%v did not have an attached name", placement)
		}

		fmt.Fprintf(stdio.Out, "%s placed at %s", itemName, placementName)
	}

	return nil
}
