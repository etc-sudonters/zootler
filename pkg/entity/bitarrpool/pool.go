package bitarrpool

import (
	"reflect"

	"sudonters/zootler/internal/bitset"
	"sudonters/zootler/pkg/entity"
)

type bitarrpool struct {
	componentBucketCount int64
	entities             []bitarrview
	table                componentTable
}

var _ entity.Pool = (*bitarrpool)(nil)
var _ entity.View = (*bitarrview)(nil)

func New(maxComponentId int64) *bitarrpool {
	var b bitarrpool
	b.componentBucketCount = bitset.Buckets(maxComponentId)
	b.entities = make([]bitarrview, 1, 128)
	(&b.table).init()
	return &b
}

func (p bitarrpool) String() string {
	return CompressedRepr{p}.String()
}

func (p *bitarrpool) Create() (entity.View, error) {
	var view bitarrview
	view.id = entity.Model(len(p.entities))
	view.comps = bitset.New(p.componentBucketCount)
	view.p = p
	p.entities = append(p.entities, view)
	return view, nil
}

// return a subset of the population that matches the provided selectors
func (p *bitarrpool) Query(qs []entity.Selector) ([]entity.View, error) {
	includeMask := bitset.New(p.componentBucketCount)
	excludeMask := bitset.New(p.componentBucketCount)
	for _, q := range qs {
		component := q.Component()
		id, ok := p.table.idType(component)
		if !ok {
			return nil, entity.ErrUnknownComponent{Type: component}
		}
		switch q.Behavior() {
		case entity.ComponentInclude:
			includeMask.Set(int64(id))
			break
		case entity.ComponentExclude:
			excludeMask.Set(int64(id))
			break
		}
	}

	if bitset.IsEmpty(includeMask) && bitset.IsEmpty(excludeMask) {
		// empty bitmask will never select anything
		return nil, entity.ErrNoEntities
	}

	var entities []entity.View
	for _, v := range p.entities[1:] {
		if includeMask.Intersect(v.comps).Eq(includeMask) &&
			v.comps.Difference(excludeMask).Eq(v.comps) {
			entities = append(entities, v)
		}
	}

	if entities == nil || len(entities) == 0 {
		return nil, entity.ErrNoEntities
	}

	return entities, nil
}

func (p *bitarrpool) Get(m entity.Model, cs ...interface{}) {
	for i := range cs {
		_ = entity.AssignComponentTo(cs[i], getComponenter(p.entities[m]))
	}
}

func (p *bitarrpool) addCompToEnt(b bitarrview, c entity.Component) error {
	row := p.table.rowFor(c)
	row.set(b.Model(), c)
	id, _ := p.table.idValue(c)
	b.comps.Set(int64(id))
	return nil
}

func (p *bitarrpool) removeCompFromEnt(b bitarrview, c entity.Component) error {
	id, ok := p.table.idValue(c)
	if !ok {
		return entity.ErrUnknownComponent{Type: reflect.TypeOf(c)}
	}
	row := p.table.row(id)
	row.unset(b.Model())
	b.comps.Clear(int64(id))
	return nil
}
