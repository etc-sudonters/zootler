package components

import (
	"sudonters/zootler/internal/entity"
	"sudonters/zootler/internal/rules/runtime"
)

type BossRoom struct{}

type HintRegion struct {
	Name, Alt string
}

type Edge struct{}

type Connection struct {
	Origin, Dest entity.Model
}

type RawLogic string

type Helper struct{}

type CompiledRule struct {
	Bytecode runtime.Chunk
}

type ExitEdge struct{}
type CheckEdge struct{}
type EventEdge struct{}
type TimePasses struct{}
type SavewarpName string
type Scene string
