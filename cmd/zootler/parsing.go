package main

import (
	"regexp"
	"sudonters/zootler/internal"
	"sudonters/zootler/internal/app"
	"sudonters/zootler/internal/rules/ast"
)

var (
	litTokenIdentRe = regexp.MustCompile("^[A-Z][A-Za-z_]+")
	litTokenStrRe   = regexp.MustCompile("^[A-Z][A-Za-z_ ]+") // allows spaces
	trickRe         = regexp.MustCompile("^logic_[a-z_]+")
)

type InstallParser struct {
	MacrosPath string
}

func (i InstallParser) Setup(z *app.Zootlr) error {
	macros, bodies := load_json_macros(i.MacrosPath)
	macros.DeclareBuiltIn("at", []string{"target", "rule"})
	macros.DeclareBuiltIn("at_dampe_time", nil)
	macros.DeclareBuiltIn("at_day", nil)
	macros.DeclareBuiltIn("at_night", nil)
	macros.DeclareBuiltIn("had_night_start", nil)
	macros.DeclareBuiltIn("here", []string{"rule"})

	sni := SpecialNameIdentifier{macros, make(map[string]struct{})}
	sni.DeclareSetting("shuffle_individual_ocarina_notes")
	astGen := ast.NewGenerator(&sni)
	initialize_script_macros(astGen, macros, bodies)

	z.AddResource(astGen)
	z.AddResource(macros)
	return nil
}

type SpecialNameIdentifier struct {
	macros   ast.MacroDecls
	settings map[string]struct{}
}

func (s SpecialNameIdentifier) Special(name string, inFuncCall bool) ast.Node {
	if litTokenIdentRe.MatchString(name) || litTokenStrRe.MatchString(name) {
		// TODO: ensure this is actually some kind of token by comparing to the store
		ident := &ast.Identifier{
			Name: string(internal.Normalize(name)),
			Kind: ast.AST_IDENT_TOK,
		}
		if inFuncCall {
			return ident
		}

		return &ast.Call{
			Callee: "has",
			Args: []ast.Node{
				ident,
				&ast.Literal{
					Value: float64(1),
					Kind:  ast.AST_LIT_NUM,
				}},
			Macro: false,
		}
	}

	if trickRe.MatchString(name) {
		ident := &ast.Identifier{
			Name: string(internal.Normalize(name)),
			Kind: ast.AST_IDENT_TRK,
		}
		if inFuncCall {
			return ident
		}

		return &ast.Call{
			Callee: "trick_enabled",
			Args:   []ast.Node{ident},
			Macro:  false,
		}
	}

	if s.IsSetting(name) {
		ident := &ast.Identifier{
			Name: string(internal.Normalize(name)),
			Kind: ast.AST_IDENT_SET,
		}
		if inFuncCall {
			return ident
		}

		return &ast.Call{
			Callee: "load_setting",
			Args:   []ast.Node{ident},
			Macro:  false,
		}
	}
	return nil
}

func (s SpecialNameIdentifier) DeclareSetting(name string) {
	if s.settings == nil {
		s.settings = make(map[string]struct{})
	}
	s.settings[name] = struct{}{}
}

func (s SpecialNameIdentifier) IsSetting(name string) bool {
	_, exists := s.settings[name]
	return exists
}

func (s *SpecialNameIdentifier) Macro(name string) bool {
	return s.macros.Exists(name)
}
