package fill

import (
	"context"
	"errors"
	"fmt"
	"sudonters/libzootr/internal/shufflequeue"
	"sudonters/libzootr/internal/skelly/bitset32"
	"sudonters/libzootr/magicbean"
	"sudonters/libzootr/zecs"
)

// not _really_ errors but do communicate why we stopped placing
var ItemPoolEmpty = errors.New("item pool empty")
var LocationPoolEmpty = errors.New("location pool empty")

type Token zecs.Entity
type Location zecs.Entity

type PlaceableTokens struct {
	*shufflequeue.FisherYatesQueue[Token]
	Eligibility map[Token]bitset32.Bitset
}

type PlaceableLocations struct {
	*shufflequeue.FisherYatesQueue[Location]
	Eligibility map[Location]bitset32.Bitset
}

type Ledger map[Location]Token

type Placement struct {
	ledger map[Location]Token
}

// see https://github.com/cjohnson57/RandomizerAlgorithms
type Random struct{}
type Forward struct{}
type Assumed struct{}

type Algorithm interface {
	Fill(context.Context, *magicbean.Exploration, *magicbean.Generation, PlaceableTokens, PlaceableLocations) (Placement, error)
}

func (this Random) Fill(
	ctx context.Context,
	xplr *magicbean.Exploration,
	gen *magicbean.Generation,
	tokens PlaceableTokens,
	locations PlaceableLocations,
) (Placement, error) {
	ocm := &gen.Ocm
	ledger := make(map[Location]Token, min(tokens.Len(), locations.Len()))
	placer := placer{ocm, ledger}
	placement := Placement{ledger}

	assertRequeued := func(token Token, location Location) {
		if err := tokens.Requeue(token); err != nil {
			panic(fmt.Errorf("failed to requeue token: %w", err))
		}
		if err := locations.Requeue(location); err != nil {
			panic(fmt.Errorf("failed to requeue location: %w", err))
		}

	}

placing:
	for {
		if tokens.Len() == 0 {
			return placement, ItemPoolEmpty
		}

		if locations.Len() == 0 {
			return placement, LocationPoolEmpty
		}

		token, _ := tokens.Dequeue()
		location, _ := locations.Dequeue()

		if !bitset32.Intersects(tokens.Eligibility[token], locations.Eligibility[location]) {
			assertRequeued(token, location)
			continue placing
		}

		if err := placer.place(token, location); err != nil {
			assertRequeued(token, location)
			return placement, err

		}
	}
}

type placer struct {
	ocm    *zecs.Ocm
	ledger Ledger
}

func (this placer) place(token Token, location Location) error {
	loc := this.ocm.Proxy(zecs.Entity(location))
	err := loc.Attach(magicbean.HoldsToken(zecs.Entity(token)))
	if err != nil {
		return fmt.Errorf("could not attach token %v at location %v: %w", token, location, err)
	}
	this.ledger[location] = token
	return nil
}
