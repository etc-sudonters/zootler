package main

import (
	"context"
	"errors"
	"sudonters/zootler/internal/components"
	"sudonters/zootler/internal/query"
	"sudonters/zootler/internal/rules/parser"
	"sudonters/zootler/internal/rules/runtime"
	"sudonters/zootler/internal/slipup"
	"sudonters/zootler/internal/table"

	"github.com/etc-sudonters/substrate/mirrors"
)

type LogicCompiler struct {
	compiled, failed uint64
}

func (l *LogicCompiler) Setup(ctx context.Context, e query.Engine) error {
	edge := new(ParsableEdge)
	q := e.CreateQuery()
	q.Load(mirrors.TypeOf[components.Name]())
	q.Load(mirrors.TypeOf[components.RawLogic]())
	edgeRules, retrieveErr := e.Retrieve(q)
	if retrieveErr != nil {
		return slipup.Describe(retrieveErr, "while preparing to compile logic")
	}

	if edgeRules.Len() == 0 {
		return errors.New("did not find any logic rules to compile")
	}

	for edgeRules.MoveNext() {
		current := edgeRules.Current()
		if initErr := edge.Init(current); initErr != nil {
			return slipup.Describe(initErr, "while initing edge rules")
		}

		ast, parseErr := parser.Parse(string(edge.RawLogic.Rule))
		if parseErr != nil {
			return slipup.Describef(parseErr, "while parsing rule for '%s'", edge.Name)
		}

		if bc, compileErr := runtime.Compile(ast); compileErr != nil {
			l.failed++
			WriteLineOut(ctx, "could not compile rule at '%s': %s", edge.Name, compileErr.Error())
		} else {
			l.compiled++
			if setErr := e.SetValues(current.Id, table.Values{components.CompiledRule{Bytecode: *bc}}); setErr != nil {
				return setErr
			}
		}
	}

	WriteLineOut(ctx, "compiled %d rules\nfailed %d rules", l.compiled, l.failed)
	return nil
}

type ParsableEdge struct {
	Name     components.Name
	RawLogic components.RawLogic
}

func (p *ParsableEdge) Init(vs *table.RowTuple) error {
	p.Name = vs.Values[0].(components.Name)
	p.RawLogic = vs.Values[1].(components.RawLogic)
	return nil
}
