package internal

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"reflect"
	"regexp"
	"strings"

	"github.com/etc-sudonters/substrate/dontio"
	"github.com/etc-sudonters/substrate/mirrors"
	"muzzammil.xyz/jsonc"
)

var alphanumonly = regexp.MustCompile("[^a-z0-9]+")

func Normalize[S ~string](s S) string {
	return alphanumonly.ReplaceAllString(strings.ToLower(string(s)), "")
}

func IsFile(e fs.DirEntry) bool {
	return e.Type()&fs.ModeType == 0
}

func ReadJsonFileStringMap(path string) (map[string]string, error) {
	t := make(map[string]string)
	raw, readErr := os.ReadFile(path)
	if readErr != nil {
		return t, readErr
	}

	err := jsonc.Unmarshal(raw, &t)
	return t, err
}

func ReadJsonFileAs[T any](path string) (T, error) {
	var t T
	raw, readErr := os.ReadFile(path)
	if readErr != nil {
		return t, readErr
	}

	err := jsonc.Unmarshal(raw, &t)
	return t, err
}

func T[E any]() reflect.Type {
	return mirrors.TypeOf[E]()
}

type Std struct{ *dontio.Std }

func (s Std) WriteLineOut(msg string, v ...any) {
	fmt.Fprintf(s.Out, msg+"\n", v...)
}

func (s Std) WriteLineErr(msg string, v ...any) {
	fmt.Fprintf(s.Err, msg+"\n", v...)
}

func WriteLineOut(ctx context.Context, tpl string, v ...any) error {
	stdio, stdErr := dontio.StdFromContext(ctx)
	if stdErr != nil {
		return stdErr
	}
	Std{stdio}.WriteLineOut(tpl, v...)
	return nil
}

func WriteLineErr(ctx context.Context, tpl string, v ...any) error {
	stdio, stdErr := dontio.StdFromContext(ctx)
	if stdErr != nil {
		return stdErr
	}
	Std{stdio}.WriteLineErr(tpl, v...)
	return nil
}
