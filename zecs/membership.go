package zecs

import (
	"sudonters/libzootr/internal/query"
	"sudonters/libzootr/internal/skelly/bitset32"
)

type Membership struct {
	row    bitset32.Bitset
	parent *Ocm
}

func IsMemberOf[V Value](this Membership) bool {
	id := query.MustAsColumnId[V](this.parent.eng)
	return bitset32.IsSet(&this.row, id)
}
