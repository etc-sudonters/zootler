package logic

import (
	"github.com/etc-sudonters/zootler/internal/rules"
	"github.com/etc-sudonters/zootler/pkg/entity"
)

type (
	Edge struct {
		Origination entity.Model
		Destination entity.Model
	}
	Name       string
	Collected  struct{}
	Trick      struct{}
	Enabled    struct{}
	Spawn      struct{}
	RawRule    string
	Location   struct{}
	Token      struct{}
	Inhabited  entity.Model
	Inhabits   entity.Model
	ParsedRule rules.Expression
)
