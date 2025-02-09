package magicbean

import (
	"math/rand/v2"
	"sudonters/libzootr/internal/settings"

	"sudonters/libzootr/mido/objects"
	"sudonters/libzootr/zecs"
)

type Generation struct {
	Ocm       zecs.Ocm
	World     ExplorableWorld
	Objects   objects.Table
	Inventory Inventory
	Rng       rand.Rand
	Settings  settings.Zootr
}
