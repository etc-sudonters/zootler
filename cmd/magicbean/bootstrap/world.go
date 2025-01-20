package bootstrap

import (
	"sudonters/zootler/magicbean"
	"sudonters/zootler/zecs"

	"github.com/etc-sudonters/substrate/skelly/graph"
)

func explorableworldfrom(ocm *zecs.Ocm) magicbean.ExplorableWorld {
	var world magicbean.ExplorableWorld
	q := ocm.Query()
	q.Build(
		zecs.Load[magicbean.RuleCompiled],
		zecs.Load[magicbean.EdgeKind],
		zecs.Load[magicbean.Connection],
	)

	rows, err := q.Execute()
	PanicWhenErr(err)

	world.Edges = make(map[magicbean.Connection]magicbean.ExplorableEdge, rows.Len())
	world.Graph = graph.WithCapacity(rows.Len() * 2)

	directed := graph.Builder{world.Graph}

	for entity, tup := range rows.All {
		trans := tup.Values[2].(magicbean.Connection)
		directed.AddEdge(graph.Origination(trans.From), graph.Destination(trans.To))
		world.Edges[trans] = magicbean.ExplorableEdge{
			Entity: entity,
			Rule:   tup.Values[0].(magicbean.RuleCompiled),
			Kind:   tup.Values[1].(magicbean.EdgeKind),
		}
	}

	return world
}
