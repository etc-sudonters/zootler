package generation

import (
	"context"
	"math/rand/v2"
	"sudonters/zootler/internal/query"
	"sudonters/zootler/internal/rules/runtime"
	"sudonters/zootler/internal/settings"
	"sudonters/zootler/internal/table"

	"github.com/etc-sudonters/substrate/skelly/graph"
)

type Generation struct {
	Id      uint64
	Seed    *SeedBuilder
	Ctx     context.Context
	Cancel  context.CancelCauseFunc
	Rngesus *rand.Rand
}

type SeedBuilder struct {
	EdgeRules map[table.RowId]runtime.Chunk
	Settings  settings.ZootrSettings
	Storage   query.Engine
	Worlds    []WorldBuilder
}

type WorldBuilder struct {
	G     graph.Builder
	State WorldState
}

type WorldState struct{}
