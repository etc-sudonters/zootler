package world

import (
	"github.com/etc-sudonters/zootler/pkg/entity"
	"github.com/etc-sudonters/zootler/pkg/logic"
)

type Pool struct {
	W Id
	entity.Pool
}

func (p Pool) Create(name logic.Name) (entity.View, error) {
	view, err := p.Pool.Create()
	if err != nil {
		return nil, err
	}

	view.Add(OriginWorld(p.W))
	view.Add(name)
	return view, nil
}
