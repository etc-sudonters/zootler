package main

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
	"sudonters/zootler/magicbeanvm"
	"sudonters/zootler/magicbeanvm/ast"
	"sudonters/zootler/magicbeanvm/code"
	"sudonters/zootler/magicbeanvm/symbols"
)

var funcName = regexp.MustCompile("[^A-Z]")

func DisassembleAll(compiled []magicbeanvm.CompiledSource) {
	slices.SortFunc(compiled, func(a, b magicbeanvm.CompiledSource) int {
		switch cmp := strings.Compare(a.OriginatingRegion, b.OriginatingRegion); cmp {
		case 0:
			return strings.Compare(a.Destination, b.Destination)
		default:
			return cmp
		}
	})

	for _, compiled := range compiled {
		fmt.Println()
		fmt.Printf("%q: %q\n", compiled.OriginatingRegion, compiled.Destination)
		fmt.Println(code.Disassemble(compiled.ByteCode.Tape))
	}
	fmt.Println()
	fmt.Printf("Disassembled %06d\n", len(compiled))
	fmt.Println()

}

func SymbolReport(symbolTable *symbols.Table) {
	size, total, aliased := symbolTable.Size(), symbolTable.RawSize(), symbolTable.AliasCount()
	fmt.Println("Symbol Report")
	fmt.Printf("ALIAS: %04d %04X\n", aliased, aliased)
	fmt.Printf("COUNT: %04d %04X\n", size, size)
	fmt.Printf("TOTAL: %04d %04X\n", total, total)
	fmt.Println()
	fmt.Println("KINDS")
	counts := make(map[symbols.Kind]int)
	for symbol := range symbolTable.RawAll {
		counts[symbol.Kind] = counts[symbol.Kind] + 1
	}
	for kind, count := range counts {
		fmt.Printf("%04d %s\n", count, kind)
	}
	fmt.Println()

}

type analysis struct {
	nodes   map[ast.Kind]int
	invokes map[string]invokecount

	invokeFinder *findinvokes
}

func newanalysis() analysis {
	return analysis{
		make(map[ast.Kind]int),
		make(map[string]invokecount),
		nil,
	}
}

func (this analysis) Report() {
	fmt.Println(this.String())
}

func (this analysis) String() string {
	var str strings.Builder
	fmt.Fprintln(&str, "INVOKE TOTALS")
	for name, item := range this.invokes {
		fmt.Fprintf(&str, "%06d\t%s\t\t%s\n", item.count, item.kind, name)
	}
	fmt.Fprintln(&str)

	fmt.Fprintln(&str, "NODE TOTALS")
	for kind, count := range this.nodes {
		fmt.Fprintf(&str, "%06d\t%s\n", count, kind)
	}
	fmt.Fprintln(&str)

	return str.String()
}

func (this analysis) register(env *magicbeanvm.CompilationEnvironment) {
	env.Analysis.PostOptimize(func(env *magicbeanvm.CompilationEnvironment) ast.Visitor {
		finder := findinvokes{env.Symbols, this.invokes}
		return ast.Visitor{Invoke: finder.Invoke}
	})
	env.Analysis.PostOptimize(func(env *magicbeanvm.CompilationEnvironment) ast.Visitor {
		return countnodes(this.nodes)
	})
}

func countnodes(counter nodecounter) ast.Visitor {
	return ast.Visitor{
		AnyOf:      thenVisit(counter.tick, ast.VisitAnyOf),
		Boolean:    thenVisit(counter.tick, ast.VisitBoolean),
		Compare:    thenVisit(counter.tick, ast.VisitCompare),
		Every:      thenVisit(counter.tick, ast.VisitEvery),
		Identifier: thenVisit(counter.tick, ast.VisitIdentifier),
		Invert:     thenVisit(counter.tick, ast.VisitInvert),
		Invoke:     thenVisit(counter.tick, ast.VisitInvoke),
		Number:     thenVisit(counter.tick, ast.VisitNumber),
		String:     thenVisit(counter.tick, ast.VisitString),
	}
}

func thenVisit[N ast.Node](do func(ast.Node), then ast.VisitFunc[N]) ast.VisitFunc[N] {
	return func(node N, visit ast.Visiting) error {
		do(node)
		return then(node, visit)
	}
}

type nodecounter map[ast.Kind]int

func (c nodecounter) tick(node ast.Node) {
	which := node.Kind()
	count := c[which]
	c[which] = count + 1
}

type invokecount struct {
	kind  symbols.Kind
	count int
}

type findinvokes struct {
	symbols  *symbols.Table
	counting map[string]invokecount
}

func (this findinvokes) Invoke(node ast.Invoke, _ ast.Visiting) error {
	symbol := ast.LookUpNodeInTable(this.symbols, node.Target)

	if symbol != nil {
		switch symbol.Kind {
		case symbols.BUILT_IN, symbols.COMPILED_FUNC, symbols.COMP_TIME, symbols.FUNCTION:
			sym := this.counting[symbol.Name]
			this.counting[symbol.Name] = invokecount{
				count: sym.count + 1,
				kind:  symbol.Kind,
			}
		}
	}

	return nil
}

type findstrings map[string]struct{}

func (this findstrings) String(node ast.String, _ ast.Visiting) error {
	this[string(node)] = struct{}{}
	return nil
}
