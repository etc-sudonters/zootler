package bitarrpool

import (
	"reflect"
	"sudonters/zootler/pkg/entity"
)

func getComponenter(b bitarrview) func(reflect.Type) (entity.Component, error) {
	return func(t reflect.Type) (entity.Component, error) {
		id, ok := b.p.table.lookup[t]
		if !ok {
			return nil, entity.ErrUnknownComponent{Type: t}
		}

		if !b.comps.Test(int64(id)) {
			return nil, entity.ErrNotAssigned
		}

		row := b.p.table.row(id)
		comp := row.get(b.id)
		if comp == nil {
			return nil, entity.ErrNilComponent{
				Entity:    b.Model(),
				Component: t,
			}
		}

		return comp, nil
	}
}
