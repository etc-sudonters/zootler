package hashpool

import (
	"reflect"

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

func (p *Pool) Query(basePopulation entity.Selector, qs ...entity.Selector) ([]entity.View, error) {
	if len(p.population) == 0 {
		p.debug("no population")
		return nil, nil
	}

	selectors := []entity.Selector{basePopulation}
	if len(qs) > 0 {
		selectors = append(selectors, qs...)
	}

	// initial population is the entire loaded pool
	var population entity.Population = entity.Population(p.All())
	componentsToLoad := make([]reflect.Type, 0, len(qs)+1)

	var behavior entity.LoadBehavior
	for _, q := range selectors {
		needed := q.Component()

		nextGenParents, ok := p.membership[needed]
		if !ok {
			continue
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
