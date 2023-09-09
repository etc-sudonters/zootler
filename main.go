package main

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"

	"github.com/etc-sudonters/zootler/internal/datastructures/stack"
	"github.com/etc-sudonters/zootler/internal/rules"
	"github.com/etc-sudonters/zootler/pkg/entity"
	"github.com/etc-sudonters/zootler/pkg/logic"
	"github.com/etc-sudonters/zootler/pkg/world"
)

var itempool map[string]int = map[string]int{
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

func main() {
	args := os.Args

	var logicFileDir string

	if len(args) > 1 {
		logicFileDir = args[1]
	} else {
		logicFileDir = "."
	}

	binPath, err := os.Executable()
	if err != nil {
		panic(err)
	}

	binDir := filepath.Dir(binPath)
	dataPath := filepath.Join(binDir, "data")

	builder := world.NewBuilder(1)
	// step 1 populate the world with all locations
	placements, err := logic.ReadLocationFile(filepath.Join(dataPath, "locations.json"))

	if err != nil {
		panic(fmt.Errorf("error while reading placements: %w", err))
	}

	if len(placements) == 0 {
		panic("no placements loaded!")
	}

	var locationMap map[string]entity.View = make(map[string]entity.View, len(placements))

	for _, placement := range placements {
		loc, _ := builder.AddEntity(logic.Name(placement.Name))
		locationMap[placement.Name] = loc
		builder.AddNode(loc)
		loc.Add(logic.Location{})

		if len(placement.Tags) > 0 {
			for _, tag := range placement.Tags {
				for _, comp := range logic.ParseComponentsFromLocationTag(tag) {
					loc.Add(comp)
				}
			}
		}
	}

	//step 2 connect everything in the graph and record the raw rules
	entries, err := os.ReadDir(logicFileDir)
	if err != nil {
		panic(err)
	}

	for _, entry := range entries {
		if filepath.Ext(entry.Name()) != "json" {
			continue
		}

		path := filepath.Join(logicFileDir, entry.Name())
		logicLocations, err := rules.ReadLogicFile(path)
		if err != nil {
			panic(fmt.Errorf("error while reading logic rules: %w", err))
		}

		for _, logicLocation := range logicLocations {
			region := locationMap[string(logicLocation.Region)]

			for name, rule := range logicLocation.Locations {
				loc := locationMap[string(name)]
				connection, err := builder.AddEdge(region, loc)
				if err != nil {
					panic(fmt.Errorf(
						"error while building %s -> %s: %w", logicLocation.Region, name, err))
				}

				connection.Add(logic.RawRule(rule))
			}
		}
	}

	// step 3 pour item pool into the mix
	items, err := logic.ReadItemFile(filepath.Join(dataPath, "items.json"))
	if err != nil {
		panic(err)
	}

	placeableItems := make([]entity.View, getItemPoolSize(itempool))

	for _, item := range items {
		count, ok := itempool[item.Name]
		if !ok {
			continue
		}

		N(count, func() {
			entity, _ := builder.AddEntity(logic.Name(item.Name))
			entity.Add(logic.Token{})
			entity.Add(item.Importance)
			entity.Add(item.Type)
			placeableItems = append(placeableItems, entity)
		})
	}

	// step 4 throw the items everywhere who fucking cares for now
	songLocations, err := builder.Pool.Query(
		entity.With[logic.Location]{},
		entity.Without[logic.Inhabited]{},
		entity.With[logic.Song]{},
		entity.Load[logic.Name]{},
	)

	if err != nil {
		panic(err)
	}

	songTokens, err := builder.Pool.Query(
		entity.With[logic.Token]{},
		entity.With[logic.Song]{},
		entity.Without[logic.Inhabited]{},
		entity.Load[logic.Name]{},
	)

	if err != nil {
		panic(err)
	}

	placeItems(
		songLocations,
		songTokens,
	)

	//step 3 get some dumb locations
	world := builder.Build()

	songs, err := world.Entities.Query(entity.With[logic.Song]{}, entity.Load[logic.Inhabited]{}, entity.Load[logic.Name]{})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Songs ended up at:\n")
	for _, check := range songs {
		var name logic.Name
		err := check.Get(&name)
		if err != nil {
			panic(err)
		}

		var inhabited logic.Inhabited
		err = check.Get(&inhabited)
		if err != nil {
			panic(err)
		}

		var songName logic.Name
		world.Entities.Get(entity.Model(inhabited), &songName)

		fmt.Printf("%s @ %s\n", name, songName)
	}
}

func N(n int, do func()) {
	for i := 0; i < n; i++ {
		do()
	}
}

func getItemPoolSize(pool map[string]int) int {
	var n int

	for _, v := range pool {
		n += v
	}

	return n
}

func placeItems(locations []entity.View, itempool []entity.View) {
	fmt.Print("placing items")
	var err error
	shuffle(locations)
	shuffle(itempool)

	L := stack.From(locations)
	I := stack.From(itempool)

	for {
		var loc entity.View
		var item entity.View
		loc, L, err = L.Pop()
		if err != nil {
			fmt.Print("No more locations, exiting")
			break
		}
		item, I, err = I.Pop()
		if err != nil {
			fmt.Print("No more items, exiting")
			break
		}

		loc.Add(logic.Inhabited(item.Model()))
	}
}

func shuffle[T any](elms []T) {
	rand.Shuffle(len(elms), func(i, j int) {
		elms[i], elms[j] = elms[j], elms[i]
	})
}
