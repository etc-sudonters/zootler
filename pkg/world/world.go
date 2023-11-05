package world

import (
	"errors"
	"fmt"
	"sudonters/zootler/internal/entity"

	"github.com/etc-sudonters/substrate/skelly/graph"
)

var ErrEntityNotConnected = errors.New("entity is not connected to other entities")
var ErrEntitiesNotConnected = errors.New("the entities are not connected")

type Name string

type World struct {
	Entities Pool
	Graph    graph.Directed
}

type Pool struct {
	entity.Pool
	ByName map[Name]entity.Model
}

func (p Pool) Create(name Name) (entity.View, error) {
	if model, ok := p.ByName[name]; ok {
		return p.Fetch(model)
	}

	view, err := p.Pool.Create()
	if err != nil {
		return nil, fmt.Errorf("failed to create entity %q: %w", name, err)
	}

	p.ByName[name] = view.Model()

	view.Add(name)
	return view, nil
}

func (w World) Edge(e Edge) (entity.View, error) {
	var conns Connections
	w.Entities.Get(e.Origination, []interface{}{&conns})
	if conns == nil {
		return nil, ErrEntityNotConnected
	}

	edgeId, ok := conns[e.Destination]
	if !ok {
		return nil, ErrEntitiesNotConnected
	}

	edge, err := w.Entities.Fetch(edgeId)
	if err != nil {
		return nil, err
	}

	return edge, nil
}
