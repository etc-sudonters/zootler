package rules

import (
	"fmt"
	"os"
	"strings"

	"muzzammil.xyz/jsonc"
)

type (
	RegionName   string
	LocationName string
	EventName    string
	SceneName    string
	HintGroup    string
	RawRule      string
)

type RawLogicLocation struct {
	Region    RegionName               `json:"region_name"`
	Locations map[LocationName]RawRule `json:"locations"`
	Exits     map[LocationName]RawRule `json:"exits"`
	Events    map[EventName]RawRule    `json:"events"`
	Scene     *SceneName               `json:"scene"`
	Hint      *HintGroup               `json:"hint"`
	Dungeon   string                   `json:"dungeon"`
	SaveWarp  string                   `json:"savewarp"`
}

func (l RawLogicLocation) String() string {
	repr := &strings.Builder{}

	fmt.Fprintf(
		repr,
		"RawLogicLocation{\n\tRegion: %s,\n\tLocationCount: %d,\n\tExitCount: %d,\n}",
		l.Region, len(l.Locations), len(l.Exits))

	return repr.String()
}

func ReadLogicFile(fp string) ([]RawLogicLocation, error) {
	contents, err := os.ReadFile(fp)
	if err != nil {
		return nil, err
	}

	var locs []RawLogicLocation
	if err := jsonc.Unmarshal(contents, &locs); err != nil {
		return nil, err
	}

	return locs, nil
}
