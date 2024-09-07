package ichiro

import (
	"path"
	"sudonters/zootler/internal/app"
	"sudonters/zootler/internal/entities"

	"github.com/etc-sudonters/substrate/slipup"
)

type DataLoader struct {
	Table    TableLoader
	DataPath string
}

func (dl *DataLoader) Setup(z *app.Zootlr) error {
	if err := dl.Table.Setup(z); err != nil {
		return slipup.Createf("while creating data table")
	}

	tokens, tokenMapErr := entities.TokenMap(z.Engine())
	if tokenMapErr != nil {
		return slipup.Describe(tokenMapErr, "while creating token map")
	}

	locations, locationMapErr := entities.LocationMap(z.Engine())
	if locationMapErr != nil {
		return slipup.Describe(locationMapErr, "while creating location map")
	}

	edges, edgeMapErr := entities.EdgeMap(z.Engine())
	if edgeMapErr != nil {
		return slipup.Describe(edgeMapErr, "while creating edge map")
	}

	if itemLoadErr := LoadDataFile[ItemComponents](
		path.Join(dl.DataPath, "items.json"),
		tokens,
	); itemLoadErr != nil {
		return slipup.Describe(itemLoadErr, "while loading item data")
	}

	if locationLoadErr := LoadDataFile[LocationComponents](
		path.Join(dl.DataPath, "locations.json"),
		locations,
	); locationLoadErr != nil {
		return slipup.Describe(locationLoadErr, "while loading item data")
	}

	app.AddResource[entities.Map[entities.Token]](tokens)
	app.AddResource[entities.Map[entities.Location]](locations)
	app.AddResource[entities.Map[entities.Edge]](edges)
	return nil
}
