package main

import (
	"context"
	"errors"
	"sudonters/zootler/internal/query"
	"sudonters/zootler/internal/slipup"
	"sudonters/zootler/internal/table"
	"sudonters/zootler/pkg/rules/parser"
	"sudonters/zootler/pkg/world/components"

	"github.com/etc-sudonters/substrate/mirrors"
)

type LogicCompiler struct{}

func (l *LogicCompiler) Configure(ctx context.Context, e query.Engine) error {
	var edge ParsableEdge
	q := e.CreateQuery()
	(&edge).Query(q)
	edgeRules, retrieveErr := e.Retrieve(q)
	if retrieveErr != nil {
		return slipup.Trace(retrieveErr, "while preparing to compile logic")
	}

	if edgeRules.Len() == 0 {
		return errors.New("did not find any logic rules to compile")
	}

	for edgeRules.MoveNext() {
		if initErr := (&edge).Init(edgeRules.Current()); initErr != nil {
			return slipup.Trace(initErr, "while initing edge rules")
		}

		if _, parseErr := parser.Parse(string(edge.RawLogic.Rule)); parseErr != nil {
			return slipup.TraceMsg(parseErr, "while parsing rule for '%s'", edge.Name)
		}
	}

	return nil
}

type ParsableEdge struct {
	Name     components.Name
	Edge     components.Edge
	RawLogic components.RawLogic
}

type Archetype interface {
	Init(*table.RowTuple)
	Query(query.Query)
}

func (p *ParsableEdge) Init(vs *table.RowTuple) error {
	p.Name = vs.Values[0].(components.Name)
	p.Edge = vs.Values[2].(components.Edge)
	p.RawLogic = vs.Values[1].(components.RawLogic)
	return nil
}

func (p *ParsableEdge) Query(q query.Query) {
	q.Load(mirrors.TypeOf[components.Name]())
	q.Load(mirrors.TypeOf[components.Edge]())
	q.Load(mirrors.TypeOf[components.RawLogic]())
}
