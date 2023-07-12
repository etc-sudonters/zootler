package hashpool

import (
	"fmt"
	"reflect"

	"github.com/etc-sudonters/zootler/internal/bag"
	"github.com/etc-sudonters/zootler/internal/datastructures/set"
	"github.com/etc-sudonters/zootler/pkg/entity"
)

func (p Pool) All() set.Hash[entity.Model] {
	return set.FromMap(p.population)
}

func (p *Pool) Get(e entity.Model, comps ...interface{}) {
	retrieve := func(t reflect.Type) (entity.Component, error) {
		ents, ok := p.membership[t]
		if !ok {
			return nil, entity.ErrNotLoaded
		}

		comp, ok := ents[e]
		if !ok {
			return nil, entity.ErrNotLoaded
		}

		return comp, nil
	}

	for i := range comps {
		assignComponentTo(comps[i], retrieve)
	}
}

func (p *Pool) Query(basePopulation entity.Component, qs ...entity.Selector) ([]entity.View, error) {
	if len(p.population) == 0 {
		p.debug("no population")
		return nil, nil
	}

	var population entity.Population
	componentsToLoad := make([]reflect.Type, 0, len(qs)+1)

	// first generation spawns, this is some subset of the total population
	{
		needed := reflect.ValueOf(basePopulation).Type()
		members, ok := p.membership[needed]
		if !ok {
			return nil, fmt.Errorf("unknown component: %s", bag.NiceTypeName(needed))
		}

		p.debug("parthenogenesis generation: %+v", members)

		population = entity.Population(set.FromMap(members))
		componentsToLoad = append(componentsToLoad, needed)
	}

	var behavior entity.LoadBehavior
	for _, q := range qs {
		needed := q.Component()

		nextGenParents, ok := p.membership[needed]
		if !ok {
			return nil, fmt.Errorf("unknown component: %s", bag.NiceTypeName(needed))
		}

		p.debug("next generation parent population: %+v", nextGenParents)

		population, behavior = q.Select(population, entity.Population(set.FromMap(nextGenParents)))

		if behavior == entity.ComponentLoad {
			componentsToLoad = append(componentsToLoad, needed)
		}

		if len(population) == 0 {
			p.debug("population went extinct")
			return nil, nil
		}
	}

	viewing := make([]entity.View, 0, len(population))

	for entity := range population {
		table := componentBuckets{}

		for _, compType := range componentsToLoad {
			table[compType] = p.membership[compType][entity]
		}

		viewing = append(viewing, &view{
			m:      entity,
			origin: p,
			loaded: table,
		})
	}

	return viewing, nil
}
