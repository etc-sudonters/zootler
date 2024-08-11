package components

import (
	"reflect"
)

type Alias struct {
	For reflect.Type
	Qty int
}
