package worldloader

import (
	"strings"
)

func EscapeName(name string) string {
	name = strings.ReplaceAll(name, "'()[]-", "")
	return strings.ReplaceAll(name, " ", "_")
}
