package entity

import (
	"reflect"

	"github.com/etc-sudonters/substrate/skelly/hashset"
)

// responsible for the total administration of a population of models
type Pool interface {
	Queryable
	Manager
}

// responsible for creation and destruction of models
type Manager interface {
	Create() (View, error)
}

// responsible for looking either individual models or creating a subset of the
// population that matches the provided selectors
type Queryable interface {
	// return a subset of the population that matches the provided selectors
	Query(f Filter) ([]View, error)
	// load the specified components from the specified model, if a component
	// isn't attached to the model its pointer should be set to nil
	Get(m Model, components []interface{})
	// return the specific model from the pool
	Fetch(m Model) (View, error)
}

type Filter struct {
	include []reflect.Type
	exclude []reflect.Type
}

func (f Filter) With() []reflect.Type {
	return f.include
}

func (f Filter) Without() []reflect.Type {
	return f.exclude
}

type FilterBuilder struct {
	include hashset.Hash[reflect.Type]
	exclude hashset.Hash[reflect.Type]
}

func (f FilterBuilder) With(t reflect.Type) FilterBuilder {
	f.include.Add(t)
	return f
}

func (f FilterBuilder) Without(t reflect.Type) FilterBuilder {
	f.exclude.Add(t)
	return f
}

func (f FilterBuilder) Build() Filter {
	return Filter{
		include: hashset.AsSlice(f.include),
		exclude: hashset.AsSlice(f.exclude),
	}
}

func (f FilterBuilder) Clone() FilterBuilder {
	return FilterBuilder{
		include: hashset.FromMap(f.include),
		exclude: hashset.FromMap(f.exclude),
	}
}
