package bootstrap

import (
	"fmt"
	"regexp"
	"strings"
	"sudonters/zootler/cmd/magicbean/z16"
	"sudonters/zootler/internal/settings"
	"sudonters/zootler/magicbean"
	"sudonters/zootler/mido"
	"sudonters/zootler/mido/ast"
	"sudonters/zootler/mido/objects"
	"sudonters/zootler/mido/optimizer"
	"sudonters/zootler/mido/symbols"
	"sudonters/zootler/zecs"
)

func loadsymbols(ocm *zecs.Ocm, syms *symbols.Table, objs *objects.Builder) error {
	batches := []tagging{
		{kind: symbols.TOKEN, tag: objects.PtrToken, q: []zecs.BuildQuery{zecs.With[magicbean.Token]}},
		{kind: symbols.REGION, tag: objects.PtrRegion, q: []zecs.BuildQuery{zecs.With[magicbean.Region]}},
		{kind: symbols.TRANSIT, tag: objects.PtrTrans, q: []zecs.BuildQuery{zecs.With[magicbean.Connection]}},
		{kind: symbols.SCRIPTED_FUNC, tag: objects.PtrFunc, q: []zecs.BuildQuery{zecs.With[magicbean.ScriptDecl]}},
	}

	for _, batch := range batches {
		batch.tagall(ocm, syms, objs)
	}

	return nil
}

type tagging struct {
	kind symbols.Kind
	tag  uint16
	q    []zecs.BuildQuery
}

func (this tagging) tagall(ocm *zecs.Ocm, syms *symbols.Table, objs *objects.Builder) {
	q := ocm.Query()
	q.Build(zecs.Load[name], this.q...)
	rows, err := q.Execute()
	PanicWhenErr(err)
	for entity, tup := range rows.All {
		name := string(tup.Values[0].(name))
		symbol := syms.Declare(name, this.kind)
		ptr := objects.PackPtr32(this.tag, uint32(entity))
		objs.AssociateSymbol(symbol, ptr)
	}

}

func loadscripts(ocm *zecs.Ocm, env *mido.CompileEnv) error {
	q := ocm.Query()
	q.Build(
		zecs.Load[name],
		zecs.Load[magicbean.ScriptDecl],
		zecs.Load[magicbean.ScriptSource],
		zecs.WithOut[magicbean.RuleParsed],
	)

	rows, rowErr := q.Execute()
	PanicWhenErr(rowErr)
	decls := make(map[string]string, rows.Len())

	for _, tup := range rows.All {
		decl := tup.Values[1].(magicbean.ScriptDecl)
		body := tup.Values[2].(magicbean.ScriptSource)
		decls[string(decl)] = string(body)
	}

	PanicWhenErr(env.BuildScriptedFuncs(decls))

	eng := ocm.Engine()
	for entity, tup := range rows.All {
		name := tup.Values[0].(name)
		script, exists := env.ScriptedFuncs.Get(string(name))
		if !exists {
			panic(fmt.Errorf("somehow scripted func %s is missing, a mystery", name))
		}
		eng.SetValues(entity, zecs.Values{magicbean.ScriptParsed{script.Body}})
	}

	return nil
}

func aliassymbols(ocm *zecs.Ocm, syms *symbols.Table) error {
	q := ocm.Query()
	q.Build(zecs.Load[name], zecs.With[magicbean.Token])
	eng := ocm.Engine()
	rows, err := q.Execute()
	PanicWhenErr(err)

	for id, tup := range rows.All {
		name := string(tup.Values[0].(name))
		original := syms.LookUpByName(name)

		switch original.Kind {
		case symbols.FUNCTION, symbols.BUILT_IN_FUNCTION, symbols.COMPILER_FUNCTION, symbols.SCRIPTED_FUNC:
			continue
		case symbols.TOKEN:
			alias := escape(name)
			syms.Alias(original, alias)
			PanicWhenErr(eng.SetValues(id, zecs.Values{magicbean.AliasingName(alias)}))
		default:
			panic(fmt.Errorf("expected to only alias function or token: %s", original))
		}

	}

	return nil
}

func installCompilerFunctions(_ *settings.Zootr) mido.ConfigureCompiler {
	return func(env *mido.CompileEnv) {
		for i, name := range settings.Names() {
			symbol := env.Symbols.Declare(name, symbols.SETTING)
			env.Objects.AssociateSymbol(
				symbol,
				objects.PackPtr32(objects.PtrSetting, uint32(i)),
			)
		}

		mido.WithCompilerFunctions(func(*mido.CompileEnv) optimizer.CompilerFunctionTable {
			return optimizer.CompilerFunctionTable{
				"load_setting":           constCompileFunc,
				"load_setting_2":         constCompileFunc,
				"compare_setting":        constCompileFunc,
				"region_has_shortcuts":   constCompileFunc,
				"is_trick_enabled":       constCompileFunc,
				"had_night_start":        constCompileFunc,
				"has_all_notes_for_song": constCompileFunc,
				"at_dampe_time":          constCompileFunc,
				"at_day":                 constCompileFunc,
				"at_night":               constCompileFunc,
			}
		})(env)
	}
}

func installConnectionGenerator(ocm *zecs.Ocm) mido.ConfigureCompiler {
	return func(env *mido.CompileEnv) {
		env.Optimize.AddOptimizer(func(ce *mido.CompileEnv) ast.Rewriter {
			var conngen ConnectionGenerator
			conngen.Nodes = z16.NewNodes(ocm)
			conngen.Tokens = z16.NewTokens(ocm)
			conngen.Symbols = ce.Symbols
			conngen.Objects = ce.Objects

			return optimizer.NewConnectionGeneration(ce.Optimize.Context, ce.Symbols, conngen)
		})

	}
}

var escaping = regexp.MustCompile("['()[\\]-]")

func escape(name string) string {
	name = escaping.ReplaceAllLiteralString(name, "")
	return strings.ReplaceAll(name, " ", "_")
}

type ConnectionGenerator struct {
	Nodes   z16.Nodes
	Tokens  z16.Tokens
	Symbols *symbols.Table
	Objects *objects.Builder
}

func (this ConnectionGenerator) AddConnectionTo(region string, rule ast.Node) (*symbols.Sym, error) {
	hash := ast.Hash(rule)
	suffix := fmt.Sprintf("#%s#%32x", region, hash)
	tokenName := magicbean.NameF("Token%s", suffix)

	if symbol := this.Symbols.LookUpByName(string(tokenName)); symbol != nil {
		return symbol, nil
	}

	token := this.Tokens.Named(tokenName)
	placement := this.Nodes.Placement(magicbean.NameF("Place%s", suffix))

	placement.Owns(token)
	ptr := objects.PackPtr32(objects.PtrToken, uint32(token.Entity()))
	token.Attach(magicbean.Event{}, ptr)

	node := this.Nodes.Region(magicbean.Name(region))
	edge := node.Has(placement)
	edge.Attach(magicbean.RuleParsed{rule})

	symbol := this.Symbols.Declare(string(tokenName), symbols.TOKEN)
	this.Objects.AssociateSymbol(symbol, ptr)
	return symbol, nil
}

func constCompileFunc([]ast.Node, ast.Rewriting) (ast.Node, error) {
	return ast.Boolean(true), nil
}
