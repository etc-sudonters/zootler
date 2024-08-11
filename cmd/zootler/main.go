package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path"
	"runtime/debug"
	"sudonters/zootler/internal/app"
	"sudonters/zootler/internal/query"
	"sudonters/zootler/internal/rules/compiler"

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

func (opts *cliOptions) init() error {
	flag.StringVar(&opts.logicDir, "l", "", "Directory where logic files are located")
	flag.StringVar(&opts.dataDir, "d", "", "Directory where data files are stored")
	flag.BoolVar(&opts.includeMq, "M", false, "Whether or not to include MQ data")
	flag.Parse()

	if opts.logicDir == "" {
		return missingRequired("-l")
	}

	if opts.dataDir == "" {
		return missingRequired("-d")
	}

	return nil
}

func main() {
	var opts cliOptions
	var appExitCode stageleft.ExitCode = stageleft.ExitSuccess
	stdio := dontio.Std{
		In:  os.Stdin,
		Out: os.Stdout,
		Err: os.Stderr,
	}
	std := std{&stdio}
	defer func() {
		os.Exit(int(appExitCode))
	}()
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				std.WriteLineErr("%s", err)
			}
			_, _ = std.Err.Write(debug.Stack())
			if appExitCode != stageleft.ExitSuccess {
				appExitCode = stageleft.AsExitCode(r, stageleft.ExitCode(126))
			}
		}
	}()

	exitWithErr := func(code stageleft.ExitCode, err error) {
		appExitCode = code
		std.WriteLineErr(err.Error())
	}

	ctx := context.Background()
	ctx = dontio.AddStdToContext(ctx, &stdio)

	if argsErr := (&opts).init(); argsErr != nil {
		exitWithErr(2, argsErr)
		return
	}

	_, appCreateErr := app.New(ctx,

		app.Setup(CreateScheme{DDL: MakeDDL()}),
		app.Setup(DataFileLoader[FileItem]{
			IncludeMQ: opts.includeMq,
			Path:      path.Join(opts.dataDir, "items.json"),
		}),
		app.Setup(DataFileLoader[FileLocation]{
			IncludeMQ: opts.includeMq,
			Path:      path.Join(opts.dataDir, "locations.json"),
			Add:       new(AttachDefaultItem),
		}),
		app.Setup(WorldFileLoader{
			IncludeMQ: opts.includeMq,
			Path:      opts.logicDir,
			Helpers:   path.Join(path.Dir(opts.logicDir), "helpers.json"),
		}),
		app.Setup(&LogicCompiler{}),
		app.Setup(DebugSetupFunc(func(ctx context.Context, storage query.Engine) error {
			c := new(compiler.ChunkBuilder)
			c.WriteOp(compiler.OP_NOP)
			c.PushConst(1337)
			c.WriteOp(compiler.OP_RETURN)
			WriteLineOut(ctx, c.Disassemble("test"))
			return nil
		})),
	)

	if appCreateErr != nil {
		exitWithErr(3, appCreateErr)
		return
	}
}
