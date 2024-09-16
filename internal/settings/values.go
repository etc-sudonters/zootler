package settings

import "github.com/etc-sudonters/substrate/slipup"

const (
	Forest Medallions = 0b00000001
	Fire              = 0b00000010
	Water             = 0b00000100
	Shadow            = 0b00001000
	Spirit            = 0b00010000
	Light             = 0b00100000
)

const (
	TrialsEnabledNone   TrialsEnabled = 0b00000000
	TrialsEnabledAll                  = 0b11111111
	TrialsEnabledRandom               = 0b10101010
	TrialsEnabledForest               = TrialsEnabled(Forest)
	TrialsEnabledFire                 = TrialsEnabled(Fire)
	TrialsEnabledWater                = TrialsEnabled(Water)
	TrialsEnabledShadow               = TrialsEnabled(Shadow)
	TrialsEnabledSpirit               = TrialsEnabled(Spirit)
	TrialsEnabledLight                = TrialsEnabled(Light)
)

const (
	KeyRingsRandom        Keyrings = 0b1010101010101010
	KeyRingsOff                    = 0b0000000000000000
	KeyringFortress                = 0b0000000000000001
	KeyringChestGame               = 0b0000000000000010
	KeyringForest                  = 0b0000000000000100
	KeyringFire                    = 0b0000000000001000
	KeyringWater                   = 0b0000000000010000
	KeyringShadow                  = 0b0000000000100000
	KeyringSpirit                  = 0b0000000001000000
	KeyringWell                    = 0b0000000010000000
	KeyringTrainingGround          = 0b0000000100000000
	KeyringGanonsCastle            = 0b0000001000000000
	KeyRingsAll                    = 0b0000001111111111
	KeyRingsGiveBossKey            = 0b0100000000000000
)

const (
	LogicGlitchess LogicRuleSet = 2
	LogicGlitched               = 4
	LogicNone                   = 8
)

const (
	ReachableAll      ReachableLocations = 2
	ReachableRequired                    = 4
)

const (
	GanonBKRemove GanonBKShuffleKind = 1 << iota
	GanonBKVanilla
	GanonBKDungeon
	GanonBKRegional
	GanonBKOverworld
	GanonBKAnyDungeon
	GanonBKKeysanity
	GanonBKOnLacs
	GanonBKStones
	GanonBKMedallions
	GanonBKDungeons
	GanonBKTokens
	GanonBKHearts
)

const (
	DungeonRewardVanilla DungeonRewardShuffle = 1 << iota
	DungeonRewardComplete
	DungeonRewardDungeon
	DungeonRewardRegional
	DungeonRewardOverworld
	DungeonRewardAnyDungeon
	DungeonRewardAnywhere
)

const (
	KeysVanilla KeyShuffle = 1 << iota
	KeysRemove
	KeysDungeon
	KeysRegional
	KeysOverworld
	KeysAnyDungeon
	KeysAnywhere
)

const (
	SilverRupeesOff              SilverRupeePouches = 0b00000000000000000000000000000000
	SilverRupeesAll                                 = 0b11111111111111111111111111111111
	SilverRupeesRandom                              = 0b11100000000000000000000000000000
	SilverRupeesDodongoCavern                       = 0b00000000000000000000000000000001
	SilverRupeesIceCavernScythe                     = 0b00000000000000000000000000000010
	SilverRupeesIceCavernPush                       = 0b00000000000000000000000000000100
	SilverRupeesWellBasement                        = 0b00000000000000000000000000010000
	SilverRupeesShadowShortcut                      = 0b00000000000000000000000000100000
	SilverRupeesShadowBlades                        = 0b00000000000000000000000001000000
	SilverRupeesShadowHugePit                       = 0b00000000000000000000000010000000
	SilverRupeesShadowSpikes                        = 0b00000000000000000000000100000000
	SilverRupeesTrainingSlopes                      = 0b00000000000000000000001000000000
	SilverRupeesTrainingLava                        = 0b00000000000000000000010000000000
	SilverRupeesTrainingWater                       = 0b00000000000000000000100000000000
	SilverRupeesSpiritTorches                       = 0b00000000000000000001000000000000
	SilverRupeesSpiritBoulders                      = 0b00000000000000000010000000000000
	SilverRupeesSpiritSunBlock                      = 0b00000000000000000100000000000000
	SilverRupeesSpiritAdultClimb                    = 0b00000000000000001000000000000000
	SilverRupeesTowerForestTrial                    = 0b00000000000000010000000000000000
	SilverRupeesTowerFireTrial                      = 0b00000000000000100000000000000000
	SilverRupeesTowerWaterTrial                     = 0b00000000000001000000000000000000
	SilverRupeesTowerShadowTrial                    = 0b00000000000010000000000000000000
	SilverRupeesTowerSpiritTrial                    = 0b00000000000100000000000000000000
	SilverRupeesTowerLightTrial                     = 0b00000000001000000000000000000000
)

const (
	MapsCompassesVanilla                  = MapsCompasses(KeysVanilla)
	MapsCompassesRemove                   = MapsCompasses(KeysRemove)
	MapsCompassesDungeon                  = MapsCompasses(KeysDungeon)
	MapsCompassesRegional                 = MapsCompasses(KeysRegional)
	MapsCompassesOverworld                = MapsCompasses(KeysOverworld)
	MapsCompassesAnyDungeon               = MapsCompasses(KeysAnyDungeon)
	MapsCompassesAnywhere                 = MapsCompasses(KeysAnywhere)
	MapsCompassesStartWith                = MapsCompasses(1 << 8)
	MapsCompassesEnhanced   MapsCompasses = 0x0F00
)

const (
	KokriForestClosed OpenForest = iota
	KokriForestOpen
	KokriForestDekuClosed
)

const (
	KakGateClosed OpenKak = iota
	KakGateLetter
	KakGateOpen
)

const (
	ZoraFountainClosed OpenZoraFountain = iota
	// move KZ for adult but not for child
	ZoraFountainOpenAdult
	// if KZ is moved for child, he's moved for adult
	ZoraFountainOpenAlways
)

func (z OpenZoraFountain) String() string {
	switch z {
	case ZoraFountainClosed:
		return "closed"
	case ZoraFountainOpenAdult:
		return "adult"
	case ZoraFountainOpenAlways:
		return "open"
	default:
		panic(slipup.Createf("unknown zora fountain setting %x", uint8(z)))
	}
}

const (
	GerudoFortressNormal GerudoForestCarpenters = iota
	GerudoFortressFast
	GerudoFortressOpen
)

const (
	ShortcutsOff    DungeonShortcuts = 0b0000000000000000
	ShortcutsAll                     = 0b0000000011111111
	ShortcutsRandom                  = 0b1010101010101010
	ShortcutsDeku                    = 0b0000000000000001
	ShortcutsCavern                  = 0b0000000000000010
	ShortcutsJabu                    = 0b0000000000000100
	ShortcutsForest                  = 0b0000000000001000
	ShortcutsFire                    = 0b0000000000010000
	ShortcutsWater                   = 0b0000000000100000
	ShortcutsShadow                  = 0b0000000001000000
	ShortcutsSpirit                  = 0b0000000010000000
)

const (
	StartAgeChild StartingAge = iota
	StartAgeAdult
	StartAgeRandom
)

const (
	// mask for upper 4 bits
	MasterQuestDungeonsNone     MasterQuestDungeons = 0b0000000000000000
	MasterQuestDungeonsAll                          = 0b1000000000000000
	MasterQuestDungeonsSpecific                     = 0b0100000000000000
	MasterQuestDungeonsCount                        = 0b0010000000000000
	MasterQuestDungeonsRandom                       = 0b0001000000000000
	// Specific Dungeons
	MasterQuestDekuTree       = 0b0000000000000001
	MasterQuestDodongoCavern  = 0b0000000000000010
	MasterQuestJabu           = 0b0000000000000100
	MasterQuestForest         = 0b0000000000001000
	MasterQuestFire           = 0b0000000000010000
	MasterQuestWater          = 0b0000000000100000
	MasterQuestShadow         = 0b0000000001000000
	MasterQuestSpirit         = 0b0000000010000000
	MasterQuestWell           = 0b0000000100000000
	MasterQuestIceCavern      = 0b0000001000000000
	MasterQuestTrainingGround = 0b0000010000000000
	MasterQuestGanonsCastle   = 0b0000100000000000
)

const (
	CompletedDungeonsNone     CompletedDungeons = 0b0000000000000000
	CompletedDungeonsSpecific                   = 0b1000000000000000
	CompletedDungeonsRewards                    = 0b0100000000000000
	CompletedDungeonsCount                      = 0b0010000000000000

	// specific dungeons
	CompletedDekuTree      = 0b0000000000000001
	CompletedDodongoCavern = 0b0000000000000010
	CompletedJabu          = 0b0000000000000100
	CompletedForest        = 0b0000000000001000
	CompletedFire          = 0b0000000000010000
	CompletedWater         = 0b0000000000100000
	CompletedShadow        = 0b0000000001000000
	CompletedSpirit        = 0b0000000010000000
)

const (
	InteriorShuffleOff    InteriorShuffle = 0
	InteriorShuffleSimple                 = 2
	InteriorShuffleAll                    = 4
)

const (
	DungeonEntranceShuffleOff    DungeonEntranceShuffle = 0
	DungeonEntranceShuffleSimple                        = 2
	DungeonEntranceShuffleAll                           = 4
)

const (
	BossShuffleOff    BossShuffle = 0
	BossShuffleSimple             = 2
	BossShuffleAll                = 4
)

const (
	SpawnVanilla     Spawn = 0
	RandomSpawn            = 0xF0F0F0F0F0F0F0F00000000000000000
	SetSpawnLocation       = 0xFFFFFFFFFFFFFFFF0000000000000000
)

const (
	ShuffleSongsOnSong ShuffleSongs = iota
	ShuffleSongsOnRewards
	ShuffleSongsAnywhere
)

const (
	// upper bits mask -- how are shops shuffled
	ShuffleShopsOff           ShuffleShops = 0
	ShuffleShopsSpecialRandom              = 0b01010101 << 8
	ShuffleShopsSpecial0                   = 0b11000000 << 8
	ShuffleShopsSpecial1                   = 0b10100000 << 8
	ShuffleShopsSpecial2                   = 0b10010000 << 8
	ShuffleShopsSpecial3                   = 0b10001000 << 8
	ShuffleShopsSpecial4                   = 0b10000100 << 8

	// lower bit mask -- do we have shop price caps
	ShuffleShopPricesRandom       ShuffleShops = 0b0000000001010101
	ShuffleShopPricesStartWallet               = 0b0000000000000011
	ShuffleShopPricesAdultWallet               = 0b0000000000000110
	ShuffleShopPricesGiantWallet               = 0b0000000000001010
	ShuffleShopPricesTycoonWallet              = 0b0000000000010010
	ShuffleShopPricesAffordable                = 0b0000000000100010
)

const (
	ShuffleGoldTokenOff       ShuffleTokens = 0
	ShuffleGoldTokenDungeons                = 1
	ShuffleGoldTokenOverworld               = 2
)

const (
	ShuffleScrubsOff         ShuffleScrubs = 0 // off off
	ShuffleScrubsUpgradeOnly               = 1 // OOTR off
	ShuffleScrubsAffordable                = 2
	ShuffleScrubsExpensive                 = 3
	ShuffleScrubsRandom                    = 4
)

const (
	ShuffleFreestandingsOff       ShuffleFreestandings = 0
	ShuffleFreestandingsDungeon                        = 1
	ShuffleFreestandingsOverworld                      = 2
)

const (
	ShufflePotsOff       ShufflePots = 0
	ShuffleEmptyPots                 = 1
	ShufflePotsDungeons              = 2
	ShufflePotsOverworld             = 4
)

const (
	ShuffleCratesOff       ShuffleCrates = 0
	ShuffleEmptyCrates                   = 1
	ShuffleCratesDungeons                = 2
	ShuffleCratesOverworld               = 4
)

const (
	ShuffleLoachRewardOff     ShuffleLoachReward = 0
	ShuffleLoachRewardVanilla                    = 1
	ShuffleLoachRewardEasy                       = 2
)

const (
	ShuffleSongPatternsOff   ShuffleSongPatterns = 0
	ShuffleSongPatternsFrogs                     = 1
	ShuffleSongPatternsWarps                     = 2
)

const (
	HintsRevealedNever  HintsRevealed = 0
	HintsRevealedMask                 = 1
	HintsRevealedStone                = 2
	HintsRevealedAlways               = 4
)

const (
	DamageMultiplierHalf   DamageMultiplier = 0
	DamageMultiplierNormal                  = 1
	DamageMultiplierDouble                  = 2
	DamageMultiplierQuad                    = 4
	DamageMultiplierOhko                    = 8
)

func (m DamageMultiplier) String() string {
	switch m {
	case DamageMultiplierHalf:
		return "half"
	case DamageMultiplierNormal:
		return "normal"
	case DamageMultiplierDouble:
		return "double"
	case DamageMultiplierQuad:
		return "quadruple"
	case DamageMultiplierOhko:
		return "ohko"
	default:
		panic(slipup.Createf("unknown damage multiple %x", uint(m)))
	}
}

const (
	BonkDamageNone   BonkDamage = 0
	BonkDamageHalf              = 1
	BonkDamageNormal            = 2
	BonkDamageDouble            = 4
	BonkDamageQuad              = 8
	BonkDamageOhko              = 16
)

func (m BonkDamage) String() string {
	switch m {
	case BonkDamageNone:
		return "none"
	case BonkDamageHalf:
		return "half"
	case BonkDamageNormal:
		return "normal"
	case BonkDamageDouble:
		return "double"
	case BonkDamageQuad:
		return "quadruple"
	case BonkDamageOhko:
		return "ohko"
	default:
		panic(slipup.Createf("unknown bonk damage %x", uint(m)))
	}
}

const (
	StartingTimeOfDayDefault StartingTimeOfDay = iota
	StartingTimeOfDayRandom
	StartingTimeOfDaySunrise
	StartingTimeOfDayMorning
	StartingTimeOfDayNoon
	StartingTimeOfDayAfternoon
	StartingTimeOfDaySunset
	StartingTimeOfDayEvening
	StartingTimeOfDayMidnight
	StartingTimeOfDayWitching
)

const (
	ItemPoolMinimal ItemPool = iota
	ItemPoolScarce
	ItemPoolDefault
	ItemPoolPlentiful
	ItemPoolLudicrous
)

const (
	IceTrapsOff IceTraps = 0
	IceTrapsNormal
	IceTrapsSomeExtraJunk
	IceTrapsAllExtraJunk
	IceTrapsAllJunk
)

const (
	AdultTradeShuffle ShuffleTradeAdult = 0b1000000000000000

	AdultTradeStartPocketEgg    = 0b0000000000000001
	AdultTradeStartPocketCucco  = 0b0000000000000010
	AdultTradeStartCojiro       = 0b0000000000000100
	AdultTradeStartOddMushroom  = 0b0000000000001000
	AdultTradeStartOddPotion    = 0b0000000000010000
	AdultTradeStartPoachersSaw  = 0b0000000000100000
	AdultTradeStartBrokenSword  = 0b0000000001000000
	AdultTradeStartPrescription = 0b0000000010000000
	AdultTradeStartEyeballFrog  = 0b0000000100000000
	AdultTradeStartEyedrops     = 0b0000001000000000
	AdultTradeStartClaimCheck   = 0b0000010000000000
)

const (
	ChildTradeShuffle  ShuffleTradeChild = 0b1000000000000000
	ChildTradeComplete                   = 0b1111111111111111

	ChildTradeStartWeirdEgg   = 0b0000000000000001
	ChildTradeStartChicken    = 0b0000000000000010
	ChildTradeStartLetter     = 0b0000000000000100
	ChildTradeStartMaskKeaton = 0b0000000000001000
	ChildTradeStartMaskSkull  = 0b0000000000010000
	ChildTradeStartMaskSpooky = 0b0000000000100000
	ChildTradeStartMaskBunny  = 0b0000000001000000
	ChildTradeStartMaskGoron  = 0b0000000010000000
	ChildTradeStartMaskZora   = 0b0000000100000000
	ChildTradeStartMaskGerudo = 0b0000001000000000
	ChildTradeStartMaskTruth  = 0b0000010000000000
)

const (
	ForestTemplePoesNone ForestTemplePoes = 0
	ForestTempleAmyMeg                    = 1
	ForestTempleJoBeth                    = 2
)

func HasFlag[F ~uint8 | ~uint16 | ~uint32 | ~uint64](value, flag F) bool {
	return value&flag == flag
}
