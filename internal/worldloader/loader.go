package worldloader

import (
	"errors"
	"iter"
	"sudonters/zootler/internal/components"
	"sudonters/zootler/internal/query"
	"sudonters/zootler/internal/table"

	"github.com/etc-sudonters/substrate/slipup"
)

type LogicLoader struct {
	L *Locations
	T *Tokens
}

func (l *LogicLoader) LoadLocations(loading iter.Seq[LogicLocation]) error {
	for loading := range loading {
		location, err := l.L.Build(loading.Name)
		if err != nil {
			return slipup.Describef(err, "while loading %s", loading.Name)
		}
		location.Attach(loading.Components())

		if edgesErr := l.loadEdges(location, loading.Edges); edgesErr != nil {
			return edgesErr
		}

		if eventsErr := l.loadEvents(loading.Events); eventsErr != nil {
			return eventsErr
		}

	}

	return nil
}

func (l *LogicLoader) loadEdges(location *LocationBuilder, edges iter.Seq[locedge]) error {
	for edge := range edges {
		destination, destErr := l.L.Build(edge.name)
		if destErr != nil {
			return slipup.Describef(destErr, "while creating destination %s", edge.name)
		}

		connection, connectErr := l.L.Connect(location, destination, edge.rule)
		if connectErr != nil {
			return slipup.Describef(connectErr, "while creating connection %s -> %s", location.name, edge.name)
		}

		if err := connection.Attach(edge.components); err != nil {
			return slipup.Describef(err, "while attaching components to %s", connection.name)
		}
	}

	return nil
}

func (l *LogicLoader) loadEvents(events map[string]string) error {
	for loc := range l.L.MustEachEdge(events) {
		l.T.item[loc.normaled] = loc.id
		if err := loc.Attach(table.Values{
			components.Advancement{},
			components.CollectableGameToken{},
			components.Inhabited(loc.id),
			components.Inhabits(loc.id),
			components.Locked{},
			components.Event{},
		}); err != nil {
			return slipup.Describef(err, "while attaching event components to %s", loc.name)
		}
	}
	return nil
}

func (l *LogicLoader) Init(eng query.Engine) error {
	locs, locsErr := NewLocations(eng)
	toks, toksErr := NewTokens(eng)

	l.L = locs
	l.T = toks
	return errors.Join(locsErr, toksErr)
}
