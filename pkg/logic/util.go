package logic

import (
	"regexp"
	"strings"
)

var _nameEscapeRe *regexp.Regexp = regexp.MustCompile(`['()\[\]-]`)

func EscapeName(name string) string {
	name = _nameEscapeRe.ReplaceAllLiteralString(name, "")
	return strings.ReplaceAll(name, " ", "_")
}

func CompressWhiteSpace[S ~string](s S) S {
	return S(strings.Join(strings.Fields(string(s)), " "))
}
