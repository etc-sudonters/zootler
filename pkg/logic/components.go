package logic

import (
	"github.com/etc-sudonters/zootler/pkg/entity"
)

const (
	// day time
	DayToD TimeOfDay = 1 << iota
	// night time
	NightToD TimeOfDay = 1 << iota
	// 6pm to 9pm, when dampe digging is available
	DampeToD TimeOfDay = 1 << iota

	Any TimeOfDay = DampeToD | NightToD | DayToD
)

type (
	Adult          struct{}
	Advancement    struct{}
	AgeExclusive   struct{}
	BossKey        Dungeon
	Check          struct{}
	Chest          struct{}
	Child          struct{}
	Cloakable      struct{}
	Collected      struct{}
	Compass        Dungeon
	CurrentDungeon Dungeon
	Dungeon        entity.Model
	DungeonReward  struct{}
	Edge           struct {
		Destination entity.Model
		Origination entity.Model
	}
	Empty               struct{}
	Enabled             bool
	EnemyDrop           struct{}
	Event               struct{}
	FreeStanding        struct{}
	GoldSkulltula       struct{}
	GossipStone         Location
	Hint                struct{}
	HintCategory        string
	HintGroup           string
	Hinted              HintCategory
	Junk                struct{}
	Location            entity.Model
	Map                 Dungeon
	Name                string
	Node                struct{}
	OcarinaButton       struct{}
	OriginDungeon       Dungeon
	Placeable           struct{}
	Pot                 struct{}
	Price               int
	Priority            int
	Progressive         struct{}
	RawRule             string
	Recovery            struct{}
	Refill              struct{}
	RestrictedPlacement struct{}
	RuleComp            Rule
	Rupee               struct{}
	Scarecrow           struct{}
	Scrub               struct{}
	Shop                struct{}
	Size                struct{}
	SmallKey            Dungeon
	Song                struct{}
	Spawn               struct{}
	Texture             struct{}
	TimeOfDay           uint8
	Token               struct{}
	Trade               struct{}
	Trick               Name
	TriforcePiece       struct{}
)
