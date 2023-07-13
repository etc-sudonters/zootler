package logic

import (
	"github.com/etc-sudonters/zootler/pkg/entity"
)

const (
	DayToD   TimeOfDay = 1 << iota
	NightToD TimeOfDay = 1 << iota
	DampeToD TimeOfDay = 1 << iota

	Any TimeOfDay = DampeToD | NightToD | DayToD
)

type (
	Name           string
	Advancement    struct{}
	Dungeon        entity.Model
	TimeOfDay      uint8
	HintCategory   string
	Adult          struct{}
	Child          struct{}
	AgeExclusive   struct{}
	OriginDungeon  Dungeon
	CurrentDungeon Dungeon
	Scarecrow      struct{}
	Chest          struct{}
	Cloakable      struct{}
	FreeStanding   struct{}
	GoldSkulltula  struct{}
	Pot            struct{}
	Texture        struct{}
	Size           struct{}
	BossKey        Dungeon
	Compass        Dungeon
	Map            Dungeon
	SmallKey       Dungeon
	DungeonReward  struct{}
	EnemyDrop      struct{}
	Event          struct{}
	Hint           struct{}
	Junk           struct{}
	Placeable      struct{}
	Price          int
	Priority       int
	Progressive    struct{}
	Recovery       struct{}
	Refill         struct{}
	Rupee          struct{}
	Scrube         struct{}
	Shop           struct{}
	SilverRupee    Region
	Region         entity.Model
	Song           struct{}
	Token          struct{}
	Trade          struct{}
	TriforcePiece  struct{}
	Warp           Region
	GossipStone    Location
	Location       entity.Model
	OcarinaButton  struct{}
	Collected      struct{}
	Empty          struct{}
	Check          struct{}
	Spawn          struct{}
	Node           struct{}
	Edge           struct {
		Destination entity.Model
		Origination entity.Model
	}
	OriginRegion        Region
	CurrentRegion       Region
	RuleComp            Rule
	RestrictedPlacement struct{}
	Hinted              HintCategory
	Trick               Name
	Enabled             struct{}
	Disabled            struct{}
)
