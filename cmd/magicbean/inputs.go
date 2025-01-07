package main

import (
	"io/fs"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"sudonters/zootler/internal"
	"sudonters/zootler/magicbeanvm"
	"sudonters/zootler/magicbeanvm/ast"
	"sudonters/zootler/magicbeanvm/symbols"

	"github.com/etc-sudonters/substrate/slipup"
)

type hardCodedRule struct{ where, logic string }
type partialRule struct {
	where, token symbols.Index
	body         ast.Node
}

var rules = []hardCodedRule{
	/*
		{"nested-and", "item and adult and dance"},
		{"nested-or", "item or adult or dance"},
		{"late-expand-at", "at('Forest Temple Outside Upper Ledge', True)"},
		{"late-expand-here", "here(logic_forest_mq_hallway_switch_boomerang and can_use(Boomerang))"},
		{"is-trick-enabled", "logic_forest_mq_hallway_switch_jumpslash"},
		{"call-func", "can_use(Hover_Boots)"},
		{"true", "True"},
		{"true-or", "True or at('Forest Temple Outside Upper Ledge', False)"},
		{"true-and", "True and at('Forest Temple Outside Upper Ledge', True)"},
		{"or-true", "at('Forest Temple Outside Upper Ledge', False) or True)"},
		{"and-true", "at('Forest Temple Outside Upper Ledge', True) and True"},
		{"false", "False"},
		{"false-or", "False or at('Forest Temple Outside Upper Ledge', False)"},
		{"false-and", "False and at('Forest Temple Outside Upper Ledge', False)"},
		{"or-false", "at('Forest Temple Outside Upper Ledge', False) or False)"},
		{"and-false", "at('Forest Temple Outside Upper Ledge', False) and False"},
		{"compare-eq-same-is-true", "same == same"},
		{"compare-nq-diff-is-true", "same != diff"},
		{"compare-eq-diff-is-false", "same == diff"},
		{"compare-nq-same-is-false", "same != same"},
		{"compare-lt", "chicken_count < 7"},
		{"compare-setting", "deadly_bonks == 'ohko'"},
		{"uses-setting", "('Triforce Piece', victory_goal_count)"},
		{"contains", "'Deku Tree' in dungeon_shortcuts"},
		{"subscript", "skipped_trials[Forest]"},
		{"float", "can_live_dmg(0.5)"},
		{"promote-standalone-token", "Progressive_Hookshot"},
		{"has-all", "has(taco, 1) and has(burrito, 1) and has(taquito, 1)"},
		{"has-any", "has(taco, 1) or has(burrito, 1) or has(taquito, 1)"},
		{"has-all-mix", "has(taco, 2) and has(burrito, 1) and has(taquito, 1) and is_adult"},
		{"has-any-mix", "has(taco, 2) or has(burrito, 1) or has(taquito, 1) or is_child"},
		{"call-helper", "can_use(Dins_Fire)"},
		{"can-use-hookshot", "can_use(Hookshot)"},
		{"can-use-goron-tunic", "can_use(Goron_Tunic)"},
		{"promote-standalone-func", "is_adult"},
		{"implicit-has", "(Spirit_Temple_Small_Key, 15)"},
		{"really-implicit-has", "Dins_Fire"},
		{"really-really-implicit-has", "'Goron Tunic'"},
		{"goron-tunic", "is_adult and ('Goron Tunic' or Buy_Goron_Tunic)"},
		{"goron-tunic", "is_adult or ('Goron Tunic' or Buy_Goron_Tunic)"},
		{"subscripts", "(skipped_trials[Forest] or 'Forest Trial Clear') and (skipped_trials[Fire] or 'Fire Trial Clear') and (skipped_trials[Water] or 'Water Trial Clear') and (skipped_trials[Shadow] or 'Shadow Trial Clear') and (skipped_trials[Spirit] or 'Spirit Trial Clear') and (skipped_trials[Light] or 'Light Trial Clear')"},
		{"logic_rules", "logic_rules == 'glitched'"},
		{"recursive-macro", "here(at('dance hall', dance))"},
		{"Big Poe Kill", "can_ride_epona and Bow and has_bottle"},
		{"nest-expansions", "at('Forest Temple Outside Upper Ledge', here(logic_forest_mq_hallway_switch_boomerang and can_use(Boomerang)))"},
	*/

	{"hmmm", "is_adult and can_trigger_lacs"},
}

func fakeLoadingRules() ([]loadingRule, error) {
	fake := make([]loadingRule, 0, len(rules))

	for idx := range rules {
		item := rules[idx]
		fake = append(fake, loadingRule{
			parent: item.where,
			name:   item.where,
			body:   item.logic,
			kind:   symbols.EVENT,
		})
	}

	return fake, nil
}

func aliasTokens(symbols *symbols.Table, funcs *ast.PartialFunctionTable, names []string) {
	for _, name := range names {
		escaped := escape(name)
		if _, exists := funcs.Get(escaped); exists {
			continue
		}
		if _, exists := funcs.Get(name); exists {
			continue
		}
		symbol := symbols.LookUpByName(name)
		symbols.Alias(symbol, escaped)
	}
}

var escaping = regexp.MustCompile("['()[\\]-]")

func escape(name string) string {
	name = escaping.ReplaceAllLiteralString(name, "")
	return strings.ReplaceAll(name, " ", "_")
}

func ReadHelpers(path string) map[string]string {
	contents, err := internal.ReadJsonFileStringMap(path)
	if err != nil {
		panic(err)
	}
	return contents
}

type loadingRule struct {
	parent, name, body string
	kind               symbols.Kind
}

func FakeSourceRules() []magicbeanvm.Source {
	source := make([]magicbeanvm.Source, len(rules))
	for i := range source {
		source[i] = magicbeanvm.Source{
			Kind:              magicbeanvm.SourceCheck,
			String:            magicbeanvm.SourceString(rules[i].logic),
			OriginatingRegion: rules[i].where,
			Destination:       rules[i].where,
		}
	}
	return source
}

func SourceRules(locations []location) (source []magicbeanvm.Source) {
	type pair struct {
		kind   magicbeanvm.SourceKind
		source map[string]string
	}

	for _, location := range locations {
		pairs := []pair{
			{magicbeanvm.SourceCheck, location.Locations},
			{magicbeanvm.SourceEvent, location.Events},
			{magicbeanvm.SourceTransit, location.Exits},
		}

		for _, pair := range pairs {
			source = slices.Concat(source, sourceRules(
				location.RegionName,
				pair.kind,
				pair.source,
			))
		}
	}

	return
}

func sourceRules(origin string, kind magicbeanvm.SourceKind, rules map[string]string) []magicbeanvm.Source {
	chunk := make([]magicbeanvm.Source, 0, len(rules))
	for destination, rule := range rules {
		chunk = append(chunk, magicbeanvm.Source{
			Kind:              kind,
			String:            magicbeanvm.SourceString(rule),
			OriginatingRegion: origin,
			Destination:       destination,
		})
	}
	return chunk
}

func readLogicFiles(logicDir string) ([]location, error) {
	var locations []location
	err := filepath.WalkDir(logicDir, func(path string, entry fs.DirEntry, err error) error {
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

		these, readErr := internal.ReadJsonFileAs[[]location](path)
		if readErr != nil {
			return slipup.Describef(readErr, "while reading file '%s'", path)
		}

		locations = slices.Concat(locations, these)
		return nil
	})

	return locations, err
}

func loaddir(logicDir string) ([]loadingRule, error) {
	var rules []loadingRule
	err := filepath.WalkDir(logicDir, func(path string, entry fs.DirEntry, err error) error {
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

		nodes, readErr := internal.ReadJsonFileAs[[]location](path)
		if readErr != nil {
			return slipup.Describef(readErr, "while reading file '%s'", path)
		}

		for _, node := range nodes {
			namedRules := []struct {
				rules map[string]string
				kind  symbols.Kind
			}{
				{node.Locations, symbols.LOCATION},
				{node.Events, symbols.EVENT},
				{node.Exits, symbols.TRANSIT},
			}

			for _, bulk := range namedRules {
				for name, body := range bulk.rules {
					rules = append(rules, loadingRule{
						parent: node.RegionName,
						name:   name,
						body:   body,
						kind:   bulk.kind,
					})
				}
			}
		}
		return nil
	})
	return rules, err
}

func loadTokensNames(path string) ([]string, error) {
	loading, err := internal.ReadJsonFileAs[[]item](path)
	if err != nil {
		return nil, slipup.Describef(err, "while loading components from '%s'", path)
	}

	names := make([]string, len(loading))
	for i := range loading {
		names[i] = loading[i].Name
	}

	return names, nil
}

type location struct {
	Events     map[string]string `json:"events"`
	Exits      map[string]string `json:"exits"`
	Locations  map[string]string `json:"locations"`
	RegionName string            `json:"region_name"`
	AltHint    string            `json:"alt_hint"`
	Hint       string            `json:"hint"`
	Dungeon    string            `json:"dungeon"`
	IsBossRoom bool              `json:"is_boss_room"`
	Savewarp   string            `json:"savewarp"`
	Scene      string            `json:"scene"`
	TimePasses bool              `json:"time_passes"`
}

type item struct {
	Name        string                 `json:"name"`
	Type        string                 `json:"type"`
	Advancement bool                   `json:"advancement"`
	Priority    bool                   `json:"priority"`
	Special     map[string]interface{} `json:"special"`
}
