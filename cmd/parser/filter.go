package main

import (
	"fmt"
	"regexp"
	"strings"

	"sudonters/zootler/internal/rules"
)

type filter struct {
	region   *regexp.Regexp
	name     *regexp.Regexp
	kind     string
	errsOnly bool
}

func (f filter) MatchRegion(r rules.RegionName) bool {
	if f.region != nil && !f.region.Match([]byte(r)) {
		return false
	}
	return true
}

func (f filter) MatchKind(kind string) bool {
	if f.kind != "" && !strings.EqualFold(f.kind, kind) {
		return false
	}
	return true
}

func (f filter) MatchSpecific(name string) bool {
	if f.name != nil && !f.name.Match([]byte(name)) {
		return false
	}

	return true
}

func parseFilter(f string) filter {
	var filt filter
	var idx int

	parseKey := func() string {
		start := idx
		for ; idx < len(f); idx++ {
			if f[idx] == '=' {
				return f[start:idx]
			}
		}

		panic("could not parse key!")
	}

	parseValue := func() string {
		start := idx
		for ; idx < len(f); idx++ {
			if f[idx] == ';' {
				return f[start:idx]
			}
		}
		return f[start:]
	}

	for ; idx < len(f); idx++ {
		key := parseKey()
		idx++ // skip =
		value := parseValue()

		switch strings.ToLower(key) {
		case "region":
			filt.region = regexp.MustCompile(value)
			break
		case "name":
			filt.name = regexp.MustCompile(value)
			break
		case "kind":
			filt.kind = value
			break
		default:
			panic(fmt.Errorf("unknown filter option %q", key))
		}
	}

	return filt
}
