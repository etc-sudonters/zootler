package entity

import (
	"github.com/etc-sudonters/rando/set"
)

type Archetype interface {
	Apply(t Tags)
}

type ArchetypeTag TagName

func (a ArchetypeTag) Apply(t Tags) {
	t.Apply(TagName(a), nil)
}

type ArchetypeFunc func(t Tags)

func (a ArchetypeFunc) Apply(t Tags) {
	a(t)
}

type Tags map[TagName]interface{}

func (t Tags) Use(a Archetype) Tags {
	a.Apply(t)
	return t
}

func (t Tags) ApplyAll(o Tags) Tags {
	for k, v := range o {
		t[k] = v
	}

	return t
}

func (t Tags) Apply(n TagName, i interface{}) Tags {
	t[n] = i
	return t
}

func (t Tags) RemoveAll(ns ...TagName) Tags {
	for _, n := range ns {
		delete(t, n)
	}
	return t
}

func (t Tags) Remove(n TagName) Tags {
	delete(t, n)
	return t
}

type TagName string

type membership map[TagName]set.Hash[Id]

type Id int

type View struct {
	Id   Id
	Tags Tags
}

func (v *View) Use(a Archetype) *View {
	v.Tags.Use(a)
	return v
}

func (v *View) Apply(n TagName, i interface{}) *View {
	v.Tags.Apply(n, i)
	return v
}

func (v *View) ApplyAll(o Tags) *View {
	v.Tags.ApplyAll(o)
	return v
}

func (v *View) Remove(n TagName) *View {
	v.Tags.Remove(n)
	return v
}

func (v *View) RemoveAll(ns ...TagName) *View {
	v.Tags.RemoveAll(ns...)
	return v
}
