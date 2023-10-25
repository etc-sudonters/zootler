package entity

import (
	"errors"
	"reflect"
)

func isTryDerefErr(e error) bool {
	var is = errors.Is(e, ErrNotAssigned) || errors.Is(e, ErrNotLoaded)
	if is {
		return true
	}
	var unknown ErrUnknownComponent
	if errors.As(e, &unknown) {
		return true
	}

	return false
}

type ComponentGetter interface {
	GetComponent(Model, reflect.Type) (Component, error)
}

func PierceComponentType(c Component) reflect.Type {
	return reflect.TypeOf(c)
}

func AssignComponentTo(entity Model, target interface{}, retrieve ComponentGetter) error {
	if target == nil {
		return ErrNilComponentPtr
	}

	rv := reflect.ValueOf(target)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return ErrNonNilPtrOnly
	}

	targetType := rv.Elem().Type()

	acquired, err := retrieve.GetComponent(entity, targetType)

	if err != nil {
		if isTryDerefErr(err) && targetType.Kind() == reflect.Pointer {
			acquired, err = retrieve.GetComponent(entity, targetType.Elem())
		}

		if err != nil {
			return err
		}
	}

	if acquired == nil {
		return ErrNilComponent{
			Entity:    entity,
			Component: targetType,
		}
	}

	acquiredValue := reflect.ValueOf(acquired)

	if acquiredValue.Kind() != reflect.Pointer && targetType.Kind() == reflect.Pointer {
		intermediate := reflect.New(acquiredValue.Type())
		intermediate.Elem().Set(acquiredValue)
		acquiredValue = intermediate
	}

	rv.Elem().Set(acquiredValue)
	return nil
}
