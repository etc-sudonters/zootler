package ioutil

import (
	"context"
	"errors"
)

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

func AttachExitCode(err error, ec ExitCode) ExitCodeError {
	return ExitCodeError{err, ec}
}

type ExitCodeError struct {
	Err  error
	Code ExitCode
}

func (ece ExitCodeError) Unwrap() error {
	return ece.Err
}

func (ece ExitCodeError) Error() string {
	return ece.Err.Error()
}

func AsExitCode(r interface{}, exit ExitCode) ExitCode {
	switch r := r.(type) {
	case ExitCode:
		exit = r
		break
	case *ExitCode:
		exit = *r
		break
	case ExitCodeError:
		exit = r.Code
		break
	case error:
		exit = GetExitCodeOr(r, exit)
		break
	case uint8:
		exit = ExitCode(r)
		break
	}

	return exit
}

func GetExitCodeOr(err error, orCode ExitCode) ExitCode {
	if code := GetExitCode(err); code != ExitUnknown {
		return code
	}
	return orCode
}

func GetExitCode(err error) ExitCode {
	var ece ExitCodeError
	ec := ExitUnknown

	if errors.As(err, &ece) {
		ec = ece.Code
	}

	return ec
}

type ExitCode uint8

var (
	ExitSuccess      ExitCode = 0
	ExitBadFlag      ExitCode = 1
	ExitUnknown      ExitCode = 99
	ExitPanic        ExitCode = 69
	ExitParseFailure ExitCode = 100
	ExitQueryFail    ExitCode = 15
)
