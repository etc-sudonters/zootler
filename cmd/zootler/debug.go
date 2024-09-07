package main

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"sudonters/zootler/internal/app"
	"sudonters/zootler/internal/query"
	"sudonters/zootler/internal/table"

	"github.com/etc-sudonters/substrate/dontio"
	"github.com/etc-sudonters/substrate/slipup"
)

func examineTable(ctx context.Context, storage query.Engine) error {
	tbl, extractTblErr := query.ExtractTable(storage)
	if extractTblErr != nil {
		return extractTblErr
	}

	dontio.WriteLineOut(ctx, "Number of columns:\t%d", len(tbl.Cols))
	dontio.WriteLineOut(ctx, "Number of rows:\t\t%d", len(tbl.Rows))
	return nil
}

type DebugSetupFunc func(ctx context.Context, storage query.Engine) error

func (d DebugSetupFunc) Setup(ctx context.Context, storage query.Engine) error {
	return d(ctx, storage)
}

type InspectTable struct {
	Columns []reflect.Type
}

func (i InspectTable) Setup(z *app.Zootlr) error {
	storage := z.Table()
	ctx := z.Ctx()
	tbl, extractTblErr := query.ExtractTable(storage)
	if extractTblErr != nil {
		return extractTblErr
	}

	dontio.WriteLineOut(ctx, "Number of columns:\t%d", len(tbl.Cols))
	dontio.WriteLineOut(ctx, "Number of rows:\t\t%d", len(tbl.Rows))
	for _, t := range i.Columns {
		id, ok := storage.ColumnIdFor(t)
		if !ok {
			return fmt.Errorf("could not find column for '%s'", t.Name())
		}
		if err := examineColumn(ctx, tbl.Cols[id]); err != nil {
			return slipup.Describef(err, "while inspecting column '%s'", t.Name())
		}
	}
	return nil
}

func examineColumn(ctx context.Context, col table.ColumnData) error {
	dontio.WriteLineOut(ctx, "Column:\t\t%s", col.Type().Name())
	dontio.WriteLineOut(ctx, "Id:\t\t%d", col.Id())
	dontio.WriteLineOut(ctx, "Population:\t%d", col.Column().Len())
	return nil
}

type InspectQuery struct {
	Exist    []reflect.Type
	NotExist []reflect.Type
}

func (iq InspectQuery) Setup(ctx context.Context, storage query.Engine) error {
	q := storage.CreateQuery()

	for i := range iq.Exist {
		t := iq.Exist[i]
		colId, exists := storage.ColumnIdFor(t)
		if !exists {
			dontio.WriteLineErr(ctx, "Unable to resolve columnId for '%s'", t)
		}
		q.Exists(colId)
	}

	for i := range iq.NotExist {
		t := iq.NotExist[i]
		colId, exists := storage.ColumnIdFor(t)
		if !exists {
			dontio.WriteLineErr(ctx, "Unable to resolve columnId for '%s'", t)
		}
		q.NotExists(colId)
	}

	results, err := storage.Retrieve(q)
	if err != nil {
		return slipup.Describe(err, "while inspecting query")
	}

	dontio.WriteLineOut(ctx, "Examining query")
	dontio.WriteLineOut(ctx, "Exists:\t\t%s", showNames(iq.Exist))
	dontio.WriteLineOut(ctx, "NotExist:\t%s", showNames(iq.NotExist))
	dontio.WriteLineOut(ctx, "Population:\t%d", results.Len())
	return nil
}

func showNames(ts []reflect.Type) string {
	var b strings.Builder
	final := len(ts) - 1
	for i := range ts {
		b.WriteString(ts[i].Name())
		if i != final {
			b.WriteRune(',')
		}
	}

	return b.String()
}
