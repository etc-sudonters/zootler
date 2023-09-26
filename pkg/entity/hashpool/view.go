package hashpool

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/etc-sudonters/zootler/internal/bag"
	"github.com/etc-sudonters/zootler/pkg/entity"
)

type view struct {
	m       entity.Model
	origin  *Pool
	loaded  map[reflect.Type]entity.Component
	session map[reflect.Type]entity.Component
}

func (v view) String() string {
	return fmt.Sprintf("View{%s}", v.m)
}

func (v view) checkDetached() {
	if v.origin == nil {
		panic("detatched view")
	}
}

func (v view) Model() entity.Model {
	v.checkDetached()
	return v.m
}

func (v view) Get(target interface{}) error {
	v.checkDetached()

	tryFind := func(t reflect.Type) (entity.Component, error) {
		v.origin.debug("target load type %s", bag.NiceTypeName(t))
		acquired, ok := v.loaded[t]
		if !ok {
			v.origin.debug("attempting to load %s from session", bag.NiceTypeName(t))
			acquired, ok = v.session[t]
			if !ok {
				return nil, entity.ErrNotLoaded
			}
		}

		return acquired, nil
	}

	err := assignComponentTo(target, tryFind)
	if err != nil {
		return err
	}
	return nil
}

func (v view) Add(target entity.Component) error {
	v.checkDetached()

	typ := reflect.TypeOf(target)
	if _, ok := v.loaded[typ]; ok {
		v.loaded[typ] = target
	} else {
		ensureTable(v.origin, target)
		if v.session == nil {
			v.session = make(map[reflect.Type]entity.Component)
		}
		v.session[typ] = target
	}

	v.origin.membership[typ][v.m] = target

	return nil
}

func (v view) Remove(target entity.Component) error {
	v.checkDetached()

	typ := reflect.TypeOf(target)

	if typ == entity.ModelComponentType {
		return errors.New("cannot remove model component")
	}

	delete(v.loaded, typ)
	delete(v.session, typ)
	removeFromTable(v.m, typ, v.origin)
	return nil
}
