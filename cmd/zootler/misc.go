package main

import (
	"regexp"
	"strings"
)

type str interface {
	~string
}

var alphaOnly = regexp.MustCompile("[^a-z]+")

func normalize[S str](s S) string {
	return alphaOnly.ReplaceAllString(strings.ToLower(s), "")
}
