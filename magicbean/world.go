package magicbean

import (
	"fmt"
	"sudonters/zootler/internal/query"
	"sudonters/zootler/internal/settings"
	"sudonters/zootler/mido/objects"

	"github.com/etc-sudonters/substrate/skelly/graph"
)

type CollectionId uint16
type age uint8

func (this age) String() string {
	switch this {
	case EXPLORE_AS_ADULT:
		return "EXPLORE_AS_ADULT"
	case EXPLORE_AS_CHILD:
		return "EXPLORE_AS_CHILD"
	default:
		panic(fmt.Errorf("unknown exploration age %d", this))
	}
}

func explorationAgeFromSettings(this *settings.Zootr) age {
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
	EXPLORE_AS_ADULT age = 0
	EXPLORE_AS_CHILD     = 1
)

type World struct {
	dependencies
	sharedstate
	nodes []nodes
}

type dependencies struct {
	DataStore   query.Engine
	Physical    *graph.Directed
	Translation *TranslationLayer
}

type WorldBuilder struct {
	dependencies
	Objects   *objects.TableBuilder
	Settings  *settings.Zootr
	RootNodes []graph.Node
}

type sharedstate struct {
	collected map[CollectionId]uint8
	builtins  map[string]objects.Index
	objects   objects.Table
	startAge  age
}

func NewWorld(from WorldBuilder) *World {
	world := new(World)
	world.dependencies.DataStore = from.DataStore
	world.dependencies.Physical = from.Physical
	world.dependencies.Translation = from.Translation

	world.nodes = make([]nodes, 2)
	world.nodes[EXPLORE_AS_ADULT] = newtracker(from.RootNodes)
	world.nodes[EXPLORE_AS_CHILD] = newtracker(from.RootNodes)

	world.collected = make(map[CollectionId]uint8, len(world.Translation.byname))
	world.startAge = explorationAgeFromSettings(from.Settings)
	world.objects = objects.NewTable(objects.BuildTableFrom(from.Objects))

	return world
}
