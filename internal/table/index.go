package table

import (
	"github.com/etc-sudonters/substrate/skelly/bitset"
)

type Index interface {
	Set(r RowId, v Value)
	Unset(r RowId)
	// this bitset is intersected / & / AND'd
	Rows(v Value) bitset.Bitset64
}
