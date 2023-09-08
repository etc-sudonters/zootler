package world

import (
	"context"
	"fmt"
	"sync"

	"github.com/etc-sudonters/zootler/internal/datastructures/graph"
	"github.com/etc-sudonters/zootler/internal/datastructures/set"
	"github.com/etc-sudonters/zootler/internal/rules"
	"github.com/etc-sudonters/zootler/pkg/entity"
	"github.com/etc-sudonters/zootler/pkg/entity/hashpool"
	"github.com/etc-sudonters/zootler/pkg/logic"
)

type edge struct {
	origin graph.Origination
	dest   graph.Destination
}

type Id int

type World struct {
	Id        Id
	Entities  Pool
	Graph     graph.Directed
	edgeCache map[edge]entity.Model
	nodeCache map[graph.Node]entity.Model
}

type Builder struct {
	id        Id
	entities  Pool
	graph     graph.Builder
	edgeCache map[edge]entity.View
	nodeCache map[graph.Node]entity.View
	nameCache map[logic.Name]entity.View
}

// caller is responsible for setting a unique id if necessary
func New(id Id) Builder {
	return Builder{
		id,
		Pool{id, hashpool.New()},
		graph.Builder{graph.New()},
		make(map[edge]entity.View),
		make(map[graph.Node]entity.View),
		make(map[logic.Name]entity.View),
	}
}

// after calling this it is no longer safe to interact with the builder
func (w *Builder) Build() World {
	edgeCache := make(map[edge]entity.Model, len(w.edgeCache))
	nodeCache := make(map[graph.Node]entity.Model, len(w.nodeCache))

	for e, v := range w.edgeCache {
		edgeCache[e] = v.Model()
	}

	for n, v := range w.nodeCache {
		nodeCache[n] = v.Model()

	}

	return World{
		Id:        w.id,
		Entities:  w.entities,
		Graph:     w.graph.G,
		edgeCache: edgeCache,
		nodeCache: nodeCache,
	}
}

// entity returned is the location we accepted
func (w *Builder) Accept(l rules.RawLogicLocation) (entity.View, error) {
	region, err := w.AddEntity(logic.Name(l.Region))
	if err != nil {
		return nil, err
	}

	if l.Hint != nil {
		region.Add(logic.HintGroup(*l.Hint))
	}

	if len(l.Locations) > 0 {
		for locName, rule := range l.Locations {
			e, _ := w.AddEntity(logic.Name(locName))
			e.Add(logic.RawRule(rule))
			w.AddEdge(logic.Name(fmt.Sprintf("%s @ %s", locName, l.Region)), region, e)
		}
	}

	return region, nil
}

func (w *Builder) AddEntity(n logic.Name) (entity.View, error) {
	if ent, ok := w.nameCache[n]; ok {
		return ent, nil
	}

	ent, err := w.entities.Create(n)
	if err != nil {
		return nil, err
	}

	w.nameCache[n] = ent
	return ent, nil
}

func (w *Builder) AddNode(v entity.View) {
	v.Add(logic.Node{})
	w.graph.AddNode(graph.Node(v.Model()))
}

func (w *Builder) AddEdge(name logic.Name, origin, destination entity.View) (entity.View, error) {
	o := graph.Origination(origin.Model())
	d := graph.Destination(destination.Model())

	if ent, ok := w.edgeCache[edge{o, d}]; ok {
		return ent, nil
	}

	if err := w.graph.AddEdge(o, d); err != nil {
		return nil, err
	}

	ent, err := w.entities.Create(name)

	if err != nil {
		return nil, err
	}

	ent.Add(logic.Edge{
		Destination: entity.Model(d),
		Origination: entity.Model(o),
	})

	w.edgeCache[edge{o, d}] = ent

	return ent, nil
}

func (w *World) FindReachableWorld(ctx context.Context) (set.Hash[graph.Node], error) {
	reachable := set.New[graph.Node]()

	bfs := graph.BreadthFirst[graph.Destination]{
		Selector: &RulesAwareSelector[graph.Destination]{
			w, graph.Successors,
		},
		Visitor: &graph.VisitSet{S: reachable},
	}

	spawns, err := w.Entities.Query(entity.With[logic.Spawn]{})
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	for _, e := range spawns {
		e := e
		wg.Add(1)
		go func() {
			defer wg.Done()
			bfs.Walk(ctx, w.Graph, graph.Node(e.Model()))
		}()
	}

	wg.Wait()
	return reachable, nil
}
