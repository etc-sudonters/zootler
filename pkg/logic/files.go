package logic

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"sudonters/zootler/internal/entity"

	"muzzammil.xyz/jsonc"
)

type (
	RegionName   string
	LocationName string
	EventName    string
	SceneName    string
	HintGroup    string
	SaveWarp     string
	Dungeon      string
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

func (l RawLogicLocation) Components() []entity.Component {
	var comps []entity.Component

	if l.Scene != nil {
		comps = append(comps, *l.Scene)
	}

	if l.Hint != nil {
		comps = append(comps, *l.Hint)
	}

	if l.Dungeon != "" {
		comps = append(comps, Dungeon(l.Dungeon))
	}

	if l.SaveWarp != "" {
		comps = append(comps, SaveWarp(l.SaveWarp))
	}

	return comps
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

	for i, l := range locs {
		for evt, rule := range l.Events {
			l.Events[evt] = compressWhiteSpace(rule)
		}

		for exit, rule := range l.Exits {
			l.Exits[exit] = compressWhiteSpace(rule)
		}

		for location, rule := range l.Locations {
			l.Locations[location] = compressWhiteSpace(rule)
		}

		locs[i] = l
	}

	return locs, nil
}

func compressWhiteSpace[S ~string](raw S) S {
	r := string(raw)
	r = trailWhiteSpace.ReplaceAllLiteralString(r, "")
	r = leadWhiteSpace.ReplaceAllLiteralString(r, "")
	return S(compressWhiteSpaceRe.ReplaceAllLiteralString(r, " "))
}

var compressWhiteSpaceRe *regexp.Regexp = regexp.MustCompile(`\s+`)
var leadWhiteSpace *regexp.Regexp = regexp.MustCompile(`^\s+`)
var trailWhiteSpace *regexp.Regexp = regexp.MustCompile(`\s+$`)
