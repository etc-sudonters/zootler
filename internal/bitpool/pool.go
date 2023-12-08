package bitpool

import (
	"fmt"
	"reflect"
	"sudonters/zootler/internal/entity"
	"sudonters/zootler/internal/entity/table"

	"github.com/etc-sudonters/substrate/mirrors"
	"github.com/etc-sudonters/substrate/skelly/bitset"
)

type bitpool struct {
	componentBucketCount int
	entities             []bitview
	table                *table.Table
}

var _ entity.Pool = (*bitpool)(nil)
var _ entity.View = bitview{}

type Settings struct {
	MaxComponentId int
	MaxEntityId    int
}

func New(s Settings) *bitpool {
	var b bitpool
	b.componentBucketCount = bitset.Buckets(s.MaxComponentId)
	b.entities = make([]bitview, 1, 128)
	b.table = table.New(s.MaxEntityId)
	return &b
}

func FromTable(tbl *table.Table, maxComponentId int) *bitpool {
	var b bitpool
	b.componentBucketCount = bitset.Buckets(maxComponentId)
	b.table = tbl
	b.entities = make([]bitview, 128)
	return &b
}

func (p *bitpool) Create() (entity.View, error) {
	var view bitview
	view.id = entity.Model(len(p.entities))
	view.comps = bitset.New(p.componentBucketCount)
	view.p = p
	p.entities = append(p.entities, view)
	return view, nil
}

// return a subset of the population that matches the provided selectors
func (p *bitpool) Query(f entity.Filter) ([]entity.View, error) {
	var filter filter
	(&filter).init(p.componentBucketCount)

	getTypeId := func(typ reflect.Type) (entity.ComponentId, error) {
		id, err := p.table.IdOf(typ)
		if err != nil {
			name := typ.Name()
			if name == "" {
				if n, ok := mirrors.TryGetLiteral(typ); ok {
					name = n
				}
			}
			return 0, fmt.Errorf("during component %s: %w", name, err)
		}

		return id, nil
	}

	for _, typ := range f.With() {
		id, err := getTypeId(typ)
		if err != nil {
			return nil, err
		}
		filter.include(id)
	}

	for _, typ := range f.Without() {
		id, err := getTypeId(typ)
		if err != nil {
			return nil, err
		}
		filter.exclude(id)
	}

	var entities []entity.View

	for _, e := range p.entities {
		e := e
		if filter.test(e) {
			entities = append(entities, e)
		}
	}

	if len(entities) == 0 {
		return nil, entity.ErrNoEntities
	}

	return entities, nil
}

func (p *bitpool) Get(m entity.Model, cs []interface{}) {
	for i := range cs {
		_ = entity.AssignComponentTo(m, cs[i], p.table.Getter())
	}
}

func (p *bitpool) Fetch(m entity.Model) (entity.View, error) {
	if int(m) >= len(p.entities) {
		return nil, entity.ErrEntityNotExist
	}

	return p.entities[int(m)], nil
}

func (p *bitpool) addCompToEnt(b bitview, c entity.Component) error {
	id := p.table.Set(b.id, c)
	b.comps.Set(int(id))
	return nil
}

func (p *bitpool) removeCompFromEnt(b bitview, c entity.Component) error {
	if id := p.table.Unset(b.id, entity.PierceComponentType(c)); id != entity.INVALID_COMPONENT {
		b.comps.Clear(int(id))
	}
	return nil
}
