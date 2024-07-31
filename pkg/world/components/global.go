package components

import "sudonters/zootler/internal/entity"

type (
	Collectable struct{}
	Collected   struct{}
	DefaultItem entity.Model
	Inhabited   entity.Model
	Inhabits    entity.Model
	Locked      struct{}
	Name        string
	Placeable   struct{} // ???: should this carry _what_ is placeable here
	Spawn       struct{}
	Trick       struct{}
)
