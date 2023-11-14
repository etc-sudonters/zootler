package bitpool

import (
	"fmt"
	"sudonters/zootler/internal/entity"
	"sudonters/zootler/internal/entity/componenttable"
	"sudonters/zootler/internal/mirrors"

	"github.com/etc-sudonters/substrate/skelly/set/bits"
)

type bitpool struct {
	componentBucketCount int
	entities             []bitview
	table                *componenttable.Table
}

var _ entity.Pool = (*bitpool)(nil)
var _ entity.View = bitview{}

type Settings struct {
	MaxComponentId int
	MaxEntityId    int
}

func New(s Settings) *bitpool {
	var b bitpool
	b.componentBucketCount = bits.Buckets(s.MaxComponentId)
	b.entities = make([]bitview, 1, 128)
	b.table = componenttable.New(s.MaxEntityId)
	return &b
}

func FromTable(tbl *componenttable.Table, maxComponentId int) *bitpool {
	var b bitpool
	b.componentBucketCount = bits.Buckets(maxComponentId)
	b.table = tbl
	b.entities = make([]bitview, 128)
	return &b
}

func (p *bitpool) Create() (entity.View, error) {
	var view bitview
	view.id = entity.Model(len(p.entities))
	view.comps = bits.New(p.componentBucketCount)
	view.p = p
	p.entities = append(p.entities, view)
	return view, nil
}

// return a subset of the population that matches the provided selectors
func (p *bitpool) Query(qs []entity.Selector) ([]entity.View, error) {
	var filter filter
	(&filter).init(p.componentBucketCount)

	for _, q := range qs {
		typ := q.Component()
		id, err := p.table.IdOf(typ)
		if err != nil {
			name := typ.Name()
			if name == "" {
				if n, ok := mirrors.TryGetLiteral(typ); ok {
					name = n
				}
			}
			return nil, fmt.Errorf("during component %s: %w", name, err)
		}

		switch q.Behavior() {
		case entity.ComponentInclude:
			filter.include(id)
			break
		case entity.ComponentExclude:
			filter.exclude(id)
			break
		default:
			return nil, fmt.Errorf("during component %s: unknown behavior: %s", typ.Name(), q.Behavior())
		}
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
