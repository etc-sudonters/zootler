package entity

import (
	"fmt"
	"reflect"
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

func (m Model) String() string {
	return fmt.Sprintf("Model{%d}", m)
}

type Component interface{}

func ComponentName(c Component) string {
	if c == nil {
		return "nil"
	}

	return reflect.TypeOf(c).Name()
}

type View interface {
	Model() Model
	/*
		var e View := ...
		var s Song
		if err := e.Get(&s); err != {
			fmt.Println(s.String())
		}
		_ALWAYS_ pass reference to what we're assiging
	*/
	Get(interface{}) error
	Add(Component) error
	Remove(Component) error
}

type Pool interface {
	Queryable
	Manager
}

type Manager interface {
	Create() (View, error)
	Delete(View) error
}
