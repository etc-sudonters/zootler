package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"

	"sudonters/zootler/internal/entity"
	"sudonters/zootler/internal/stack"
	"sudonters/zootler/pkg/logic"
	"sudonters/zootler/pkg/world"
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
		entity.With[logic.Location]{},
		entity.With[logic.Song]{},
	)

	if err != nil {
		panic(err)
	}

	tokens, err := builder.Pool.Query(
		entity.With[logic.Song]{},
		entity.With[logic.Token]{},
	)

	if err != nil {
		panic(err)
	}

	world := builder.Build()

	placeItems(
		locations,
		tokens,
	)

	placedSongs, err := world.Entities.Query(
		entity.With[logic.Token]{},
		entity.With[logic.Inhabits]{},
		entity.With[logic.Song]{},
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
	var err error

	L := stack.From(locations)
	I := stack.From(itempool)
	shuffle(L)
	shuffle(I)

	for {
		var loc entity.View
		var token entity.View

		loc, L, err = L.Pop()
		if err != nil {
			break
		}
		token, I, err = I.Pop()
		if err != nil {
			break
		}

		placeItem(loc, token)

		var inh logic.Inhabited
		loc.Get(&inh)
	}
}

func placeItem(loc, token entity.View) {
	var locName logic.Name
	var itemName logic.Name
	loc.Get(&locName)
	token.Get(&itemName)

	loc.Add(logic.Inhabited(token.Model()))
	token.Add(logic.Inhabits(loc.Model()))
}

func shuffle[T any, E ~[]T](elms E) {
	rand.Shuffle(len(elms), func(i, j int) {
		elms[i], elms[j] = elms[j], elms[i]
	})
}
