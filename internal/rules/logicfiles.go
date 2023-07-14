package rules

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"unicode/utf8"
)

type LogicLocation struct {
	Region    string            `json:"region_name"`
	Locations map[string]string `json:"locations"`
	Exits     map[string]string `json:"exits"`
	Events    map[string]string `json:"events"`
}

func ReadLogicFile(fp string) ([]LogicLocation, error) {
	contents, err := os.ReadFile(fp)
	if err != nil {
		return nil, err
	}

	var locs []LogicLocation
	if err := json.Unmarshal(contents, &locs); err != nil {
		return nil, err
	}

	return locs, nil
}

func LexAllLocationRules(locs []LogicLocation) error {
	var allErrs []error

	for _, loc := range locs {
		for check, rule := range loc.Locations {
			name := fmt.Sprintf("%s: %s", loc.Region, check)
			if err := lexEntire(name, rule); err != nil {
				allErrs = append(allErrs, err)
			}
		}

		for exit, rule := range loc.Exits {
			if err := lexEntire(exit, rule); err != nil {
				allErrs = append(allErrs, err)
			}
		}

		for event, rule := range loc.Events {
			name := fmt.Sprintf("%s %s", loc.Region, event)
			if err := lexEntire(name, rule); err != nil {
				allErrs = append(allErrs, err)
			}
		}
	}

	if allErrs != nil {
		return errors.Join(allErrs...)
	}

	return nil
}

func lexEntire(name, rule string) error {
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
