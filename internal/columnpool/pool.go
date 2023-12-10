package columnpool

import (
	"reflect"
	"sudonters/zootler/internal/entity"
	"sudonters/zootler/internal/table"

	"github.com/etc-sudonters/substrate/mirrors"
	"github.com/etc-sudonters/substrate/skelly/bitset"
	"github.com/etc-sudonters/substrate/stageleft"
)

type record struct {
	id    entity.Model
	comps bitset.Bitset64
}

type ColumnEntityPool struct {
	tbl      table.Table
	typmap   mirrors.TypeMap
	entities []record

	defaultColumn table.ColumnFactory
}

// return a subset of the population that matches the provided filter
func (c *ColumnEntityPool) Query(f any) ([]entity.View, error) {
	panic("not implemented") // TODO: Implement
}

// load the specified components from the specified model, if a component
// isn't attached to the model its pointer should be set to nil
func (c *ColumnEntityPool) Get(m entity.Model, components []interface{}) {
	panic("not implemented") // TODO: Implement
}

func (c *ColumnEntityPool) Create() (entity.View, error) {
	panic("not implemented") // TODO: Implement
}

type columnEntity struct {
	id         table.RowId
	components bitset.Bitset64
	mapper     map[reflect.Type]int
	columns    []entity.Component

	p *ColumnEntityPool
}

func (cols *columnEntity) Model() entity.Model {
	return entity.Model(cols.id)
}

func (cols *columnEntity) Get(c interface{}) error {
	panic(stageleft.NotImplErr)
}

func (cols *columnEntity) Add(c entity.Component) error {
	panic("not implemented") // TODO: Implement
}

func (cols *columnEntity) Remove(c entity.Component) error {
	panic("not implemented") // TODO: Implement
}
