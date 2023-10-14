package main

import (
	"context"
	"flag"
	"fmt"
	"os"

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
	var exit *ioutil.ExitCode = &ioutil.ExitSuccess
	defer func() { os.Exit(int(*exit)) }()
	stdio := ioutil.Std{
		In:  os.Stdin,
		Out: os.Stdout,
		Err: os.Stdout,
	}

	ctx := context.Background()
	ctx = ioutil.AddStdToContext(ctx, &stdio)
	ctx = ioutil.AddExitCodeToContext(ctx, exit)

	flag.StringVar(&opts.logicDir, "l", "", "Directory where logic files are located")
	flag.StringVar(&opts.dataDir, "d", "", "Directory where data files are stored")
	flag.Parse()

	if cliErr := opts.validate(); cliErr != nil {
		fmt.Fprintf(stdio.Err, "%s\n", cliErr.Error())
		flag.PrintDefaults()
		*exit = ioutil.ExitBadFlag
		return
	}

	builder := world.NewBuilder(1)
	buildFromDataDir(ctx, builder, opts.dataDir)
	buildFromLogicDir(ctx, builder, opts.logicDir)
}

type missingRequired string // option name

func (arg missingRequired) Error() string {
	return fmt.Sprintf("%s is required", string(arg))
}

func buildFromLogicDir(ctx context.Context, b *world.Builder, path string) error {
	panic("not impled!")
}

func buildFromDataDir(ctx context.Context, b *world.Builder, path string) error {
	panic("not impled")
}
