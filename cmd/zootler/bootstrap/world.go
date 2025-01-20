package bootstrap

import (
	"sudonters/zootler/internal/skelly/graph32"
	"sudonters/zootler/magicbean"
	"sudonters/zootler/zecs"
)

func explorableworldfrom(ocm *zecs.Ocm) magicbean.ExplorableWorld {
	var world magicbean.ExplorableWorld
	q := ocm.Query()
	q.Build(
		zecs.Load[magicbean.RuleCompiled],
		zecs.Load[magicbean.EdgeKind],
		zecs.Load[magicbean.Connection],
		zecs.Load[magicbean.Name],
	)

	rows, err := q.Execute()
	PanicWhenErr(err)

	world.Edges = make(map[magicbean.Connection]magicbean.ExplorableEdge, rows.Len())
	world.Graph = graph32.WithCapacity(rows.Len() * 2)

	directed := graph32.Builder{world.Graph}

	for entity, tup := range rows.All {
		trans := tup.Values[2].(magicbean.Connection)
		directed.AddEdge(graph32.Node(trans.From), graph32.Node(trans.To))
		world.Edges[trans] = magicbean.ExplorableEdge{
			Entity: entity,
			Rule:   tup.Values[0].(magicbean.RuleCompiled),
			Kind:   tup.Values[1].(magicbean.EdgeKind),
			Name:   tup.Values[3].(magicbean.Name),
		}
	}

	return world
}
