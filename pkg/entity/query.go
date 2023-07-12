package entity

import (
	"errors"
	"reflect"

	"github.com/etc-sudonters/zootler/internal/bag"
	"github.com/etc-sudonters/zootler/internal/datastructures/set"
)

var ErrNotLoaded = errors.New("not loaded")

type Population set.Hash[Model]

type LoadBehavior int

const (
	ComponentLoad LoadBehavior = iota
	ComponentIgnore
)

type Queryable interface {
	All() set.Hash[Model]
	Query(Component, ...Selector) ([]View, error)
	Get(Model, ...interface{})
}

type Selector interface {
	// which component pool to select against
	Component() reflect.Type
	// creates a new population based on the current population
	// the secondary return value controls how/if a component is attached
	Select(current, candidates Population) (Population, LoadBehavior)
}

type Includeable interface {
	Component | *Component
}

type DebugSelector struct {
	F func(string, ...any)
	S Selector
}

type componentFromGeneric[T Includeable] struct{}

func (i componentFromGeneric[T]) Component() reflect.Type {
	var t T
	return reflect.TypeOf(t)
}

// filters entities with an arbitrary entity set
// when Op is called the current generation is passed as the first arg
type Arbitrary struct {
	Elems set.Hash[Model]
	Op    ArbitraryOp
}

type ArbitraryOp set.Operation[Model, Population, set.Hash[Model]]

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
type Load[T Includeable] struct {
	componentFromGeneric[T]
}

// entities with this component are excluded
type Without[T Includeable] struct {
	componentFromGeneric[T]
}

// filter entities to ones with this component, but do not load it
type With[T Includeable] struct {
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
