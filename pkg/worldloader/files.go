package worldloader

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sudonters/zootler/internal/entity"
	"sudonters/zootler/internal/reitertools"
	"sudonters/zootler/pkg/logic"
	"sudonters/zootler/pkg/world"

	"muzzammil.xyz/jsonc"
)

type FileSystemLoader struct {
	LogicDirectory string
	DataDirectory  string
}

func loadLocations(ctx context.Context, f FileSystemLoader, b *world.Builder) error {
	locations, err := ReadLocationFile(ctx, path.Join(f.DataDirectory, "locations.json"))
	if err != nil {
		return fmt.Errorf("while loading locations: %w", err)
	}

	for _, loc := range locations {
		ent, err := b.Pool.Create(world.Name(loc.Name))
		if err != nil {
			return fmt.Errorf("failed to create %q: %w", loc.Name, err)
		}

		ent.Add(logic.Location{})
		ent.Add(logic.GetAllLocationComponents(loc))
		b.AddNode(ent)
	}
	return nil
}

func loadItems(ctx context.Context, f FileSystemLoader, b *world.Builder) error {
	items, err := ReadItemFile(ctx, path.Join(f.DataDirectory, "items.json"))
	if err != nil {
		return fmt.Errorf("while loading items: %w", err)
	}

	for _, item := range items {
		ent, err := b.Pool.Create(world.Name(item.Name))
		if err != nil {
			return fmt.Errorf("failed to create %q: %w", item.Name, err)
		}

		ent.Add(logic.Token{})
		ent.Add(item.Importance)
		ent.Add(item.Components)
		ent.Add(item.Type)
	}
	return nil
}

func loadConnections(ctx context.Context, f FileSystemLoader, b *world.Builder) error {
	iter, err := ReadLogicDirectory(ctx, f.LogicDirectory)
	if err != nil {
		return fmt.Errorf("failed to read dir %q: %w", f.DataDirectory, err)
	}

	for iter.MoveNext() {
		loc := iter.Current()
		ent, err := b.Pool.Create(world.Name(loc.Region))
		if err != nil {
			return fmt.Errorf("while adding region %q: %w", loc.Region, err)
		}

		ent.Add(loc.Components())
		b.AddNode(ent)

		for event, rule := range loc.Events {
			evt, err := b.Pool.Create(world.Name(event))
			if err != nil {
				return fmt.Errorf("while creating event %q: %w", event, err)
			}
			evt.Add(logic.Token{})
			edge, err := b.AddEdge(ent, evt)
			if err != nil {
				return fmt.Errorf("while linking %q to event %q: %w", loc.Region, event, err)
			}
			edge.Add(rule)
		}

		for exit, rule := range loc.Exits {
			ext, err := b.Pool.Create(world.Name(exit))
			edge, err := b.AddEdge(ent, ext)
			if err != nil {
				return fmt.Errorf("while adding edge from %q to %q: %w", loc.Region, exit, err)
			}
			edge.Add(rule)
		}

		for check, rule := range loc.Locations {
			chk, err := b.Pool.Create(world.Name(check))
			if err != nil {
				return fmt.Errorf("while creating check %q at %q: %w", loc.Region, check, err)
			}
			edge, err := b.AddEdge(ent, chk)
			if err != nil {
				return fmt.Errorf("while linking check %q to %q: %w", check, loc.Region, err)
			}
			edge.Add(rule)
		}
	}

	return nil
}

func (f FileSystemLoader) LoadInto(ctx context.Context, b *world.Builder) error {
	if err := loadItems(ctx, f, b); err != nil {
		return err
	}

	if err := loadLocations(ctx, f, b); err != nil {
		return err
	}

	if err := loadConnections(ctx, f, b); err != nil {
		return err
	}

	return nil
}

func ReadLocationFile(ctx context.Context, filename string) ([]logic.PlacementLocation, error) {
	_, data, err := jsonc.ReadFromFile(filename)
	if err != nil {
		return nil, err
	}

	var locations []logic.PlacementLocation
	if err := json.Unmarshal(data, &locations); err != nil {
		return nil, err
	}

	return locations, nil
}

func ReadItemFile(ctx context.Context, filename string) ([]logic.PlacementItem, error) {
	_, data, err := jsonc.ReadFromFile(filename)
	if err != nil {
		return nil, err
	}

	var items []item
	if err := jsonc.Unmarshal(data, &items); err != nil {
		return nil, err
	}

	Items := make([]logic.PlacementItem, len(items))
	for i, item := range items {
		Items[i] = logic.PlacementItem{
			Name:       item.Name,
			Type:       componentFromItemType(item.Type),
			Importance: importanceFrom(item.Progressive),
		}

		var comps []entity.Component

		if details := item.Details; details != nil {
			if junk := details.Junk; junk != nil {
				comps = append(comps, logic.Junk{}, logic.Count(*junk))
			}

			if progressive := details.Progressive; progressive != nil {
				comps = append(comps, logic.Count(*progressive))
			}

			if details.Bottle {
				comps = append(comps, logic.Bottle{})
			}

			if s := details.Shop; s != nil {
				comps = append(comps, logic.ShopObject(*s))
			}

			if p := details.Price; p != nil {
				comps = append(comps, logic.Price(*p))
			}

			if details.Stone {
				comps = append(comps, logic.SpiritualStone{})
			}

			if details.Medallion {
				comps = append(comps, logic.Medallion{})
			}

			if details.Trade {
				comps = append(comps, logic.Trade{})
			}

			Items[i].Components = comps
		}
	}

	return Items, nil
}

func ReadLogicDirectory(ctx context.Context, directory string) (reitertools.Iterator[logic.RawLogicLocation], error) {
	entries, err := os.ReadDir(directory)
	if err != nil {
		return nil, fmt.Errorf("while reading directory %q: %w", directory, err)
	}

	iter := reitertools.Filter(
		reitertools.SliceIter(entries),
		func(entry os.DirEntry, _ int) bool {
			return strings.HasSuffix(entry.Name(), "json") && !strings.Contains(entry.Name(), "Helpers.json")
		})

	return reitertools.Flatten(iter, func(e os.DirEntry) reitertools.Iterator[logic.RawLogicLocation] {
		contents, err := os.ReadFile(filepath.Join(directory, e.Name()))
		if err != nil {
			panic(err)
		}

		var locs []logic.RawLogicLocation
		if err := jsonc.Unmarshal(contents, &locs); err != nil {
			panic(err)
		}

		return reitertools.SliceIter(locs)
	}), nil
}

func componentFromItemType(t string) entity.Component {
	switch t {
	case "BossKey":
		return logic.BossKey{}
	case "Compass":
		return logic.Compass{}
	case "Drop":
		return logic.Drop{}
	case "DungeonReward":
		return logic.DungeonReward{}
	case "Event":
		return logic.Event{}
	case "GanonBossKey":
		return logic.GanonBossKey{}
	case "HideoutSmallKey":
		return logic.HideoutSmallKey{}
	case "Item":
		return logic.Item{}
	case "Map":
		return logic.Map{}
	case "Refill":
		return logic.Refill{}
	case "Shop":
		return logic.Shop{}
	case "SmallKey":
		return logic.SmallKey{}
	case "Song":
		return logic.Song{}
	case "Token":
		return logic.GoldSkulltulaToken{}
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

func importanceFrom(b *bool) logic.Importance {
	if b == nil {
		return logic.ImportanceJunk
	} else if !(*b) {
		return logic.ImportancePriority
	} else {
		return logic.ImportanceAdvancement
	}
}
