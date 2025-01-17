package z2

import (
	"sudonters/zootler/internal/query"
	"sudonters/zootler/internal/table"
)

type AttachingComponents struct {
	v table.Values
}

func (this *AttachingComponents) Add(v table.Value) {
	this.v = append(this.v, v)
}

func (this *AttachingComponents) Components() table.Values {
	return this.v
}

type proxy struct {
	q  query.Engine
	id Entity
}

func (this *proxy) Attach(v ...table.Value) error {
	return this.AttachAll(table.Values(v))
}

func (this *proxy) AttachAll(vs table.Values) error {
	return this.q.SetValues(table.RowId(this.id), vs)
}

func (this *proxy) Entity() Entity {
	return this.id
}

func Tracked[T TrackableValue](eng query.Engine) TrackedEntities[T] {
	return TrackedEntities[T]{
		q:        eng,
		entities: make(map[T]Entity),
	}
}

func IntoNamed(t TrackedEntities[Name]) NamedEntities {
	return NamedEntities{t}
}

func Named(eng query.Engine) NamedEntities {
	return IntoNamed(Tracked[Name](eng))
}

type NamedEntities struct {
	TrackedEntities[Name]
}

type TrackableValue interface {
	table.Value
	comparable
}

type TrackedEntities[T TrackableValue] struct {
	q        query.Engine
	entities map[T]Entity
}

func (this TrackedEntities[T]) Entity(key T) proxy {
	if id, exists := this.entities[key]; exists {
		return proxy{this.q, id}
	}

	row, err := this.q.InsertRow(key)
	if err != nil {
		panic(err)
	}
	this.entities[key] = Entity(row)
	return proxy{this.q, Entity(row)}
}
