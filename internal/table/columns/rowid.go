package columns

import (
	"sudonters/zootler/internal/table"
)

type RowId struct {
}

func (i RowId) Get(e table.RowId) table.Value {
	return e
}

func (i RowId) Set(e table.RowId, c table.Value) {
}

func (i RowId) Unset(e table.RowId) {
}
