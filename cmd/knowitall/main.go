package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime/debug"

	"github.com/etc-sudonters/substrate/dontio"
	"github.com/etc-sudonters/substrate/stageleft"
)

func main() {
	var opts cliOptions
	var appExitCode stageleft.ExitCode = stageleft.ExitSuccess
	realStd := dontio.Std{
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
				realStd.WriteLineErr("%s", err)
			} else if str, ok := r.(string); ok {
				realStd.WriteLineErr("%s", str)
			}
			_, _ = realStd.Err.Write(debug.Stack())
			if appExitCode != stageleft.ExitSuccess {
				appExitCode = stageleft.AsExitCode(r, stageleft.ExitCode(126))
			}
		}
	}()

	ctx := context.Background()
	if argsErr := (&opts).init(); argsErr != nil {
		appExitCode = 2
		realStd.WriteLineErr(argsErr.Error())
		return
	}
	appStd := dontio.Std{}
	cleanup, err := redirectAppStd(&appStd, &opts)
	if cleanup != nil {
		defer cleanup()
	}
	if err != nil {
		fmt.Fprintf(realStd.Err, "Failed to redirect application std{in,out,err}\n%v\n", err)
		appExitCode = 3
		return
	}
	ctx = dontio.AddStdToContext(ctx, &appStd)

	appExitCode = runMain(ctx, &appStd, &opts)
	return
}

func redirectAppStd(std *dontio.Std, opts *cliOptions) (func(), error) {
	logDirErr := os.Mkdir(opts.logDir, 0777)
	if logDirErr != nil && !os.IsExist(logDirErr) {
		return nil, fmt.Errorf("failed to initialize log dir %s: %w", opts.logDir, logDirErr)
	}

	std.In = dontio.AlwaysErrReader{io.ErrUnexpectedEOF}
	return dontio.FileStd(std, opts.logDir)
}

type missingRequired string // option name

func (arg missingRequired) Error() string {
	return fmt.Sprintf("%s is required", string(arg))
}

type cliOptions struct {
	logDir string
}

func (opts *cliOptions) init() error {
	var path string
	var pathErr error

	flag.StringVar(&path, "l", "", "Directory open log files in")
	flag.Parse()

	if path == "" {
		path = ".logs"
	}
	opts.logDir, pathErr = filepath.Abs(path)
	if pathErr != nil {
		pathErr = fmt.Errorf("failed to initialize log dir %q: %w", path, pathErr)
	}

	if opts.logDir == "" {
		return fmt.Errorf("failed to initialize log dir: empty path!")
	}

	return pathErr
}
