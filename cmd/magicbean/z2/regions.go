package z2

import (
	"io/fs"
	"path/filepath"
	"sudonters/zootler/internal"

	"github.com/etc-sudonters/substrate/slipup"
)

type Region struct {
	proxy
	Name
}

type Edge struct {
	proxy
	Connection
}

type Regions struct {
	Regions     NamedEntities
	Connections TrackedEntities[Connection]
}

func (this *Regions) RegionNamed(name Name) Region {
	proxy := this.Regions.Entity(name)
	return Region{proxy, name}
}

func (this *Regions) Connect(from, to Region, kind ConnectionKind) Edge {
	connection := Connection{from.id, to.id}
	edge := this.Connections.Entity(connection)
	return Edge{edge, connection}
}

type RegionLoader struct {
	Regions
	Tokens Tokens
}

func (this *RegionLoader) Load(raw region) {
	region := this.RegionNamed(Name(raw.RegionName))

	for exit, rule := range raw.Exits {
		exit := this.RegionNamed(Name(exit))
		edge := this.Connect(region, exit, ConnectionExit)
		edge.Attach(StringSource(rule))
	}

	for location, rule := range raw.Locations {
		location := this.RegionNamed(Name(location))
		location.Attach(EmptyPlacement{})
		edge := this.Connect(region, location, ConnectionCheck)
		edge.Attach(StringSource(rule))
	}

	for event, rule := range raw.Events {
		eventName := Name(event)
		location := this.RegionNamed(NameF("%s %s", raw.RegionName, eventName))
		token := this.Tokens.Entity(eventName)
		edge := this.Connect(region, location, ConnectionCheck)
		edge.Attach(StringSource(rule))
		location.Attach(HoldsToken(token.Entity()), FixedPlacement{}, Generated{})
		token.Attach(HeldAt(location.Entity()), FixedPlacement{})
	}

	var attachments attachments

	if raw.Hint != "" {
		attachments.add(HintRegion(raw.Hint))
	}

	if raw.AltHint != "" {
		attachments.add(AltHintRegion(raw.AltHint))
	}

	if raw.Dungeon != "" {
		attachments.add(DungeonName(raw.Dungeon))
	}

	if raw.IsBossRoom {
		attachments.add(IsBossRoom{})
	}

	if raw.Savewarp != "" {
		attachments.add(Savewarp(raw.Savewarp))
	}

	if raw.Scene != "" {
		attachments.add(Scene(raw.Scene))
	}

	if raw.TimePasses {
		attachments.add(TimePassess{})
	}

	region.AttachAll(attachments.v)
}

func LoadRegionsFromFile(loader *RegionLoader, path string) error {
	these, readErr := internal.ReadJsonFileAs[[]region](path)
	if readErr != nil {
		return slipup.Describef(readErr, "while reading file '%s'", path)
	}
	for i := range these {
		loader.Load(these[i])
	}
	return nil
}

func LoadRegionsFromDirectory(loader *RegionLoader, dir string) error {
	return filepath.WalkDir(dir, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return slipup.Describe(err, "logic directory walk called with err")
		}

		info, err := entry.Info()
		if err != nil || info.Mode() != (^fs.ModeType)&info.Mode() {
			// either we couldn't get the info, which doesn't bode well
			// or it's some kind of not file thing which we also don't want
			return nil
		}

		if ext := filepath.Ext(path); ext != ".json" {
			return nil
		}

		return LoadRegionsFromFile(loader, path)
	})
}
