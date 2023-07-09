package graph

import "errors"

var ErrNodeNotFound = errors.New("node not found")

type Builder struct {
	G Model
}

func (b *Builder) AddNode() Node {
	newId := Node(len(b.G.nodes))
	b.G.nodes = append(b.G.nodes, make(Neighbors))
	b.G.inEdges[newId] = make(Neighbors)
	return newId
}

func (b *Builder) AddEdge(o Origination, d Destination) error {
	origin := Node(o)
	destination := Node(d)
	if !b.G.canNodeExist(origin) || !b.G.canNodeExist(destination) {
		return ErrNodeNotFound
	}

	b.G.nodes[origin].Add(destination)
	b.G.inEdges[destination].Add(origin)
	return nil
}
