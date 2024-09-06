package main

import (
	"io/fs"
	"iter"
	"path/filepath"
	"strings"
	"sudonters/zootler/internal"
	"sudonters/zootler/internal/rules/ast"
	"sudonters/zootler/internal/rules/parser"
)

type AstEdge struct {
	RuleEdge
	Ast ast.Node
}

type RuleEdge struct {
	Origin, Dest, Rule, Kind string
}

type fileloc struct {
	Checks map[string]string `json:"locations"`
	Exits  map[string]string `json:"exits"`
	Name   string            `json:"region_name"`
}

type AstAllRuleEdges struct {
	AllEdgeRulesFrom
}

func (a AstAllRuleEdges) Iter(yield func(AstEdge) bool) {
	var aster ast.Ast
	for edge := range a.AllEdgeRulesFrom.Iter {
		pt, ptErr := parser.Parse(edge.Rule)
		if ptErr != nil {
			panic(ptErr)
		}
		ast, err := parser.Transform(&aster, pt)
		if err != nil {
			panic(err)
		}

		if !yield(AstEdge{edge, ast}) {
			return
		}
	}
}

type AllEdgeRulesFrom struct {
	Path string
}

func (a AllEdgeRulesFrom) Iter(yield func(RuleEdge) bool) {
	filepath.Walk(a.Path, func(path string, info fs.FileInfo, err error) error {
		ext := filepath.Ext(path)
		if ext != ".json" {
			return nil
		}

		locs, err := internal.ReadJsonFileAs[[]fileloc](path)
		if err != nil {
			return err
		}

		for edge := range all_rules_from(locs) {
			if !yield(edge) {
				return nil
			}
		}

		return nil
	})
}

func all_rules_from(locs []fileloc) iter.Seq[RuleEdge] {
	return func(yield func(RuleEdge) bool) {
		var edge RuleEdge
		for _, loc := range locs {
			edge.Origin = loc.Name
			edge.Kind = "CHECK"
			for name, rule := range loc.Checks {
				edge.Dest = name
				edge.Rule = strings.Trim(rule, " \n")
				if !yield(edge) {
					return
				}
			}

			edge.Kind = "EXIT"
			for name, rule := range loc.Exits {
				edge.Dest = name
				edge.Rule = strings.Trim(rule, " \n")
				if !yield(edge) {
					return
				}
			}
		}
	}
}
