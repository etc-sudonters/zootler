package world

import (
	"fmt"

	"github.com/etc-sudonters/zootler/internal/datastructures/graph"
	"github.com/etc-sudonters/zootler/pkg/entity"
	"github.com/etc-sudonters/zootler/pkg/entity/hashpool"
	"github.com/etc-sudonters/zootler/pkg/logic"
)

type Builder struct {
	id        Id
	Pool      Pool
	graph     graph.Builder
	edgeCache map[edge]entity.View
	nodeCache map[graph.Node]entity.View
}

// caller is responsible for setting a unique id if necessary
func NewBuilder(id Id) Builder {
	return Builder{
		id,
		Pool{id, hashpool.New()},
		graph.Builder{graph.New()},
		make(map[edge]entity.View),
		make(map[graph.Node]entity.View),
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
		Entities:  w.Pool,
		Graph:     w.graph.G,
		edgeCache: edgeCache,
		nodeCache: nodeCache,
	}
}

func (w *Builder) AddEntity(n logic.Name) (entity.View, error) {
	ent, err := w.Pool.Create(n)
	if err != nil {
		return nil, err
	}

	return ent, nil
}

func (w *Builder) AddNode(v entity.View) {
	w.graph.AddNode(graph.Node(v.Model()))
}

func (w *Builder) AddEdge(origin, destination entity.View) (entity.View, error) {
	var oName logic.Name
	var dName logic.Name
	origin.Get(&oName)
	destination.Get(&dName)

	name := logic.Name(fmt.Sprintf("%s -> %s", oName, dName))

	o := graph.Origination(origin.Model())
	d := graph.Destination(destination.Model())

	if ent, ok := w.edgeCache[edge{o, d}]; ok {
		return ent, nil
	}

	if err := w.graph.AddEdge(o, d); err != nil {
		return nil, err
	}

	ent, err := w.Pool.Create(name)

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
