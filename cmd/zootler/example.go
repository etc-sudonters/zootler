package main

import (
	"context"
	"errors"
	"sudonters/zootler/internal/app"
	"sudonters/zootler/internal/components"
	"sudonters/zootler/internal/query"
	"sudonters/zootler/internal/rules/runtime"
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
	c := new(runtime.ChunkBuilder)
	c.LoadConst(runtime.StackValueOrPanic(1))
	c.LoadConst(runtime.StackValueOrPanic(0))
	c.Equal()
	jmpTrue := c.JumpIfTrue()
	c.LoadConst(runtime.StackValueOrPanic(2))
	c.LoadConst(runtime.StackValueOrPanic(1))
	c.Equal()
	jumpFalse := c.JumpIfFalse()
	jmp1Target, _ := c.LoadConst(runtime.StackValueOrPanic(true))
	convergeJump := c.UnconditionalJump()
	jmp2Target, _ := c.LoadConst(runtime.StackValueOrPanic(false))
	converge, _ := c.LoadConst(runtime.StackValueOrPanic(3.14))
	c.LoadIdentifier("frank")
	c.DumpStack()
	c.SetReturn()
	c.Return()

	c.PatchJump(jmpTrue, jmp1Target)
	c.PatchJump(jumpFalse, jmp2Target)
	c.PatchJump(convergeJump, converge)

	WriteLineOut(z.Ctx(), c.Disassemble("test"))
	WriteLineOut(z.Ctx(), "%s\n", c.Ops)

	env := runtime.NewEnv()
	env.Set("frank", runtime.StackValueOrPanic(3.14))

	vm := runtime.CreateVM(env)
	vm.Debug(true)
	result, runErr := vm.Run(z.Ctx(), &c.Chunk)
	WriteLineOut(z.Ctx(), "vm dump:\n%#v", result)
	if runErr == nil {
		WriteLineOut(z.Ctx(), "result:\t%#v", result.Unwrap())
	}
	return runErr
}
