package eligibity

import (
	"sudonters/libzootr/components"
	"sudonters/libzootr/internal"
	"sudonters/libzootr/internal/skelly/bitset32"
	"sudonters/libzootr/zecs"
)

func CreateAll(ocm *zecs.Ocm, flagging Flagging) error {
	for token := range zecs.IterEntities[components.TokenMarker](ocm) {
		proxy := ocm.Proxy(token)
		membership, err := proxy.Membership()
		internal.PanicOnError(err)
		criteria := Critera{bitset32.Bitset{}}
		flagging.FlagToken(membership, &criteria)
	}

	for location := range zecs.IterEntities[components.PlacementLocationMarker](ocm) {
		proxy := ocm.Proxy(location)
		membership, err := proxy.Membership()
		internal.PanicOnError(err)
		criteria := Critera{bitset32.Bitset{}}
		flagging.FlagLocation(membership, &criteria)
	}

	return nil
}

type Criterion uint32
type Critera struct {
	set bitset32.Bitset
}

func (this *Critera) Set(c Criterion) {
	bitset32.Set(&this.set, c)
}

type Flagging interface {
	FlagToken(zecs.Membership, *Critera)
	FlagLocation(zecs.Membership, *Critera)
}

const (
	_ Criterion = iota
	Any
	Major
	Advancement
	Junk
	DungeonReward
	SkullToken
	Song
	Shop

	Anywhere = Any
	Anything = Any
)
