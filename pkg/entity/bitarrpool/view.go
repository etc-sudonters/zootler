package bitarrpool

import (
	"fmt"
	"sudonters/zootler/internal/bitset"
	"sudonters/zootler/pkg/entity"
)

type bitarrview struct {
	id    entity.Model
	comps bitset.Bitset64
	p     *bitarrpool
}

func (v bitarrview) mask(m bitset.Bitset64) bitset.Bitset64 {
	return m.Intersect(v.comps)
}

func (v bitarrview) isFullyMasked(mask bitset.Bitset64) bool {
	return v.mask(mask).Eq(mask)
}

func (b bitarrview) String() string {
	return fmt.Sprintf("bitview{ %d: %s }", b.id, b.comps)
}

func (b bitarrview) Model() entity.Model {
	return b.id
}

func (b bitarrview) Get(w interface{}) error {
	return entity.AssignComponentTo(w, getComponenter(b))
}

func (b bitarrview) addMany(cs ...entity.Component) error {
	for _, c := range cs {
		if err := b.Add(c); err != nil {
			return err
		}
	}
	return nil
}

func (b bitarrview) Add(c entity.Component) error {
	if many, ok := c.([]entity.Component); ok {
		return b.addMany(many...)
	}

	return b.p.addCompToEnt(b, c)
}

func (b bitarrview) removeMany(cs ...entity.Component) error {
	for _, c := range cs {
		if err := b.Remove(c); err != nil {
			return err
		}
	}
	return nil
}

func (b bitarrview) Remove(c entity.Component) error {
	if many, ok := c.([]entity.Component); ok {
		return b.removeMany(many...)
	}

	return b.p.removeCompFromEnt(b, c)
}
