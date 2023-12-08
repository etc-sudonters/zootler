package table

import (
	"sudonters/zootler/internal/entity"

	"github.com/etc-sudonters/substrate/skelly/bitset"
)

type Index interface {
	Column() ColumnId
	Set(e entity.Model, c entity.Component)
	Unset(e entity.Model, c entity.Component)
	// lower is better
	Estimate(c entity.Component) uint64
	// this bitset is intersected / & / AND'd
	Rows(c entity.Component) bitset.Bitset64
}
