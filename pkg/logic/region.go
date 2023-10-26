package logic

import (
	"sudonters/zootler/internal/entity"
)

type Collection struct {
	Name       string
	Rule       RawRule
	Components []entity.Component
	Vanilla    string
}

type Transit struct {
	Exit string
	Rule RawRule
}

type Region struct {
	Name        string
	Collections []Collection
	Transits    []Transit
	Components  []entity.Component
}
