package indicies

import (
	"sudonters/zootler/internal/table"

	"github.com/etc-sudonters/substrate/skelly/bitset"
)

func NewBitset64(k int) Bitset64 {
	return Bitset64{bitset.New(k)}
}

// matches every entity inserted into the column, used for fast membership tests
type Bitset64 struct {
	entities bitset.Bitset64
}

func (b Bitset64) Set(e table.RowId, c table.Value) {
	b.entities.Set(int(e))
}

func (b Bitset64) Unset(e table.RowId, c table.Value) {
	b.entities.Clear(int(e))
}

func (b Bitset64) Matches(c table.Value) int {
	return b.entities.Len()
}

func (b Bitset64) Get(c table.Value) bitset.Bitset64 {
	return bitset.Copy(b.entities)
}
