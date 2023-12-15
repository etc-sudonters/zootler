package entity

import (
	"reflect"
)

type ComponentId uint64

const INVALID_COMPONENT ComponentId = 0

// arbitrary attachments to a Model
type Component interface{}
type ComponentType reflect.Type

func ComponentName(c Component) string {
	if c == nil {
		return "nil"
	}

	return PierceComponentType(c).Name()
}
