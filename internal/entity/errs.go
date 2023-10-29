package entity

import (
	"errors"
	"fmt"
	"reflect"
)

var ErrNoMoreIds = errors.New("no more ids available")
var ErrEntityNotExist = errors.New("entity does not exist")
var ErrNoEntities = errors.New("no entities")
var ErrNotLoaded = errors.New("not loaded")
var ErrNotAssigned = errors.New("not assigned")
var ErrNonNilPtrOnly = errors.New("non-nil pointers only")
var ErrNilComponentPtr = errors.New("nil pointer to component")

type ErrUnknownComponent struct {
	reflect.Type
}

func (u ErrUnknownComponent) Error() string {
	return fmt.Sprintf("unknown component: %s", u.Name())
}

type ErrNilComponent struct {
	Entity    Model
	Component reflect.Type
}

func (n ErrNilComponent) Error() string {
	return fmt.Sprintf("nil component: %s on %d", n.Component.Name(), n.Entity)
}
