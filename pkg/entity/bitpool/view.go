package bitpool

import (
	"fmt"
	"sudonters/zootler/pkg/entity"

	"github.com/etc-sudonters/substrate/skelly/set/bits"
)

type bitview struct {
	id    entity.Model
	comps bits.Bitset64
	p     *bitpool
}

func (b bitview) String() string {
	return fmt.Sprintf("bitview{ %d: %s }", b.id, b.comps)
}

func (b bitview) Model() entity.Model {
	return b.id
}

func (b bitview) Get(w interface{}) error {
	return entity.AssignComponentTo(b.id, w, b.p.table.Getter())
}

func (b bitview) addMany(cs ...entity.Component) error {
	for _, c := range cs {
		if err := b.Add(c); err != nil {
			return err
		}
	}
	return nil
}

func (b bitview) Add(c entity.Component) error {
	if many, ok := c.([]entity.Component); ok {
		return b.addMany(many...)
	}

	return b.p.addCompToEnt(b, c)
}

func (b bitview) removeMany(cs ...entity.Component) error {
	for _, c := range cs {
		if err := b.Remove(c); err != nil {
			return err
		}
	}
	return nil
}

func (b bitview) Remove(c entity.Component) error {
	if many, ok := c.([]entity.Component); ok {
		return b.removeMany(many...)
	}

	return b.p.removeCompFromEnt(b, c)
}
