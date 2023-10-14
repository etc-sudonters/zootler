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
		entity.With[logic.Name]{},
		entity.With[logic.Location]{},
		entity.Without[logic.Beehive]{},
		entity.Without[logic.Cow]{},
		entity.Without[logic.Crate]{},
		entity.Without[logic.Drop]{},
		entity.Without[logic.Flying]{},
		entity.Without[logic.GoldSkulltula]{},
		entity.Without[logic.Pot]{},
		entity.Without[logic.RupeeTower]{},
		entity.Without[logic.Shop]{},
		entity.Without[logic.SmallCrate]{},
		entity.Without[logic.Hint]{},
		entity.Without[logic.HintStone]{},
	)

	if err != nil {
		panic(err)
	}

	tokens, err := builder.Pool.Query(
		entity.With[logic.Name]{},
		entity.With[logic.Token]{},
	)

	if err != nil {
		panic("bruno")
	}

	world := builder.Build()

	placeItems(
		locations,
		tokens,
	)

	placedSongs, err := world.Entities.Query(
		entity.With[logic.Token]{},
		entity.With[logic.Name]{},
		entity.With[logic.Inhabits]{},
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

		var inhabits logic.Inhabits
		err = placedSong.Get(&inhabits)
		if err != nil {
			panic(err)
		}

		var locationName logic.Name
		world.Entities.Get(entity.Model(inhabits), &locationName)

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
		var token entity.View

		loc, L, err = L.Pop()
		if err != nil {
			fmt.Print("No more locations, exiting")
			break
		}
		token, I, err = I.Pop()
		if err != nil {
			fmt.Print("No more items, exiting")
			break
		}

		placeItem(loc, token)

		var inh logic.Inhabited
		loc.Get(&inh)
		fmt.Printf("%+v", inh)
	}
}

func placeItem(loc, token entity.View) {
	var locName logic.Name
	var itemName logic.Name
	loc.Get(&locName)
	token.Get(&itemName)

	fmt.Printf("Placing %s at '%s'\n", itemName, locName)
	loc.Add(logic.Inhabited(token.Model()))
	token.Add(logic.Inhabits(loc.Model()))
}

func shuffle[T any, E ~[]T](elms E) {
	rand.Shuffle(len(elms), func(i, j int) {
		elms[i], elms[j] = elms[j], elms[i]
	})
}
