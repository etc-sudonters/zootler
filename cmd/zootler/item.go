package main

import (
	"fmt"
	"sudonters/zootler/internal/table"
	"sudonters/zootler/pkg/world/components"
)

type FileItem struct {
	Name        string                 `json:"name"`
	Type        string                 `json:"type"`
	Advancement bool                   `json:"advancement"`
	Priority    bool                   `json:"priority"`
	Special     map[string]interface{} `json:"special"`
}

func (item FileItem) TableValues() (table.Values, error) {
	values := []table.Value{
		components.Name(item.Name),
		components.Token{},
	}

	if kind, err := item.kind(); err != nil {
		return nil, err
	} else {
		values = append(values, kind)
	}

	/*
		if item.Advancement {
			values = append(values, components.Advancement{})
		}
		if item.Priority {
			values = append(values, components.Priority{})
		}
	*/

	if special, err := item.special(); err != nil {
		return nil, err
	} else {
		values = append(values, special...)
	}

	return table.Values(values), nil
}

func (item FileItem) kind() (table.Value, error) {
	switch normalize(item.Type) {
	case "bosskey":
		return components.BossKey{}, nil
	case "compass":
		return components.Compass{}, nil
	case "drop":
		return components.Drop{}, nil
	case "dungeonreward":
		return components.DungeonReward{}, nil
	case "event":
		return components.Event{}, nil
	case "ganonbosskey":
		return components.GanonBossKey{}, nil
	case "hideoutsmallkey":
		return components.HideoutSmallKey{}, nil
	case "item":
		return components.Item{}, nil
	case "map":
		return components.Map{}, nil
	case "refill":
		return components.Refill{}, nil
	case "shop":
		return components.Shop{}, nil
	case "silverrupee":
		return components.SilverRupee{}, nil
	case "smallkey":
		return components.SmallKey{}, nil
	case "song":
		return components.Song{}, nil
	case "tcgsmallkey":
		return components.TCGSmallKey{}, nil
	case "token":
		return components.GoldSkulltulaToken{}, nil
	}
	return nil, fmt.Errorf("unknown item type '%s'", item.Type)
}

func (item FileItem) special() (table.Values, error) {
	var values table.Values

	if price, ok := item.Special["price"]; ok {
		if price, ok := price.(float64); ok {
			values = append(values, components.Price(price))
		}
	}

	if _, ok := item.Special["bottle"]; ok {
		values = append(values, components.Bottle{})
	}

	if _, ok := item.Special["ocarina_button"]; ok {
		values = append(values, components.OcarinaButton{})
	}

	if _, ok := item.Special["junk"]; ok {
		values = append(values, components.Junk{})
	}

	if _, ok := item.Special["medallion"]; ok {
		values = append(values, components.Medallion{})
	}

	if _, ok := item.Special["stone"]; ok {
		values = append(values, components.SpiritualStone{})
	}

	if _, ok := item.Special["trade"]; ok {
		values = append(values, components.Trade{})
	}

	return values, nil
}
