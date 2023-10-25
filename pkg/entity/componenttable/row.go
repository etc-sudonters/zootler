package componenttable

import "sudonters/zootler/pkg/entity"

type Row struct {
	id         entity.ComponentId
	components []entity.Component
}

func (r *Row) Init(id entity.ComponentId) {
	r.id = id
	r.components = make([]entity.Component, 0)
}

func (row *Row) Set(e entity.Model, c entity.Component) {
	row.EnsureSize(int(e))
	row.components[e] = c
}

func (row *Row) Unset(e entity.Model) {
	if len(row.components) < int(e) {
		return
	}

	row.components[e] = nil
}

func (row Row) Get(e entity.Model) entity.Component {
	if len(row.components) <= int(e) {
		return nil
	}

	return row.components[e]
}

func (row *Row) EnsureSize(n int) {
	if len(row.components) > n {
		return
	}

	expaded := make([]entity.Component, n+1, n*2)
	copy(expaded, row.components)
	row.components = expaded
}
