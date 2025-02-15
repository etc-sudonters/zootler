package settings

import (
	"errors"
	"fmt"
	"math/bits"
)

type LogicSetting uint8

func (this LogicSetting) String() string {
	switch this {
	case LogicNone:
		return "none"
	case LogicGlitched:
		return "glitched"
	case LogicGlitchless:
		return "glitchless"
	default:
		panic("unreachable")
	}
}

type LocationsReachable uint8

func (this LocationsReachable) String() string {
	switch this {
	case ReachableAll:
		return "all"
	case ReachableGoals:
		return "goal"
	case ReachableNecessary:
		return "beatable"
	default:
		panic("unreachable")
	}
}

type ConditionedAmount uint64
type ConditionKind uint32

type BossShuffle uint8
type InteriorShuffle uint8
type OverworldShuffle uint8
type DungeonDoorShuffle uint8

type Enum uint8

type Flags uint64
type Flag uint64

func CountFlags[U ~uint64](union U) int {
	return bits.OnesCount64(uint64(union))
}

func HasFlag[U ~uint64](union U, flag U) bool {
	return uint64(union)&uint64(flag) == uint64(flag)
}

type NavigationShuffle uint8

func EncodeConditionedAmount(kind ConditionKind, qty uint32) ConditionedAmount {
	return ConditionedAmount(uint64(kind)<<32 | uint64(qty))
}

func (this ConditionedAmount) Kind() ConditionKind {
	return ConditionKind(uint64(this) >> 32)
}

func (this ConditionedAmount) Amount() uint32 {
	return uint32(this)
}

func (this ConditionedAmount) Decode() (ConditionKind, uint32) {
	kind := ConditionKind(this >> 32)
	qty := uint32(kind)
	return kind, qty
}

type TrialFlag uint64

const (
	LogicGlitchless LogicSetting = iota
	LogicNone
	LogicGlitched

	ReachableAll LocationsReachable = iota
	ReachableGoals
	ReachableNecessary

	CondUnitialized ConditionKind = iota
	CondDefault
	CondMedallions
	CondStones
	CondRewards
	CondTokens
	CondHearts
	CondOpen
	CondVanilla
	CondTriforce

	TrialForest TrialFlag = 1 << iota
	TrialFire
	TrialWater
	TrialShadow
	TrialSpirit
	TrialLight
	TrialAll = TrialForest | TrialFire | TrialWater | TrialShadow | TrialSpirit | TrialLight
)

func (this ConditionKind) String() string {
	switch this {
	case CondUnitialized:
		panic(errors.New("uninitialized condition"))
	case CondMedallions:
		return "medallions"
	case CondStones:
		return "stones"
	case CondRewards:
		return "rewards"
	case CondTokens:
		return "tokens"
	case CondHearts:
		return "hearts"
	case CondDefault:
		return "default"
	case CondOpen:
		return "open"
	case CondVanilla:
		return "vanilla"
	case CondTriforce:
		return "triforce"
	default:
		panic(fmt.Errorf("unknown condition flag %x", uint8(this)))
	}
}

func ParseCondition(raw string) (ConditionKind, error) {
	switch raw {
	case "medallions":
		return CondMedallions, nil
	case "stones":
		return CondStones, nil
	case "rewards":
		return CondRewards, nil
	case "tokens":
		return CondTokens, nil
	case "hearts":
		return CondHearts, nil
	case "default":
		return CondDefault, nil
	case "open":
		return CondOpen, nil
	case "vanilla":
		return CondVanilla, nil
	case "triforce":
		return CondTriforce, nil
	default:
		return CondUnitialized, fmt.Errorf("unknown condition %q", raw)
	}
}

func ConditionFrom(raw string, qty uint32) (ConditionedAmount, error) {
	cond, err := ParseCondition(raw)
	if err != nil {
		return 0, err
	}

	return EncodeConditionedAmount(cond, qty), nil
}

type OpenForest uint8

const (
	KokriForestClosed OpenForest = iota
	KokriForestOpen
	KokriForestClosedDeku
)

func (this OpenForest) String() string {
	switch this {
	case KokriForestOpen:
		return "open"
	case KokriForestClosedDeku:
		return "closed_deku"
	case KokriForestClosed:
		return "closed"
	default:
		panic("unreachable")
	}
}

type OpenKakarikoGate uint8

const (
	KakGateClosed OpenKakarikoGate = iota
	KakGateOpen
	KakGateZelda
)

func (this OpenKakarikoGate) String() string {
	switch this {
	case KakGateOpen:
		return "open"
	case KakGateZelda:
		return "zelda"
	case KakGateClosed:
		return "closed"
	default:
		panic("unreachable")
	}
}

type OpenZoraFountain uint8

const (
	ZoraFountainClosed OpenZoraFountain = iota
	ZoraFountainOpenAdult
	ZoraFountainOpen
)

func (this OpenZoraFountain) String() string {
	switch this {
	case ZoraFountainClosed:
		return "closed"
	case ZoraFountainOpenAdult:
		return "adult"
	case ZoraFountainOpen:
		return "open"
	default:
		panic("unreachable")
	}
}

type GerudoFortressCarpenterRescue uint8

const (
	RescueAllCarpenters GerudoFortressCarpenterRescue = iota
	RescueOneCarpenters
	RescueZeroCarpenters
)

func (this GerudoFortressCarpenterRescue) String() string {
	switch this {
	case RescueAllCarpenters:
		return "normal"
	case RescueOneCarpenters:
		return "fast"
	case RescueZeroCarpenters:
		return "open"
	default:
		panic("unreachable")
	}
}

type ShuffleScrub uint8

const (
	ShuffleUpgradeScrub ShuffleScrub = iota
	ShuffleScrubsAffordable
	ShuffleScrubsExpensive
	ShuffleScrubsRandomPrices
)

func (this ShuffleScrub) String() string {
	switch this {
	case ShuffleUpgradeScrub:
		return "off"
	case ShuffleScrubsAffordable:
		return "low"
	case ShuffleScrubsExpensive:
		return "regular"
	case ShuffleScrubsRandomPrices:
		return "random"
	default:
		panic("unreachable")
	}
}

type PartitionedShuffle uint8

type ShuffleSkullTokens PartitionedShuffle
type ShufflePots PartitionedShuffle
type ShuffleCrates PartitionedShuffle
type ShuffleFreestanding PartitionedShuffle

const (
	ShufflePartionOff PartitionedShuffle = iota
	ShufflePartitionDungeons
	ShufflePartitionOverworld
	ShufflePartionAll
)

func (this PartitionedShuffle) String() string {
	switch this {
	case ShufflePartionOff:
		return "off"
	case ShufflePartitionDungeons:
		return "dungeons"
	case ShufflePartitionOverworld:
		return "overworld"
	case ShufflePartionAll:
		return "all"
	default:
		panic("unreachable")
	}

}

type ShuffleDungeonRewards uint8

const (
	ShuffleDungeonRewardsReward ShuffleDungeonRewards = iota
	ShuffleDungeonRewardsVanilla
	ShuffleDungeonRewardsOwnDungeon
	ShuffleDungeonRewardsRegional
	ShuffleDungeonRewardsOverworld
	ShuffleDungeonRewardsAnyDungeon
	ShuffleDungeonRewardsAnywhere
)

func (this ShuffleDungeonRewards) String() string {

	switch this {
	case ShuffleDungeonRewardsVanilla:
		return "vanilla"
	case ShuffleDungeonRewardsReward:
		return "reward"
	case ShuffleDungeonRewardsOwnDungeon:
		return "dungeon"
	case ShuffleDungeonRewardsRegional:
		return "regional"
	case ShuffleDungeonRewardsOverworld:
		return "overworld"
	case ShuffleDungeonRewardsAnyDungeon:
		return "any_dungeon"
	case ShuffleDungeonRewardsAnywhere:
		return "anywhere"
	default:
		panic("unreachable")
	}
}

type ShuffleKeys uint8

const (
	ShuffleKeyOwnDungeon ShuffleKeys = iota
	ShuffleKeysVanilla
	ShuffleKeysRemove
	ShuffleKeyRegional
	ShuffleKeyOverworld
	ShuffleKeyAnyDungeon
	ShuffleKeysAnywhere
)

func (this ShuffleKeys) String() string {
	switch this {
	case ShuffleKeysVanilla:
		return "vanilla"
	case ShuffleKeysRemove:
		return "remove"
	case ShuffleKeyOwnDungeon:
		return "dungeon"
	case ShuffleKeyOverworld:
		return "overworld"
	case ShuffleKeyAnyDungeon:
		return "any_dungeon"
	case ShuffleKeysAnywhere:
		return "keysanity"
	default:
		panic("unreachable")
	}
}

type GanonBossKeyShuffle uint8

const (
	GanonBossKeyInDungeon GanonBossKeyShuffle = iota
	GanonBossKeyRemove
	GanonBossKeyVanilla
	GanonBossKeyRegional
	GanonBossKeyOverworld
	GanonBossKeyAnyDungeon
	GanonBossKeyAnywhere
	GanonBossKeyOnLacs
	GanonBossKeyStones
	GanonBossKeyMedallions
	GanonBossKeyTokens
	GanonBossKeyHearts
)

func (this GanonBossKeyShuffle) String() string {
	switch this {
	case GanonBossKeyInDungeon:
		return "dungeon"
	case GanonBossKeyRemove:
		return "remove"
	case GanonBossKeyVanilla:
		return "vanilla"
	case GanonBossKeyRegional:
		return "regional"
	case GanonBossKeyOverworld:
		return "overworld"
	case GanonBossKeyAnyDungeon:
		return "any_dungeon"
	case GanonBossKeyAnywhere:
		return "keysanity"
	case GanonBossKeyOnLacs:
		return "on_lacs"
	case GanonBossKeyStones:
		return "stones"
	case GanonBossKeyMedallions:
		return "medallions"
	case GanonBossKeyTokens:
		return "tokens"
	case GanonBossKeyHearts:
		return "hearts"
	default:
		panic("unreachable")
	}

}

func ParseGanonBossKeyShuffle(raw string) (GanonBossKeyShuffle, error) {
	switch raw {
	case "dungeon":
		return GanonBossKeyInDungeon, nil
	case "remove":
		return GanonBossKeyRemove, nil
	case "vanilla":
		return GanonBossKeyVanilla, nil
	case "regional":
		return GanonBossKeyRegional, nil
	case "overworld":
		return GanonBossKeyOverworld, nil
	case "any_dungeon":
		return GanonBossKeyAnyDungeon, nil
	case "keysanity":
		return GanonBossKeyAnywhere, nil
	case "on_lacs":
		return GanonBossKeyOnLacs, nil
	case "stones":
		return GanonBossKeyStones, nil
	case "medallions":
		return GanonBossKeyMedallions, nil
	case "tokens":
		return GanonBossKeyTokens, nil
	case "hearts":
		return GanonBossKeyHearts, nil
	default:
		return GanonBossKeyInDungeon, fmt.Errorf("unknown ganon boss key shuffle: %q", raw)
	}
}

type HintsRevealed uint8

const (
	HintsRevealedNever HintsRevealed = iota
	HintsRevealedWithMask
	HintsRevealedWithStone
	HintsRevealedAlways
)

func (this HintsRevealed) String() string {
	switch this {
	case HintsRevealedNever:
		return "none"
	case HintsRevealedWithMask:
		return "mask"
	case HintsRevealedWithStone:
		return "agony"
	case HintsRevealedAlways:
		return "always"
	default:
		panic("unreachable")
	}
}

type DamageMultiplier uint8
type DeadlyBonks DamageMultiplier

const (
	DamageMultiplierNormal DamageMultiplier = iota
	DamageMultiplierHalf
	DamageMultiplierDouble
	DamageMultiplierQuadruple
	DamageMultiplierOHKO
)

func (this DamageMultiplier) String() string {
	switch this {
	case DamageMultiplierNormal:
		return "normal"
	case DamageMultiplierHalf:
		return "half"
	case DamageMultiplierDouble:
		return "double"
	case DamageMultiplierQuadruple:
		return "quadruple"
	case DamageMultiplierOHKO:
		return "ohko"
	default:
		panic("unreachable")
	}
}

func (this DeadlyBonks) String() string {
	return DamageMultiplier(this).String()
}

type GerudoFortressHeartPiece uint

const (
	GerudoFortressHeartPieceRemove GerudoFortressHeartPiece = iota
	GerudoFortressHeartPieceVanilla
	GerudoFortressHeartPieceShuffle
)

func (this GerudoFortressHeartPiece) String() string {
	switch this {
	case GerudoFortressHeartPieceRemove:
		return "remove"
	case GerudoFortressHeartPieceVanilla:
		return "vanilla"
	case GerudoFortressHeartPieceShuffle:
		return "shuffle"
	default:
		panic("unreachable")
	}
}

type ShufflingFlags uint64

const (
	ShuffleEmptyPots ShufflingFlags = 1 << iota
	ShuffleEmptyCrates
	ShuffleCows
	ShuffleOcarinaNotes
	ShuffleBeehives
	ShuffleWonderItems
	ShuffleKokiriSword
	ShuffleOcarinas
	ShuffleGerudoCard
	ShuffleBeans
	ShuffleExpensiveMerchants
	ShuffleFrogRupees
	ShuffleLoachReward
)

type TimeOfDay uint8

const (
	_ TimeOfDay = iota
	TimeOfDaySunrise
	TimeOfDayMorning
	TimeOfDayNoon
	TimeOfDayAfternoon
	TimeOfDaySunset
	TimeOfDayEvening
	TimeOfDayMidnight
	TimeOfDayWitching
)

func (this TimeOfDay) IsNight() bool {
	switch this {
	case TimeOfDayEvening,
		TimeOfDaySunset,
		TimeOfDayMidnight,
		TimeOfDayWitching:
		return true
	default:
		return false
	}

}

type ChildTradeItems uint64

type AdultTradeItems uint64

const (
	AdultTradePocketEgg AdultTradeItems = 1 << iota
	AdultTradePocketCucco
	AdultTradeOddMushroom
	AdultTradeOddPotion
	AdultTradePoachersSaw
	AdultTradeBrokenSword
	AdultTradePrescription
	AdultTradeEyeballFrog
	AdultTradeEyedrops
	AdultTradeClaimCheck
	AdultTradeItemsAll = AdultTradePocketEgg | AdultTradePocketCucco |
		AdultTradeOddMushroom | AdultTradeOddPotion |
		AdultTradePoachersSaw | AdultTradeBrokenSword |
		AdultTradePrescription | AdultTradeEyeballFrog |
		AdultTradeEyedrops | AdultTradeClaimCheck
)

type StartAge bool

const (
	StartAgeAdult StartAge = true
	StartAgeChild StartAge = false
)

type LocationFlags uint64

const (
	LocationSkipRauruReward LocationFlags = 1 << iota
	LocationSkipChildZelda
	LocationsFreeScarecrow
	LocationsPlantBeans
	LocationsCompleteMaskQuest
)

type ShuffleSongs uint8

const (
	ShuffleSongsOnSongs ShuffleSongs = iota
	ShuffleSongsOnDungeonRewards
	ShuffleSongsAnywhere
)

type ShuffleShop uint8

const (
	ShuffleShopsOff ShuffleShop = iota
	ShuffleShopsZero
	ShuffleShopsOne
	ShuffleShopTwo
	ShuffleShopThree
	ShuffleShopFour
)

type ShuffleMapCompass uint8

const (
	ShuffleMapCompassDungeon ShuffleMapCompass = iota
	ShuffleMapCompassRemove
	ShuffleMapCompassStartWith
	ShuffleMapCompassVanilla
	ShuffleMapCompassRegional
	ShuffleMapCompassOverworld
	ShuffleMapCompassAnyDungeon
	ShuffleMapCompassAnywhere
)

type ConnectionFlag uint64

const (
	ConnectionOpenDoorOfTime ConnectionFlag = 1 << iota
	ConnectionShuffleHideoutEntrances
)
