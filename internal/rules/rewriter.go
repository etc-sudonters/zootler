package rules

import (
	"sudonters/zootler/internal/query"
	"sudonters/zootler/internal/rules/runtime"
	"sudonters/zootler/internal/table"
)

type AstRewriter struct {
	Env              *runtime.ExecutionEnvironment
	Funcs            *runtime.FuncNamespace
	ProgressiveItems ProgressiveItems
}

type ProgressiveItems map[string]table.RowId

func (p *ProgressiveItems) Init(eng query.Engine) error {
	newItems := make(map[string]table.RowId, 2048)

	q := eng.CreateQuery()

	*p = newItems
	return nil
}
