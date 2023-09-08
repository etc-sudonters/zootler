package rules

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"unicode/utf8"

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

func LexAllLocationRules(loc RawLogicLocation) error {
	var allErrs []error

	for check, rule := range loc.Locations {
		name := fmt.Sprintf("%s: %s", loc.Region, check)
		if err := lexEntire(name, rule); err != nil {
			allErrs = append(allErrs, err)
		}
	}

	for exit, rule := range loc.Exits {
		if err := lexEntire(string(exit), rule); err != nil {
			allErrs = append(allErrs, err)
		}
	}

	for event, rule := range loc.Events {
		name := fmt.Sprintf("%s %s", loc.Region, event)
		if err := lexEntire(name, rule); err != nil {
			allErrs = append(allErrs, err)
		}
	}

	if allErrs != nil {
		return errors.Join(allErrs...)
	}

	return nil
}

func lexEntire(name string, rawRule RawRule) error {
	var rule string = string(rawRule)
	l := lex(name, rule)
	runeCount := utf8.RuneCountInString(rule)
	i := 0
	collected := []item{}

	for {
		if i > runeCount {
			repr := &strings.Builder{}
			fmt.Fprintf(repr, "lexer %s spinning, killing test", l.name)
			fmt.Fprintf(repr, "collected tokens %+v", collected)
			fmt.Fprintf(repr, "lexer state %+v", l)
			return errors.New(repr.String())
		}

		item := l.nextItem()
		if item.typ == itemEof {
			break
		}

		collected = append(collected, item)

		// collect err but not eof
		if item.typ == itemErr {
			break
		}

		i++
	}

	lastTok := collected[len(collected)-1]
	if lastTok.typ == itemErr {
		repr := &strings.Builder{}
		fmt.Fprintf(repr, "failed to parse %s\n", name)
		fmt.Fprintf(repr, "err: %s", lastTok)
		fmt.Fprintf(repr, "rule snippet: %.50q", rule)
		return errors.New(repr.String())
	}

	return nil
}
