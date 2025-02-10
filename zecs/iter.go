package zecs

import "iter"

func IterEntities[V Value](ocm *Ocm, qs ...BuildQuery) iter.Seq[Entity] {
	return func(yield func(Entity) bool) {
		q := ocm.Query()
		q.Build(With[V], qs...)

		for row, _ := range q.Rows {
			if !yield(Entity(row)) {
				return
			}
		}
	}
}
