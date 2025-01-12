package magicbeanvm

import (
	"sudonters/zootler/internal/components"
	"sudonters/zootler/internal/query"
	"sudonters/zootler/internal/table"

	"github.com/etc-sudonters/substrate/skelly/bitset"
)

type Exploration struct {
	active  bool
	age     ExplorationAge
	reached bitset.Bitset64
	workset bitset.Bitset64
	world   *World
}

func NewExploration(age ExplorationAge, world *World) Exploration {
	var explore Exploration
	explore.world = world
	explore.age = age

	roots := explore.world.RootNodes()
	if len(roots) == 0 {
		panic("no root nodes declared in physical world")
	}

	for _, root := range roots {
		bitset.Set(&explore.workset, root)
	}

	return explore
}

func (this Exploration) QuantityFor(id CollectionId) uint8 {
	return this.world.collected[id]
}

func (this Exploration) IsExploringAs(age ExplorationAge) bool {
	return this.age == age
}

func (this Exploration) IsStartingAge() bool {
	return this.age == explorationAgeFromSettings(this.world.Settings)
}

func (this Exploration) HasBottle() bool {
	store := this.world.DataStore
	q := store.CreateQuery()
	q.Load(query.MustAsColumnId[CollectionId](store))
	q.Exists(query.MustAsColumnId[components.IsBottle](store))

	bottles, err := store.Retrieve(q)
	if err != nil {
		panic(err)
	}

	for _, bottle := range bottles.All {
		id := bottle.Values[0].(CollectionId)
		if this.world.collected[id] > 0 {
			return true
		}
	}

	return false
}

func (this Exploration) HasDungeonRewards(count int) bool {
	total := this.sumCollectionFor(query.MustAsColumnId[components.DungeonReward](this.world.DataStore))
	return total >= count
}

func (this Exploration) HasHearts(count int) bool {
	store := this.world.DataStore
	q := store.CreateQuery()
	q.Load(query.MustAsColumnId[CollectionId](store))
	q.Load(query.MustAsColumnId[components.PieceOfHeart](store))
	hearts, err := store.Retrieve(q)
	if err != nil {
		panic(err)
	}

	var total int
	for _, heart := range hearts.All {
		id := heart.Values[0].(CollectionId)
		mul := heart.Values[1].(components.PieceOfHeart)

		total += int(mul) * int(this.world.collected[id])
	}

	return total >= count
}

func (this Exploration) HasMedallions(count int) bool {
	total := this.sumCollectionFor(query.MustAsColumnId[components.Medallion](this.world.DataStore))
	return total >= count
}

func (this Exploration) HasStones(count int) bool {
	total := this.sumCollectionFor(query.MustAsColumnId[components.SpiritualStone](this.world.DataStore))
	return total >= count
}

func (this Exploration) sumCollectionFor(kind table.ColumnId) int {
	store := this.world.DataStore
	q := store.CreateQuery()
	q.Load(query.MustAsColumnId[CollectionId](store))
	q.Exists(kind)
	rows, err := store.Retrieve(q)
	if err != nil {
		panic(err)
	}

	var sum int
	for _, token := range rows.All {
		id := token.Values[0].(CollectionId)
		sum += int(this.world.collected[id])
	}
	return sum
}
