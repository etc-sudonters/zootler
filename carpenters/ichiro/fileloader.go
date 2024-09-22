package ichiro

import (
	"sudonters/zootler/internal"
	"sudonters/zootler/internal/components"
	"sudonters/zootler/internal/entities"
	"sudonters/zootler/internal/table"

	"github.com/etc-sudonters/substrate/slipup"
)

type ComponentLoader interface {
	EntityName() components.Name
	AsComponents() table.Values
}

func LoadDataFile[C ComponentLoader, E entities.Entity, M entities.Map[E]](path string, entities M) error {
	loaders, err := internal.ReadJsonFileAs[[]C](path)
	if err != nil {
		return slipup.Describef(err, "while loading components from '%s'", path)
	}

	for _, loader := range loaders {
		name := loader.EntityName()
		entity, entityErr := entities.Entity(name)
		if entityErr != nil {
			return slipup.Describef(entityErr, "while retrieving entity '%s'", name)
		}
		if addErr := entity.AddComponents(loader.AsComponents()); addErr != nil {
			return slipup.Describef(entityErr, "while populating entity '%s'", name)
		}
	}

	return nil
}
