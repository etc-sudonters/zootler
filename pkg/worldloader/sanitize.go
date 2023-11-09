package worldloader

import (
	"regexp"
	"strings"
)

var _nameEscapeRe *regexp.Regexp = regexp.MustCompile(`['()\[\]-]`)

func EscapeName(name string) string {
	name = _nameEscapeRe.ReplaceAllLiteralString(name, "")
	return strings.ReplaceAll(name, " ", "_")
}
