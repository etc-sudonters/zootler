package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"runtime/debug"
	"sudonters/zootler/internal/query"
	"sudonters/zootler/internal/table"
	"sudonters/zootler/internal/table/columns"
	"sudonters/zootler/pkg/world/components"

	"github.com/etc-sudonters/substrate/dontio"
	"github.com/etc-sudonters/substrate/mirrors"
	"github.com/etc-sudonters/substrate/stageleft"
	"muzzammil.xyz/jsonc"
)

type missingRequired string // option name

func (arg missingRequired) Error() string {
	return fmt.Sprintf("%s is required", string(arg))
}

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

	storage := query.NewEngine()

	storage.CreateColumn(table.BuildColumnOf[components.Name](columns.NewSliceColumn()))
	storage.CreateColumn(table.BuildColumnOf[components.Song](columns.NewHashMap()))
	storage.CreateColumn(table.BuildColumnOf[components.Location](columns.NewBit[components.Location]()))
	storage.CreateColumn(table.BuildColumnOf[components.Token](columns.NewBit[components.Token]()))

	loadLocations("inputs/data/locations.json", storage)
	loadItems("inputs/data/items.json", storage)

	q := storage.CreateQuery()
	q.Exists(mirrors.TypeOf[components.Location]())
	allLocs, err := storage.Retrieve(q)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(stdio.Out, "Count of all locations: %d\n", allLocs.Len())

	q = storage.CreateQuery()
	q.Exists(mirrors.TypeOf[components.Location]())
	q.Exists(mirrors.TypeOf[components.Song]())
	songLocs, err := storage.Retrieve(q)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(stdio.Out, "Count of Song locations: %d\n", songLocs.Len())

	q = storage.CreateQuery()

	q = storage.CreateQuery()
	q.Exists(mirrors.TypeOf[components.Location]())
	q.NotExists(mirrors.TypeOf[components.Song]())
	notSongLocs, err := storage.Retrieve(q)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(stdio.Out, "Count of Song locations: %d\n", notSongLocs.Len())

	q = storage.CreateQuery()
	q.Exists(mirrors.TypeOf[components.Token]())
	allToks, err := storage.Retrieve(q)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(stdio.Out, "Count of all tokens: %d\n", allToks.Len())

	q = storage.CreateQuery()
	q.Exists(mirrors.TypeOf[components.Token]())
	q.Exists(mirrors.TypeOf[components.Song]())
	songToks, err := storage.Retrieve(q)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(stdio.Out, "Count of Song tokens: %d\n", songToks.Len())

	q = storage.CreateQuery()

	q = storage.CreateQuery()
	q.Exists(mirrors.TypeOf[components.Token]())
	q.NotExists(mirrors.TypeOf[components.Song]())
	notSongToks, err := storage.Retrieve(q)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(stdio.Out, "Count of Song tokens: %d\n", notSongToks.Len())
}

func loadLocations(path string, storage query.Engine) {
	raw, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var locs []struct {
		Name string `json:"name"`
		Type string `json:"type"`
	}

	if err := jsonc.Unmarshal(raw, &locs); err != nil {
		panic(err)
	}

	var song components.Song
	var location components.Location

	for _, l := range locs {
		id, err := storage.InsertRow(components.Name(l.Name), location)
		if err != nil {
			panic(err)
		}
		if l.Type == "Song" {
			storage.SetValues(id, table.Values{song})
		}
	}
}

func loadItems(path string, storage query.Engine) {
	raw, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var items []struct {
		Name string `json:"name"`
		Type string `json:"type"`
	}

	var tok components.Token
	var song components.Song

	if err := jsonc.Unmarshal(raw, &items); err != nil {
		panic(err)
	}

	for _, item := range items {
		id, err := storage.InsertRow(components.Name(item.Name), tok)
		if err != nil {
			panic(err)
		}
		if item.Type == "Song" {
			storage.SetValues(id, table.Values{song})
		}
	}

}
