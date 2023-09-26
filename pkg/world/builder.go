package world

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/etc-sudonters/zootler/internal/graph"
	"github.com/etc-sudonters/zootler/internal/rules"
	"github.com/etc-sudonters/zootler/pkg/entity"
	"github.com/etc-sudonters/zootler/pkg/entity/hashpool"
	"github.com/etc-sudonters/zootler/pkg/logic"
)

type Builder struct {
	id        Id
	Pool      Pool
	graph     graph.Builder
	edgeCache map[edge]entity.View
	nodeCache map[graph.Node]entity.View
}

// caller is responsible for setting a unique id if necessary
func NewBuilder(id Id) *Builder {
	return &Builder{
		id,
		Pool{id, hashpool.New()},
		graph.Builder{graph.New()},
		make(map[edge]entity.View),
		make(map[graph.Node]entity.View),
	}
}

// after calling this it is no longer safe to interact with the builder
func (w *Builder) Build() World {
	edgeCache := make(map[edge]entity.Model, len(w.edgeCache))
	nodeCache := make(map[graph.Node]entity.Model, len(w.nodeCache))

	for e, v := range w.edgeCache {
		edgeCache[e] = v.Model()
	}

	for n, v := range w.nodeCache {
		nodeCache[n] = v.Model()

	}

	return World{
		Id:        w.id,
		Entities:  w.Pool,
		Graph:     w.graph.G,
		edgeCache: edgeCache,
		nodeCache: nodeCache,
	}
}

func (w *Builder) AddEntity(n logic.Name) (entity.View, error) {
	ent, err := w.Pool.Create(n)
	if err != nil {
		return nil, err
	}

	return ent, nil
}

func (w *Builder) AddNode(v entity.View) {
	w.graph.AddNode(graph.Node(v.Model()))
}

func (w *Builder) AddEdge(origin, destination entity.View) (entity.View, error) {
	var oName logic.Name
	var dName logic.Name
	origin.Get(&oName)
	destination.Get(&dName)

	name := logic.Name(fmt.Sprintf("%s -> %s", oName, dName))

	o := graph.Origination(origin.Model())
	d := graph.Destination(destination.Model())

	if ent, ok := w.edgeCache[edge{o, d}]; ok {
		return ent, nil
	}

	if err := w.graph.AddEdge(o, d); err != nil {
		return nil, err
	}

	ent, err := w.Pool.Create(name)

	if err != nil {
		return nil, err
	}

	ent.Add(logic.Edge{
		Destination: entity.Model(d),
		Origination: entity.Model(o),
	})

	w.edgeCache[edge{o, d}] = ent

	return ent, nil
}

func BuildFromDirectory(world Id, dataDir, logicDir string) (*Builder, error) {
	b := NewBuilder(world)

	if err := loadItems(b, dataDir); err != nil {
		return nil, fmt.Errorf("while loading items: %w", err)
	}

	quickLocations, err := loadLocations(b, dataDir)
	if err != nil {
		return nil, fmt.Errorf("while loading locations: %w", err)
	}

	if err := loadLogic(b, logicDir, quickLocations); err != nil {
		return nil, fmt.Errorf("while loading logic: %w", err)
	}

	return b, nil
}

func loadLocations(b *Builder, d string) (map[string]entity.View, error) {
	locations, err := logic.ReadLocationFile(filepath.Join(d, "locations.json"))
	if err != nil {
		return nil, err
	}

	quickLocation := make(map[string]entity.View, len(locations))

	for _, l := range locations {
		loc, _ := b.AddEntity(logic.Name(l.Name))
		quickLocation[l.Name] = loc
		b.AddNode(loc)
		loc.Add(logic.Location{})

		for _, tag := range l.Tags {
			for _, comp := range logic.ParseComponentsFromLocationTag(tag) {
				loc.Add(comp)
			}
		}

		for _, typ := range logic.ParseComponentsFromLocationType(l.Type) {
			loc.Add(typ)
		}
	}

	return quickLocation, nil
}

func loadLogic(b *Builder, logicDir string, quickLocations map[string]entity.View) error {
	entries, err := os.ReadDir(logicDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if filepath.Ext(entry.Name()) != "json" {
			continue
		}

		path := filepath.Join(logicDir, entry.Name())
		logicLocations, err := rules.ReadLogicFile(path)
		if err != nil {
			return err
		}

		for _, logicLocation := range logicLocations {
			region := quickLocations[string(logicLocation.Region)]

			for name, rule := range logicLocation.Locations {
				loc := quickLocations[string(name)]
				connection, err := b.AddEdge(region, loc)
				if err != nil {
					return fmt.Errorf("error while building %s -> %s: %w", logicLocation.Region, name, err)
				}

				connection.Add(logic.RawRule(rule))
			}
		}
	}
	return nil
}

func loadItems(b *Builder, dataDir string) error {
	items, err := logic.ReadItemFile(filepath.Join(dataDir, "items.json"))
	if err != nil {
		panic(err)
	}

	placeableItems := make([]entity.View, len(defaultItemPool))

	for _, item := range items {
		count, ok := defaultItemPool[item.Name]
		if !ok {
			continue
		}

		for i := 0; i < count; i++ {
			entity, _ := b.AddEntity(logic.Name(item.Name))
			entity.Add(logic.Token{})
			entity.Add(item.Importance)
			entity.Add(item.Type)
			placeableItems = append(placeableItems, entity)
		}
	}

	return nil
}

var defaultItemPool map[string]int = map[string]int{
	"Arrows (10)":                          8,
	"Arrows (30)":                          6,
	"Arrows (5)":                           3,
	"Biggoron Sword":                       1,
	"Bolero of Fire":                       1,
	"Bomb Bag":                             3,
	"Bombchus (10)":                        3,
	"Bombchus (20)":                        1,
	"Bombchus (5)":                         1,
	"Bombs (10)":                           2,
	"Bombs (20)":                           2,
	"Bombs (5)":                            8,
	"Boomerang":                            1,
	"Bottle with Blue Potion":              1,
	"Bottle with Red Potion":               2,
	"Bow":                                  3,
	"Claim Check":                          1,
	"Deku Nut Capacity":                    2,
	"Deku Nuts (10)":                       1,
	"Deku Nuts (5)":                        4,
	"Deku Seeds (30)":                      4,
	"Deku Shield":                          4,
	"Deku Stick (1)":                       3,
	"Deku Stick Capacity":                  2,
	"Dins Fire":                            1,
	"Double Defense":                       1,
	"Eponas Song":                          1,
	"Farores Wind":                         1,
	"Fire Arrows":                          1,
	"Goron Tunic":                          1,
	"Heart Container":                      8,
	"Hover Boots":                          1,
	"Hylian Shield":                        2,
	"Ice Arrows":                           1,
	"Iron Boots":                           1,
	"Kokiri Sword":                         1,
	"Lens of Truth":                        1,
	"Light Arrows":                         1,
	"Magic Meter":                          2,
	"Megaton Hammer":                       1,
	"Minuet of Forest":                     1,
	"Mirror Shield":                        1,
	"Nayrus Love":                          1,
	"Nocturne of Shadow":                   1,
	"Piece of Heart (Treasure Chest Game)": 1,
	"Piece of Heart":                       35,
	"Prelude of Light":                     1,
	"Progressive Hookshot":                 2,
	"Progressive Scale":                    2,
	"Progressive Strength Upgrade":         3,
	"Progressive Wallet":                   2,
	"Recovery Heart":                       11,
	"Requiem of Spirit":                    1,
	"Rupee (1)":                            1,
	"Rupees (20)":                          6,
	"Rupees (200)":                         6,
	"Rupees (5)":                           23,
	"Rupees (50)":                          7,
	"Rutos Letter":                         1,
	"Sarias Song":                          1,
	"Serenade of Water":                    1,
	"Slingshot":                            3,
	"Song of Storms":                       1,
	"Song of Time":                         1,
	"Stone of Agony":                       1,
	"Suns Song":                            1,
	"Zeldas Lullaby":                       1,
	"Zora Tunic":                           1,
}
