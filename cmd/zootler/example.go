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
	stdio, _ := dontio.StdFromContext(ctx)
	q := storage.CreateQuery()
	q.Exists(mirrors.TypeOf[components.Location]())
	allLocs, err := storage.Retrieve(q)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(stdio.Out, "Count of all locations: %d\n", allLocs.Len())

	q = storage.CreateQuery()
	q.Exists(mirrors.TypeOf[components.Location]())
	q.Exists(mirrors.TypeOf[components.Song]())
	songLocs, err := storage.Retrieve(q)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(stdio.Out, "Count of Song locations: %d\n", songLocs.Len())

	q = storage.CreateQuery()

	q = storage.CreateQuery()
	q.Exists(mirrors.TypeOf[components.Location]())
	q.NotExists(mirrors.TypeOf[components.Song]())
	notSongLocs, err := storage.Retrieve(q)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(stdio.Out, "Count of not Song locations: %d\n", notSongLocs.Len())

	q = storage.CreateQuery()
	q.Exists(mirrors.TypeOf[components.CollectableGameToken]())
	allToks, err := storage.Retrieve(q)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(stdio.Out, "Count of all collectable tokens: %d\n", allToks.Len())

	q = storage.CreateQuery()
	q.Exists(mirrors.TypeOf[components.CollectableGameToken]())
	q.Exists(mirrors.TypeOf[components.Song]())
	songToks, err := storage.Retrieve(q)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(stdio.Out, "Count of Song tokens: %d\n", songToks.Len())

	q = storage.CreateQuery()

	q = storage.CreateQuery()
	q.Exists(mirrors.TypeOf[components.CollectableGameToken]())
	q.NotExists(mirrors.TypeOf[components.Song]())
	notSongToks, err := storage.Retrieve(q)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(stdio.Out, "Count of not Song tokens: %d\n", notSongToks.Len())
	return nil
}
