package main

import (
	"regexp"
	"strings"
)

var alphaOnly = regexp.MustCompile("[^a-z]+")

func normalize(s string) string {
	return alphaOnly.ReplaceAllString(strings.ToLower(s), "")
}
