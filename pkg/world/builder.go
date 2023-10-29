package world

import (
	"errors"
	"fmt"

	"sudonters/zootler/internal/entity"
	"sudonters/zootler/internal/entity/bitpool"

	"github.com/etc-sudonters/substrate/reiterate"
	"github.com/etc-sudonters/substrate/skelly/graph"
)

var ErrOriginNotLoaded = errors.New("origin model is not loaded in graph")
var ErrUntrackedGraphChange = errors.New("world graph has untracked changes")

type Builder struct {
	Pool      Pool
	graph     graph.Builder
	nameCache map[Name]entity.View
}

func NewBuilder() *Builder {
	return &Builder{
		Pool{bitpool.New(bitpool.Settings{
			MaxComponentId: 300,
			MaxEntityId:    2000,
		})},
		graph.Builder{G: graph.New()},
		make(map[Name]entity.View, 128),
	}
}

// after calling this it is no longer safe to interact with the builder
func (w *Builder) Build() World {
	return World{
		Entities: w.Pool,
		Graph:    w.graph.G,
	}
}

// unique names means we can forward declare entities w/o worry
func (w *Builder) AddEntity(n Name) (entity.View, error) {
	if ent, ok := w.nameCache[n]; ok {
		return ent, nil
	}

	ent, err := w.Pool.Create(n)
	if err != nil {
		return nil, err
	}

	w.nameCache[n] = ent
	return ent, nil
}

func (w *Builder) AddNode(v entity.View) {
	w.graph.AddNode(graph.Node(v.Model()))
}

func (w *Builder) AddEdge(origin, destination entity.View) (entity.View, error) {
	// we require that origin exist in the graph already rather than just adding it
	successors, err := w.graph.G.Successors(graph.Node(origin.Model()))
	if err != nil {
		if errors.Is(err, graph.ErrOriginNotFound) {
			return nil, ErrOriginNotLoaded
		}
		return nil, err
	}

	var conns Connections

	if err := origin.Get(&conns); err != nil {
		if errors.Is(err, entity.ErrNotLoaded) {
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

	edge, err := w.Pool.Create(nameEdgeFrom(origin, destination))
	if err != nil {
		panic(err)
	}

	conns[destination.Model()] = edge.Model()
	if err := w.graph.AddEdge(graph.Origination(origin.Model()), graph.Destination(destination.Model())); err != nil {
		panic(err)
	}

	edge.Add(Edge{
		Origination: origin.Model(),
		Destination: origin.Model(),
	})

	return edge, nil
}

func nameEdgeFrom(origin, destination entity.View) Name {
	var from Name
	var to Name

	if err := origin.Get(&from); err != nil {
		panic(err)
	}
	if err := destination.Get(&to); err != nil {
		panic(err)
	}

	return Name(fmt.Sprintf("%s -> %s", from, to))
}
