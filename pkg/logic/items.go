package logic

import (
	"sudonters/zootler/internal/entity"
)

//components
type (
	BossKey            struct{}
	Bottle             struct{}
	Compass            struct{}
	Count              float64
	Drop               struct{}
	DungeonReward      struct{}
	Event              struct{}
	GanonBossKey       struct{}
	GoldSkulltulaToken struct{}
	HideoutSmallKey    struct{}
	Item               struct{}
	Junk               struct{}
	Map                struct{}
	Medallion          struct{}
	Price              float64
	Refill             struct{}
	ShopObject         float64
	SmallKey           struct{}
	SpiritualStone     struct{}
	Trade              struct{}
)

func (c BossKey) String() string            { return "BossKey" }
func (c Compass) String() string            { return "Compass" }
func (c Drop) String() string               { return "Drop" }
func (c DungeonReward) String() string      { return "DungeonReward" }
func (c Event) String() string              { return "Event" }
func (c GanonBossKey) String() string       { return "GanonBossKey" }
func (c HideoutSmallKey) String() string    { return "HideoutSmallKey" }
func (c Item) String() string               { return "Item" }
func (c Map) String() string                { return "Map" }
func (c Refill) String() string             { return "Refill" }
func (c Shop) String() string               { return "Shop" }
func (c SmallKey) String() string           { return "SmallKey" }
func (c Song) String() string               { return "Song" }
func (c GoldSkulltulaToken) String() string { return "Token" }

type PlacementItem struct {
	Name       string
	Type       entity.Component
	Importance Importance
	Components []entity.Component
}

type Importance int

const (
	ImportanceJunk        Importance = 1 << iota
	ImportancePriority    Importance = 1 << iota
	ImportanceAdvancement Importance = 1 << iota
)
