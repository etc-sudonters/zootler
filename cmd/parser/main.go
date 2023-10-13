package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/etc-sudonters/zootler/internal/console"
	"github.com/etc-sudonters/zootler/internal/rules"
	"muzzammil.xyz/jsonc"
)

var parseErrorColor console.ForegroundColor = 141
var compressWhiteSpaceRe *regexp.Regexp = regexp.MustCompile(`\s+`)
var leadWhiteSpace *regexp.Regexp = regexp.MustCompile(`^\s+`)
var trailWhiteSpace *regexp.Regexp = regexp.MustCompile(`\s+$`)

func main() {
	var logicFileDir string
	var errsOnly bool
	var rawFilter string
	var filt filter
	var showHelpers bool = false
	var pretty bool = false

	flag.StringVar(&logicFileDir, "l", "", "Directory where logic files are located")
	flag.StringVar(&rawFilter, "f", "", "Look it's complicated")
	flag.BoolVar(&errsOnly, "E", false, "set to only display errors")
	flag.BoolVar(&showHelpers, "H", false, "")
	flag.BoolVar(&pretty, "P", false, "")
	flag.Parse()

	if logicFileDir == "" {
		fmt.Fprint(os.Stderr, "-l is required")
		os.Exit(2)
	}

	if rawFilter != "" {
		filt = parseFilter(rawFilter)
	}

	filt.errsOnly = errsOnly
	loadLogic(logicFileDir, filt, pretty)

	if showHelpers {
		loadHelpers(logicFileDir, filt, newDual, showHelpers, pretty)
	}
}

func loadHelpers(logicDir string, filt filter, f func() *dualVisit, display bool, pretty bool) {
	contents, err := os.ReadFile(filepath.Join(logicDir, "LogicHelpers.json"))
	if err != nil {
		panic(err)
	}

	var helpers map[string]string

	if err := jsonc.Unmarshal(contents, &helpers); err != nil {
		panic(err)
	}

	for name, helper := range helpers {
		helper = compressWhiteSpace(helper)
		l := rules.NewLexer(name, helper)
		p := rules.NewParser(l)
		totalRule, err := p.ParseTotalRule()
		if err != nil {
			fmt.Fprintf(os.Stdout, "Name:\t%s\n", name)
			fmt.Fprintf(os.Stdout, "FAILED TO PARSE: %s\n", helper)
			fmt.Fprintf(os.Stdout, "ERROR: %s\n", err.Error())
			continue
		}

		if display && !filt.errsOnly {
			v := f()
			totalRule.Rule.Visit(v)
			pretty := v.pretty.b.String()
			single := v.single.b.String()
			fmt.Fprintf(os.Stdout, "Name:\t%s\n", name)
			fmt.Fprintf(os.Stdout, "Helper:\t%s\n%s\n", single, pretty)
		}

	}
}

func loadLogic(logicDir string, filt filter, pretty bool) {
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
			ParseAllChecks(logicLocation, newDual, filt, pretty)
		}
	}
}
func newDual() *dualVisit {
	return &dualVisit{
		pretty: newFancy(),
		single: newSingleLine(),
	}
}

func ParseAllChecks(loc rules.RawLogicLocation, f func() *dualVisit, filt filter, pretty bool) {
	parseAll("Event", loc.Events, loc.Region, f, filt, pretty)
	parseAll("Check", loc.Locations, loc.Region, f, filt, pretty)
	parseAll("Exit", loc.Exits, loc.Region, f, filt, pretty)
}

func parseAll[E ~string, R ~string, M map[E]R, N ~string, F func() *dualVisit](ctx string, m M, region N, f F, filt filter, prettiness bool) {
	if !filt.MatchKind(ctx) {
		return
	}

	for check, rule := range m {
		if !filt.MatchSpecific(string(check)) {
			continue
		}
		rule = R(compressWhiteSpace(string(rule)))
		name := fmt.Sprintf("%s: %s: %s", ctx, region, check)
		l := rules.NewLexer(name, string(rule))
		p := rules.NewParser(l)
		totalRule, err := p.ParseTotalRule()
		if err != nil {
			fmt.Fprint(os.Stdout, "Failed to parse rule\n")
			fmt.Fprintf(os.Stdout, "Region:\t%s\nName:\t%s\nKind:\t%s\n", region, check, ctx)
			fmt.Fprintf(os.Stdout, "Raw:\t%s\n", rule)
			fmt.Fprintf(os.Stdout, "ERROR: %s\n\n", err.Error())
			continue
		}

		if !filt.errsOnly {
			v := f()
			totalRule.Rule.Visit(v)
			single := v.single.b.String()
			fmt.Fprintf(os.Stdout, "Region:\t%s\nName:\t%s\nKind:\t%s\n", region, check, ctx)
			fmt.Fprintf(os.Stdout, "Raw:\t%s\n", rule)
			fmt.Fprintf(os.Stdout, "Rule:\t%s\n", single)
			if prettiness {
				pretty := v.pretty.b.String()
				fmt.Fprintf(os.Stdout, "%s\n", pretty)
			}
			fmt.Fprint(os.Stdout, "\n")
		}
	}
}

const errColor console.BackgroundColor = 210

type dualVisit struct {
	pretty *FancyAstWriter
	single *singleLine
}

func (d dualVisit) VisitAttrAccess(n *rules.AttrAccess) {
	d.pretty.VisitAttrAccess(n)
	d.single.VisitAttrAccess(n)
}
func (d dualVisit) VisitBinOp(n *rules.BinOp) {
	d.pretty.VisitBinOp(n)
	d.single.VisitBinOp(n)
}
func (d dualVisit) VisitBoolOp(n *rules.BoolOp) {
	d.pretty.VisitBoolOp(n)
	d.single.VisitBoolOp(n)
}
func (d dualVisit) VisitBoolean(n *rules.Boolean) {
	d.pretty.VisitBoolean(n)
	d.single.VisitBoolean(n)
}
func (d dualVisit) VisitCall(n *rules.Call) {
	d.pretty.VisitCall(n)
	d.single.VisitCall(n)
}
func (d dualVisit) VisitIdentifier(n *rules.Identifier) {
	d.pretty.VisitIdentifier(n)
	d.single.VisitIdentifier(n)
}
func (d dualVisit) VisitNumber(n *rules.Number) {
	d.pretty.VisitNumber(n)
	d.single.VisitNumber(n)
}
func (d dualVisit) VisitString(n *rules.String) {
	d.pretty.VisitString(n)
	d.single.VisitString(n)
}
func (d dualVisit) VisitSubscript(n *rules.Subscript) {
	d.pretty.VisitSubscript(n)
	d.single.VisitSubscript(n)
}
func (d dualVisit) VisitTuple(n *rules.Tuple) {
	d.pretty.VisitTuple(n)
	d.single.VisitTuple(n)
}
func (d dualVisit) VisitUnary(n *rules.UnaryOp) {
	d.pretty.VisitUnary(n)
	d.single.VisitUnary(n)
}

func compressWhiteSpace(r string) string {
	r = trailWhiteSpace.ReplaceAllLiteralString(r, "")
	r = leadWhiteSpace.ReplaceAllLiteralString(r, "")
	return compressWhiteSpaceRe.ReplaceAllLiteralString(r, " ")
}
