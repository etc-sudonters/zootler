package main

import (
	"fmt"
	"sudonters/zootler/internal/query"
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

func (item FileItem) GetName() components.Name {
	return components.Name(item.Name)
}

func (item FileItem) AddComponents(rid table.RowId, storage query.Engine) error {
	if err := storage.SetValues(rid, table.Values{components.CollectableGameToken{}}); err != nil {
		return err
	}

	if kindErr := item.kind(rid, storage); kindErr != nil {
		return kindErr
	}

	if item.Advancement {
		if err := storage.SetValues(rid, table.Values{components.Advancement{}}); err != nil {
			return err
		}
	}
	if item.Priority {
		if err := storage.SetValues(rid, table.Values{components.Priority{}}); err != nil {
			return err
		}
	}

	if specialErr := item.special(rid, storage); specialErr != nil {
		return specialErr
	}

	return nil
}

func (item FileItem) kind(rid table.RowId, storage query.Engine) error {
	switch item.Type {
	case "BossKey":
		return storage.SetValues(rid, table.Values{components.BossKey{}})
	case "Compass":
		return storage.SetValues(rid, table.Values{components.Compass{}})
	case "Drop":
		return storage.SetValues(rid, table.Values{components.Drop{}})
	case "DungeonReward":
		return storage.SetValues(rid, table.Values{components.DungeonReward{}})
	case "Event":
		return storage.SetValues(rid, table.Values{components.Event{}})
	case "GanonBossKey":
		return storage.SetValues(rid, table.Values{components.GanonBossKey{}})
	case "HideoutSmallKey":
		return storage.SetValues(rid, table.Values{components.HideoutSmallKey{}})
	case "Item":
		return storage.SetValues(rid, table.Values{components.Item{}})
	case "Map":
		return storage.SetValues(rid, table.Values{components.Map{}})
	case "Refill":
		return storage.SetValues(rid, table.Values{components.Refill{}})
	case "Shop":
		return storage.SetValues(rid, table.Values{components.Shop{}})
	case "SilverRupee":
		return storage.SetValues(rid, table.Values{components.SilverRupee{}})
	case "SmallKey":
		return storage.SetValues(rid, table.Values{components.SmallKey{}})
	case "Song":
		return storage.SetValues(rid, table.Values{components.Song{}})
	case "TCGSmallKey":
		return storage.SetValues(rid, table.Values{components.TCGSmallKey{}})
	case "Token":
		return storage.SetValues(rid, table.Values{components.GoldSkulltulaToken{}})
	default:
		return fmt.Errorf("unknown item type '%s'", item.Type)
	}

}

func (item FileItem) special(rid table.RowId, storage query.Engine) error {
	if price, ok := item.Special["price"]; ok {
		if price, ok := price.(float64); ok {
			if err := storage.SetValues(rid, table.Values{components.Price(price)}); err != nil {
				return err
			}
		}
	}

	if _, ok := item.Special["bottle"]; ok {
		if err := storage.SetValues(rid, table.Values{components.Bottle{}}); err != nil {
			return err
		}
	}

	if _, ok := item.Special["ocarina_button"]; ok {
		if err := storage.SetValues(rid, table.Values{components.OcarinaButton{}}); err != nil {
			return err
		}
	}

	if _, ok := item.Special["junk"]; ok {
		if err := storage.SetValues(rid, table.Values{components.Junk{}}); err != nil {
			return err
		}
	}

	if _, ok := item.Special["medallion"]; ok {
		if err := storage.SetValues(rid, table.Values{components.Medallion{}}); err != nil {
			return err
		}
	}

	if _, ok := item.Special["stone"]; ok {
		if err := storage.SetValues(rid, table.Values{components.SpiritualStone{}}); err != nil {
			return err
		}
	}

	if _, ok := item.Special["trade"]; ok {
		if err := storage.SetValues(rid, table.Values{components.Trade{}}); err != nil {
			return err
		}
	}

	return nil
}
