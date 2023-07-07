package entity

import "github.com/etc-sudonters/rando/set"

type Pool struct {
	entities   map[Id]Tags
	membership membership
	nextId     Id
}

func PoolWithCapacity(capacity int) *Pool {
	return &Pool{
		entities:   make(map[Id]Tags, capacity),
		membership: make(membership, capacity),
		nextId:     0,
	}
}

func (p Pool) Add(name TagName) *View {
	thisId := p.nextId
	tags := make(Tags)
	tags["id"] = thisId
	p.entities[thisId] = tags
	p.nextId += 1
	return &View{
		Id:   thisId,
		Tags: tags,
	}
}

func (p *Pool) Index() {
	m := make(membership, len(p.membership))

	for id, tags := range p.entities {
		for tag := range tags {
			members, exists := m[tag]
			if !exists {
				members = set.New[Id]()
			}
			members.Add(id)
			m[tag] = members
		}
	}

	p.membership = m
}

func (p Pool) Query(qs ...Query) []*View {
	matched := p.membership["id"]

	for _, q := range qs {
		tag, operation := q()
		if tag == "id" {
			continue
		}

		members, exists := p.membership[tag]
		if !exists || len(members) == 0 {
			return nil
		}

		matched = operation(matched, members)

		if len(matched) == 0 {
			return nil
		}
	}

	view := make([]*View, 0, len(matched))

	for id := range matched {
		view = append(view, &View{
			Id:   id,
			Tags: p.entities[id],
		})
	}

	return view
}
