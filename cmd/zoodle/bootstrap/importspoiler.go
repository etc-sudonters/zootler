package bootstrap

import (
	"context"
	"fmt"
	"io"
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

func readSetting(name string, these *settings.Model, obj *json.ObjectParser) error {
	var err error

	tradeShuffle := these.Logic.Trade.AdultItems
	switch name {
	case "adult_trade_start":
		for i, name := range json.ReadStringArray(obj, &err) {
			switch name {
			case "Pocket Egg":
				tradeShuffle |= settings.AdultTradeStartPocketEgg
			case "Pocket Cucco":
				tradeShuffle |= settings.AdultTradeStartPocketCucco
			case "Odd Mushroom":
				tradeShuffle |= settings.AdultTradeStartOddMushroom
			case "Odd Potion":
				tradeShuffle |= settings.AdultTradeStartOddPotion
			case "Poachers Saw":
				tradeShuffle |= settings.AdultTradeStartPoachersSaw
			case "Broken Sword":
				tradeShuffle |= settings.AdultTradeStartBrokenSword
			case "Prescription":
				tradeShuffle |= settings.AdultTradeStartPrescription
			case "Eyeball Frog":
				tradeShuffle |= settings.AdultTradeStartEyeballFrog
			case "Eyedrops":
				tradeShuffle |= settings.AdultTradeStartEyedrops
			case "Claim Check":
				tradeShuffle |= settings.AdultTradeStartClaimCheck
			default:
				return fmt.Errorf("unknown adult trade shuffle item %q at index %d", name, i)
			}
		}
		these.Logic.Trade.AdultItems = tradeShuffle
		return err

	default:
		return obj.DiscardValue()
	}
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
