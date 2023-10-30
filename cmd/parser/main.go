package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"sudonters/zootler/pkg/logic"
	"sudonters/zootler/pkg/rulesparser"

	"github.com/etc-sudonters/substrate/dontio"
	"muzzammil.xyz/jsonc"
)

var parseErrorColor dontio.ForegroundColor = 141

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
		loadHelpers(logicFileDir, filt, showHelpers, pretty)
	}
}

func loadHelpers(logicDir string, filt filter, display bool, pretty bool) {
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
		l := rulesparser.NewRulesLexer(helper)
		p := rulesparser.NewRulesParser(l)
		rule, err := p.Parse()
		if err != nil {
			fmt.Fprintf(os.Stdout, "Name:\t%s\n", name)
			fmt.Fprintf(os.Stdout, "FAILED TO PARSE: %s\n", helper)
			fmt.Fprintf(os.Stdout, "ERROR: %s\n", err.Error())
			continue
		}

		if display && !filt.errsOnly {
			fancy := newFancy()
			single := newSingleLine()
			rule.Visit(manyVisitors(fancy, single))
			fmt.Fprintf(os.Stdout, "Name:\t%s\n", name)
			fmt.Fprintf(os.Stdout, "Helper:\t%s\n%s\n", single.b.String(), fancy.b.String())
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
		logicLocations, err := logic.ReadLogicFile(path)
		if err != nil {
			panic(err)
		}

		for _, logicLocation := range logicLocations {
			if !filt.MatchRegion(logicLocation.Region) {
				continue
			}
			ParseAllChecks(logicLocation, filt, pretty)
		}
	}
}

func ParseAllChecks(loc logic.RawLogicLocation, filt filter, pretty bool) {
	parseAll("Event", loc.Events, loc.Region, filt, pretty)
	parseAll("Check", loc.Locations, loc.Region, filt, pretty)
	parseAll("Exit", loc.Exits, loc.Region, filt, pretty)
}

func parseAll[E ~string, R ~string, M map[E]R, N ~string](ctx string, m M, region N, filt filter, prettiness bool) {
	if !filt.MatchKind(ctx) {
		return
	}

	for check, rule := range m {
		if !filt.MatchSpecific(string(check)) {
			continue
		}
		rule = R(compressWhiteSpace(string(rule)))
		l := rulesparser.NewRulesLexer(string(rule))
		p := rulesparser.NewRulesParser(l)
		totalRule, err := p.Parse()
		if err != nil {
			fmt.Fprint(os.Stdout, "Failed to parse rule\n")
			fmt.Fprintf(os.Stdout, "Region:\t%s\nName:\t%s\nKind:\t%s\n", region, check, ctx)
			fmt.Fprintf(os.Stdout, "Raw:\t%s\n", rule)
			fmt.Fprintf(os.Stdout, "ERROR: %s\n\n", err.Error())
			continue
		}

		if !filt.errsOnly {
			fancy := newFancy()
			single := newSingleLine()
			totalRule.Visit(manyVisitors(fancy, single))
			fmt.Fprintf(os.Stdout, "Region:\t%s\nName:\t%s\nKind:\t%s\n", region, check, ctx)
			fmt.Fprintf(os.Stdout, "Raw:\t%s\n", rule)
			fmt.Fprintf(os.Stdout, "Rule:\t%s\n", single.b.String())
			if prettiness {
				fmt.Fprintf(os.Stdout, "%s\n", fancy.b.String())
			}
			fmt.Fprint(os.Stdout, "\n")
		}
	}
}

const errColor dontio.BackgroundColor = 210

func manyVisitors(v ...rulesparser.RuleVisitor) rulesparser.RuleVisitor {
	return manyVisit{v}
}

type manyVisit struct {
	visitors []rulesparser.RuleVisitor
}

func (m manyVisit) visit(n rulesparser.Expression) {
	for _, v := range m.visitors {
		if v == nil {
			continue
		}
		n.Visit(v)
	}
}

func (m manyVisit) VisitAttrAccess(n *rulesparser.AttrAccess) {
	m.visit(n)
}
func (m manyVisit) VisitBinOp(n *rulesparser.BinOp) {
	m.visit(n)
}
func (m manyVisit) VisitBoolOp(n *rulesparser.BoolOp) {
	m.visit(n)
}
func (m manyVisit) VisitBoolean(n *rulesparser.Boolean) {
	m.visit(n)
}
func (m manyVisit) VisitCall(n *rulesparser.Call) {
	m.visit(n)
}
func (m manyVisit) VisitIdentifier(n *rulesparser.Identifier) {
	m.visit(n)
}
func (m manyVisit) VisitNumber(n *rulesparser.Number) {
	m.visit(n)
}
func (m manyVisit) VisitString(n *rulesparser.String) {
	m.visit(n)
}
func (m manyVisit) VisitSubscript(n *rulesparser.Subscript) {
	m.visit(n)
}
func (m manyVisit) VisitTuple(n *rulesparser.Tuple) {
	m.visit(n)
}
func (m manyVisit) VisitUnary(n *rulesparser.UnaryOp) {
	m.visit(n)
}

func compressWhiteSpace(s string) string { return s }
