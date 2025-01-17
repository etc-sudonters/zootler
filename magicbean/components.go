package magicbean

import (
	"fmt"
	"sudonters/zootler/mido/ast"
	"sudonters/zootler/mido/compiler"
	"sudonters/zootler/mido/objects"
	"sudonters/zootler/zecs"
)

type Name string

func NameF(tpl string, v ...any) Name {
	return Name(fmt.Sprintf(tpl, v...))
}

type Transit struct{ From, To zecs.Entity }
type Region struct{}
type Placement struct{}
type DefaultPlacement zecs.Entity
type Token struct{}
type Fixed struct{}

type ScriptDecl string
type ScriptBody string

type StringRule string
type ParsedRule ast.Node
type OptimizedRule ast.Node
type CompiledRule compiler.Bytecode

type HeldAt zecs.Entity
type HoldsToken zecs.Entity
type Empty struct{}
type Generated struct{}
type Ptr objects.Object

type Collectable struct{}
type Location struct{}
type EdgeKind uint8
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
	_             EdgeKind = 0
	EdgeTransit            = 0x69
	EdgePlacement          = 0xBB

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
