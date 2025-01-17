package bootstrap

import (
	"fmt"
	"regexp"
	"sudonters/zootler/internal/settings"
	"sudonters/zootler/magicbean"
	"sudonters/zootler/mido"
	"sudonters/zootler/mido/ast"
	"sudonters/zootler/mido/objects"
	"sudonters/zootler/mido/symbols"
	"sudonters/zootler/zecs"
)

func loadsymbols(ocm *zecs.Ocm, syms *symbols.Table, objs *objects.Builder) error {
	batches := []tagging{
		{kind: symbols.TOKEN, tag: objects.PtrToken, q: []zecs.BuildQuery{zecs.With[magicbean.Token]}},
		{kind: symbols.REGION, tag: objects.PtrRegion, q: []zecs.BuildQuery{zecs.With[magicbean.Region]}},
		{kind: symbols.PLACEMENT, tag: objects.PtrPlace, q: []zecs.BuildQuery{zecs.With[magicbean.Placement]}},
		{kind: symbols.TRANSIT, tag: objects.PtrTrans, q: []zecs.BuildQuery{zecs.With[magicbean.Transit]}},
	}

	for _, batch := range batches {
		batch.tagall(ocm, syms, objs)
	}

	return nil
}

type tagging struct {
	kind symbols.Kind
	tag  objects.PtrTag
	q    []zecs.BuildQuery
}

func (this tagging) tagall(ocm *zecs.Ocm, syms *symbols.Table, objs *objects.Builder) {
	q := ocm.Query()
	q.Build(zecs.Load[name], this.q...)
	for entity, tup := range q.Rows {
		name := string(tup.Values[0].(name))
		symbol := syms.Declare(name, this.kind)
		objs.AssociateSymbol(symbol, objects.PackTaggedPtr32(this.tag, uint32(entity)))
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
	panicWhenErr(rowErr)
	decls := make(map[string]string, rows.Len())

	for _, tup := range rows.All {
		decl := tup.Values[1].(magicbean.ScriptDecl)
		body := tup.Values[2].(magicbean.ScriptSource)
		decls[string(decl)] = string(body)
	}

	panicWhenErr(env.BuildScriptedFuncs(decls))

	eng := ocm.Engine()
	for entity, tup := range rows.All {
		name := tup.Values[0].(name)
		script, exists := env.ScriptedFuncs.Get(string(name))
		if !exists {
			panic(fmt.Errorf("somehow scripted func %s is missing, a mystery", name))
		}
		eng.SetValues(entity, zecs.Values{magicbean.ScriptParsed(script.Body)})
	}

	return nil
}

func aliassymbols(ocm *zecs.Ocm, syms *symbols.Table, scripts *ast.ScriptedFunctions) error {
	q := ocm.Query()
	q.Build(zecs.Load[name], zecs.With[magicbean.Token])
	eng := ocm.Engine()

	for id, tup := range q.Rows {
		name := string(tup.Values[0].(name))
		alias := escape(name)

		if _, exists := scripts.Get(alias); exists {
			continue
		}
		if _, exists := scripts.Get(name); exists {
			continue
		}

		original := syms.LookUpByName(name)
		syms.Alias(original, alias)
		panicWhenErr(eng.SetValues(id, zecs.Values{magicbean.AliasingName(alias)}))
	}

	return nil
}

func installSettings(_ *settings.Zootr) mido.ConfigureCompiler {
	return func(_ *mido.CompileEnv) {}
}

func installConnectionGenerator(_ *zecs.Ocm) mido.ConfigureCompiler {
	return func(_ *mido.CompileEnv) {}
}

var escaping = regexp.MustCompile("['()[\\]-]")

func escape(name string) string {
	name = escaping.ReplaceAllLiteralString(name, "")
	return strings.ReplaceAll(name, " ", "_")
}
