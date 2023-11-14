package filler

import (
	"context"
	"errors"

	"sudonters/zootler/internal/entity"
	"sudonters/zootler/pkg/logic"
	"sudonters/zootler/pkg/world"
	"sudonters/zootler/pkg/world/components"

	"github.com/etc-sudonters/substrate/bag"
	"github.com/etc-sudonters/substrate/skelly/queue"
	"github.com/etc-sudonters/substrate/stageleft"
)

type ConstGoal bool

func (c ConstGoal) Reachable(context.Context, world.World) (bool, error) {
	return bool(c), nil
}

type Goal interface {
	Reachable(context.Context, world.World) (bool, error)
}

type Filler interface {
	Fill(context.Context, world.World, Goal) error
}

type AssumedFill struct {
	Locations []entity.Selector
	Items     []entity.Selector
}

func (a *AssumedFill) Fill(ctx context.Context, w world.World, g Goal) error {
	var err error = nil
	var locs []entity.View
	var items []entity.View

	var filt []entity.Selector

	filt = make([]entity.Selector, len(a.Locations)+1)
	filt[0] = entity.With[components.Location]{}
	copy(filt[1:], a.Locations)

	locs, err = w.Entities.Pool.Query(filt)
	if err != nil {
		return stageleft.AttachExitCode(err, stageleft.ExitCode(99))
	}

	filt = make([]entity.Selector, len(a.Items)+1)
	filt[0] = entity.With[components.Token]{}
	copy(filt[1:], a.Items)
	items, err = w.Entities.Pool.Query(filt)
	if err != nil {
		return stageleft.AttachExitCode(err, stageleft.ExitCode(99))
	}

	L := queue.From(locs)
	I := queue.From(items)

	var solved bool
	maxTries := L.Len() * I.Len()

	for i := 0; i <= maxTries; i++ {
		if L.Len() == 0 || I.Len() == 0 {
			break
		}

		if err != nil {
			return err
		}

		var loc entity.View
		var item entity.View

		bag.Shuffle(*L)
		bag.Shuffle(*I)
		loc, err = L.Pop()
		if err != nil {
			return err
		}
		item, err = I.Pop()
		if err != nil {
			return err
		}

		loc.Add(logic.Inhabited(item.Model()))
		item.Add(logic.Inhabits(loc.Model()))

		solved, err = g.Reachable(ctx, w)
		if err != nil {
			return err
		}

		if !solved {
			loc.Remove(logic.Inhabited(item.Model()))
			item.Remove(logic.Inhabits(loc.Model()))
			L.Push(loc)
			I.Push(item)
		}
	}

	if !solved {
		err = errors.New("could not solve placement")
	}

	return err
}
