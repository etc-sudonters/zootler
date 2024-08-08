package main

import (
	"context"
	"sudonters/zootler/internal/query"

	"github.com/etc-sudonters/substrate/dontio"
)

type InspectTable struct{}

func (i InspectTable) Configure(ctx context.Context, storage query.Engine) error {
	stdio, dontioErr := dontio.StdFromContext(ctx)
	if dontioErr != nil {
		return dontioErr
	}

	std := std{stdio}

	tbl, extractTblErr := query.ExtractTable(storage)
	if extractTblErr != nil {
		return extractTblErr
	}

	std.WriteLineOut("Number of columns:\t%d", len(tbl.Cols))
	std.WriteLineOut("Number of rows:\t\t%d", len(tbl.Rows))
	return nil
}

type Func func(context.Context, query.Engine) error

func (f Func) Configure(ctx context.Context, storage query.Engine) error {
	return f(ctx, storage)
}
