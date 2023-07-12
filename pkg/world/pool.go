package world

import (
	"github.com/etc-sudonters/zootler/pkg/entity"
)

type Pool struct {
	W Id
	entity.Pool
}

func (p Pool) Create(name Name) (entity.View, error) {
	view, err := p.Pool.Create()
	if err != nil {
		return nil, err
	}

	view.Add(OriginWorld(p.W))
	view.Add(name)
	return view, nil
}
