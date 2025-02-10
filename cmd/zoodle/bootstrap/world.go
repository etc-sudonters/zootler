package bootstrap

import (
	"sudonters/libzootr/components"
	"sudonters/libzootr/internal/skelly/graph32"
	"sudonters/libzootr/magicbean"
	"sudonters/libzootr/zecs"
)

func explorableworldfrom(ocm *zecs.Ocm) magicbean.ExplorableWorld {
	var world magicbean.ExplorableWorld
	q := ocm.Query()
	q.Build(
		zecs.Load[components.RuleCompiled],
		zecs.Load[components.EdgeKind],
		zecs.Load[components.Connection],
		zecs.Load[components.Name],
		zecs.Optional[components.RuleSource],
	)

	rows, err := q.Execute()
	PanicWhenErr(err)

	world.Edges = make(map[components.Connection]magicbean.ExplorableEdge, rows.Len())
	world.Graph = graph32.WithCapacity(rows.Len() * 2)
	directed := graph32.Builder{Graph: &world.Graph}
	roots := zecs.EntitiesMatching(ocm, zecs.With[components.WorldGraphRoot])
	if len(roots) == 0 {
		panic("no graph roots loaded")
	}
	for _, root := range roots {
		directed.AddRoot(graph32.Node(root))
	}

	for entity, tup := range rows.All {
		trans := tup.Values[2].(components.Connection)
		directed.AddEdge(graph32.Node(trans.From), graph32.Node(trans.To))
		edge := magicbean.ExplorableEdge{
			Entity: entity,
			Rule:   tup.Values[0].(components.RuleCompiled),
			Kind:   tup.Values[1].(components.EdgeKind),
			Name:   tup.Values[3].(components.Name),
		}

		src := tup.Values[4]
		if src != nil {
			edge.Src = src.(components.RuleSource)
		}

		world.Edges[trans] = edge
	}

	return world
}
