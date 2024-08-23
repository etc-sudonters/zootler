package main

import (
	"errors"
	"fmt"
	"strings"
	"sudonters/zootler/internal/app"
	"sudonters/zootler/internal/components"
	"sudonters/zootler/internal/query"
	"sudonters/zootler/internal/rules/parser"
	"sudonters/zootler/internal/rules/runtime"
	"sudonters/zootler/internal/table"

	"github.com/etc-sudonters/substrate/dontio"
	"github.com/etc-sudonters/substrate/slipup"
)

type LogicCompiler struct {
	compiled, failed uint64
}

func (l *LogicCompiler) Setup(z *app.Zootlr) error {
	ctx := z.Ctx()
	e := z.Engine()
	compiler := app.GetResource[runtime.Compiler](z)
	if compiler == nil {
		return slipup.Createf("expected compiler resource to be available")
	}
	edge := new(ParsableEdge)
	q := e.CreateQuery()
	q.Load(query.MustAsColumnId[components.Name](e))
	q.Load(query.MustAsColumnId[components.RawLogic](e))
	edgeRules, retrieveErr := e.Retrieve(q)
	if retrieveErr != nil {
		return slipup.Describe(retrieveErr, "while preparing to compile logic")
	}

	if edgeRules.Len() == 0 {
		return errors.New("did not find any logic rules to compile")
	}

    allIdentitifers := make(identifiers, 256)
	for id, tup := range edgeRules.All {
		if initErr := edge.Init(id, tup); initErr != nil {
			return slipup.Describe(initErr, "while initing edge rules")
		}

		ast, parseErr := parser.Parse(string(edge.RawLogic.Rule))
		if parseErr != nil {
			return slipup.Describef(parseErr, "while parsing rule for '%s'", edge.Name)
		}

        if visitErr := parser.Visit(allIdentitifers, ast); visitErr != nil {
            return slipup.Describef(visitErr, "while gathering identifiers from '%s' rule: '%s'", edge.Name, edge.RawLogic.Rule)
        }

		if bc, compileErr := runtime.CompileEdgeRule(&compiler.Res, ast); compileErr != nil {
			l.failed++
			dontio.WriteLineOut(ctx, "could not compile rule at '%s': %s", edge.Name, compileErr.Error())
		} else {
			l.compiled++
			if setErr := e.SetValues(id, table.Values{components.CompiledRule{Bytecode: *bc}}); setErr != nil {
				return setErr
			}
		}
	}

    identifierString := &strings.Builder{}

    for id := range allIdentitifers {
        fmt.Fprintf(identifierString, "%s\n", id)
    }

	dontio.WriteLineOut(ctx, "compiled %d rules\nfailed %d rules", l.compiled, l.failed)
    dontio.WriteLineOut(ctx, "all identifiers found:\n%s", identifierString)
	return nil
}

type ParsableEdge struct {
	Name     components.Name
	RawLogic components.RawLogic
}

func (p *ParsableEdge) Init(_ table.RowId, tup table.ValueTuple) error {
	p.Name = tup.Values[0].(components.Name)
	p.RawLogic = tup.Values[1].(components.RawLogic)
	return nil
}

type empty = struct{}
type identifiers map[string]struct{}

func (i identifiers) VisitBinOp(ast *parser.BinOp) error {
    lhErr := parser.Visit(i, ast.Left)
    rhErr := parser.Visit(i, ast.Right)
	return errors.Join(lhErr, rhErr)
}

func (i identifiers) VisitBoolOp(ast *parser.BoolOp) error {
    lhErr := parser.Visit(i, ast.Left)
    rhErr := parser.Visit(i, ast.Right)
	return errors.Join(lhErr, rhErr)
}

func (i identifiers) VisitCall(ast *parser.Call) error {
    var errs []error

    if callErr := parser.Visit(i, ast.Callee); callErr != nil {
        errs = append(errs, callErr)
    }

    for _, arg := range ast.Args {
        if err := parser.Visit(i, arg); err != nil {
            errs = append(errs, err)
        }
    }

    return errors.Join(errs...)
}

func (i identifiers) VisitIdentifier(ast *parser.Identifier) error {
    i[ast.Value] = empty{}
    return nil
}

func (i identifiers) VisitSubscript(ast *parser.Subscript) error {
    idxErr := parser.Visit(i, ast.Index)
    targetErr := parser.Visit(i, ast.Target)
    return errors.Join(idxErr, targetErr)
}

func (i identifiers) VisitTuple(ast *parser.Tuple) error {
	var errs []error
	for _, elm := range ast.Elems {
		if err := parser.Visit(i, elm); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

func (i identifiers) VisitUnary(ast *parser.UnaryOp) error {
	return parser.Visit(i, ast.Target)
}

func (i identifiers) VisitLiteral(ast *parser.Literal) error {
	return nil
}
