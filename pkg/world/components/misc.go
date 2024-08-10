package components

import (
	"sudonters/zootler/internal/entity"
)

type BossRoom struct{}

type HintRegion struct {
	Name, Alt string
}

type Edge struct {
	Origin entity.Model
	Dest   entity.Model
}

type RawLogic struct {
	Rule string
}

type Helper struct{}

type CompiledRule struct {
	Code []uint8
}
