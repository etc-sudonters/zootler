package hashpool

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/etc-sudonters/rando/entity"
	"github.com/etc-sudonters/rando/set"
)

var NotLoaded = errors.New("not loaded")

type entityBucket map[entity.Model]entity.Component
type componentBucket map[reflect.Type]entity.Component
type storage map[reflect.Type]entityBucket

var modelType = reflect.TypeOf(entity.Model(0))

func interfaceToComponentType(i interface{}) (reflect.Type, error) {
	return nil, nil
}

type Pool struct {
	population set.Hash[entity.Model]
	membership storage
	nextModel  entity.Model
	debug      func(string, ...any)
}

func EnsureTable(p *Pool, component entity.Component) error {
	componentType := reflect.TypeOf(component)
	if _, ok := p.membership[componentType]; !ok {
		p.membership[componentType] = make(entityBucket)
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

func (p Pool) All() set.Hash[entity.Model] {
	return set.FromMap(p.population)
}

func (p *Pool) createEasy() view {
	thisModel := p.nextModel
	p.nextModel++

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

func (p *Pool) Query(qs ...entity.Selector) ([]entity.View, error) {
	if len(qs) == 0 {
		p.debug("no queries")
		return nil, nil
	}

	if len(p.population) == 0 {
		p.debug("no population")
		return nil, nil
	}

	population := map[entity.Model]componentBucket{}

	// first generation is the entire population
	for member, component := range p.membership[entity.ModelComponentType] {
		population[member] = componentBucket{
			entity.ModelComponentType: component,
		}
	}

	if population == nil {
		p.debug("nil population!")
	}

	p.debug("starting with total population: %+v", p.membership[entity.ModelComponentType])
	p.addToPopulation(
		population,
		entity.ModelComponentType,
		p.membership[entity.ModelComponentType],
	)

	for _, q := range qs {
		needed := q.Component()

		nextGenParents, ok := p.membership[needed]
		if !ok {
			return nil, fmt.Errorf("unknown component: %s", entity.NiceTypeName(needed))
		}

		p.debug("next generation parent population: %+v", nextGenParents)

		nextGeneration := q.Select(set.FromMap(population), nextGenParents)
		set.DiscardUsing(population, set.FromMap(nextGeneration))

		if len(population) == 0 {
			p.debug("population went extinct")
			return nil, nil
		}

		p.addToPopulation(population, needed, nextGenParents)
	}

	viewing := make([]entity.View, 0, len(population))

	for entity, table := range population {
		viewing = append(viewing, &view{
			m:      entity,
			origin: p,
			loaded: table,
		})
	}

	return viewing, nil
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

type view struct {
	m       entity.Model
	origin  *Pool
	loaded  map[reflect.Type]entity.Component
	session map[reflect.Type]entity.Component
}

func (v view) checkDetached() {
	if v.origin == nil {
		panic("detatched view")
	}
}

func (v view) Model() entity.Model {
	v.checkDetached()
	return v.m
}

func (v view) Get(target interface{}) error {
	v.checkDetached()

	if target == nil {
		return errors.New("nil component reference")
	}

	value := reflect.ValueOf(target)
	typ := value.Type()

	if typ.Kind() != reflect.Pointer || value.IsNil() {
		return errors.New("non-nil pointers only")
	}

	targetType := typ.Elem()

	tryFind := func(t reflect.Type) (entity.Component, error) {
		v.origin.debug("target load type %s", entity.NiceTypeName(t))
		acquired, ok := v.loaded[t]
		if !ok {
			v.origin.debug("attempting to load %s from session", entity.NiceTypeName(t))
			acquired, ok = v.session[t]
			if !ok {
				return nil, NotLoaded
			}
		}

		return acquired, nil
	}

	acquired, err := tryFind(targetType)
	if err != nil {
		if errors.Is(err, NotLoaded) && targetType.Kind() == reflect.Pointer {
			v.origin.debug("pointer to pointer? try dereferencing once")
			acquired, err = tryFind(targetType.Elem())
		}

		if err != nil {
			return err
		}
	}

	if acquired == nil {
		panic(fmt.Sprintf("nil component loaded for %s on model{%v}", entity.NiceTypeName(targetType), v.m))
	}
	v.origin.debug("acquired %+v", acquired)

	acquiredValue := reflect.ValueOf(acquired)
	v.origin.debug("acquired value's type: %s", acquiredValue.Type())

	if acquiredValue.Kind() != reflect.Pointer && targetType.Kind() == reflect.Pointer {
		intermediate := reflect.New(acquiredValue.Type())
		intermediate.Elem().Set(acquiredValue)
		acquiredValue = intermediate
	}

	value.Elem().Set(acquiredValue)
	return nil
}

func (v view) Add(target entity.Component) error {
	v.checkDetached()

	typ := reflect.TypeOf(target)
	if _, ok := v.loaded[typ]; ok {
		v.loaded[typ] = target
	} else {
		if err := EnsureTable(v.origin, target); err != nil {
			return err
		}
		v.session[typ] = target
	}

	v.origin.membership[typ][v.m] = target

	return nil
}

func (v view) Remove(target entity.Component) error {
	v.checkDetached()

	typ := reflect.TypeOf(target)

	if typ == entity.ModelComponentType {
		return errors.New("cannot remove model component")
	}

	delete(v.loaded, typ)
	delete(v.session, typ)
	removeFromTable(v.m, typ, v.origin)
	return nil
}

var _ entity.Pool = (*Pool)(nil)

func removeFromTable(entity entity.Model, compType reflect.Type, origin *Pool) {
	delete(origin.membership[compType], entity)
}

func (p *Pool) addToPopulation(
	population map[entity.Model]componentBucket,
	component reflect.Type,
	nextGeneration map[entity.Model]entity.Component,
) {
	for member := range population {
		if population[member] == nil {
			population[member] = make(map[reflect.Type]entity.Component)
		}
	}

	p.debug("population size: %d", len(population))
	p.debug("%+v", population)

	for k, v := range nextGeneration {
		p.debug("succeeding model{%v} into next generation", k)
		population[k][component] = v
	}
}
