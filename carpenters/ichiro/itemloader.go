package ichiro

import (
	"sudonters/zootler/internal/components"
	"sudonters/zootler/internal/table"

	"github.com/etc-sudonters/substrate/slipup"
)

type TokenComponents struct {
	Name        string                 `json:"name"`
	Type        string                 `json:"type"`
	Advancement bool                   `json:"advancement"`
	Priority    bool                   `json:"priority"`
	Special     map[string]interface{} `json:"special"`
}

func (i TokenComponents) EntityName() components.Name {
	return components.Name(i.Name)
}

func (i TokenComponents) AsComponents() table.Values {
	vt := table.Values{i.kind(), components.CollectableGameToken{}}
	if i.Advancement {
		vt = append(vt, components.Advancement{})
	}
	if i.Priority {
		vt = append(vt, components.Priority{})
	}
	return i.special(vt)
}

func (i TokenComponents) kind() table.Value {
	switch i.Type {
	case "BossKey":
		return components.BossKey{}
	case "Compass":
		return components.Compass{}
	case "Drop":
		return components.Drop{}
	case "DungeonReward":
		return components.DungeonReward{}
	case "Event":
		return components.Event{}
	case "GanonBossKey":
		return components.GanonBossKey{}
	case "HideoutSmallKey":
		return components.HideoutSmallKey{}
	case "Item":
		return components.Item{}
	case "Map":
		return components.Map{}
	case "Refill":
		return components.Refill{}
	case "Shop":
		return components.Shop{}
	case "SilverRupee":
		return components.SilverRupee{}
	case "SmallKey":
		return components.SmallKey{}
	case "Song":
		return components.Song{}
	case "TCGSmallKey":
		return components.TCGSmallKey{}
	case "Token":
		return components.GoldSkulltulaToken{}
	default:
		panic(slipup.Createf("unknown item type '%s'", i.Type))
	}
}

func (i TokenComponents) special(vt table.Values) table.Values {
	if price, ok := i.Special["price"]; ok {
		if price, ok := price.(float64); ok {
			vt = append(vt, components.Price(price))
		}
	}

	if _, ok := i.Special["bottle"]; ok {
		vt = append(vt, components.Bottle{})
	}

	if _, ok := i.Special["ocarina_button"]; ok {
		vt = append(vt, components.OcarinaButton{})
	}

	if _, ok := i.Special["junk"]; ok {
		vt = append(vt, components.Junk{})
	}

	if _, ok := i.Special["medallion"]; ok {
		vt = append(vt, components.Medallion{})
	}

	if _, ok := i.Special["stone"]; ok {
		vt = append(vt, components.SpiritualStone{})
	}

	if _, ok := i.Special["trade"]; ok {
		vt = append(vt, components.Trade{})
	}

	return vt
}
