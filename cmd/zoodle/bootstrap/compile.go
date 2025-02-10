package bootstrap

import (
	"sudonters/libzootr/components"
	"sudonters/libzootr/internal/query"
	"sudonters/libzootr/internal/table"
	"sudonters/libzootr/mido"
	"sudonters/libzootr/mido/optimizer"
	"sudonters/libzootr/zecs"
)

func parseall(ocm *zecs.Ocm, codegen *mido.CodeGen) error {
	q := ocm.Query()
	q.Build(
		zecs.Load[components.RuleSource],
		zecs.With[components.Connection],
		zecs.WithOut[components.RuleParsed],
	)

	for ent, tup := range q.Rows {
		entity := ocm.Proxy(ent)
		source := tup.Values[0].(components.RuleSource)

		parsed, err := codegen.Parse(string(source))
		PanicWhenErr(err)
		entity.Attach(components.RuleParsed{parsed})
	}

	return nil
}

func optimizeall(ocm *zecs.Ocm, codegen *mido.CodeGen) error {
	eng := ocm.Engine()
	unoptimized := ocm.Query()
	unoptimized.Build(
		zecs.Load[components.RuleParsed],
		zecs.Load[components.Connection],
		zecs.WithOut[components.RuleOptimized],
	)

	for {
		rows, err := unoptimized.Execute()
		PanicWhenErr(err)
		if rows.Len() == 0 {
			break
		}

		for ent, tup := range rows.All {
			entity := ocm.Proxy(ent)
			parsed := tup.Values[0].(components.RuleParsed)
			edge := tup.Values[1].(components.Connection)

			parent, parentErr := eng.GetValues(
				edge.From, table.ColumnIds{
					query.MustAsColumnId[components.Name](eng),
				},
			)
			PanicWhenErr(parentErr)
			optimizer.SetCurrentLocation(codegen.Context, string(parent.Values[0].(components.Name)))
			optimized, optimizeErr := codegen.Optimize(parsed.Node)
			PanicWhenErr(optimizeErr)
			entity.Attach(components.RuleOptimized{optimized})
		}
	}

	return nil
}

func compileall(ocm *zecs.Ocm, codegen *mido.CodeGen) error {
	uncompiled := ocm.Query()
	uncompiled.Build(
		zecs.Load[components.RuleOptimized],
		zecs.With[components.Connection],
		zecs.WithOut[components.RuleCompiled],
	)

	for ent, tup := range uncompiled.Rows {
		entity := ocm.Proxy(ent)
		compiling := tup.Values[0].(components.RuleOptimized)

		bytecode, err := codegen.Compile(compiling.Node)
		PanicWhenErr(err)
		entity.Attach(components.RuleCompiled(bytecode))
	}

	return nil
}
