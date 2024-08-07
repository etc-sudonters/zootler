package main

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sudonters/zootler/internal/query"
	"sudonters/zootler/internal/table"
	"sudonters/zootler/pkg/world/components"

	"github.com/etc-sudonters/substrate/dontio"
	"github.com/etc-sudonters/substrate/mirrors"
)

func T[E any]() reflect.Type {
	return mirrors.TypeOf[E]()
}

func example(ctx context.Context, storage query.Engine) error {
	stdio, stdErr := dontio.StdFromContext(ctx)
	std := std{stdio}
	if stdErr != nil {
		return stdErr
	}

	func() {
		q := storage.CreateQuery()
		q.Exists(T[components.Location]())
		allLocs, err := storage.Retrieve(q)
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(stdio.Out, "Count of all locations: %d\n", allLocs.Len())
	}()

	func() {
		q := storage.CreateQuery()
		q.Exists(T[components.Location]())
		q.Exists(T[components.Song]())
		songLocs, err := storage.Retrieve(q)
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(stdio.Out, "Count of Song locations: %d\n", songLocs.Len())
	}()

	func() {
		q := storage.CreateQuery()
		q.Exists(T[components.Location]())
		q.NotExists(T[components.Song]())
		notSongLocs, err := storage.Retrieve(q)
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(stdio.Out, "Count of not Song locations: %d\n", notSongLocs.Len())
	}()

	func() {
		q := storage.CreateQuery()
		q.Exists(T[components.CollectableGameToken]())
		allToks, err := storage.Retrieve(q)
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(stdio.Out, "Count of all collectable tokens: %d\n", allToks.Len())
	}()

	func() {
		q := storage.CreateQuery()
		q.Exists(T[components.CollectableGameToken]())
		q.Exists(T[components.Song]())
		songToks, err := storage.Retrieve(q)
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(stdio.Out, "Count of Song tokens: %d\n", songToks.Len())
	}()

	func() {
		q := storage.CreateQuery()
		q.Exists(T[components.CollectableGameToken]())
		q.NotExists(T[components.Song]())
		notSongToks, err := storage.Retrieve(q)
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(stdio.Out, "Count of not Song tokens: %d\n", notSongToks.Len())
	}()

	func() {
		q := storage.CreateQuery()
		q.NotExists(T[components.Location]())
		q.Exists(T[components.Song]())
		q.Load(T[components.Name]())
		songNames, err := storage.Retrieve(q)
		if err != nil {
			panic(err)
		}

		for songNames.MoveNext() {
			row := songNames.Current()
			name := row.Values[0].(components.Name)
			std.WriteLineOut("now playing '%s' (%d)", name, row.Id)
		}
	}()

	func() {
		lookupName := "Spirit Medallion"
		l := storage.CreateLookup()
		l.Load(T[components.Medallion]())
		l.Load(T[components.DungeonReward]())
		l.Load(T[components.Advancement]())
		l.Load(T[components.Pot]())

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
				fmt.Fprintf(stdio.Out, "Loaded column '%s' for '%s': %v\n", medallion.Cols[i].T.Name(), lookupName, medallion.Values[i])
			}
		}
	}()

	if err := func() error {
		l := storage.CreateQuery()
		l.Load(T[components.Name]())

		entries, err := storage.Retrieve(l)
		if err != nil {
			panic(err)
		} else if entries.Len() == 0 {
			return errors.New("did not find any rows!")
		}

		std.WriteLineOut("found %d rows", entries.Len())
		for entries.MoveNext() {
			row := entries.Current()
			h := new(hintable)
			h.init(row)
		}

		return nil
	}(); err != nil {
		return err
	}

	return nil

}

type hintable struct {
	Rule components.RawLogic
	Name components.Name
}

func (h *hintable) init(r *table.RowTuple) error {
	m := r.ColumnMap()
	name, nameErr := table.Extract[components.Name](m)
	if nameErr != nil {
		return nameErr
	}
	h.Name = *name
	return nil
}
