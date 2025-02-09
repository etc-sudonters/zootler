package eligibity

import (
	"sudonters/libzootr/internal/skelly/bitset32"
	"sudonters/libzootr/zecs"
)

func CreateAll(ocm *zecs.Ocm, criteria Criteria) error {
	// todo: loop all items/locations and create an eligibity set for each
	return nil
}

type Criterion uint32

type Criteria interface {
	Create(zecs.Proxy) bitset32.Bitset
}

const (
	_ Criterion = iota
	Major
	Advancement
	Junk
	DungeonReward
	SkullToken
	Song
	Shop
)
