package components

import (
	"sudonters/zootler/internal/entity"
	"sudonters/zootler/internal/rules/runtime"
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
	Bytecode runtime.Chunk
}

type ExitEdge struct{}
type CheckEdge struct{}
type EventEdge struct{}
type TimePasses struct{}
type SavewarpName string