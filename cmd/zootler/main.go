package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path"
	"runtime/debug"
	"sudonters/zootler/internal/app"

	"github.com/etc-sudonters/substrate/dontio"
	"github.com/etc-sudonters/substrate/stageleft"
)

type missingRequired string // option name

func (arg missingRequired) Error() string {
	return fmt.Sprintf("%s is required", string(arg))
}

type cliOptions struct {
	logicDir  string
	dataDir   string
	includeMq bool
}

func (opts *cliOptions) init() {
	flag.StringVar(&opts.logicDir, "l", "", "Directory where logic files are located")
	flag.StringVar(&opts.dataDir, "d", "", "Directory where data files are stored")
	flag.BoolVar(&opts.includeMq, "M", false, "Whether or not to include MQ data")
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

type std struct{ *dontio.Std }

func (s std) WriteLineOut(msg string, v ...any) {
	fmt.Fprintf(s.Out, msg+"\n", v...)
}

func main() {
	var opts cliOptions
	var appExitCode stageleft.ExitCode = stageleft.ExitSuccess
	stdio := dontio.Std{
		In:  os.Stdin,
		Out: os.Stdout,
		Err: os.Stderr,
	}
	defer func() {
		os.Exit(int(appExitCode))
	}()
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				fmt.Fprintf(stdio.Err, "%s\n", err)
			}
			_, _ = stdio.Err.Write(debug.Stack())
			if appExitCode != stageleft.ExitSuccess {
				appExitCode = stageleft.AsExitCode(r, stageleft.ExitCode(126))
			}
		}
	}()

	exitWithErr := func(code stageleft.ExitCode, err error) {
		appExitCode = code
		fmt.Fprintf(stdio.Err, "%s\n", err.Error())
	}

	ctx := context.Background()
	ctx = dontio.AddStdToContext(ctx, &stdio)

	(&opts).init()

	if cliErr := opts.validate(); cliErr != nil {
		exitWithErr(2, cliErr)
		return
	}

	app, err := app.NewApp(ctx,
		app.Storage(CreateScheme{DDL: MakeDDL()}),
		app.Storage(DataFileLoader[FileItem]{
			IncludeMQ: opts.includeMq,
			Path:      path.Join(opts.dataDir, "items.json"),
		}),
		app.Storage(DataFileLoader[FileLocation]{
			IncludeMQ: opts.includeMq,
			Path:      path.Join(opts.dataDir, "locations.json"),
		}),
		app.Storage(WorldFileLoader{
			IncludeMQ: opts.includeMq,
			Path:      opts.logicDir,
			Helpers:   path.Join(path.Dir(opts.logicDir), "helpers.json"),
		}),
	)

	if err != nil {
		exitWithErr(3, err)
		return
	}

	if err := example(app.Ctx(), app.Engine()); err != nil {
		exitWithErr(4, err)
		return
	}
}
