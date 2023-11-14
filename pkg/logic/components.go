package logic

import (
	"reflect"
	"sudonters/zootler/internal/entity"
	"sudonters/zootler/internal/mirrors"
	"sudonters/zootler/pkg/rules/ast"
)

type (
	ParsedRule struct {
		R ast.Expression
	}
	Edge struct {
		Origination entity.Model
		Destination entity.Model
	}
	Collected struct{}
	Trick     struct{}
	Enabled   struct{}
	Spawn     struct{}
	RawRule   string
	Inhabited entity.Model
	Inhabits  entity.Model
)

type TypedStringSelector struct {
	TypedStrs mirrors.TypedStrings
}

func (t TypedStringSelector) With(s string) entity.Selector {
	return selector{
		literal:  s,
		typ:      t.TypedStrs.Typed(s),
		behavior: entity.ComponentInclude,
	}
}

func (t TypedStringSelector) Without(s string) entity.Selector {
	return selector{
		literal:  s,
		typ:      t.TypedStrs.Typed(s),
		behavior: entity.ComponentExclude,
	}

}

type selector struct {
	literal  string
	typ      reflect.Type
	behavior entity.LoadBehavior
}

func (s selector) Component() reflect.Type {
	return s.typ
}

func (s selector) Behavior() entity.LoadBehavior {
	return s.behavior
}
