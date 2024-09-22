package shiro

import (
	"strings"
	"sudonters/zootler/icearrow/compiler"
	"sudonters/zootler/internal"
	"sudonters/zootler/internal/settings"

	"github.com/etc-sudonters/substrate/slipup"
)

func intoIntrinsics(settings *settings.ZootrSettings) intrinsics {
	var i intrinsics
	i.settings = settings
	i.bools = make(map[string]bool)

	for trick, enabled := range i.settings.Tricks.Enabled {
		i.bools[string(internal.Normalize(trick))] = enabled
	}

	return i
}

type settingResolver = func(compiler.Invocation, *compiler.Symbol, *compiler.SymbolTable) (bool, bool)

func intoIntrinsicFunc(resolver settingResolver) compiler.Intrinsic {
	return func(ct compiler.Invocation, sym *compiler.Symbol, st *compiler.SymbolTable) (compiler.CompileTree, error) {
		result, fine := resolver(ct, sym, st)
		if !fine {
			return ct, slipup.Createf("something weird with this call %#v", ct)
		}

		immed := compiler.Immediate{Value: result}

		if result {
			immed.Kind = compiler.CT_IMMED_TRUE
		} else {
			immed.Kind = compiler.CT_IMMED_FALSE
		}

		return immed, nil
	}

}

type intrinsics struct {
	settings *settings.ZootrSettings
	bools    map[string]bool
}

func (si *intrinsics) CompareEq(ct compiler.Invocation, sym *compiler.Symbol, st *compiler.SymbolTable) (bool, bool) {
	name := settingname(ct.Args[0].(compiler.Load), st)
	switch name {
	case "deadly_bonks", "deadlybonks":
		comperand := st.String(ct.Args[1].(compiler.Load).Id)
		flags := map[string]settings.BonkDamage{
			"ohko":      settings.BonkDamageOhko,
			"quadruple": settings.BonkDamageQuad,
		}
		compareTo, exists := flags[comperand.Value]
		if exists {
			return si.settings.Damage.Bonk == compareTo, true
		}
		panic(slipup.Createf("unknown bonk damage value %q", comperand.Value))

	case "damage_multiplier", "damagemultiplier":
		comperand := st.String(ct.Args[1].(compiler.Load).Id)
		flags := map[string]settings.DamageMultiplier{
			"ohko":      settings.DamageMultiplierOhko,
			"quadruple": settings.DamageMultiplierQuad,
		}
		compareTo, exists := flags[comperand.Value]
		if exists {
			return si.settings.Damage.Multiplier == compareTo, true
		}
		panic(slipup.Createf("unknown damage multiplier value %q", comperand.Value))

	case "bridge":
		comperand := st.String(ct.Args[1].(compiler.Load).Id)
		cond, _ := settings.Decode(si.settings.BridgeCondition)
		conds := map[string]settings.Condition{
			"open":       settings.CondOpen,
			"vanilla":    settings.CondDefault,
			"stones":     settings.CondStones,
			"medallions": settings.CondMedallions,
			"dungeons":   settings.CondRewards,
			"tokens":     settings.CondTokens,
			"hearts":     settings.CondHearts,
		}
		compareTo, exists := conds[comperand.Value]
		if exists {
			return cond == compareTo, true
		}
		panic(slipup.Createf("unknown bridge condition value %q", comperand.Value))

	case "logic_rules", "logicrules":
		comperand := st.String(ct.Args[1].(compiler.Load).Id)
		desired, exists := map[string]settings.LogicRuleSet{
			"glitchless": settings.LogicGlitchess,
			"glitched":   settings.LogicGlitched,
			"none":       settings.LogicNone,
		}[comperand.Value]
		if !exists {
			panic(slipup.Createf("unknown logic set value %s", comperand.Value))
		}
		return si.settings.LogicRules == desired, true

	case "open_forest", "openforest":
		comperand := st.String(ct.Args[1].(compiler.Load).Id)
		forestState, exists := map[string]settings.OpenForest{
			"closed_deku": settings.KokriForestDekuClosed,
			"closed":      settings.KokriForestClosed,
			"open":        settings.KokriForestOpen,
		}[comperand.Value]
		if !exists {
			panic(slipup.Createf("unknown open forest value %s", comperand.Value))
		}
		return si.settings.Locations.KokriForest == forestState, true

	case "gerudo_fortress", "gerudofortress":
		comperand := st.String(ct.Args[1].(compiler.Load).Id)
		fortressState, exists := map[string]settings.GerudoFortress{
			"normal": settings.GerudoFortressNormal,
			"fast":   settings.GerudoFortressFast,
			"open":   settings.GerudoFortressOpen,
		}[comperand.Value]
		if !exists {
			panic(slipup.Createf("unknown gerudo fortress value %s", comperand.Value))
		}
		return si.settings.Locations.GerudoFortress == fortressState, true

	case "shuffle_pots", "shufflepots":
		comperand := st.String(ct.Args[1].(compiler.Load).Id)
		shuffling, exists := map[string]settings.ShufflePots{
			"off": settings.ShufflePotsOff,
		}[comperand.Value]
		if !exists {
			panic(slipup.Createf("unknown shuffle pots value %s", comperand.Value))
		}
		return settings.HasFlag(si.settings.Shuffling.Pots, shuffling), true

	case "open_kakariko", "openkakariko":
		comperand := st.String(ct.Args[1].(compiler.Load).Id)
		kakState, exists := map[string]settings.OpenKak{
			"open":   settings.KakGateOpen,
			"closed": settings.KakGateClosed,
			"zelda":  settings.KakGateLetter,
		}[comperand.Value]
		if !exists {
			panic(slipup.Createf("unknown gerudo fortress value %s", comperand.Value))
		}
		return si.settings.Locations.Kakariko == kakState, true

	case "selected_adult_trade_item", "selectedadulttradeitem":
		comperand := st.String(ct.Args[1].(compiler.Load).Id)
		tradeItem, exists := map[string]settings.ShuffleTradeAdult{
			"Odd Potion":   settings.AdultTradeStartOddPotion,
			"Poachers Saw": settings.AdultTradeStartPoachersSaw,
			"Broken Sword": settings.AdultTradeStartBrokenSword,
			"Prescription": settings.AdultTradeStartPrescription,
			"Eyeball Frog": settings.AdultTradeStartEyeballFrog,
			"Eyedrops":     settings.AdultTradeStartEyedrops,
			"Claim Check":  settings.AdultTradeStartClaimCheck,
		}[comperand.Value]
		if !exists {
			panic(slipup.Createf("unknown adult trade item %s", comperand.Value))
		}
		return si.settings.Trades.Adult == tradeItem, true

	case "shuffle_scrubs", "shufflescrubs":
		comperand := st.String(ct.Args[1].(compiler.Load).Id)
		if comperand.Value != "off" {
			return si.settings.Shuffling.Scrubs == settings.ShuffleScrubsOff ||
				si.settings.Shuffling.Scrubs == settings.ShuffleScrubsUpgradeOnly, true
		}
		return si.settings.Shuffling.Scrubs != settings.ShuffleScrubsOff &&
			si.settings.Shuffling.Scrubs != settings.ShuffleScrubsUpgradeOnly, true

	case "shuffle_overworld_entrances", "shuffleoverworldentrances":
		comperand := st.String(ct.Args[1].(compiler.Load).Id)
		return si.settings.Entrances.Overworld && comperand.Value == "off", true

	case "shuffle_tcgkeys", "shuffletcgkeys":
		comperand := st.String(ct.Args[1].(compiler.Load).Id)
		choice, exists := map[string]settings.KeyShuffle{
			"vanilla": settings.KeysVanilla,
		}[comperand.Value]
		if !exists {
			panic(slipup.Createf("unknown tcg key shuffle value %s", comperand.Value))
		}
		return si.settings.KeyShuffle.ChestGameKeys == choice, true

	case "shuffle_ganon_bosskey", "shuffleganonbosskey":
		comperand := st.String(ct.Args[1].(compiler.Load).Id)
		choice, exists := map[string]settings.GanonBKShuffleKind{
			"remove":      settings.GanonBKRemove,
			"vanilla":     settings.GanonBKVanilla,
			"dungeon":     settings.GanonBKDungeon,
			"regional":    settings.GanonBKRegional,
			"overworld":   settings.GanonBKOverworld,
			"any_dungeon": settings.GanonBKAnyDungeon,
			"keysanity":   settings.GanonBKKeysanity,
			"on_lacs":     settings.GanonBKOnLacs,
			"stones":      settings.GanonBKStones,
			"medallions":  settings.GanonBKMedallions,
			"dungeons":    settings.GanonBKDungeonRewards,
			"tokens":      settings.GanonBKTokens,
			"hearts":      settings.GanonBKHearts,
			"triforce":    settings.GanonBKTriforcePieces,
		}[comperand.Value]
		if exists {
			return si.settings.KeyShuffle.GanonBKShuffle == choice, true
		}
		panic(slipup.Createf("unknown ganon bosskey placement value %q", comperand.Value))

	case "starting_age", "startingage":
		comperand := st.Symbol(ct.Args[1].(compiler.Load).Id)
		switch comperand.Name {
		case "age":
			return true, true // TODO: this is a runtime compare
		case "adult":
			return si.settings.Spawns.StartingAge == settings.StartAgeAdult, true
		case "child":
			return si.settings.Spawns.StartingAge == settings.StartAgeChild, true
		default:
			panic(slipup.Createf("unknown starting age value %q", comperand.Name))
		}

	case "lacs_condition", "lacscondition":
		comperand := st.String(ct.Args[1].(compiler.Load).Id)
		cond, _ := settings.Decode(si.settings.LacsCondition)
		conds := map[string]settings.Condition{
			"vanilla":    settings.CondDefault,
			"stones":     settings.CondStones,
			"medallions": settings.CondMedallions,
			"dungeons":   settings.CondRewards,
			"tokens":     settings.CondTokens,
			"hearts":     settings.CondHearts,
		}
		compareTo, exists := conds[comperand.Value]
		if exists {
			return cond == compareTo, true
		}
		panic(slipup.Createf("unknown bridge condition value %q", comperand.Value))

	case "zora_fountain", "zorafountain":
		comperand := st.String(ct.Args[1].(compiler.Load).Id)
		fountainState, exists := map[string]settings.OpenZoraFountain{
			"adult":  settings.ZoraFountainOpenAdult,
			"closed": settings.ZoraFountainClosed,
			"open":   settings.ZoraFountainOpenAlways,
		}[comperand.Value]
		if !exists {
			panic(slipup.Createf("unknown open forest value %s", comperand.Value))
		}
		return si.settings.Locations.ZoraFountain == fountainState, true

	default:
		panic(slipup.Createf("%s %q\n%#v", sym.Name, name, ct))
	}
}

func (si *intrinsics) CompareNq(ct compiler.Invocation, sym *compiler.Symbol, st *compiler.SymbolTable) (bool, bool) {
	result, exists := si.CompareEq(ct, sym, st)
	return !result, exists
}

func (si *intrinsics) CompareLt(ct compiler.Invocation, sym *compiler.Symbol, st *compiler.SymbolTable) (bool, bool) {
	name := settingname(ct.Args[0].(compiler.Load), st)
	switch name {
	case "chicken_count", "chickencount":
		comperand := st.Const(ct.Args[1].(compiler.Load).Id)
		return si.settings.Minigames.KakChickens < uint8(comperand.Value), true
	default:
		panic(slipup.Createf("compare_lt %#v\n%#v", sym, ct))
	}
}

func (si *intrinsics) InvertHasShortcuts(ct compiler.Invocation, sym *compiler.Symbol, st *compiler.SymbolTable) (bool, bool) {
	res, exists := si.HasShortcuts(ct, sym, st)
	return !res, exists
}

func (si *intrinsics) HasShortcuts(ct compiler.Invocation, sym *compiler.Symbol, st *compiler.SymbolTable) (bool, bool) {
	arg := ct.Args[0].(compiler.Load)
	dungeon := st.String(arg.Id)

	if dungeon.Value == "King Dodongo Boss Room" {
		// todo: actually locate this and determine if shortcuts are on
		return false, true
	}

	flag, exists := map[string]settings.DungeonShortcuts{
		"Deku Tree":        settings.ShortcutsDeku,
		"Dodongos Cavern":  settings.ShortcutsCavern,
		"Jabu Jabus Belly": settings.ShortcutsJabu,
		"Forest Temple":    settings.ShortcutsForest,
		"Fire Temple":      settings.ShortcutsFire,
		"Water Temple":     settings.ShortcutsWater,
		"Shadow Temple":    settings.ShortcutsShadow,
		"Spirit Temple":    settings.ShortcutsSpirit,
	}[dungeon.Value]

	if !exists {
		panic(slipup.Createf("did not expect dungeon value %q", dungeon.Value))
	}

	return settings.Has(si.settings.Dungeons.Shortcuts, flag), true

}

func (si *intrinsics) IsTrialSkipped(ct compiler.Invocation, sym *compiler.Symbol, st *compiler.SymbolTable) (bool, bool) {
	name := settingname(ct.Args[0].(compiler.Load), st)
	trial, exists := map[string]settings.TrialsEnabled{
		"fire":   settings.TrialsEnabledFire,
		"forest": settings.TrialsEnabledForest,
		"light":  settings.TrialsEnabledLight,
		"shadow": settings.TrialsEnabledShadow,
		"spirit": settings.TrialsEnabledSpirit,
		"water":  settings.TrialsEnabledWater,
	}[name]

	if !exists {
		panic(slipup.Createf("unknown trial %q", name))
	}

	return !settings.HasFlag(si.settings.Dungeons.Trials, trial), true
}

func (si *intrinsics) LoadSetting(ct compiler.Invocation, sym *compiler.Symbol, st *compiler.SymbolTable) (bool, bool) {
	inverting := strings.HasPrefix("invert", sym.Name)
	name := settingname(ct.Args[0].(compiler.Load), st)
	var answer bool
	switch name {
	case "free_bombchu_drops", "freebombchudrops":
		answer = si.settings.FreeBombchuDrops
		break
	case "shuffle_individual_ocarina_notes", "shuffleindividualocarinanotes":
		answer = si.settings.Shuffling.OcarinaNotes
		break
	case "fix_broken_drops", "fixbrokendrops":
		answer = si.settings.FixBrokenDrops
		break
	case "plant_beans", "plantbeans":
		answer = si.settings.Starting.PlantBeans
		break
	case "shuffle_dungeon_entrances", "shuffledungeonentrances":
		answer = si.settings.Entrances.DungeonEntrances != settings.DungeonEntranceShuffleOff
		break
	case "shuffle_empty_pots", "shuffleemptypots":
		answer = settings.HasFlag(si.settings.Shuffling.Pots, settings.ShuffleEmptyPots)
		break
	case "disable_trade_revert", "disabletraderevert":
		answer = si.settings.Trades.DisableRevert
		break
	case "adult_trade_shuffle", "adulttradeshuffle":
		answer = si.settings.Trades.Adult != settings.AdultTradeShuffleDisabled
		break
	case "warp_songs", "warpsongs":
		answer = si.settings.Entrances.WarpSongs
		break
	case "shuffle_silver_rupees", "shufflesilverrupees":
		answer = si.settings.KeyShuffle.SilverRupees != settings.KeysVanilla
		break
	case "shuffle_expensive_merchants", "shuffleexpensivemerchants":
		answer = si.settings.Shuffling.ExpensiveMerchants
		break
	case "entrance_shuffle", "entranceshuffle":
		answer = si.settings.Entrances.ShufflingAny() || si.settings.Spawns.Randomized()
		break
	case "shuffle_interior_entrances", "shuffleinteriorentrances":
		answer = si.settings.Entrances.Interior != settings.InteriorShuffleOff
		break
	case "skip_child_zelda", "skipchildzelda":
		answer = true // TODO we need to check starting items and child shuffle
		break
	case "free_scarecrow", "freescarecrow":
		answer = si.settings.Starting.Scarecrow
		break
	case "complete_mask_quest", "completemaskquest":
		answer = si.settings.Starting.CompleteMaskQuest
		break
	case "skip_reward_from_rauru", "skiprewardfromrauru":
		answer = si.settings.Starting.RauruReward
		break
	case "open_door_of_time", "opendooroftime":
		answer = si.settings.Locations.OpenDoorOfTime
		break
	default:
		panic(slipup.Createf("loading (inverting? %t) setting %q\n%#v", inverting, name, ct))
	}

	if inverting {
		answer = !answer
	}

	return answer, true
}

func (si *intrinsics) IsTrickEnabled(ct compiler.Invocation, sym *compiler.Symbol, st *compiler.SymbolTable) (bool, bool) {
	trick := st.Symbol(ct.Args[0].(compiler.Load).Id)
	enabled := si.bools[trick.Name]
	if enabled {
		return true, true
	}
	return false, true
}

func settingname(loading compiler.Load, st *compiler.SymbolTable) string {
	var name string
	switch loading.Kind {
	case compiler.CT_LOAD_IDENT:
		setSym := st.Symbol(loading.Id)
		name = setSym.Name
		break
	case compiler.CT_LOAD_STR:
		str := st.String(loading.Id)
		name = str.Value
		break
	default:
		panic("tried to load neither ident or string for setting name")
	}
	return name
}
