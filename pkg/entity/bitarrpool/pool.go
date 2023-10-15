package bitarrpool

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"sudonters/zootler/internal/bag"
	"sudonters/zootler/internal/bitset"
	"sudonters/zootler/pkg/entity"
)

func unknown(t reflect.Type, b *bitarrpool) error {
	return fmt.Errorf("unknown component: %q\navailable:\n%s", bag.NiceTypeName(t), bpSummary{*b})
}

type componentId int64

const (
	INVALID_COMPONENT componentId  = 0
	INVALID_ENTITY    entity.Model = 0
)

type componentTable []componentRow   // idx'd by compId
type componentRow []entity.Component // idx'd by entId

type bitarrpool struct {
	componentBucketCount int64

	entities  []*bitarrview
	maxCompId componentId
	compIds   map[reflect.Type]componentId
	table     componentTable
}

var ErrNoMoreIds = errors.New("no more ids available")
var ErrNoEntities = errors.New("no entities")
var _ entity.Pool = (*bitarrpool)(nil)
var _ entity.View = (*bitarrview)(nil)

func New(maxComponentId int64) *bitarrpool {
	return &bitarrpool{
		componentBucketCount: bitset.Buckets(maxComponentId),

		/* entites and table are initialized with 1 existing member so we can
		* skip needing to figure out if we accidentally have an id of 0. We
		* can't, and that's a nil pointer
		 */

		entities: make([]*bitarrview, 1, 128),

		maxCompId: componentId(maxComponentId),
		compIds:   make(map[reflect.Type]componentId, 128),
		table:     make(componentTable, 1, 128),
	}
}

func (b *bitarrpool) addToTable(newComponent componentId) {
	capacity := bag.Max(len(b.entities)*2, 128)
	b.table = append(b.table, make(componentRow, len(b.entities), capacity))
}

func (b *bitarrpool) ensureRowSize(c componentId) {
	currentRow := b.table[c]
	if len(currentRow) < len(b.entities) {
		newRow := make([]entity.Component, len(b.entities))
		copy(newRow, currentRow)
		b.table[c] = newRow
	}
}

func (b *bitarrpool) componentId(t reflect.Type) componentId {
	var id componentId
	if id, ok := b.compIds[t]; ok {
		return id
	}

	id = componentId(len(b.compIds) + 1)
	b.compIds[t] = id
	b.addToTable(id)
	return id
}

func (b *bitarrpool) Create() (entity.View, error) {
	comps := bitset.New(b.componentBucketCount)

	view := &bitarrview{
		id:    entity.Model(len(b.entities) + 1),
		comps: &comps,
		p:     b,
	}

	b.entities = append(b.entities, view)
	return view, nil
}

// return a subset of the population that matches the provided selectors
func (b *bitarrpool) Query(q entity.Selector, qs ...entity.Selector) ([]entity.View, error) {
	componentId := b.componentId(q.Component())
	filter := entityFilter{id: componentId, behavior: q.Behavior()}
	entities := bag.Filter(b.entities, filter.Test)

	for _, q = range qs {
		componentId = b.componentId(q.Component())
		filter = entityFilter{id: componentId, behavior: q.Behavior()}
		entities = bag.Filter(b.entities, filter.Test)
		if len(entities) == 0 {
			return nil, ErrNoEntities
		}
	}

	view := make([]entity.View, len(entities))
	for idx, ent := range entities {
		view[idx] = ent
	}
	return view, nil
}

func (b *bitarrpool) Get(m entity.Model, cs ...interface{}) {
	for i := range cs {
		_ = entity.AssignComponentTo(cs[i], tryFindCompFor(b, m))
	}
}

type bpSummary struct {
	bitarrpool
}

func (b bpSummary) String() string {
	return fmt.Sprintf("bitpool{ k: %d, maxId: %d, nextId: %d, all: %s, tests: %s }", b.entityBucketCount, b.maxEntId, b.nextEntId, b.entities, b.summarizeTable())
}

func (b bpSummary) summarizeTable() string {
	repr := &strings.Builder{}
	repr.WriteRune('{')
	repr.WriteRune(' ')
	for typ, _ := range b.tests {
		fmt.Fprintf(repr, "%q", typ)
		repr.WriteRune(' ')
	}
	repr.WriteRune('}')
	return repr.String()
}

type bitarrview struct {
	id    entity.Model
	comps *bitset.Bitset64
	p     *bitarrpool
}

func (b bitarrview) String() string {
	return fmt.Sprintf("bitview{ %d: %s }", b.id, b.comps)
}

func (b *bitarrview) Model() entity.Model {
	return b.id
}

func (b *bitarrview) Get(w interface{}) error {
	return entity.AssignComponentTo(w, tryFindCompFor(b.p, b.id))
}

func (b *bitarrview) addMany(cs ...entity.Component) error {
	for _, c := range cs {
		if err := b.Add(c); err != nil {
			return err
		}
	}
	return nil
}

func (b *bitarrview) Add(c entity.Component) error {
	if many, ok := c.([]entity.Component); ok {
		return b.addMany(many...)
	}

	return b.p.addCompToEnt(b, c)
}

func (b *bitarrview) Remove(c entity.Component) error {
	if many, ok := c.([]entity.Component); ok {
		return b.removeMany(many...)
	}

	return b.p.removeCompFromEnt(b, c)
}

type entityFilter struct {
	id       componentId
	behavior entity.LoadBehavior
}

func (e entityFilter) Test(ent *bitarrview) bool {
	switch e.behavior {
	case entity.ComponentOptional:
		return true
	case entity.ComponentInclude:
		return ent.comps.Test(int64(e.id))
	case entity.ComponentExclude:
		return !ent.comps.Test(int64(e.id))
	}
	return false
}
