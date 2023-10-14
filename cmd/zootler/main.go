package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"runtime/debug"

	"github.com/etc-sudonters/zootler/internal/errs"
	"github.com/etc-sudonters/zootler/internal/ioutil"
	"github.com/etc-sudonters/zootler/pkg/world"
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
			stdio.Err.Write(debug.Stack())
			exit = ioutil.AsExitCode(r, ioutil.ExitPanic)
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

	builder := world.NewBuilder(1)
	if err := buildFromDataDir(ctx, builder, opts.dataDir); err != nil {
		exit = ioutil.GetExitCodeOr(err, ioutil.ExitCode(2))
		fmt.Fprintf(stdio.Err, "Error while parsing files: %s\n", err.Error())
		return
	}
	if err := buildFromLogicDir(ctx, builder, opts.logicDir); err != nil {
		exit = ioutil.GetExitCodeOr(err, ioutil.ExitCode(2))
		fmt.Fprintf(stdio.Err, "Error while parsing files: %s\n", err.Error())
		return
	}
}

type missingRequired string // option name

func (arg missingRequired) Error() string {
	return fmt.Sprintf("%s is required", string(arg))
}

func buildFromLogicDir(ctx context.Context, b *world.Builder, path string) error {
	return errs.NotImplErr
}

func buildFromDataDir(ctx context.Context, b *world.Builder, path string) error {
	return errs.NotImplErr
}
