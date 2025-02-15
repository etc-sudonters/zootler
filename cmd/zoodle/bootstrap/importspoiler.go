package bootstrap

import (
	"context"
	"fmt"
	"io"
	"strings"
	"sudonters/libzootr/components"
	"sudonters/libzootr/internal/json"
	"sudonters/libzootr/magicbean/tracking"
	"sudonters/libzootr/settings"

	"github.com/etc-sudonters/substrate/dontio"
)

func malformed(err error) error {
	return fmt.Errorf("malformed spoilers document: %w", err)
}

func LoadSpoilerData(
	ctx context.Context,
	r io.Reader,
	these *settings.Model,
	nodes *tracking.Nodes,
	tokens *tracking.Tokens) error {

	if std, err := dontio.StdFromContext(ctx); err == nil {
		std.WriteLineOut("Loading spoiler data...")
	}

	reader := json.NewParser(json.NewScanner(r))
	obj, err := reader.ReadObject()
	if err != nil {
		return malformed(err)
	}

	for obj.More() {
		property, err := obj.ReadPropertyName()
		if err != nil {
			return malformed(err)
		}

		switch property {
		case "settings":
			obj, err := obj.ReadObject()
			if err != nil {
				return malformed(err)
			}
			if err := CopySettings(these, obj); err != nil {
				return malformed(err)
			}
		case "locations":
			obj, err := obj.ReadObject()
			if err != nil {
				return malformed(err)
			}
			if err := CopyPlacements(tokens, nodes, obj); err != nil {
				return malformed(err)
			}
		case ":skipped_locations":
			obj, err := obj.ReadObject()
			if err != nil {
				return malformed(err)
			}
			if err := CopySkippedLocations(nodes, obj); err != nil {
				return malformed(err)
			}
		default:
			if err := obj.DiscardValue(); err != nil {
				return malformed(err)
			}
		}
	}

	if err := obj.ReadEnd(); err != nil {
		return malformed(err)
	}

	return nil
}

func CopySettings(these *settings.Model, obj *json.ObjectParser) error {
	for obj.More() {
		setting, err := obj.ReadPropertyName()
		if err != nil {
			return err
		}
		if err := readSetting(setting, these, obj); err != nil {
			return fmt.Errorf("while reading setting %q: %w", setting, err)
		}
	}
	return obj.ReadEnd()
}

func readSetting(propertyName string, these *settings.Model, obj *json.ObjectParser) error {
	var err error

	switch propertyName {
	case "logic_rules":
		value, err := obj.ReadString()
		if err != nil {
			return fmt.Errorf("logic_rules: valid settings: none, glitches, glitchless: %w")
		}
		switch value {
		case "none":
			these.Logic.Set = settings.LogicNone
		case "glitched":
			these.Logic.Set = settings.LogicGlitched
		case "glitchless":
			these.Logic.Set = settings.LogicGlitchless
		default:
			return fmt.Errorf("unknown logic_rules value: %q", value)
		}
	case "reachable_locations":
		value, err := obj.ReadString()
		if err != nil {
			return fmt.Errorf("reachable_locations: valid settings: ... : %w", value)
		}
		switch value {
		case "all":
			these.Logic.Locations.Reachability = settings.ReachableAll
		case "goals":
			these.Logic.Locations.Reachability = settings.ReachableGoals
		case "beatable":
			these.Logic.Locations.Reachability = settings.ReachableNecessary
		default:
			return fmt.Errorf("reachable_locations: unknown setting: %q", value)
		}
	case "triforce_hunt":
		value, err := obj.ReadBool()
		these.Logic.WinConditions.TriforceHunt = value
		return err
	case "lacs_condition":
		val, readErr := obj.ReadString()
		if readErr != nil {
			return readErr
		}
		cond, parseErr := settings.ParseCondition(val)
		if parseErr != nil {
			return parseErr
		}
		_, qty := these.Logic.WinConditions.Lacs.Decode()
		these.Logic.WinConditions.Lacs = settings.EncodeConditionedAmount(cond, qty)
	case "bridge":
		val, readErr := obj.ReadString()
		if readErr != nil {
			return readErr
		}
		cond, parseErr := settings.ParseCondition(val)
		if parseErr != nil {
			return parseErr
		}
		_, qty := these.Logic.WinConditions.Bridge.Decode()
		these.Logic.WinConditions.Bridge = settings.EncodeConditionedAmount(cond, qty)
	case "trials":
		// this only tells how many trials _MIGHT_ be enabled the spoiler
		// already has these pinned down in its trials top level key
		count, err := obj.ReadInt()
		if err != nil {
			return fmt.Errorf("expected number for trials: %w", err)
		}
		if count == 0 {
			these.Logic.WinConditions.Trials = 0
		}
	case "shuffle_ganon_bosskey":
		val, readErr := obj.ReadString()
		if readErr != nil {
			return fmt.Errorf("expected string for shuffle_ganon_bosskey: %w", readErr)
		}
		cond, err := settings.ParseGanonBossKeyShuffle(val)
		if err != nil {
			return err
		}
		these.Logic.Dungeon.GanonBossKeyShuffle = cond
	case "ganon_bosskey_medallions",
		"ganon_bosskey_stones",
		"ganon_bosskey_rewards",
		"ganon_bosskey_hearts",
		"ganon_bosskey_tokens":
		if these.Logic.WinConditions.GanonBossKey != 0 {
			return fmt.Errorf("ganon_bosskey condition is already set")
		}

		qty, qtyErr := obj.ReadInt()
		if qtyErr != nil {
			fmt.Errorf("expected number for %s: %w", propertyName, qtyErr)
		}

		condStr, _ := strings.CutPrefix(propertyName, "ganon_bosskey_")
		cond, condErr := settings.ConditionFrom(condStr, uint32(qty))
		if condErr != nil {
			return fmt.Errorf("could not parse %s: %w", propertyName, condErr)
		}
		these.Logic.WinConditions.GanonBossKey = cond
	case "adult_trade_start":
		var adultTradeItems settings.AdultTradeItems
		for i, propertyName := range json.ReadStringArray(obj, &err) {
			switch propertyName {
			case "Pocket Egg":
				adultTradeItems |= settings.AdultTradePocketEgg
			case "Pocket Cucco":
				adultTradeItems |= settings.AdultTradePocketCucco
			case "Odd Mushroom":
				adultTradeItems |= settings.AdultTradeOddMushroom
			case "Odd Potion":
				adultTradeItems |= settings.AdultTradeOddPotion
			case "Poachers Saw":
				adultTradeItems |= settings.AdultTradePoachersSaw
			case "Broken Sword":
				adultTradeItems |= settings.AdultTradeBrokenSword
			case "Prescription":
				adultTradeItems |= settings.AdultTradePrescription
			case "Eyeball Frog":
				adultTradeItems |= settings.AdultTradeEyeballFrog
			case "Eyedrops":
				adultTradeItems |= settings.AdultTradeEyedrops
			case "Claim Check":
				adultTradeItems |= settings.AdultTradeClaimCheck
			default:
				return fmt.Errorf("unknown adult trade shuffle item %q at index %d", propertyName, i)
			}
		}
		these.Logic.Trade.AdultItems = adultTradeItems
		return err

	default:
		return obj.DiscardValue()
	}

	return nil
}

func CopyPlacements(tokens *tracking.Tokens, nodes *tracking.Nodes, obj *json.ObjectParser) error {
	var err error

	for node, placed := range json.ReadObjectProperties(obj, readPlacedItem, &err) {
		node := nodes.Placement(components.Name(node))
		token := tokens.MustGet(components.Name(placed.name))
		if placed.price > -1 {
			token.Attach(components.Price(placed.price))
		}
		node.Holding(token)
	}
	return err
}

func CopySkippedLocations(nodes *tracking.Nodes, obj *json.ObjectParser) error {
	var err error
	for skipped := range obj.Keys(&err) {
		node := nodes.Placement(components.Name(skipped))
		node.Attach(components.Collected{})
	}

	return err
}

type placed struct {
	name  string
	price int
}

func readPlacedItem(obj *json.ObjectParser) (placed, error) {
	var placed placed
	placed.price = -1
	switch kind := obj.Current().Kind; kind {
	case json.STRING:
		var readErr error
		placed.name, readErr = obj.ReadString()
		if readErr != nil {
			return placed, readErr
		}
		break
	case json.OBJ_OPEN:
		inner, readErr := obj.ReadObject()
		if readErr != nil {
			return placed, readErr
		}
		for inner.More() {
			prop, err := obj.ReadPropertyName()
			if err != nil {
				return placed, err
			}
			switch prop {
			case "item":
				placed.name, readErr = obj.ReadString()
				if readErr != nil {
					return placed, err
				}
				break
			case "price":
				placed.price, readErr = obj.ReadInt()
				if readErr != nil {
					return placed, readErr
				}
				break
			default:
				return placed, fmt.Errorf("unexpected property: %q", prop)
			}
		}

	}
	return placed, nil
}
