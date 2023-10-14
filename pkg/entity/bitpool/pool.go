package bitpool

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/etc-sudonters/zootler/internal/bitset"
	"github.com/etc-sudonters/zootler/pkg/entity"
)

var ErrNoMoreIds = errors.New("no more ids available")
var ErrNoEntities = errors.New("no entities")
var _ entity.Pool = (*bitpool)(nil)
var _ entity.View = (*bitview)(nil)
var _ entity.Population = bitpop{}

func New(maxId int64) *bitpool {
	buckets := bitset.Buckets(maxId)
	b := bitset.New(buckets)
	return &bitpool{
		k:      buckets,
		maxId:  entity.Model(maxId),
		nextId: 0,
		all:    &b,
		tests:  make(map[reflect.Type]*bitset.Bitset64),
		table:  make(map[reflect.Type]map[entity.Model]entity.Component),
	}
}

// return a subset of the population that matches the provided selectors
func (b *bitpool) Query(q entity.Selector, qs ...entity.Selector) ([]entity.View, error) {
	if b.nextId == 0 {
		return nil, ErrNoEntities
	}

	test := q.Component()
	members, ok := b.tests[test]
	if !ok {
		return nil, entity.UnknownComponent(test)
	}

	population := q.Select(bitpop{*b.all}, bitpop{*members})

	for _, q = range qs {
		test = q.Component()
		members, ok = b.tests[test]
		if !ok {
			return nil, entity.UnknownComponent(test)
		}
		population = q.Select(population, bitpop{*members})

		if bitset.Empty((population.(bitpop)).b) {
			return nil, ErrNoEntities
		}
	}

	var entities []entity.View
	p := (population.(bitpop)).b

	for i := int64(0); i < int64(b.nextId); i++ {
		if p.Test(i) {
			e := bitview{entity.Model(i), b}
			entities = append(entities, &e)
		}
	}

	return entities, nil
}

func (b *bitpool) Get(m entity.Model, cs ...interface{}) {
	for i := range cs {
		entity.AssignComponentTo(cs[i], tryFindCompFor(b, m))
	}
}

func (b *bitpool) Create() (entity.View, error) {
	if b.maxId == b.nextId {
		return nil, ErrNoMoreIds
	}
	v := bitview{id: b.nextId, p: b}
	b.nextId++
	b.all.Set(int64(v.id))
	v.Add(v.id)
	return &v, nil
}

type bitpool struct {
	k      int64
	maxId  entity.Model
	nextId entity.Model
	all    *bitset.Bitset64
	tests  map[reflect.Type]*bitset.Bitset64
	table  map[reflect.Type]map[entity.Model]entity.Component
}

type bitview struct {
	id entity.Model
	p  *bitpool
}

func (b bitview) String() string {
	return fmt.Sprintf("bitview(%d)", b.id)
}

func (b *bitview) Model() entity.Model {
	return entity.Model(b.id)
}

func (b *bitview) Get(w interface{}) error {
	return entity.AssignComponentTo(w, tryFindCompFor(b.p, b.id))
}

func (b *bitview) Add(c entity.Component) error {
	typ := reflect.TypeOf(c)
	ensureTables(b.p, typ)
	b.p.tests[typ].Set(int64(b.id))
	b.p.table[typ][b.id] = c
	return nil
}

func (b *bitview) Remove(c entity.Component) error {
	typ := reflect.TypeOf(c)
	members, ok := b.p.tests[typ]
	if !ok {
		return nil
	}

	members.Clear(int64(b.id))
	delete(b.p.table[typ], b.id)
	return nil
}

func ensureTables(b *bitpool, t reflect.Type) {
	if _, ok := b.tests[t]; !ok {
		set := bitset.New(b.k)
		b.tests[t] = &set
		b.table[t] = make(map[entity.Model]entity.Component)
	}
}

func existsInCompTable(b *bitpool, t reflect.Type, m entity.Model) bool {
	tab, ok := b.tests[t]
	if !ok {
		return false
	}
	return tab.Test(int64(m))
}

func tryFindCompFor(p *bitpool, m entity.Model) func(reflect.Type) (entity.Component, error) {
	return func(t reflect.Type) (entity.Component, error) {
		if !existsInCompTable(p, t, m) {
			return nil, entity.ErrNotAssigned
		}
		return p.table[t][m], nil
	}
}

type bitpop struct {
	b bitset.Bitset64
}

func (b bitpop) Difference(o entity.Population) entity.Population {
	return bitpop{b.b.Difference((o.(bitpop).b))}
}

func (b bitpop) Intersect(e entity.Population) entity.Population {
	return bitpop{b.b.Intersect((e.(bitpop).b))}
}

func (b bitpop) Union(e entity.Population) entity.Population {
	return bitpop{b.b.Union((e.(bitpop)).b)}
}
