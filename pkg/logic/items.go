package logic

import (
	"os"

	"github.com/etc-sudonters/zootler/pkg/entity"
	"muzzammil.xyz/jsonc"
)

//components
type (
	BossKey            struct{}
	Compass            struct{}
	Drop               struct{}
	DungeonReward      struct{}
	Event              struct{}
	GanonBossKey       struct{}
	HideoutSmallKey    struct{}
	Item               struct{}
	Map                struct{}
	Refill             struct{}
	SmallKey           struct{}
	GoldSkulltulaToken struct{}
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

type item struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Progressive *bool  `json:"progressive"`
}

type PlacementItem struct {
	Name       string
	Type       entity.Component
	Importance Importance
}

type Importance int

func importanceFrom(b *bool) Importance {
	if b == nil {
		return Junk
	} else if !(*b) {
		return Priority
	} else {
		return Advancement
	}
}

const (
	Junk        Importance = 1 << iota
	Priority    Importance = 1 << iota
	Advancement Importance = 1 << iota
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
	}

	return Items, nil
}
