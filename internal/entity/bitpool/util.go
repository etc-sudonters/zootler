package bitpool

import (
	"errors"
	"sudonters/zootler/internal/entity"
	"sudonters/zootler/internal/entity/componenttable"
)

var ErrNoComponentTable = errors.New("pool does not have a component table")

func ExtractComponentTable(p entity.Pool) (*componenttable.Table, error) {
	b, ok := p.(*bitpool)
	if !ok {
		return nil, ErrNoComponentTable
	}

	return b.table, nil
}
