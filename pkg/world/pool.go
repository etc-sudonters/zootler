package world

import (
	"fmt"

	"sudonters/zootler/pkg/entity"
	"sudonters/zootler/pkg/logic"
)

type Pool struct {
	W Id
	entity.Pool
}

func (p Pool) Create(name logic.Name) (entity.View, error) {
	view, err := p.Pool.Create()
	if err != nil {
		return nil, fmt.Errorf("failed to create entity %q: %w", name, err)
	}

	view.Add(OriginWorld(p.W))
	view.Add(name)
	return view, nil
}
