package logic

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/etc-sudonters/zootler/pkg/entity"
	"muzzammil.xyz/jsonc"
)

//components
type (
	BossKey            struct{}
	Bottle             struct{}
	Compass            struct{}
	Count              float64
	Drop               struct{}
	DungeonReward      struct{}
	Event              struct{}
	GanonBossKey       struct{}
	GoldSkulltulaToken struct{}
	HideoutSmallKey    struct{}
	Item               struct{}
	Junk               struct{}
	Map                struct{}
	Medallion          struct{}
	Price              float64
	Refill             struct{}
	ShopObject         float64
	SmallKey           struct{}
	SpiritualStone     struct{}
	Trade              struct{}
)

func (c BossKey) String() string            { return "BossKey" }
func (c Compass) String() string            { return "Compass" }
func (c Drop) String() string               { return "Drop" }
func (c DungeonReward) String() string      { return "DungeonReward" }
func (c Event) String() string              { return "Event" }
func (c GanonBossKey) String() string       { return "GanonBossKey" }
func (c HideoutSmallKey) String() string    { return "HideoutSmallKey" }
func (c Item) String() string               { return "Item" }
func (c Map) String() string                { return "Map" }
func (c Refill) String() string             { return "Refill" }
func (c Shop) String() string               { return "Shop" }
func (c SmallKey) String() string           { return "SmallKey" }
func (c Song) String() string               { return "Song" }
func (c GoldSkulltulaToken) String() string { return "Token" }

func componentFromItemType(t string) entity.Component {
	switch t {
	case "BossKey":
		return BossKey{}
	case "Compass":
		return Compass{}
	case "Drop":
		return Drop{}
	case "DungeonReward":
		return DungeonReward{}
	case "Event":
		return Event{}
	case "GanonBossKey":
		return GanonBossKey{}
	case "HideoutSmallKey":
		return HideoutSmallKey{}
	case "Item":
		return Item{}
	case "Map":
		return Map{}
	case "Refill":
		return Refill{}
	case "Shop":
		return Shop{}
	case "SmallKey":
		return SmallKey{}
	case "Song":
		return Song{}
	case "Token":
		return GoldSkulltulaToken{}
	}

	return nil
}

type itemDetails struct {
	Junk        *float64
	Progressive *float64
	Bottle      bool
	Shop        *float64
	Price       *float64
	Stone       bool
	Medallion   bool
	Trade       bool
	Alias       []itemAlias
}

type itemAlias struct {
	For   string
	Count int
}

func (i *itemDetails) UnmarshalJSON(data []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	for k, v := range raw {
		switch strings.ToLower(k) {
		case "junk":
			junk := v.(float64)
			i.Junk = &junk
			break
		case "progressive":
			prog := v.(float64)
			i.Progressive = &prog
			break
		case "bottle":
			i.Bottle = true
			break
		case "shop_object":
			shop := v.(float64)
			i.Shop = &shop
			break
		case "price":
			price := v.(float64)
			i.Price = &price
			break
		case "stone":
			i.Stone = true
			break
		case "medallion":
			i.Medallion = true
			break
		case "trade":
			i.Trade = true
			break
		case "alias":
			aliases := v.([]interface{})

			for x := 0; x < len(aliases)-1; x += 2 {
				aliasFor := aliases[x].(string)
				amount := aliases[x+1].(float64)

				alias := itemAlias{
					For:   aliasFor,
					Count: int(amount),
				}

				i.Alias = append(i.Alias, alias)
			}
		}
	}

	return nil
}

type item struct {
	Name        string       `json:"name"`
	Type        string       `json:"type"`
	Progressive *bool        `json:"progressive"`
	Details     *itemDetails `json:"special"`
}

type PlacementItem struct {
	Name       string
	Type       entity.Component
	Importance Importance
	Components []entity.Component
}

type Importance int

func importanceFrom(b *bool) Importance {
	if b == nil {
		return ImportanceJunk
	} else if !(*b) {
		return ImportancePriority
	} else {
		return ImportanceAdvancement
	}
}

const (
	ImportanceJunk        Importance = 1 << iota
	ImportancePriority    Importance = 1 << iota
	ImportanceAdvancement Importance = 1 << iota
)

func ReadItemFile(fp string) ([]PlacementItem, error) {

	contents, err := os.ReadFile(fp)
	if err != nil {
		return nil, err
	}

	var items []item
	if err := jsonc.Unmarshal(contents, &items); err != nil {
		return nil, err
	}

	var Items []PlacementItem = make([]PlacementItem, len(items))
	for i, item := range items {
		Items[i] = PlacementItem{
			Name:       item.Name,
			Type:       componentFromItemType(item.Type),
			Importance: importanceFrom(item.Progressive),
		}

		var comps []entity.Component

		if details := item.Details; details != nil {
			if junk := details.Junk; junk != nil {
				comps = append(comps, Junk{}, Count(*junk))
			}

			if progressive := details.Progressive; progressive != nil {
				comps = append(comps, Count(*progressive))
			}

			if details.Bottle {
				comps = append(comps, Bottle{})
			}

			if s := details.Shop; s != nil {
				comps = append(comps, ShopObject(*s))
			}

			if p := details.Price; p != nil {
				comps = append(comps, Price(*p))
			}

			if details.Stone {
				comps = append(comps, SpiritualStone{})
			}

			if details.Medallion {
				comps = append(comps, Medallion{})
			}

			if details.Trade {
				comps = append(comps, Trade{})
			}

			Items[i].Components = comps
		}
	}

	return Items, nil
}
