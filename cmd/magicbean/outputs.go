package main

import (
	"regexp"
	"sudonters/zootler/magicbeanvm/ast"
	"sudonters/zootler/magicbeanvm/symbols"
)

var funcName = regexp.MustCompile("[^A-Z]")

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

type symcount struct {
	kind  symbols.Kind
	count int
}

type findinvokes struct {
	symbols  *symbols.Table
	counting map[string]symcount
	has      map[string]int
}

func (this findinvokes) Invoke(node ast.Invoke, _ ast.Visiting) error {
	symbol := ast.LookUpNodeInTable(this.symbols, node.Target)

	if symbol != nil {
		switch symbol.Kind {
		case symbols.BUILT_IN:
			sym := this.counting[symbol.Name]
			this.counting[symbol.Name] = symcount{
				count: sym.count + 1,
				kind:  symbol.Kind,
			}
		}
		switch {
		case symbol.Name == "has":
			var name string
			qty, ok := node.Args[1].(ast.Number)
			if !ok {
				qty = 1
			}
			switch what := node.Args[0].(type) {
			case ast.Identifier:
				name = what.Symbol.Name
			case ast.String:
				symbol := this.symbols.LookUpByName(string(what))
				name = symbol.Name
			default:
				panic("unreachable...?")
			}
			this.has[name] = int(qty) + this.has[name]
		case symbol.Name == "has_every" || symbol.Name == "has_anyof":
			for i := range node.Args {
				var name string
				switch what := node.Args[i].(type) {
				case ast.Identifier:
					name = what.Symbol.Name
				case ast.String:
					symbol := this.symbols.LookUpByName(string(what))
					name = symbol.Name
				default:
					panic("unreachable...?")
				}
				this.has[name] = this.has[name] + 1
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
