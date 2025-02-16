package magicbean

import (
	"fmt"
	"github.com/etc-sudonters/substrate/skelly/bitset32"
	"sudonters/libzootr/components"
	"sudonters/libzootr/zecs"
)

func CollectTokensFrom(
	ocm *zecs.Ocm,
	from bitset32.Bitset,
	into Inventory,
) error {

	q := ocm.Query()
	q.Build(
		zecs.Load[components.HoldsToken],
		zecs.With[components.PlacementLocationMarker],
		zecs.WithOut[components.Collected],
		zecs.FromSubset(&from),
	)

	rows, err := q.Execute()
	if err != nil {
		return err
	}

	for entity, tup := range rows.All {
		token := tup.Values[0].(components.HoldsToken)
		into.CollectOne(zecs.Entity(token))

		proxy := ocm.Proxy(entity)
		err := proxy.Attach(components.Collected{})
		if err != nil {
			return fmt.Errorf("while marking %v as collected: %w", entity, err)
		}
	}

	return nil
}

func CollectStartingItems(generation *Generation) error {
	tokens := &generation.Tokens
	these := generation.Settings
	inventory := generation.Inventory

	for name, qty := range these.Logic.Spawns.Items {
		token := tokens.MustGet(components.Name(name))
		inventory.Collect(token.Entity(), qty)
	}

	ocm := &generation.Ocm
	skipped := zecs.Bitset32Matching(ocm, zecs.With[components.Skipped])
	return CollectTokensFrom(ocm, skipped, inventory)
}
