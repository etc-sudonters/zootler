package world

import (
	"sudonters/zootler/internal/rules"
	"sudonters/zootler/pkg/entity"
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
