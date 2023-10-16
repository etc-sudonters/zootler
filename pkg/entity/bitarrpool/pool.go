package bitarrpool

import (
	"reflect"

	"sudonters/zootler/internal/bitset"
	"sudonters/zootler/pkg/entity"
)

type bitarrpool struct {
	componentBucketCount int64

	entities  []bitarrview
	maxCompId componentId
	compIds   map[reflect.Type]componentId
	table     componentTable
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

func (p *bitarrpool) Create() (entity.View, error) {
	var view bitarrview
	view.id = entity.Model(len(p.entities))
	view.comps = bitset.New(p.componentBucketCount)
	view.p = p
	p.entities = append(p.entities, view)
	return &view, nil
}

// return a subset of the population that matches the provided selectors
func (p *bitarrpool) Query(qs []entity.Selector) ([]entity.View, error) {
	mask := bitset.New(p.componentBucketCount)
	for _, q := range qs {
		component := q.Component()
		id, ok := p.table.id(component)
		if !ok {
			return nil, entity.ErrUnknownComponent{Type: component}
		}
		switch q.Behavior() {
		case entity.ComponentInclude:
			mask.Set(int64(id))
			break
		case entity.ComponentExclude:
			mask.Clear(int64(id))
			break
		}
	}

	if bitset.IsEmpty(mask) {
		// empty bitmask will never select anything
		return nil, entity.ErrNoEntities
	}

	var entities []entity.View
	for _, v := range p.entities {
		if v.mask(mask).Eq(mask) {
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
	return nil
}

func (p *bitarrpool) removeCompFromEnt(b bitarrview, c entity.Component) error {
	id, ok := p.table.id(c)
	if !ok {
		return entity.ErrUnknownComponent{Type: reflect.TypeOf(c)}
	}
	row := p.table.row(id)
	row.unset(b.Model())
	return nil
}
