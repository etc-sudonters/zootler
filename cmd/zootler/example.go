package main

import (
	"io/fs"
	"iter"
	"path/filepath"
	"regexp"
	"sudonters/zootler/internal"
	"sudonters/zootler/internal/app"
	"sudonters/zootler/internal/rules/ast"
	"sudonters/zootler/internal/rules/parser"
	"sudonters/zootler/internal/rules/visitor"
	"sudonters/zootler/internal/rules/zasm"

	"github.com/etc-sudonters/substrate/dontio"
	"github.com/etc-sudonters/substrate/slipup"
)

func panic_not_nil(err error) {
	if err != nil {
		panic(err)
	}
}

type fileloc struct {
	Checks map[string]string `json:"locations"`
	Exits  map[string]string `json:"exits"`
	Name   string            `json:"region_name"`
}

type AllTheRulesFrom struct {
	Path string
}

var (
	litTokenIdentRe = regexp.MustCompile("^[A-Z][A-Za-z_]+")
	litTokenStrRe   = regexp.MustCompile("^[A-Z][A-Za-z_ ]+") // allows spaces
	trickRe         = regexp.MustCompile("^logic_[a-z_]+")
)

type SpecialNameIdentifier ast.MacroDecls

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
			Callee: "load_setting_2",
			Args: []ast.Node{
				&ast.Identifier{
					Name: "tricks",
					Kind: ast.AST_IDENT_SET,
				},
				ident,
			},
			Macro: false,
		}
	}

	// if s.IsSettingName(name) ...
	return nil
}

type MacroExpander struct {
	Macros        ast.MacroDecls
	CurrentOrigin string
}

func (x MacroExpander) Expand(node *ast.Call, scope zasm.AssemblerScope) (zasm.Instructions, error) {
	decl, _ := x.Macros.Get(node.Callee)
	if !decl.IsBuiltIn() {
		return x.expandScript(node, decl, scope)
	}
	return nil, nil
}

func (x MacroExpander) expandScript(call *ast.Call, decl ast.MacroDecl, scope zasm.AssemblerScope) (zasm.Instructions, error) {
	argErr := scope.InParentScope(func() error {
		for idx, param := range decl.Params {
			scratch := zasm.Scratch()
			if err := scope.Assembler.AssembleInto(&scratch, call.Args[idx]); err != nil {
				return slipup.Describef(err, "macro arg %d", idx)
			}
			scope.Replacements[param] = scratch.Instr
		}
		return nil
	})
	if argErr != nil {
		return nil, argErr
	}

	replacement := zasm.Scratch()
	if err := scope.Assembler.AssembleInto(&replacement, decl.Body); err != nil {
		return nil, err
	}
	return replacement.Instr, nil
}

func (s SpecialNameIdentifier) Macro(name string) bool {
	return ast.MacroDecls(s).Exists(name)
}

func (r AllTheRulesFrom) Setup(z *app.Zootlr) error {
	ctx := z.Ctx()
	macros, bodies := load_json_macros(filepath.Join(r.Path, "..", "helpers.json"))
	xpander := MacroExpander{
		Macros: macros,
	}
	assembly := zasm.NewAssembly()
	assembler := zasm.Assembler{
		Data:   zasm.NewDataBuilder(),
		Macros: xpander,
		Functions: map[string]int{
			"has":                    2,
			"has_all_notes_for_song": 1,
			"has_hearts":             1,
			"has_medallions":         1,
			"has_stones":             1,
			"load_setting":           1,
			"load_setting_2":         2,
			"region_has_shortcuts":   1,
		},
	}

	macros.DeclareBuiltIn("at", []string{"target", "rule"})
	macros.DeclareBuiltIn("at_dampe_time", nil)
	macros.DeclareBuiltIn("at_day", nil)
	macros.DeclareBuiltIn("at_night", nil)
	macros.DeclareBuiltIn("had_night_start", nil)
	macros.DeclareBuiltIn("here", []string{"rule"})
	astGen := ast.NewGenerator(SpecialNameIdentifier(macros))
	initialize_script_macros(astGen, macros, bodies)

	filepath.Walk(r.Path, func(path string, info fs.FileInfo, err error) error {
		dontio.WriteLineOut(ctx, "at logic dir file %s", path)
		ext := filepath.Ext(path)
		if ext != ".json" {
			dontio.WriteLineOut(ctx, "skipping non json file: %s", path)
			return nil
		}

		locs, err := internal.ReadJsonFileAs[[]fileloc](path)
		panic_not_nil(err)

		for current, rule := range all_rules_from(locs) {
			dontio.WriteLineOut(ctx, rule)
			pt, parseErr := parser.Parse(rule)
			if parseErr != nil {
				return slipup.Describe(parseErr, "while parsing")
			}

			astGen.ClearStacks()
			ast, astErr := visitor.Transform(astGen, pt)
			if astErr != nil {
				return slipup.Describe(astErr, "pt -> ast on rule")
			}

			asmBlock, blockErr := assembly.Block(current)
			if blockErr != nil {
				panic(blockErr)
			}
			xpander.CurrentOrigin = current
			if err := assembler.AssembleInto(&asmBlock, ast); err != nil {
				panic(err)
			}
		}

		return nil
	})

	assembly.Load(assembler.Data)
	dontio.WriteLineOut(ctx, "%#v", assembly.Data)
	return nil
}

func all_rules_from(locs []fileloc) iter.Seq2[string, string] {
	return func(yield func(string, string) bool) {
		for _, loc := range locs {
			for _, rule := range loc.Checks {
				if !yield(loc.Name, rule) {
					return
				}
			}

			for _, rule := range loc.Exits {
				if !yield(loc.Name, rule) {
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
		ast, astErr := visitor.Transform[ast.Node](gen, body)
		if astErr != nil {
			panic(astErr)
		}
		decls.Initialize(which, ast)
	}
}
