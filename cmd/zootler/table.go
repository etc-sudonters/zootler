package main

import (
	"context"
	"sudonters/zootler/internal/query"
)

func examineTable(ctx context.Context, storage query.Engine) error {
	tbl, extractTblErr := query.ExtractTable(storage)
	if extractTblErr != nil {
		return extractTblErr
	}

	WriteLineOut(ctx, "Number of columns:\t%d", len(tbl.Cols))
	WriteLineOut(ctx, "Number of rows:\t\t%d", len(tbl.Rows))
	return nil
}
