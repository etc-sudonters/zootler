package world

import (
	"sudonters/zootler/internal/entity"
	"sudonters/zootler/pkg/world/settings"
)

type Connections map[entity.Model]entity.Model // destination -> edge

type Edge struct {
	Origination entity.Model
	Destination entity.Model
}

func BuildLocationPool(settings.SeedSettings) map[string]int {
	return nil
}
