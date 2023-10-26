package logic

import (
	"sudonters/zootler/internal/entity"
	"sudonters/zootler/pkg/rulesparser"
)

type (
	ParsedRule struct {
		R rulesparser.Expression
	}
	Edge struct {
		Origination entity.Model
		Destination entity.Model
	}
	Name      string
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
