package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/etc-sudonters/zootler/internal/rules"
)

func main() {
	var logicFileDir string

	flag.StringVar(&logicFileDir, "l", "", "Directory where logic files are located")
	flag.Parse()

	if logicFileDir == "" {
		fmt.Fprint(os.Stderr, "-l is required")
		os.Exit(2)
	}

	fmt.Fprintf(os.Stdout, "Reading logic from %s\n", logicFileDir)

	loadLogic(logicFileDir)
}

func loadLogic(logicDir string) {
	entries, err := os.ReadDir(logicDir)
	if err != nil {
		panic(err)
	}

	for _, entry := range entries {
		if !strings.HasSuffix(entry.Name(), "json") {
			fmt.Fprintf(os.Stdout, "discarding %s\n", entry.Name())
			continue
		}

		fmt.Fprintf(os.Stdout, "reading %s\n", entry.Name())
		path := filepath.Join(logicDir, entry.Name())
		logicLocations, err := rules.ReadLogicFile(path)
		if err != nil {
			panic(err)
		}

		for _, logicLocation := range logicLocations {
			err := ParseAllChecks(logicLocation, newFancy)
			if err != nil {
				fmt.Fprint(os.Stderr, err.Error())
			}
		}
	}
}

func ParseAllChecks(loc rules.RawLogicLocation, f func() *FancyAstWriter) error {
	var allErrs []error

	parseAll("Event", loc.Events, loc.Region, allErrs, f)
	parseAll("Check", loc.Locations, loc.Region, allErrs, f)
	parseAll("Exit", loc.Exits, loc.Region, allErrs, f)

	if len(allErrs) != 0 {
		return errors.Join(allErrs...)
	}

	return nil
}

func parseAll[E ~string, R ~string, M map[E]R, N ~string, F func() *FancyAstWriter](ctx string, m M, region N, allErrs []error, f F) {
	r := bufio.NewScanner(os.Stdin)
	for check, rule := range m {
		name := fmt.Sprintf("%s: %s: %s", ctx, region, check)
		l := rules.NewLexer(name, string(rule))
		p := rules.NewParser(l)

		rule, err := p.ParseTotalRule()
		if err != nil {
			allErrs = append(allErrs, fmt.Errorf("failed on %s: %w", name, err))
		}

		fmt.Fprintf(os.Stdout, "Region:\t%s\nName:\t%s\nKind:\t%s\n", region, check, ctx)
		v := f()
		rule.Rule.Visit(v)

		fmt.Fprintf(os.Stdout, "Rule:\n%s\n", v.b.String())

		r.Scan()
	}
}
