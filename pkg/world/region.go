package world

import (
	"github.com/etc-sudonters/zootler/internal/rules"
	"github.com/etc-sudonters/zootler/pkg/entity"
)

type Collection struct {
	Name       string
	Rule       rules.RawRule
	Components []entity.Component
	Vanilla    string
}

type Transit struct {
	Exit string
	Rule rules.RawRule
}

type Region struct {
	Name        string
	Collections []Collection
	Transits    []Transit
	Components  []entity.Component
}
