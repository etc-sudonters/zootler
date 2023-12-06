package world

import (
	"errors"
	"fmt"

	"sudonters/zootler/internal/entity"
	"sudonters/zootler/internal/entity/bitpool"
	"sudonters/zootler/internal/entity/table"
	"sudonters/zootler/pkg/world/components"

	"github.com/etc-sudonters/substrate/mirrors"
	"github.com/etc-sudonters/substrate/reiterate"
	"github.com/etc-sudonters/substrate/skelly/graph"
)

var ErrOriginNotLoaded = errors.New("origin model is not loaded in graph")
var ErrUntrackedGraphChange = errors.New("world graph has untracked changes")

type FromName string
type ToName string

type Builder struct {
	Pool       WorldPool
	Graph      graph.Builder
	NameCache  map[components.Name]entity.View
	TypedStrs  mirrors.TypedStrings
	Components *table.Table
}

func DefaultBuilder() *Builder {
	tbl := table.New(10000)
	pool := bitpool.FromTable(tbl, 600)
	return NewBuilder(pool, tbl)
}

func NewBuilder(pool entity.Pool, tbl *table.Table) *Builder {
	return &Builder{
		Pool:       WorldPool{pool},
		Graph:      graph.Builder{G: graph.New()},
		NameCache:  make(map[components.Name]entity.View, 128),
		TypedStrs:  mirrors.NewTypedStrings(),
		Components: tbl,
	}
}

// after calling this it is no longer safe to interact with the builder
func (w *Builder) Build() World {
	return World{
		Entities: w.Pool,
		Graph:    w.Graph.G,
	}
}

// unique names means we can forward declare entities w/o worry
func (w *Builder) Entity(n components.Name) (entity.View, error) {
	if ent, ok := w.NameCache[n]; ok {
		return ent, nil
	}

	ent, err := w.Pool.Create(n)
	if err != nil {
		return nil, err
	}

	w.NameCache[n] = ent
	return ent, nil
}

func (w *Builder) Node(v entity.View) {
	w.Graph.AddNode(graph.Node(v.Model()))
}

func (w *Builder) Edge(origin, destination entity.View) (entity.View, error) {
	// we require that origin exist in the graph already rather than just adding it
	successors, err := w.Graph.G.Successors(graph.Node(origin.Model()))
	if err != nil {
		if errors.Is(err, graph.ErrOriginNotFound) {
			return nil, ErrOriginNotLoaded
		}
		return nil, err
	}

	var conns Connections

	if err := origin.Get(&conns); err != nil {
		if errors.Is(err, entity.ErrNotLoaded) || errors.Is(err, entity.ErrNotAssigned) {
			conns = make(Connections, 4)
			origin.Add(conns)
		} else {
			return nil, err
		}
	}

	if reiterate.Contains(graph.Destination(destination.Model()), successors) {
		edgeId, ok := conns[destination.Model()]
		if !ok {
			panic(ErrUntrackedGraphChange)
		}

		edge, err := w.Pool.Fetch(edgeId)
		if err != nil {
			panic(err)
		}
		return edge, nil
	}

	edgeName, fromName, toName := namesForEdge(origin, destination)
	edge, err := w.Pool.Create(edgeName)
	if err != nil {
		panic(err)
	}

	conns[destination.Model()] = edge.Model()
	if err := w.Graph.AddEdge(graph.Origination(origin.Model()), graph.Destination(destination.Model())); err != nil {
		panic(err)
	}
	archetype := edgeArchetype{
		origination:     origin.Model(),
		destination:     destination.Model(),
		originationName: fromName,
		destinationName: toName,
	}

	if err := archetype.Apply(edge); err != nil {
		panic(err)
	}

	return edge, nil
}

type edgeArchetype struct {
	origination     entity.Model
	originationName FromName
	destination     entity.Model
	destinationName ToName
}

func (e edgeArchetype) Apply(entity entity.View) error {
	if err := entity.Add(Edge{
		Origination: e.origination,
		Destination: e.destination,
	}); err != nil {
		return err
	}
	if err := entity.Add(e.originationName); err != nil {
		return err
	}
	if err := entity.Add(e.destinationName); err != nil {
		return err
	}
	return nil
}

func namesForEdge(origin, destination entity.View) (components.Name, FromName, ToName) {
	var from components.Name
	var to components.Name

	if err := origin.Get(&from); err != nil {
		panic(err)
	}
	if err := destination.Get(&to); err != nil {
		panic(err)
	}

	return components.Name(fmt.Sprintf("%s -> %s", from, to)), FromName(from), ToName(to)
}
