package entity

import (
	"reflect"

	"github.com/etc-sudonters/zootler/internal/bag"
	"github.com/etc-sudonters/zootler/internal/set"
)

// the population, or a subset thereof, from a pool
type Population set.Hash[Model]

// determines component loading behavior for a selector
type LoadBehavior int

const (
	// the associated component should be loaded
	ComponentLoad LoadBehavior = iota
	// the associated component should not be loaded
	ComponentIgnore
)

// responsible for looking either individual models or creating a subset of the
// population that matches the provided selectors
type Queryable interface {
	// return every model currently loaded
	All() set.Hash[Model]
	// return a subset of the population that matches the provided selectors
	Query(Selector, ...Selector) ([]View, error)
	// load the specified components from the specified model, if a component
	// isn't attached to the model its pointer should be set to nil
	Get(Model, ...interface{})
}

type Selector interface {
	// which component pool to select against
	Component() reflect.Type
	// creates a new population based on the current population
	// the secondary return value controls how/if a component is attached
	Select(current, candidates Population) (Population, LoadBehavior)
}

type includable interface {
	Component | *Component
}

// something funky happening?
type DebugSelector struct {
	F func(string, ...any)
	S Selector
}

type componentFromGeneric[T includable] struct{}

func (i componentFromGeneric[T]) Component() reflect.Type {
	var t T
	return reflect.TypeOf(t)
}

// filters entities with an arbitrary entity set
// when Op is called the current generation is passed as the first arg
type Arbitrary struct {
	Elems Population
	Op    ArbitraryOp
}

type ArbitraryOp set.Operation[Model, Population, Population]

func (a Arbitrary) Component() reflect.Type {
	return ModelComponentType
}

func (a Arbitrary) Select(
	currentGeneration Population,
	_ Population,
) (Population, LoadBehavior) {
	return Population(a.Op(currentGeneration, a.Elems)), ComponentIgnore
}

// makes the specified component available, entities w/o this component are excluded
type Load[T includable] struct {
	componentFromGeneric[T]
}

// entities with this component are excluded
type Without[T includable] struct {
	componentFromGeneric[T]
}

// filter entities to ones with this component, but do not load it
type With[T includable] struct {
	componentFromGeneric[T]
}

func (l Load[T]) Select(
	currentGeneration Population,
	candidates Population,
) (Population, LoadBehavior) {
	return Population(
			set.Intersection(candidates, currentGeneration),
		),
		ComponentLoad
}

func (e Without[T]) Select(
	currentGeneration Population,
	candidates Population,
) (Population, LoadBehavior) {
	return Population(
			set.Difference(candidates, currentGeneration),
		),
		ComponentIgnore
}

func (w With[T]) Select(
	currentGeneration Population,
	candidates Population,
) (Population, LoadBehavior) {
	return Population(
			set.Intersection(candidates, currentGeneration),
		),
		ComponentIgnore
}

func (d DebugSelector) Component() reflect.Type {
	target := d.S.Component()
	d.F("selecting against %s\n", bag.NiceTypeName(target))
	return target
}

func (d DebugSelector) Select(
	currentGeneration Population,
	candidates Population,
) (Population, LoadBehavior) {
	d.F("current generation %+v\n", currentGeneration)
	d.F("candidates %+v\n", candidates)
	population, behavior := d.S.Select(currentGeneration, candidates)
	d.F("next generation %+v", population)
	d.F("behavior %d", behavior)
	return population, behavior
}
