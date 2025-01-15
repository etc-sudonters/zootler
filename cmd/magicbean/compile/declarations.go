package compile

import (
	"regexp"
	"strings"
	"sudonters/zootler/cmd/magicbean/z2"
	"sudonters/zootler/internal/query"
	"sudonters/zootler/internal/table"
	"sudonters/zootler/mido"
	"sudonters/zootler/mido/ast"
	"sudonters/zootler/mido/objects"
	"sudonters/zootler/mido/symbols"
)

func AliasTokens(symbols *symbols.Table, funcs *ast.PartialFunctionTable, names []string) {
	for _, name := range names {
		escaped := escape(name)
		if _, exists := funcs.Get(escaped); exists {
			continue
		}
		if _, exists := funcs.Get(name); exists {
			continue
		}
		symbol := symbols.LookUpByName(name)
		symbols.Alias(symbol, escaped)
	}
}

var escaping = regexp.MustCompile("['()[\\]-]")

func escape(name string) string {
	name = escaping.ReplaceAllLiteralString(name, "")
	return strings.ReplaceAll(name, " ", "_")
}

func DeclareCompilerSymbolsFrom(engine query.Engine) mido.ConfigureCompiler {
	return func(env *mido.CompileEnv) {
		declaresymbols[z2.Collectable](engine, env.Symbols, env.Objects, symbols.TOKEN, objects.PtrToken)
		declaresymbols[z2.Location](engine, env.Symbols, env.Objects, symbols.LOCATION, objects.PtrLoc)
	}
}

func declaresymbols[T table.Value](engine query.Engine, syms *symbols.Table, objs *objects.Builder, kind symbols.Kind, tag objects.PtrTag) {
	q := z2.CreateQuery(engine)
	q.Build(
		z2.QueryWith[T],
		z2.QueryLoad[z2.Name],
	)

	rows, err := q.Execute()
	if err != nil {
		panic(err)
	}

	for id, tup := range rows.All {
		name := string(tup.Values[0].(z2.Name))
		symbol := syms.Declare(name, kind)
		objs.AssociateSymbol(symbol, objects.PackTaggedPtr32(tag, uint32(id)))
	}
}
