package entity

import (
	"errors"
	"fmt"
	"reflect"

	"sudonters/zootler/internal/bag"
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

func AssignComponentTo(target interface{}, retrieve func(reflect.Type) (Component, error)) error {
	if target == nil {
		return ErrNilComponentPtr
	}

	value := reflect.ValueOf(target)
	typ := value.Type()

	if typ.Kind() != reflect.Pointer || value.IsNil() {
		return ErrNonNilPtrOnly
	}

	targetType := typ.Elem()

	acquired, err := retrieve(targetType)

	if err != nil {
		if isTryDerefErr(err) && targetType.Kind() == reflect.Pointer {
			acquired, err = retrieve(targetType.Elem())
		}

		if err != nil {
			return err
		}
	}

	if acquired == nil {
		panic(
			fmt.Sprintf(
				"retrieved nil component for %s", bag.NiceTypeName(targetType),
			))
	}
	acquiredValue := reflect.ValueOf(acquired)

	if acquiredValue.Kind() != reflect.Pointer && targetType.Kind() == reflect.Pointer {
		intermediate := reflect.New(acquiredValue.Type())
		intermediate.Elem().Set(acquiredValue)
		acquiredValue = intermediate
	}

	value.Elem().Set(acquiredValue)
	return nil
}
