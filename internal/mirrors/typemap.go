package mirrors

import (
	"errors"
	"reflect"
)

var ErrExistsAlready = errors.New("present already")
var ErrNoId = errors.New("no id present")

type TypeId uint64

type TypeMap map[reflect.Type]TypeId

func (m TypeMap) Add(t reflect.Type) TypeId {
	if id, ok := m[t]; ok {
		return id
	}
	id := TypeId(len(m))
	m[t] = id
	return id
}

func (m TypeMap) IdOf(t reflect.Type) (TypeId, error) {
	id, ok := m[t]
	if !ok {
		return id, ErrNoId
	}
	return id, nil
}

func IdOf[T any](m TypeMap) (TypeId, error) {
	return m.IdOf(TypeOf[T]())
}

func Add[T any](m TypeMap) TypeId {
	return m.Add(TypeOf[T]())
}
