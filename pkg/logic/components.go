package logic

import (
	"sudonters/zootler/internal/entity"
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
	Location  struct{}
	Token     struct{}
	Inhabited entity.Model
	Inhabits  entity.Model
)
