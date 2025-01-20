package bootstrap

import (
	"sudonters/zootler/internal/query"
	"sudonters/zootler/internal/table"
	"sudonters/zootler/magicbean"
	"sudonters/zootler/mido"
	"sudonters/zootler/mido/optimizer"
	"sudonters/zootler/zecs"
)

func parseall(ocm *zecs.Ocm, codegen *mido.CodeGen) error {
	q := ocm.Query()
	q.Build(
		zecs.Load[magicbean.RuleSource],
		zecs.With[magicbean.Connection],
		zecs.WithOut[magicbean.RuleParsed],
	)

	for ent, tup := range q.Rows {
		entity := ocm.Proxy(ent)
		source := tup.Values[0].(magicbean.RuleSource)

		parsed, err := codegen.Parse(string(source))
		PanicWhenErr(err)
		entity.Attach(magicbean.RuleParsed{parsed})
	}

	return nil
}

func optimizeall(ocm *zecs.Ocm, codegen *mido.CodeGen) error {
	eng := ocm.Engine()
	unoptimized := ocm.Query()
	unoptimized.Build(
		zecs.Load[magicbean.RuleParsed],
		zecs.Load[magicbean.Connection],
		zecs.WithOut[magicbean.RuleOptimized],
	)

	for {
		rows, err := unoptimized.Execute()
		PanicWhenErr(err)
		if rows.Len() == 0 {
			break
		}

		for ent, tup := range rows.All {
			entity := ocm.Proxy(ent)
			parsed := tup.Values[0].(magicbean.RuleParsed)
			edge := tup.Values[1].(magicbean.Connection)

			parent, parentErr := eng.GetValues(
				edge.From, table.ColumnIds{
					query.MustAsColumnId[magicbean.Name](eng),
				},
			)
			PanicWhenErr(parentErr)
			optimizer.SetCurrentLocation(codegen.Context, string(parent.Values[0].(magicbean.Name)))
			optimized, optimizeErr := codegen.Optimize(parsed.Node)
			PanicWhenErr(optimizeErr)
			entity.Attach(magicbean.RuleOptimized{optimized})
		}
	}

	return nil
}

func compileall(ocm *zecs.Ocm, codegen *mido.CodeGen) error {
	uncompiled := ocm.Query()
	uncompiled.Build(
		zecs.Load[magicbean.RuleOptimized],
		zecs.With[magicbean.Connection],
		zecs.WithOut[magicbean.RuleCompiled],
	)

	for ent, tup := range uncompiled.Rows {
		entity := ocm.Proxy(ent)
		compiling := tup.Values[0].(magicbean.RuleOptimized)

		bytecode, err := codegen.Compile(compiling.Node)
		PanicWhenErr(err)
		entity.Attach(magicbean.RuleCompiled(bytecode))
	}

	return nil
}
