package hashpool

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/etc-sudonters/zootler/internal/bag"
	"github.com/etc-sudonters/zootler/pkg/entity"
)

func assignComponentTo(target interface{}, retrieve func(reflect.Type) (entity.Component, error)) error {
	if target == nil {
		return errors.New("nil component reference")
	}

	value := reflect.ValueOf(target)
	typ := value.Type()

	if typ.Kind() != reflect.Pointer || value.IsNil() {
		return errors.New("non-nil pointers only")
	}

	targetType := typ.Elem()

	acquired, err := retrieve(targetType)

	if err != nil {
		if errors.Is(err, entity.ErrNotLoaded) && targetType.Kind() == reflect.Pointer {
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
