package components

// break import cycles

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
	Song               struct{}
	Location           struct{}
	Token              struct{}
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
func (c SmallKey) String() string           { return "SmallKey" }
func (c Song) String() string               { return "Song" }
func (c GoldSkulltulaToken) String() string { return "Token" }
func (c Location) String() string           { return "Location" }
func (c Token) String() string              { return "Token" }
