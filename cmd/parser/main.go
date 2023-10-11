package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/etc-sudonters/zootler/internal/console"
	"github.com/etc-sudonters/zootler/internal/rules"
	"muzzammil.xyz/jsonc"
)

var parseErrorColor console.ForegroundColor = 141

func main() {
	var logicFileDir string
	var errsOnly bool
	var rawFilter string
	var filt filter
	var helpersOnly bool

	flag.StringVar(&logicFileDir, "l", "", "Directory where logic files are located")
	flag.StringVar(&rawFilter, "f", "", "Look it's complicated")
	flag.BoolVar(&errsOnly, "E", false, "set to only display errors")
	flag.BoolVar(&helpersOnly, "H", false, "")
	flag.Parse()

	if logicFileDir == "" {
		fmt.Fprint(os.Stderr, "-l is required")
		os.Exit(2)
	}

	if rawFilter != "" {
		filt = parseFilter(rawFilter)
	}

	filt.errsOnly = errsOnly
	if !helpersOnly {
		loadLogic(logicFileDir, filt)
	}

	loadHelpers(logicFileDir, filt, newFancy)
}

func loadHelpers(logicDir string, filt filter, f func() *FancyAstWriter) {
	contents, err := os.ReadFile(filepath.Join(logicDir, "LogicHelpers.json"))
	if err != nil {
		panic(err)
	}

	var helpers map[string]string

	if err := jsonc.Unmarshal(contents, &helpers); err != nil {
		panic(err)
	}

	for name, helper := range helpers {
		l := rules.NewLexer(name, helper)
		p := rules.NewParser(l)
		totalRule, err := p.ParseTotalRule()
		if err != nil {
			switch e := err.(type) {
			case rules.InvalidToken:
				helper = highlightParseError(int(e.Have.Pos), helper)
				break
			case rules.UnexpectedToken:
				helper = highlightParseError(int(e.Have.Pos), helper)
				break
			}

			fmt.Fprintf(os.Stdout, "Name:\t%s\n", name)
			fmt.Fprintf(os.Stdout, "FAILED TO PARSE: %s\n", helper)
			fmt.Fprintf(os.Stdout, "ERROR: %s\n", err.Error())
			continue
		}

		if !filt.errsOnly {
			v := f()
			totalRule.Rule.Visit(v)
			fmt.Fprintf(os.Stdout, "Name:\t%s\n", name)
			fmt.Fprintf(os.Stdout, "Helper:\n%s\n", v.b.String())
		}

	}
}

func loadLogic(logicDir string, filt filter) {
	entries, err := os.ReadDir(logicDir)
	if err != nil {
		panic(err)
	}

	for _, entry := range entries {
		if !strings.HasSuffix(entry.Name(), "json") || strings.Contains(entry.Name(), "Helpers.json") {
			continue
		}

		path := filepath.Join(logicDir, entry.Name())
		logicLocations, err := rules.ReadLogicFile(path)
		if err != nil {
			panic(err)
		}

		for _, logicLocation := range logicLocations {
			if !filt.MatchRegion(logicLocation.Region) {
				continue
			}
			ParseAllChecks(logicLocation, newFancy, filt)
		}
	}
}

func ParseAllChecks(loc rules.RawLogicLocation, f func() *FancyAstWriter, filt filter) {
	parseAll("Event", loc.Events, loc.Region, f, filt)
	parseAll("Check", loc.Locations, loc.Region, f, filt)
	parseAll("Exit", loc.Exits, loc.Region, f, filt)
}

func parseAll[E ~string, R ~string, M map[E]R, N ~string, F func() *FancyAstWriter](ctx string, m M, region N, f F, filt filter) {
	if !filt.MatchKind(ctx) {
		return
	}

	for check, rule := range m {
		if !filt.MatchSpecific(string(check)) {
			continue
		}
		name := fmt.Sprintf("%s: %s: %s", ctx, region, check)
		l := rules.NewLexer(name, string(rule))
		p := rules.NewParser(l)
		totalRule, err := p.ParseTotalRule()
		if err != nil {
			switch e := err.(type) {
			case rules.InvalidToken:
				rule = R(highlightParseError(int(e.Have.Pos), string(rule)))
				break
			case rules.UnexpectedToken:
				rule = R(highlightParseError(int(e.Have.Pos), string(rule)))
				break
			}

			fmt.Fprintf(os.Stdout, "Region:\t%s\nName:\t%s\nKind:\t%s\n", region, check, ctx)
			fmt.Fprintf(os.Stdout, "FAILED TO PARSE: %s\n", rule)
			fmt.Fprintf(os.Stdout, "ERROR: %s\n", err.Error())
			continue
		}

		if !filt.errsOnly {
			v := f()
			totalRule.Rule.Visit(v)
			fmt.Fprintf(os.Stdout, "Region:\t%s\nName:\t%s\nKind:\t%s\n", region, check, ctx)
			fmt.Fprintf(os.Stdout, "Rule:\n%s\n", v.b.String())
		}
	}
}

func highlightParseError(where int, rule string) string {
	newRule := &strings.Builder{}
	newRule.WriteString(rule)
	newRule.WriteRune('\n')
	newRule.WriteString(strings.Repeat(" ", where-5))
	newRule.WriteString(parseErrorColor.Paint(strings.Repeat("~", 4)))
	newRule.WriteString(parseErrorColor.Paint("^"))
	return newRule.String()
}
