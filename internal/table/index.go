package table

import "sudonters/zootler/internal/skelly/bitset"

type Index interface {
	Set(r RowId, v Value)
	Unset(r RowId)
	// this bitset is intersected / & / AND'd
	Rows(v Value) bitset.Bitset32
}
