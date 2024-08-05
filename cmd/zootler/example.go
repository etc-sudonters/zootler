package main

import (
	"context"
	"fmt"
	"sudonters/zootler/internal/query"
	"sudonters/zootler/pkg/world/components"

	"github.com/etc-sudonters/substrate/dontio"
	"github.com/etc-sudonters/substrate/mirrors"
)

func example(ctx context.Context, storage query.Engine) error {
	stdio, stdErr := dontio.StdFromContext(ctx)
	if stdErr != nil {
		return stdErr
	}

	func() {
		q := storage.CreateQuery()
		q.Exists(mirrors.TypeOf[components.Location]())
		allLocs, err := storage.Retrieve(q)
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(stdio.Out, "Count of all locations: %d\n", allLocs.Len())
	}()

	func() {
		q := storage.CreateQuery()
		q.Exists(mirrors.TypeOf[components.Location]())
		q.Exists(mirrors.TypeOf[components.Song]())
		songLocs, err := storage.Retrieve(q)
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(stdio.Out, "Count of Song locations: %d\n", songLocs.Len())
	}()

	func() {
		q := storage.CreateQuery()
		q.Exists(mirrors.TypeOf[components.Location]())
		q.NotExists(mirrors.TypeOf[components.Song]())
		notSongLocs, err := storage.Retrieve(q)
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(stdio.Out, "Count of not Song locations: %d\n", notSongLocs.Len())
	}()

	func() {
		q := storage.CreateQuery()
		q.Exists(mirrors.TypeOf[components.CollectableGameToken]())
		allToks, err := storage.Retrieve(q)
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(stdio.Out, "Count of all collectable tokens: %d\n", allToks.Len())
	}()

	func() {
		q := storage.CreateQuery()
		q.Exists(mirrors.TypeOf[components.CollectableGameToken]())
		q.Exists(mirrors.TypeOf[components.Song]())
		songToks, err := storage.Retrieve(q)
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(stdio.Out, "Count of Song tokens: %d\n", songToks.Len())
	}()

	func() {
		q := storage.CreateQuery()
		q.Exists(mirrors.TypeOf[components.CollectableGameToken]())
		q.NotExists(mirrors.TypeOf[components.Song]())
		notSongToks, err := storage.Retrieve(q)
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(stdio.Out, "Count of not Song tokens: %d\n", notSongToks.Len())
	}()

	func() {
		lookupName := "Spirit Medallion"
		l := storage.CreateLookup()
		l.Load(mirrors.TypeOf[components.Medallion]())
		l.Load(mirrors.TypeOf[components.DungeonReward]())
		l.Load(mirrors.TypeOf[components.Advancement]())
		l.Load(mirrors.TypeOf[components.Pot]())

		l.Lookup(components.Name(lookupName))
		med, err := storage.Lookup(l)
		if err != nil {
			panic(err)
		}
		foundMed := med.Len() == 1
		fmt.Fprintf(stdio.Out, "Found %s? %t\n", lookupName, foundMed)

		if foundMed {
			med.MoveNext()
			medallion := med.Current()

			for i := range medallion.Cols {
				fmt.Fprintf(stdio.Out, "Loaded column '%d' for '%s': %v\n", medallion.Cols[i], lookupName, medallion.Values[i])
			}
		}
	}()

	return nil

}
