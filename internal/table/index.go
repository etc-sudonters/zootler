package table

import "sudonters/libzootr/internal/skelly/bitset32"

type Index interface {
	Set(r RowId, v Value)
	Unset(r RowId)
	// this bitset is intersected / & / AND'd
	Rows(v Value) bitset32.Bitset
}
