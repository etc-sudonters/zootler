package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"

	"github.com/etc-sudonters/zootler/internal/stack"
	"github.com/etc-sudonters/zootler/pkg/entity"
	"github.com/etc-sudonters/zootler/pkg/logic"
	"github.com/etc-sudonters/zootler/pkg/world"
)

func main() {
	var logicFileDir string
	var dataDir string

	flag.StringVar(&logicFileDir, "l", "", "Directory where logic files are located")
	flag.StringVar(&dataDir, "d", "", "Directory where data files are stored")
	flag.Parse()

	if logicFileDir == "" || dataDir == "" {
		fmt.Fprint(os.Stderr, "-l and -d are both required!\n")
		os.Exit(2)
	}

	builder, err := world.BuildFromDirectory(1, dataDir, logicFileDir)
	if err != nil {
		panic(err)
	}

	locations, err := builder.Pool.Query(
		entity.Load[logic.Name]{},
		entity.With[logic.Location]{},
		entity.Without[logic.MasterQuest]{},
		entity.Without[logic.Inhabited]{},
		entity.Without[logic.Cow]{},
		entity.Without[logic.Beehive]{},
		entity.Without[logic.Cow]{},
		entity.Without[logic.Shop]{},
		entity.Without[logic.Crate]{},
		entity.Without[logic.Pot]{},
		entity.Without[logic.Flying]{},
		entity.Without[logic.RupeeTower]{},
		entity.Without[logic.GoldSkulltula]{},
		entity.Without[logic.SmallCrate]{},
		entity.Without[logic.Refill]{},
		entity.Without[logic.Drop]{},
		entity.Without[logic.RupeeTower]{},
	)

	if err != nil {
		panic("oh no!")
	}

	tokens, err := builder.Pool.Query(
		entity.Load[logic.Name]{},
		entity.With[logic.Token]{},
		entity.Without[logic.Inhabited]{},
	)

	if err != nil {
		panic("bruno")
	}

	placeItems(
		locations,
		tokens,
	)

	world := builder.Build()

	placedSongs, err := world.Entities.Query(
		entity.With[logic.Token]{},
		entity.Load[logic.Inhabited]{},
		entity.Load[logic.Name]{},
	)

	if err != nil {
		panic(err)
	}

	fmt.Printf("Songs ended up at:\n")
	for _, placedSong := range placedSongs {
		var name logic.Name
		err := placedSong.Get(&name)
		if err != nil {
			panic(err)
		}

		var inhabited logic.Inhabited
		err = placedSong.Get(&inhabited)
		if err != nil {
			panic(err)
		}

		var locationName logic.Name
		world.Entities.Get(entity.Model(inhabited), &locationName)

		fmt.Printf("%s @ %s\n", name, locationName)
	}
}

func placeItems(locations []entity.View, itempool []entity.View) {
	fmt.Print("placing items")
	var err error

	fmt.Printf("Placing %d items in %d locations\n", len(itempool), len(locations))

	L := stack.From(locations)
	I := stack.From(itempool)
	shuffle(L)
	shuffle(I)

	for {
		var loc entity.View
		var item entity.View
		var locName logic.Name
		var itemName logic.Name

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

		loc.Get(&locName)
		item.Get(&itemName)

		fmt.Printf("Placing %s at '%s'\n", itemName, locName)

		loc.Add(logic.Inhabited(item.Model()))
	}
}

func shuffle[T any, E ~[]T](elms E) {
	rand.Shuffle(len(elms), func(i, j int) {
		elms[i], elms[j] = elms[j], elms[i]
	})
}
