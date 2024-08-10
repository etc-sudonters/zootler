package main

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

var alphaOnly = regexp.MustCompile("[^a-z]+")

func normalize[S ~string](s S) string {
	return alphaOnly.ReplaceAllString(strings.ToLower(string(s)), "")
}

func IsFile(e fs.DirEntry) bool {
	return e.Type()&fs.ModeType == 0
}

func ReadJsonFile[T any](path string) (T, error) {
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

type BasicSampler struct {
	period, Modulo uint
	Do             func(v ...any) error
}

func (b *BasicSampler) Sample(v ...any) error {
	b.period++
	if b.period%b.Modulo == 0 {
		return b.Do(v...)
	}

	return nil
}

type std struct{ *dontio.Std }

func (s std) WriteLineOut(msg string, v ...any) {
	fmt.Fprintf(s.Out, msg+"\n", v...)
}

func (s std) WriteLineErr(msg string, v ...any) {
	fmt.Fprintf(s.Err, msg+"\n", v...)
}

func WriteLineOut(ctx context.Context, tpl string, v ...any) error {
	stdio, stdErr := dontio.StdFromContext(ctx)
	if stdErr != nil {
		return stdErr
	}
	std{stdio}.WriteLineOut(tpl, v...)
	return nil
}

func WriteLineErr(ctx context.Context, tpl string, v ...any) error {
	stdio, stdErr := dontio.StdFromContext(ctx)
	if stdErr != nil {
		return stdErr
	}
	std{stdio}.WriteLineErr(tpl, v...)
	return nil
}
