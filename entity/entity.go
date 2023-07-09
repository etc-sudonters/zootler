package entity

import (
	"fmt"
	"reflect"

	"github.com/etc-sudonters/rando/set"
)

// used to fill in the blanks when a component is optional
var ComponentType reflect.Type
var OptionalComponent Component
var OptionalComponentType reflect.Type
var ModelComponentType reflect.Type

func init() {
	optional := optionalComponent{}
	OptionalComponent = Component(optional)
	OptionalComponentType = reflect.TypeOf(optional)
	ComponentType = reflect.TypeOf([]Component{}).Elem()
	ModelComponentType = reflect.TypeOf(Model(0))
}

type optionalComponent struct{}

type ModelName string
type Model uint

type Component interface{}

func ComponentName(c Component) string {
	if c == nil {
		return "nil"
	}

	return fmt.Sprintf("%s", reflect.TypeOf(c).Name())
}

type View interface {
	Model() Model
	/*
		var e View := ...
		var s Song
		if err := e.Get(&s); err != {
			fmt.Println(s.String())
		}
	*/
	Get(interface{}) error
	Add(Component) error
	Remove(Component) error
}

type Pool interface {
	Queryable
	Manager
}

type Queryable interface {
	Query(Component, ...Selector) ([]View, error)
}

type Manager interface {
	All() set.Hash[Model]
	Create() (View, error)
	Delete(View) error
}

type Selector interface {
	Component() reflect.Type
	Select(set.Hash[Model], map[Model]Component) map[Model]Component
}

type Includeable interface {
	Component | *Component
}

type DebugSelector struct {
	Debug func(string, ...any)
	Selector
}

func (d DebugSelector) Component() reflect.Type {
	target := d.Selector.Component()
	d.Debug("selecting against %s\n", NiceTypeName(target))
	return target
}

func (d DebugSelector) Select(currentGeneration set.Hash[Model], nextGenerationParents map[Model]Component) map[Model]Component {
	d.Debug("current generation %+v\n", currentGeneration)
	d.Debug("next generation parents %+v\n", set.FromMap(nextGenerationParents))
	nextGeneration := d.Selector.Select(currentGeneration, nextGenerationParents)
	d.Debug("next generation %+v", nextGeneration)
	return nextGeneration
}

type Include[T Includeable] struct {
	t T
}

func (i Include[T]) Component() reflect.Type {
	return reflect.TypeOf(i.t)
}

func (i Include[T]) Select(
	currentGeneration set.Hash[Model], nextGenerationParents map[Model]Component,
) map[Model]Component {
	maximumNextGenerationPopulation := len(currentGeneration)
	if len(currentGeneration) > len(nextGenerationParents) {
		maximumNextGenerationPopulation = len(nextGenerationParents)
	}

	nextGeneration := make(map[Model]Component, maximumNextGenerationPopulation)

	for parent, component := range nextGenerationParents {
		if currentGeneration.Exists(parent) {
			nextGeneration[parent] = component
		}
	}

	return nextGeneration
}

func NiceTypeName(t reflect.Type) string {
	if t.Kind() != reflect.Pointer {
		return t.Name()
	}

	t = t.Elem()
	return fmt.Sprintf("&%s", t.Name())
}
