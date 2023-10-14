package entity

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/etc-sudonters/zootler/internal/bag"
)

var ErrNotLoaded = errors.New("not loaded")
var ErrNotAssigned = errors.New("not assigned")

// a member of a pool's population
type Model uint

func (m Model) String() string {
	return fmt.Sprintf("Model{%d}", m)
}

// arbitrary attachments to a Model
type Component interface{}

func ComponentName(c Component) string {
	if c == nil {
		return "nil"
	}

	return bag.NiceTypeName(reflect.TypeOf(c))
}

// a mutable reference to a pool population member
// a view may be created with some components loaded from the pool
type View interface {
	// which population member we are
	Model() Model
	/*
		Attempts to retrieve a loaded component. A view may choose to load
		components on demand but otherwise should ensure the pointer passed is
		set to nil and return an ErrNotLoaded if the component is simply not
		found in the loaded components.

		The argument passed must be a pointer to a value of the sought actual
		component type. If the sought after component is registered as
		`*MyComponent` then this method must be passed a `**MyComponent`.

			var e View := ...
			var s Song
			var l *Location
			if err := e.Get(&s); err != nil {
				fmt.Println(s.String())
			}
			if err := e.Get(&l); err != nil
				fmt.Println(s.String())
			}


		Ideally we could define this as `Get[T any](*T) error` instead
	*/
	Get(interface{}) error
	// attaches a component to this model, this component is retrievable via Get
	// the component is retrievable via the type housed behind the Component
	Add(Component) error
	// removes a component from this model even if it is unloaded
	// the value passed is not considered, only the type housed behind the
	// Component interface
	Remove(Component) error
}

// responsible for the total administration of a population of models
type Pool interface {
	Queryable
	Manager
}

// responsible for creation and destruction of models
type Manager interface {
	Create() (View, error)
}
