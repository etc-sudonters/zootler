package table

import (
	"sudonters/zootler/internal/entity"

	"github.com/etc-sudonters/substrate/skelly/bitset"
)

type Index interface {
	Column() ColumnId
	Set(e entity.Model, c entity.Component)
	Unset(e entity.Model, c entity.Component)
	Len() int
	Matches(c entity.Component) int
	// this bitset is intersected / & / AND'd
	Get(c entity.Component) bitset.Bitset64
}
