package components

import "sudonters/zootler/internal/entity"

type (
	Collectable     struct{}
	Collected       int
	DefaultItemName string
	DefaultItem     entity.Model
	Inhabited       entity.Model
	Inhabits        entity.Model
	Locked          struct{}
	Placeable       struct{} // ???: should this carry _what_ is placeable here
	Spawn           struct{}
	Trick           struct{}
	Name            string
)
