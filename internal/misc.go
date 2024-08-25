package internal

import (
	"github.com/etc-sudonters/substrate/mirrors"
	"github.com/etc-sudonters/substrate/slipup"
	"io/fs"
	"muzzammil.xyz/jsonc"
	"os"
	"reflect"
	"regexp"
	"strings"
)

var idcharsonly = regexp.MustCompile("[^a-z0-9]+")

type NormalizedStr string

func Normalize[S ~string](s S) NormalizedStr {
	return NormalizedStr(idcharsonly.ReplaceAllString(strings.ToLower(string(s)), ""))
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

func TypeAssert[T any](a any) (t T, err error) {
	t, cast := a.(T)
	if !cast {
		err = slipup.Createf("failed to cast %v to %s", a, mirrors.TypeOf[T]().Name())
	}
	err = nil
	return
}
