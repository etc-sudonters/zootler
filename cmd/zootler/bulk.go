package main

import (
	"io/fs"
	"iter"
	"path/filepath"
	"strings"
	"sudonters/zootler/internal"
	"sudonters/zootler/internal/rules/ast"
	"sudonters/zootler/internal/rules/parser"
	"sudonters/zootler/internal/rules/visitor"
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
	Aster *ast.AstGenerator
}

func (a AstAllRuleEdges) Iter(yield func(AstEdge) bool) {
	for edge := range a.AllEdgeRulesFrom.Iter {
		pt, ptErr := parser.Parse(edge.Rule)
		if ptErr != nil {
			panic(ptErr)
		}
		ast, err := visitor.Transform(a.Aster, pt)
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

func load_json_macros(path string) (ast.MacroDecls, map[string]parser.Expression) {
	helpers, err := internal.ReadJsonFileStringMap(path)
	if err != nil {
		panic(err)
	}

	macros := ast.NewMacros(len(helpers))
	macroBodies := make(map[string]parser.Expression, len(helpers))

	for decl, body := range helpers {
		decl := parser.MustParse(decl)
		switch decl.Type() {
		case parser.ExprCall:
			decl := parser.MustAssertAs[*parser.Call](decl)
			macro := ast.DeclareFromParseTree(macros, decl)

			macroBodies[macro.Name] = parser.MustParse(body)
			break
		case parser.ExprIdentifier:
			ident := parser.MustAssertAs[*parser.Identifier](decl)
			macro := macros.Declare(ident.Value, nil)
			macroBodies[macro.Name] = parser.MustParse(body)
			break
		default:
			panic("macro decl wasn't call-ish or ident-ish")
		}
	}

	return macros, macroBodies
}

func initialize_script_macros(gen *ast.AstGenerator, decls ast.MacroDecls, bodies map[string]parser.Expression) {
	for which, body := range bodies {
		gen.ClearStacks()
		ast, astErr := visitor.Transform(gen, body)
		if astErr != nil {
			panic(astErr)
		}
		decls.Initialize(which, ast)
	}
}
