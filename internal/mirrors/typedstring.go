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

func (t TypedStrings) Typed(s string) reflect.Type {
	if typ, ok := t.strings[s]; ok {
		return typ
	}

	typ := reflect.StructOf([]reflect.StructField{
		{
			Name: "TypedString",
			Type: TypeOf[int](),
			Tag:  reflect.StructTag(fmt.Sprintf(`literal:"%s"`, s)),
		},
	})
	t.strings[s] = typ
	return typ
}

func (t TypedStrings) InstanceOf(s string) any {
	return reflect.New(t.Typed(s)).Interface()
}

func TryGetLiteral(t reflect.Type) (string, bool) {
	field, ok := t.FieldByName("TypedString")
	if !ok {
		return "", false
	}

	literal := field.Tag.Get("literal")
	if literal == "" {
		return "", false
	}
	return literal, true
}
