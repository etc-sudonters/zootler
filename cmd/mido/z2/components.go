package z2

import (
	"fmt"
	"sudonters/zootler/internal/table"
	"sudonters/zootler/mido/ast"
	"sudonters/zootler/mido/compiler"
	"sudonters/zootler/mido/objects"
)

type HeldAt Entity
type HoldsToken Entity
type FixedPlacement struct{}
type EmptyPlacement struct{}
type Generated struct{}
type Ptr objects.Object

type Entity table.RowId
type Collectable struct{}
type Location struct{}
type Name string
type Connection struct{ From, To Entity }
type StringSource string
type ParsedSource ast.Node
type OptimizedSource ast.Node
type CompiledSource compiler.Bytecode
type ConnectionKind uint8
type HintRegion string
type AltHintRegion string
type DungeonName string
type IsBossRoom struct{}
type Savewarp string
type Scene string
type TimePassess struct{}
type CollectableType string
type CollectablePriority uint8

const (
	_               ConnectionKind = 0
	ConnectionCheck                = 0xBB
	ConnectionEvent                = 0xA0
	ConnectionExit                 = 0x69

	PriorityJunk        CollectablePriority = 0
	PriorityMajor                           = 0xE0
	PriorityAdvancement                     = 0xF0
)

type BossKey struct{}
type Compass struct{}
type Drop struct{}
type DungeonReward struct{}
type Event struct{}
type GanonBossKey struct{}
type HideoutSmallKey struct{}
type HideoutSmallKeyRing struct{}
type Item struct{}
type Map struct{}
type Refill struct{}
type Shop struct{}
type SilverRupee struct{}
type SmallKey struct{}
type SmallKeyRing struct{}
type Song struct{}
type TCGSmallKey struct{}
type TCGSmallKeyRing struct{}
type GoldSkulltulaToken struct{}

func NameF(tpl string, v ...any) Name {
	return Name(fmt.Sprintf(tpl, v...))
}
