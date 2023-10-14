package entity

import (
	"reflect"

	"github.com/etc-sudonters/zootler/internal/bag"
)

var _ Selector = With[interface{}]{}
var _ Selector = Without[interface{}]{}
var _ Selector = Optional[interface{}]{}
var _ Selector = DebugSelector{}

func UnknownComponent(t reflect.Type) unknownComponent {
	return unknownComponent(bag.NiceTypeName(t))
}

type unknownComponent string

func (u unknownComponent) Error() string {
	return "unknown component " + string(u)
}

type Population interface {
	Difference(Population) Population
	Intersect(Population) Population
	Union(Population) Population
}

// responsible for looking either individual models or creating a subset of the
// population that matches the provided selectors
type Queryable interface {
	// return a subset of the population that matches the provided selectors
	Query(Selector, ...Selector) ([]View, error)
	// load the specified components from the specified model, if a component
	// isn't attached to the model its pointer should be set to nil
	Get(Model, ...interface{})
}

type Selector interface {
	Component() reflect.Type
	Select(current, candidates Population) Population
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

// entities with this component are excluded
type Without[T includable] struct {
	componentFromGeneric[T]
}

func (e Without[T]) Select(current, next Population) Population {
	return current.Difference(next)
}

// filter entities to ones with this component, but do not load it
type With[T includable] struct {
	componentFromGeneric[T]
}

func (w With[T]) Select(current, next Population) Population {
	return current.Intersect(next)
}

type Optional[T includable] struct {
	componentFromGeneric[T]
}

func (o Optional[T]) Select(current, next Population) Population {
	return current.Union(next)
}

func (d DebugSelector) Component() reflect.Type {
	target := d.S.Component()
	d.F("selecting against %s\n", bag.NiceTypeName(target))
	return target
}

func (d DebugSelector) Select(currentGeneration Population, candidates Population) Population {
	d.F("current generation %+v\n", currentGeneration)
	d.F("candidates %+v\n", candidates)
	population := d.S.Select(currentGeneration, candidates)
	d.F("next generation %+v", population)
	return population
}
