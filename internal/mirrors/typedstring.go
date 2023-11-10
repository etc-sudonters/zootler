package mirrors

import (
	"fmt"
	"reflect"
)

type TypedStrings struct {
	strings map[string]reflect.Type
}

func NewTypedStrings() TypedStrings {
	return TypedStrings{strings: make(map[string]reflect.Type)}
}

func (t TypedStrings) Typed(s string) any {
	if typ, ok := t.strings[s]; ok {
		return reflect.New(typ)
	}

	typ := reflect.StructOf([]reflect.StructField{
		{
			Name: "TypedString",
			Type: TypeOf[int](),
			Tag:  reflect.StructTag(fmt.Sprintf(`literal:"%s"`, s)),
		},
	})
	t.strings[s] = typ
	return reflect.New(typ).Interface()
}
