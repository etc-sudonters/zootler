package world

import (
	"github.com/etc-sudonters/zootler/entity"
	"github.com/etc-sudonters/zootler/graph"
	"github.com/etc-sudonters/zootler/logic"
)

type TimeOfDay uint8

const (
	DayToD   TimeOfDay = 1 << iota
	NightToD TimeOfDay = 1 << iota
	DampeToD TimeOfDay = 1 << iota

	Any TimeOfDay = DampeToD | NightToD | DayToD
)

type HintCategory string

const (
	// Age
	AdultExclusiveComponent entity.TagName = "adult-exclusive"
	ChildExclusiveComponent entity.TagName = "child-exclusive"

	// Time of Day
	TimeOfDayComponent entity.TagName = "time-of-day"

	// dungeon
	OriginDungeon      entity.TagName = "origin-dungeon"  // where did we come from
	CurrentDungeon     entity.TagName = "current-dungeon" // where did we go
	ScarecrowComponent entity.TagName = "scarecrow"       // cotton eye joe

	// presentation
	ChestComponent         entity.TagName = "chest"
	CloakableComponent     entity.TagName = "cloakable"
	FreestandingComponent  entity.TagName = "freestanding"
	GoldSkulltulaComponent entity.TagName = "gold-skulltula"
	PotComponent           entity.TagName = "pot"

	// CSMC
	TextureComponent entity.TagName = "texture"
	SizeComponent    entity.TagName = "size"

	// pick up kinds
	BossKeyComponent       entity.TagName = "boss-key"
	CompassComponent       entity.TagName = "compass"
	MapComponent           entity.TagName = "map"
	SmallKeyComponent      entity.TagName = "small-key"
	BottlableComponent     entity.TagName = "bottleable"     // bugs, poes, fairies
	DungeonReward          entity.TagName = "dungeon-reward" // medallion, stone
	EnemyDropComponent     entity.TagName = "enemy-drop"     // deku sticks
	EventComponent         entity.TagName = "event"          // "show mido the sword and shield"
	HintComponent          entity.TagName = "hint"           // read from a gossip stone or frogs
	JunkComponent          entity.TagName = "junk"           // blupee, recovery heart
	PlaceableComponent     entity.TagName = "placeable"      // this token is available for placing
	GroupComponent         entity.TagName = "grouping"       // silver rupee pouch, keyring
	GroupableComponent     entity.TagName = "groupable"      // silver rupee pouch, keyring
	PriceComponent         entity.TagName = "price"          // token and check can both have this
	PriorityComponent      entity.TagName = "priority"       // ice arrows, stone of agony
	ProgressiveComponent   entity.TagName = "progressive"    // hookshot, hookshot
	RecoveryComponent      entity.TagName = "recovery"       // recovery heart, heart piece, heart container
	RefillComponent        entity.TagName = "refill"         // deku nuts, magic drop, recovery heart
	RupeeComponent         entity.TagName = "rupee"          // blupee
	ScrubComponent         entity.TagName = "scrub"          // this barbie is a deku scrub location
	ShopComponent          entity.TagName = "shop"           // this barbie is a shelf space in a shop
	SilverRupeeComponent   entity.TagName = "silver-rupee"   // e.g. ice cavern scythe room
	SongComponent          entity.TagName = "song"           // sun song, Serenade
	TokenComponent         entity.TagName = "token"          // anything we collected, can be item but also event
	TradeComponent         entity.TagName = "trade-quest"    // this token or check is associated to a trade quest
	TrapComponent          entity.TagName = "trap"           // this barbie is actually an ice trap on pick up
	TriforcePiece          entity.TagName = "triforce-piece" // triforce hunt, could be interesting
	WarpComponent          entity.TagName = "warp"           // this song or location is associated to a warp
	GossipStoneComponent   entity.TagName = "gossip-stone"   // everything else points at a fixed, specific location
	OcarinaButtonComponent entity.TagName = "ocarina-button" // ^ < > v A

	// state related
	/*
		Collected means slightly different things depending on where
		it is encountered:
		- Alongside Token: This token belongs to us currently
		- Alongside Check: We came here and took the item that was here
		- Alongside Location: There are no empty checks neighboring this location
	*/
	CollectedComponent entity.TagName = "collected"
	EmptyComponent     entity.TagName = "empty"    // this check has not been filled yet
	CheckComponent     entity.TagName = "check"    // anything that gives us a token
	LocationComponent  entity.TagName = "location" // the immediate playable game area
	RegionComponent    entity.TagName = "region"   // the geographic area we're in

	//graph
	SpawnComponent      entity.TagName = "spawn-node"
	SpawnAge            entity.TagName = "spawn-age"
	DestNodeComponent   entity.TagName = "dest-node"
	EdgeComponent       entity.TagName = "edge"
	NodeComponent       entity.TagName = "node"
	OriginNodeComponent entity.TagName = "origin-node"

	//entrances
	TransitComponent entity.TagName = "transit"
	InteriorEntrance entity.TagName = "interior-entrance"
	DungeonEntrance  entity.TagName = "dungeon-entrance"

	//logic/placement
	OriginRegion          entity.TagName = "origin-region"
	OriginWorldComponent  entity.TagName = "origin-world"
	CurrentWorldComponent entity.TagName = "current-world"
	RuleComponent         entity.TagName = "rule"
	RestrictedPlacement   entity.TagName = "restricted-placement"
	HintedComponent       entity.TagName = "hinted"
)

const (
	// archetypes
	CollectedArchetype      = entity.ArchetypeTag(CollectedComponent)
	EmptyArchetype          = entity.ArchetypeTag(EmptyComponent)
	HintArchetype           = entity.ArchetypeTag(HintComponent)
	LocationArchetype       = entity.ArchetypeTag(LocationComponent)
	PlaceableArchetype      = entity.ArchetypeTag(PlaceableComponent)
	RegionArchetype         = entity.ArchetypeTag(RegionComponent)
	TokenArchetype          = entity.ArchetypeTag(TokenComponent)
	TransitArchetype        = entity.ArchetypeTag(TransitComponent)
	AdultExclusiveArchetype = entity.ArchetypeTag(AdultExclusiveComponent)
	ChildExclusiveArchetype = entity.ArchetypeTag(ChildExclusiveComponent)
)

func OriginWorldArchetype(w Id) entity.ArchetypeFunc {
	return entity.SingleTagArchetype(OriginWorldComponent, w)
}

func NodeArchetype(n graph.Node) entity.ArchetypeFunc {
	return entity.SingleTagArchetype(NodeComponent, n)
}

func EdgeArchetype(o graph.Origination, d graph.Destination) entity.ArchetypeFunc {
	return func(t entity.Tags) {
		t.Apply(OriginNodeComponent, o).
			Apply(DestNodeComponent, d).
			Use(entity.ArchetypeTag(EdgeComponent))
	}
}

func RuleArchetype(r logic.Rule) entity.ArchetypeFunc {
	return func(t entity.Tags) {
		t.Apply(RuleComponent, r)
	}
}

func TransitRuleArchetype(r logic.Rule) entity.Archetype {
	return entity.BundleArchetypes(RuleArchetype(r), TransitArchetype)
}

func PickupRuleArchetype(r logic.Rule) entity.Archetype {
	return entity.BundleArchetypes(RuleArchetype(r), PickupArchectype)
}

func HintRuleArchetype(r logic.Rule) entity.Archetype {
	return entity.BundleArchetypes(RuleArchetype(r), HintArchetype)
}

func TimeOfDayArchetype(t TimeOfDay) entity.Archetype {
	return entity.SingleTagArchetype(TimeOfDayComponent, t)
}
