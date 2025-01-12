package magicbeanvm

import (
	"fmt"
	"sudonters/zootler/internal/query"
	"sudonters/zootler/internal/settings"
	"sudonters/zootler/midologic/objects"

	"github.com/etc-sudonters/substrate/skelly/graph"
)

type CollectionId uint16
type ExplorationAge uint8

func explorationAgeFromSettings(this *settings.Zootr) ExplorationAge {
	switch this.Spawns.StartingAge {
	case settings.StartAgeAdult:
		return EXPLORE_AS_ADULT
	case settings.StartAgeChild:
		return EXPLORE_AS_CHILD
	default:
		panic(fmt.Errorf("unknown starting age %v", this.Spawns.StartingAge))
	}
}

const (
	EXPLORE_AS_ADULT ExplorationAge = 0
	EXPLORE_AS_CHILD                = 1
)

type World struct {
	Settings    *settings.Zootr
	Physical    graph.Directed
	DataStore   query.Engine
	Translation TranslationLayer

	explorations []Exploration
	collected    []uint8
}

func (this *World) Init() {
	this.explorations = make([]Exploration, 2)
	this.explorations[EXPLORE_AS_ADULT] = NewExploration(EXPLORE_AS_ADULT, this)
	this.explorations[EXPLORE_AS_CHILD] = NewExploration(EXPLORE_AS_CHILD, this)
	this.explorations[explorationAgeFromSettings(this.Settings)].active = true

	var largestCollectId CollectionId
	for _, id := range this.Translation.byname {
		if id > largestCollectId {
			largestCollectId = id
		}
	}

	this.collected = make([]uint8, int(largestCollectId))
}

func (this *World) QuantityByName(name string) uint8 {
	id := this.Translation.FromName(name)
	return this.collected[id]
}

func (this *World) QuantityByPointer(ptr objects.Ptr) uint8 {
	id := this.Translation.FromPtr(ptr)
	return this.collected[id]
}

func (w *World) RootNodes() []graph.Node {
	panic("not implemented")
}
