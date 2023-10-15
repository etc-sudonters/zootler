package entity

import (
	"fmt"
	"reflect"
)

var _ Selector = With[interface{}]{}
var _ Selector = Without[interface{}]{}
var _ Selector = Optional[interface{}]{}
var _ Selector = DebugSelector{}

type LoadBehavior uint

const (
	_ LoadBehavior = iota
	ComponentInclude
	ComponentExclude
	ComponentOptional
)

func (l LoadBehavior) String() string {
	switch uint(l) {
	case 1:
		return "Include"
	case 2:
		return "Exclude"
	case 3:
		return "Optional"
	default:
		panic(fmt.Errorf("unknown load behavior %d", l))
	}
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
	Behavior() LoadBehavior
}

type includable interface {
	Component | *Component
}

// something funky happening?
type DebugSelector struct {
	F func(string, ...any)
	Selector
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

func (w Without[T]) Behavior() LoadBehavior {
	return ComponentExclude
}

// filter entities to ones with this component, but do not load it
type With[T includable] struct {
	componentFromGeneric[T]
}

func (w With[T]) Behavior() LoadBehavior {
	return ComponentInclude
}

type Optional[T includable] struct {
	componentFromGeneric[T]
}

func (o Optional[T]) Behavior() LoadBehavior {
	return ComponentOptional
}

func (d DebugSelector) Component() reflect.Type {
	target := d.Selector.Component()
	behavior := d.Selector.Behavior()
	d.F("%s %s", target, behavior)
	return target
}
