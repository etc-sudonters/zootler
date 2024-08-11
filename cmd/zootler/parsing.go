package main

import (
	"context"
	"errors"
	"sudonters/zootler/internal/components"
	"sudonters/zootler/internal/query"
	"sudonters/zootler/internal/rules/compiler"
	"sudonters/zootler/internal/rules/parser"
	"sudonters/zootler/internal/slipup"
	"sudonters/zootler/internal/table"

	"github.com/etc-sudonters/substrate/mirrors"
)

type LogicCompiler struct{}

func (l *LogicCompiler) Setup(ctx context.Context, e query.Engine) error {
	var compiledEdges int
	defer func() {
		WriteLineOut(ctx, "compiled %d rules", compiledEdges)
	}()

	edge := new(ParsableEdge)
	q := e.CreateQuery()
	q.Load(mirrors.TypeOf[components.Name]())
	q.Load(mirrors.TypeOf[components.RawLogic]())
	edgeRules, retrieveErr := e.Retrieve(q)
	if retrieveErr != nil {
		return slipup.Trace(retrieveErr, "while preparing to compile logic")
	}

	if edgeRules.Len() == 0 {
		return errors.New("did not find any logic rules to compile")
	}

	WriteLineOut(ctx, "found %d rules to compile", edgeRules.Len())
	for edgeRules.MoveNext() {
		compiledEdges++
		current := edgeRules.Current()
		if initErr := edge.Init(current); initErr != nil {
			return slipup.Trace(initErr, "while initing edge rules")
		}

		ast, parseErr := parser.Parse(string(edge.RawLogic.Rule))
		if parseErr != nil {
			return slipup.TraceMsg(parseErr, "while parsing rule for '%s'", edge.Name)
		}

		bc, compileErr := compiler.Compile(ast)
		if compileErr != nil {
			return slipup.TraceMsg(compileErr, "while compiling rule for '%s'", edge.Name)
		}

		if setErr := e.SetValues(current.Id, table.Values{components.CompiledRule{Bytecode: bc}}); setErr != nil {
			return setErr
		}
	}

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
