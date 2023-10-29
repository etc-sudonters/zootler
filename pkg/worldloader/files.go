package worldloader

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
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

func (f FileSystemLoader) LoadInto(ctx context.Context, b *world.Builder) error {
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

	items, err := ReadItemFile(ctx, path.Join(f.DataDirectory, "items.json"))
	if err != nil {
		return fmt.Errorf("while loading items: %w", err)
	}

	for _, item := range items {
		ent, err := b.Pool.Create(world.Name(item.Name))
		if err != nil {
			return fmt.Errorf("failed to create %q: %w", item.Name, err)
		}

		ent.Add([]entity.Component{
			logic.Token{},
			item.Type,
			item.Importance,
		})

		ent.Add(item.Components)
	}

	iter, err := ReadLogicDirectory(ctx, f.DataDirectory)
	if err != nil {
		return fmt.Errorf("failed to read dir %q: %w", f.DataDirectory, err)
	}

	for iter.MoveNext() {
		loc := iter.Current()
		ent, err := b.Pool.Create(world.Name(loc.Region))
		if err != nil {
			// we should already exist, this is really anomalous
			panic(fmt.Errorf("entity %q should have existed already: %w", loc.Region, err))
		}

		ent.Add(loc.Components())

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

	var items []logic.PlacementItem
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, err
	}

	return items, nil
}

func ReadLogicDirectory(ctx context.Context, directory string) (reitertools.Iterator[logic.RawLogicLocation], error) {
	entries, err := os.ReadDir(directory)
	if err != nil {
		return nil, fmt.Errorf("while reading directory %q: %w", directory, err)
	}

	iter := reitertools.Filter(
		reitertools.SliceIter(entries),
		func(e os.DirEntry) bool {
			return filepath.Ext(e.Name()) == "json"
		})

	return reitertools.Flatten(iter, func(e os.DirEntry) reitertools.Iterator[logic.RawLogicLocation] {
		_, data, err := jsonc.ReadFromFile(filepath.Join(directory, e.Name()))
		if err != nil {
			panic(err)
		}
		var locs []logic.RawLogicLocation
		if err := json.Unmarshal(data, &locs); err != nil {
			panic(err)
		}

		return reitertools.SliceIter(locs)
	}), nil
}
