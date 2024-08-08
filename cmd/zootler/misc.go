package main

import (
	"io/fs"
	"os"
	"reflect"
	"regexp"
	"strings"

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
