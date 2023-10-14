package ioutil

import "context"

var exitkey ctxkey = "exitcode"

func AddExitCodeToContext(ctx context.Context, e *ExitCode) context.Context {
	return context.WithValue(ctx, exitkey, e)
}

func ExitCodeFromCtx(ctx context.Context) (*ExitCode, error) {
	v := ctx.Value(exitkey)
	if v == nil {
		return nil, notInCtx("exit code")
	}
	return v.(*ExitCode), nil
}

type ExitCode uint8

var (
	ExitSuccess ExitCode = 0
	ExitBadFlag ExitCode = 1
)
