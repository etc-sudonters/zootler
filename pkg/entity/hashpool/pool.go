package hashpool

import (
	"reflect"

	"github.com/etc-sudonters/zootler/internal/datastructures/set"
	"github.com/etc-sudonters/zootler/pkg/entity"
)

var _ entity.Pool = (*Pool)(nil)

type entityBuckets map[entity.Model]entity.Component
type componentBuckets map[reflect.Type]entity.Component
type storage map[reflect.Type]entityBuckets

type Pool struct {
	population set.Hash[entity.Model]
	membership storage
	lastModel  entity.Model
	debug      func(string, ...any)
}

func EnsureTable(p *Pool, component entity.Component) error {
	componentType := reflect.TypeOf(component)
	if _, ok := p.membership[componentType]; !ok {
		p.membership[componentType] = make(entityBuckets)
	}

	return nil
}

func New() (*Pool, error) {
	p := &Pool{
		population: make(set.Hash[entity.Model]),
		membership: make(storage),
		debug:      func(s string, a ...any) {},
	}

	err := EnsureTable(p, entity.Model(0))

	return p, err
}

func (p *Pool) createEasy() view {
	p.lastModel++
	thisModel := p.lastModel

	v := view{
		m:       thisModel,
		origin:  p,
		loaded:  make(map[reflect.Type]entity.Component),
		session: make(map[reflect.Type]entity.Component),
	}

	v.Add(thisModel)
	v.origin.population.Add(thisModel)

	return v
}

func (p *Pool) Create() (entity.View, error) {
	return p.createEasy(), nil
}

func (p *Pool) Delete(v entity.View) error {
	m := v.(view)
	model := m.m
	delete(p.population, model)

	for _, members := range p.membership {
		delete(members, model)
	}

	m.loaded = nil
	m.session = nil
	m.origin = nil

	return nil
}

func removeFromTable(entity entity.Model, compType reflect.Type, origin *Pool) {
	delete(origin.membership[compType], entity)
}
