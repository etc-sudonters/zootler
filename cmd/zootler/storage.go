package main

import (
	"os"
	"sudonters/zootler/internal/query"
	"sudonters/zootler/internal/table"
	"sudonters/zootler/internal/table/columns"
	"sudonters/zootler/pkg/world/components"

	"muzzammil.xyz/jsonc"
)

type IntoComponents interface {
	GetName() components.Name
	AddComponents(table.RowId, query.Engine) error
}

type DataFileLoader[T IntoComponents] string

func (l DataFileLoader[T]) Configure(storage query.Engine) error {
	raw, readErr := os.ReadFile(string(l))
	if readErr != nil {
		return readErr
	}

	var items []T

	if err := jsonc.Unmarshal(raw, &items); err != nil {
		return err
	}

	for _, item := range items {
		row, insertErr := storage.InsertRow(item.GetName())
		if insertErr != nil {
			return insertErr
		}

		if valuesErr := item.AddComponents(row, storage); valuesErr != nil {
			return valuesErr
		}
	}

	return nil
}

type CreateScheme struct {
	DDL []DDL
}
type DDL func() *table.ColumnBuilder

func (cs CreateScheme) Configure(storage query.Engine) error {
	for _, ddl := range cs.DDL {
		if _, err := storage.CreateColumn(ddl()); err != nil {
			return err
		}
	}

	return nil
}

func BitColumnOf[T any]() *table.ColumnBuilder {
	var t T
	return table.BuildColumnOf[T](columns.NewBit(t))
}

func MapColumn[T any]() *table.ColumnBuilder {
	return table.BuildColumnOf[T](columns.NewMap())
}

func SliceColumn[T any]() *table.ColumnBuilder {
	return table.BuildColumnOf[T](columns.NewSlice())
}
