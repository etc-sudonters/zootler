package main

import (
	"context"
	"fmt"
	"reflect"
	"sudonters/zootler/internal/query"
	"sudonters/zootler/internal/slipup"
	"sudonters/zootler/internal/table"
)

type InspectTable struct {
	Columns []reflect.Type
}

func (i InspectTable) Configure(ctx context.Context, storage query.Engine) error {
	tbl, extractTblErr := query.ExtractTable(storage)
	if extractTblErr != nil {
		return extractTblErr
	}

	WriteLineOut(ctx, "Number of columns:\t%d", len(tbl.Cols))
	WriteLineOut(ctx, "Number of rows:\t\t%d", len(tbl.Rows))
	for _, t := range i.Columns {
		id, ok := query.ColumnIdFor(storage, t)
		if !ok {
			return fmt.Errorf("could not find column for '%s'", t.Name())
		}
		if err := examineColumn(ctx, tbl.Cols[id]); err != nil {
			return slipup.TraceMsg(err, "while inspecting column '%s'", t.Name())
		}
	}
	return nil
}

func examineColumn(ctx context.Context, col table.ColumnData) error {
	WriteLineOut(ctx, "Column:\t\t%s", col.Type().Name())
	WriteLineOut(ctx, "Id:\t\t%d", col.Id())
	WriteLineOut(ctx, "Population:\t%d", col.Column().Len())
	return nil
}
