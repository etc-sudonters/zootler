package entity

import (
	"fmt"
	"reflect"
)

var _ Selector = With[interface{}]{}
var _ Selector = Without[interface{}]{}
var _ Selector = DebugSelector{}

type LoadBehavior uint

const (
	_ LoadBehavior = iota
	ComponentInclude
	ComponentExclude
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

type Selector interface {
	Component() reflect.Type
	Behavior() LoadBehavior
}

// something funky happening?
type DebugSelector struct {
	F func(string, ...any)
	Selector
}

type componentFromGeneric[T Component] struct{}

func (i componentFromGeneric[T]) Component() reflect.Type {
	var t T
	return reflect.TypeOf(t)
}

// entities with this component are excluded
type Without[T Component] struct {
	componentFromGeneric[T]
}

func (w Without[T]) Behavior() LoadBehavior {
	return ComponentExclude
}

// filter entities to ones with this component, but do not load it
type With[T Component] struct {
	componentFromGeneric[T]
}

func (w With[T]) Behavior() LoadBehavior {
	return ComponentInclude
}

func (d DebugSelector) Component() reflect.Type {
	target := d.Selector.Component()
	behavior := d.Selector.Behavior()
	d.F("%s %s", target, behavior)
	return target
}
