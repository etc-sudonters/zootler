package filter

import "sudonters/zootler/internal/entity"

type location struct {
	b entity.FilterBuilder
}

func (l location) Build() entity.Filter {
	return l.b.Build()
}
