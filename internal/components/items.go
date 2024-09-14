package components

// break import cycles

type (
	BossKey              struct{}
	Bottle               struct{}
	Compass              struct{}
	Count                float64
	Drop                 struct{}
	DungeonReward        struct{}
	Event                struct{}
	GanonBossKey         struct{}
	GoldSkulltulaToken   struct{}
	HideoutSmallKey      struct{}
	Item                 struct{}
	Junk                 struct{}
	Map                  struct{}
	Medallion            struct{}
	Price                float64
	Refill               struct{}
	ShopObject           float64
	SmallKey             struct{}
	SpiritualStone       struct{}
	Trade                struct{}
	Location             struct{}
	CollectableGameToken struct{}
	Song                 struct{}

	TCGSmallKey struct{}

	Advancement struct{}
	Priority    struct{}

	AnonymousEvent struct{}
)

func (c BossKey) String() string              { return "Boss Key" }
func (c CollectableGameToken) String() string { return "Collectable Token" }
func (c Compass) String() string              { return "Compass" }
func (c Drop) String() string                 { return "Drop" }
func (c DungeonReward) String() string        { return "Dungeon Reward" }
func (c Event) String() string                { return "Event" }
func (c GanonBossKey) String() string         { return "Ganon Boss Key" }
func (c GoldSkulltulaToken) String() string   { return "Gold Skulltula Token" }
func (c HideoutSmallKey) String() string      { return "Hideout Small Key" }
func (c Item) String() string                 { return "Item" }
func (c Location) String() string             { return "Location" }
func (c Map) String() string                  { return "Map" }
func (c Refill) String() string               { return "Refill" }
func (c SmallKey) String() string             { return "Small Key" }
func (c Song) String() string                 { return "Song" }
