package main

import (
	"regexp"
	"sudonters/zootler/magicbeanvm/ast"
	"sudonters/zootler/magicbeanvm/symbols"
)

var funcName = regexp.MustCompile("[^A-Z]")

type findinvokes struct {
	tbl   *symbols.Table
	found map[string]symbols.Kind
}

func (this findinvokes) Invoke(node ast.Invoke, _ ast.Visiting) error {
	symbol := ast.LookUpNodeInTable(this.tbl, node.Target)

	if symbol != nil && funcName.MatchString(symbol.Name) {
		switch symbol.Kind {
		case symbols.UNKNOWN, symbols.BUILT_IN:
			this.found[symbol.Name] = symbol.Kind
		}
	}

	return nil
}

type findstrings map[string]struct{}

func (this findstrings) String(node ast.String, _ ast.Visiting) error {
	this[string(node)] = struct{}{}
	return nil
}
