package main

import (
	"context"
	"errors"
	"sudonters/zootler/internal/app"
	"sudonters/zootler/internal/components"
	"sudonters/zootler/internal/query"
	"sudonters/zootler/internal/rules/bytecode"
	"sudonters/zootler/internal/rules/vm"
	"sudonters/zootler/internal/table"
)

func example(ctx context.Context, storage query.Engine) error {
	func() {
		q := storage.CreateQuery()
		q.Exists(T[components.Location]())
		allLocs, err := storage.Retrieve(q)
		if err != nil {
			panic(err)
		}
		WriteLineOut(ctx, "Count of all locations: %d", allLocs.Len())
	}()

	func() {
		q := storage.CreateQuery()
		q.Exists(T[components.Location]())
		q.Exists(T[components.Song]())
		songLocs, err := storage.Retrieve(q)
		if err != nil {
			panic(err)
		}
		WriteLineOut(ctx, "Count of Song locations: %d", songLocs.Len())
	}()

	func() {
		q := storage.CreateQuery()
		q.Exists(T[components.Location]())
		q.NotExists(T[components.Song]())
		notSongLocs, err := storage.Retrieve(q)
		if err != nil {
			panic(err)
		}
		WriteLineOut(ctx, "Count of not Song locations: %d", notSongLocs.Len())
	}()

	func() {
		q := storage.CreateQuery()
		q.Exists(T[components.CollectableGameToken]())
		allToks, err := storage.Retrieve(q)
		if err != nil {
			panic(err)
		}
		WriteLineOut(ctx, "Count of all collectable tokens: %d", allToks.Len())
	}()

	func() {
		q := storage.CreateQuery()
		q.Exists(T[components.CollectableGameToken]())
		q.Exists(T[components.Song]())
		songToks, err := storage.Retrieve(q)
		if err != nil {
			panic(err)
		}
		WriteLineOut(ctx, "Count of Song tokens: %d", songToks.Len())
	}()

	func() {
		q := storage.CreateQuery()
		q.Exists(T[components.CollectableGameToken]())
		q.NotExists(T[components.Song]())
		notSongToks, err := storage.Retrieve(q)
		if err != nil {
			panic(err)
		}
		WriteLineOut(ctx, "Count of not Song tokens: %d", notSongToks.Len())
	}()

	func() {
		q := storage.CreateQuery()
		q.Load(T[components.Name]())
		q.Load(T[components.Song]())
		q.Load(T[components.CollectableGameToken]())
		things, err := storage.Retrieve(q)
		if err != nil {
			panic(err)
		}

		WriteLineOut(ctx, "size: %d", things.Len())
		for things.MoveNext() {
			row := things.Current()
			name := row.Values[0].(components.Name)
			WriteLineOut(ctx, "'%s' (%d)", name, row.Id)
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
		WriteLineOut(ctx, "Found %s? %t", lookupName, foundMed)

		if foundMed {
			med.MoveNext()
			medallion := med.Current()

			for i := range medallion.Cols {
				WriteLineOut(ctx, "Loaded column '%s' for '%s': %v", medallion.Cols[i].T.Name(), lookupName, medallion.Values[i])
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

		WriteLineOut(ctx, "found %d rows", entries.Len())
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

func manualProgram(z *app.Zootlr) error {
	c := new(bytecode.ChunkBuilder)
	c.PushConst(bytecode.ValueFromFloat(1))
	c.PushConst(bytecode.ValueFromFloat(0))
	c.Equal()
	c.PushConst(bytecode.ValueFromFloat(2))
	c.PushConst(bytecode.ValueFromFloat(1))
	c.NotEqual()
	c.Or()
	c.Dup()
	_, jmpFalse := c.JumpFalse()
	c.PushConst(bytecode.ValueFromBool(true))
	c.Rotate()
	_, jmpTrue := c.UnconditionalJump()
	c.PushConst(bytecode.ValueFromBool(false))
	jmpTrueTarget := c.Return()
	c.PatchJump(jmpFalse, bytecode.PC(jmpTrue))
	c.PatchJump(jmpTrue, jmpTrueTarget)

	runtime, runErr := vm.Evaluate(z.Ctx(), &c.Chunk)

	WriteLineOut(z.Ctx(), c.Disassemble("test"))
	WriteLineOut(z.Ctx(), "%s\n", c.Ops)
	WriteLineOut(z.Ctx(), "vm dump:\n%#v", runtime)
	if runErr == nil {
		WriteLineOut(z.Ctx(), "result:\t%#v", runtime.Result().Unwrap())
	}
	return runErr
}
