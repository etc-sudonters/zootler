package settings

// enum members are arranged so a 0 value can be interpreted as default
// which is context dependent

type LogicRuleSet uint8

const (
	_ LogicRuleSet = iota
	LogicGlitchless
	LogicGlitched
	LogicNone
)

type KokiriForest uint8

const (
	_ KokiriForest = iota
	KokiriForestOpen
	KokiriForestClosed
	KokiriForestClosedDeku
)

type KakarikoGate uint8

const (
	_ KakarikoGate = iota
	KakarikoGateOpenGate
	KakarikoGateClosedGate
	KakarikoGateLetterOpensGate
)

type DoorOfTime uint8

const (
	_ DoorOfTime = iota
	DoorOfTimeOpen
	DoorOfTimeClosed
)

type ZorasFountain uint8

const (
	_ ZorasFountain = iota
	FountainClosed
	FountainAdultOpen
	FountainOpen
)

type FortressCarpenters uint8

const (
	_ FortressCarpenters = iota
	FortressAllCarpenters
	FortressOneCarpenter
	FortressNoCarpenters
)

type BridgeKind uint8

const (
	_ BridgeKind = iota
	BridgeOpen
	BridgeVanilla
	BridgeStones
	BridgeMedallions
	BridgeDungeonRewards
	BridgeSkulls
	BridgeHearts
)

type BridgeRequirement struct {
	Kind   BridgeKind
	Amount uint8
}

type TowerTrialCount uint8

type StartingAge uint8

const (
	_ StartingAge = iota
	StartingAdult
	StartingChild
)

type SongShuffle uint8

const (
	_ SongShuffle = iota
	ShuffleSongLocations
	ShuffleSongOnDungeonRewards
	ShuffleSongsAnywhere
)

type ShopShuffle uint8

const (
	_            ShopShuffle = iota
	ShopShuffle0             // lol this sucks but special cases aren't special
	ShopShuffle1
	ShopShuffle2
	ShopShuffle3
	ShopShuffle4
)

type GoldTokenShuffle uint8 // flags

const (
	TokenShuffleDungeons = 1 << iota
	TokenShuffleOverworld
)

type ScrubShuffle uint8

const (
	_ ScrubShuffle = iota
	ScrubShuffleAffordable
	ScrubShuffleExpensive
	ScrubShuffleRandom
)

type ChildTradeQuest uint

const (
	_ ChildTradeQuest = iota
	ChildTradeVanilla
	ChildTradeShuffleEgg
	ChildTradeSkipZeldaMeeting
)

type PotShuffle uint8 // flags

const (
	PotShuffleOverworld = 1 << iota
	PotShuffleDungeon
)

type CrateShuffle uint8 // flags

const (
	CrateShuffleOverworld = 1 << iota
	CrateShuffleDungeon
)

type CowShuffle uint8

const (
	_ CowShuffle = iota
	CowShuffleAll
)

type BeehiveShuffle uint8

const (
	_ BeehiveShuffle = iota
	BeehiveShuffleAll
)

type KokriSwordShuffle uint8

const (
	_ KokriSwordShuffle = iota
	KokriSwordShuffleAnywhere
)

type OcarinaShuffle uint8

const (
	_ OcarinaShuffle = iota
	OcarinaShuffleAnywhere
)

type GerudoCardShuffle uint8

const (
	_ GerudoCardShuffle = iota
	GerudoCardShuffleAnywhere
)

type MagicBeanShuffle uint8

const (
	_ MagicBeanShuffle = iota
	MagicBeanShuffleBag
)

type RepeatMerchantShuffle uint8 // flags

const (
	MerchantShuffleMedigoron RepeatMerchantShuffle = 1 << iota
	MerchantShuffleCarpet
)

type FrogRupeeShuffle uint8

const (
	_ FrogRupeeShuffle = iota
	FrogRupeesAnywhere
)

type MapsAndCompassesShuffle uint8

const (
	_ MapsAndCompassesShuffle = iota
	MapsAndCompassesNone
	MapsAndCompassesBeginWith
	MapsAndCompassesVanilla
	MapsAndCompassesOwnDungeon
	MapsAndCompassesRegional
	MapsAndCompassesOverworld
	MapsAndCompassesAnyDungeon
	MapsAndCompassesAnywhere
)

type KeyShuffle uint8
type SmallKeyShuffle KeyShuffle
type BossKeyShuffle KeyShuffle
type TowerBossKeyShuffle KeyShuffle

const (
	_ KeyShuffle = iota
	KeysRemove
	KeysVanilla
	KeysOwnDungeon
	KeysRegion
	KeysOverworld
	KeysAnyDungeon
	KeysAnywhere
)

type MapsAndCompassesExtras uint8 // flags

const (
	MapsAndCompassesKnowRewards MapsAndCompassesExtras = 1 << iota
)

type ItemPool uint8

const (
	_ ItemPool = iota
	ItemPoolLudicrous
	ItemPoolPlentiful
	ItemPoolBalanced
	ItemPoolScarce
	ItemPoolMinimal
)

type AdultTradeItems uint16 // flags

const (
	AdultTradeEgg AdultTradeItems = 1 << iota
	AdultTradeCucco
	AdultTradeCojiro
	AdultTradeMushroom
	AdultTradeSaw
	AdultTradeBrokenSword
	AdultTradePrescription
	AdultTradeFrog
	AdultTradeEyeDrops
	AdultTradeClaimCheck
)

type ChestGameKeyShuffle uint

const (
	_ ChestGameKeyShuffle = iota
	ChestGameKeysVanilla
)
