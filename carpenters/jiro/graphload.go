package jiro

import (
	"slices"
	"sudonters/zootler/internal/components"
	"sudonters/zootler/internal/entities"
	"sudonters/zootler/internal/entity"
	"sudonters/zootler/internal/table"
	"sudonters/zootler/internal/world"

	"github.com/etc-sudonters/substrate/skelly/graph"
	"github.com/etc-sudonters/substrate/slipup"
)

type graphloader struct {
	locs entities.Locations
	edge entities.Edges
	grph graph.Builder
	root world.Root
}

func (l *graphloader) load(wn worldnode) error {
	origin, originErr := l.locs.Entity(components.Name(wn.RegionName))
	paniconerr(originErr)
	if err := origin.AddComponents(wn.AsComponents()); err != nil {
		paniconerr(err)
	}

	if origin.Name() == "Root" {
		l.root = world.Root(origin.Id())
	}

	for destName, nodeEdge := range wn.Edges {
		dest, destErr := l.locs.Entity(destName)
		paniconerr(destErr)
		_, edgeErr := l.connect(nodeEdge.name, origin, dest, nodeEdge.kind, nodeEdge.rule)
		paniconerr(edgeErr)
	}

	return nil
}

func (l *graphloader) connect(name components.Name, origin, dest entities.Location, kind edgekind, rule components.RawLogic) (entities.Edge, error) {
	edge, edgeErr := l.edge.Entity(name)
	if edgeErr != nil {
		return entities.Edge{}, slipup.Describef(edgeErr, "edge %s", name)
	}
	l.grph.AddEdge(graph.Origination(origin.Id()), graph.Destination(dest.Id()))

	edge.Stash("origin", string(origin.Name()))
	edge.Stash("dest", string(dest.Name()))
	edge.StashRawRule(rule)
	comps := slices.Concat(kind.ascomponents(), table.Values{
		rule,
		components.Connection{
			Origin: entity.Model(origin.Id()),
			Dest:   entity.Model(dest.Id()),
		},
	})

	if err := edge.AddComponents(comps); err != nil {
		return edge, slipup.Describef(err, "while adding components to '%s'", name)
	}

	return edge, nil
}

func paniconerr(e error) {
	if e != nil {
		panic(e)
	}
}
