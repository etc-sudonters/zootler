package magicbean

import (
	"math/rand/v2"
	"sudonters/zootler/internal/settings"

	"sudonters/zootler/mido/objects"
	"sudonters/zootler/zecs"
)

type Generation struct {
	Ocm       zecs.Ocm
	World     ExplorableWorld
	Objects   objects.Table
	Inventory Inventory
	Rng       rand.Rand
	Settings  settings.Zootr
}
