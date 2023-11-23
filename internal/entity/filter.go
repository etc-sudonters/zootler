package entity

import (
	"reflect"

	"github.com/etc-sudonters/substrate/skelly/hashset"
)

type FilterOption func(FilterBuilder) FilterBuilder

func BuildFilter(opts ...FilterOption) FilterBuilder {
	return FilterBuilder{}.Configure(opts...)
}

// used to narrow the entity population
type Filter struct {
	include []reflect.Type
	exclude []reflect.Type
}

func (f Filter) With() []reflect.Type {
	return f.include
}

func (f Filter) Without() []reflect.Type {
	return f.exclude
}

// constructs a Filter
type FilterBuilder struct {
	include hashset.Hash[reflect.Type]
	exclude hashset.Hash[reflect.Type]
}

func (f FilterBuilder) Configure(opts ...FilterOption) FilterBuilder {
	for _, o := range opts {
		f = o(f)
	}
	return f
}

func (f FilterBuilder) With(t reflect.Type) FilterBuilder {
	if f.include == nil {
		f.include = make(hashset.Hash[reflect.Type])
	}
	f.include.Add(t)
	return f
}

func (f FilterBuilder) Without(t reflect.Type) FilterBuilder {
	if f.exclude == nil {
		f.exclude = make(hashset.Hash[reflect.Type])
	}
	f.exclude.Add(t)
	return f
}

func (f FilterBuilder) Build() Filter {
	return Filter{
		include: hashset.AsSlice(f.include),
		exclude: hashset.AsSlice(f.exclude),
	}
}

func (f FilterBuilder) Clone() FilterBuilder {
	return FilterBuilder{
		include: hashset.FromMap(f.include),
		exclude: hashset.FromMap(f.exclude),
	}
}

func (f FilterBuilder) Combine(o FilterBuilder) FilterBuilder {
	return FilterBuilder{
		include: hashset.Union(f.include, o.include),
		exclude: hashset.Union(f.exclude, o.exclude),
	}
}

func (f FilterBuilder) Invert() FilterBuilder {
	return FilterBuilder{
		include: hashset.FromMap(f.exclude),
		exclude: hashset.FromMap(f.include),
	}
}
