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
type DungeonDoorShuffle uint8

type Enum uint8

type Flags uint64
type Flag uint64

func (this Flags) Count() int {
	return bits.OnesCount64(uint64(this))
}

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

func (this ConditionedAmount) Decode() (ConditionKind, int) {
	kind := ConditionKind(this >> 32)
	qty := int(uint32(kind))
	return kind, qty
}

const (
	LogicNone       LogicSetting = 1
	LogicGlitched   LogicSetting = 2
	LogicGlitchless LogicSetting = 3

	ReachableAll       LocationsReachable = 1
	ReachableGoals     LocationsReachable = 2
	ReachableNecessary LocationsReachable = 3

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

	_ Flag = iota
	Forest
	Fire
	Water
	Shadow
	Spirit
	Light
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

type OpenForest uint8

const (
	_ OpenForest = iota
	KokriForestOpen
	KokriForestClosedDeku
	KokriForestClosed
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
	_ OpenKakarikoGate = iota
	KakGateOpen
	KakGateZelda
	KakGateClosed
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
	_ OpenZoraFountain = iota
	ZoraFountainClosed
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
	_ GerudoFortressCarpenterRescue = iota
	RescueAllCarpenters
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
	_ ShuffleScrub = iota
	ShuffleUpgradeScrub
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
	ShuffleDungeonRewardsVanilla ShuffleDungeonRewards = iota
	ShuffleDungeonRewardsReward
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
	ShuffleKeysVanilla ShuffleKeys = iota
	ShuffleKeysRemove
	ShuffleKeyOwnDungeon
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
	AdultTradeStartPocketEgg AdultTradeItems = 1 << iota
	AdultTradeStartPocketCucco
	AdultTradeStartOddMushroom
	AdultTradeStartOddPotion
	AdultTradeStartPoachersSaw
	AdultTradeStartBrokenSword
	AdultTradeStartPrescription
	AdultTradeStartEyeballFrog
	AdultTradeStartEyedrops
	AdultTradeStartClaimCheck
)

type StartAge bool

const (
	StartAgeAdult StartAge = true
	StartAgeChild StartAge = false
)
