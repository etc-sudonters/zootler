package table

import (
	"sudonters/zootler/internal/entity"

	"github.com/etc-sudonters/substrate/skelly/bitset"
)

type Index interface {
	Set(e entity.Model, c entity.Component)
	Unset(e entity.Model, c entity.Component)
	Len() int
	Matches(c entity.Component) int
	Get(c entity.Component) bitset.Bitset64
}
