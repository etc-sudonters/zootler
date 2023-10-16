package bitarrpool

import (
	"reflect"
	"sudonters/zootler/pkg/entity"
)

type componentId int64

const (
	INVALID_COMPONENT componentId  = 0
	INVALID_ENTITY    entity.Model = 0
)

type componentTable struct {
	rows   []*componentRow
	lookup map[reflect.Type]componentId
}

func (t *componentTable) rowFor(c entity.Component) *componentRow {
	id, ok := t.id(c)
	if !ok {
		return t.addrow(reflect.TypeOf(c))
	}

	return t.rows[id]
}

func (t componentTable) id(c entity.Component) (componentId, bool) {
	id, ok := t.lookup[reflect.TypeOf(c)]
	if !ok {
		return INVALID_COMPONENT, false
	}
	return id, true
}

func (c componentTable) row(r componentId) *componentRow {
	return c.rows[r]
}

func (c *componentTable) addrow(t reflect.Type) *componentRow {
	row := make(componentRow, 1, 128)
	c.lookup[t] = componentId(len(c.rows))
	c.rows = append(c.rows, &row)
	return &row
}

func (c *componentTable) init() {
	c.rows = make([]*componentRow, 1, 128)
	c.lookup = make(map[reflect.Type]componentId, 128)
	c.lookup[nil] = 0
}

type componentRow []entity.Component // idx'd by entId

func (r *componentRow) set(e entity.Model, c entity.Component) {
	r.ensureSize(int(e))
	(*r)[e] = c
}

func (r *componentRow) unset(e entity.Model) {
	if len(*r) < int(e) {
		return
	}

	(*r)[e] = nil
}

func (r componentRow) get(e entity.Model) entity.Component {
	if len(r) < int(e) {
		return nil
	}

	return r[e]
}

func (r *componentRow) ensureSize(n int) {
	if len(*r) > n {
		return
	}

	row := make(componentRow, n, n*2)
	copy(row, *r)
	*r = row
}
