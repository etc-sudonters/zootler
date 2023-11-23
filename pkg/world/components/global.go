package components

import "sudonters/zootler/internal/entity"

type (
	Name      string
	Collected struct{}
	Trick     struct{}
	Spawn     struct{}
	Inhabited entity.Model
	Inhabits  entity.Model
	Locked    struct{}
)
