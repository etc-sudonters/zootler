package main

import (
	"io/fs"
	"path/filepath"
	"regexp"
	"strings"
	"sudonters/zootler/internal"
	"sudonters/zootler/magicbeanvm/ast"
	"sudonters/zootler/magicbeanvm/symbols"

	"github.com/etc-sudonters/substrate/slipup"
)

type rule struct{ where, logic string }
type partialRule struct {
	where, token symbols.Index
	body         ast.Node
}

var rules = []rule{
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
}

func aliasTokens(table *symbols.Table, names []string) {
	for _, name := range names {
		symbol := table.LookUpByName(name)
		table.Alias(symbol, escape(name))
	}
}

var escaping = regexp.MustCompile("['()[\\]-]")

func escape(name string) string {
	name = escaping.ReplaceAllLiteralString(name, "")
	return strings.ReplaceAll(name, " ", "_")
}

var settings = []string{
	"adult_trade_shuffle",
	"big_poe_count",
	"blue_fire_arrows",
	"bridge",
	"bridge_hearts",
	"bridge_medallions",
	"bridge_rewards",
	"bridge_stones",
	"bridge_tokens",
	"chicken_count",
	"complete_mask_quest",
	"damage_multiplier",
	"deadly_bonks",
	"disable_trade_revert",
	"dungeon_shortcuts",
	"entrance_shuffle",
	"fix_broken_drops",
	"free_bombchu_drops",
	"free_scarecrow",
	"ganon_bosskey_hearts",
	"ganon_bosskey_medallions",
	"ganon_bosskey_rewards",
	"ganon_bosskey_stones",
	"ganon_bosskey_tokens",
	"ganon_bosskey_tokens_hearts",
	"ganon_bosskey_tokens_medallions",
	"ganon_bosskey_tokens_rewards",
	"ganon_bosskey_tokens_stones",
	"ganon_bosskey_tokens_tokens",
	"gerudo_fortress",
	"hints",
	"keysanity",
	"lacs_condition",
	"lacs_hearts",
	"lacs_hearts",
	"lacs_medallions",
	"lacs_medallions",
	"lacs_rewards",
	"lacs_stones",
	"lacs_tokens",
	"logic_rules",
	"open_door_of_time",
	"open_forest",
	"open_kakariko",
	"plant_beans",
	"selected_adult_trade_item",
	"shuffle_dungeon_entrances",
	"shuffle_empty_pots",
	"shuffle_expensive_merchants",
	"shuffle_ganon_bosskey",
	"shuffle_gerudo_fortress_heart_piece",
	"shuffle_hideout_entrances",
	"shuffle_individual_ocarina_notes",
	"shuffle_interior_entrances",
	"shuffle_overworld_entrances",
	"shuffle_pots",
	"shuffle_scrubs",
	"shuffle_silver_rupees",
	"shuffle_tcgkeys",
	"skip_child_zelda",
	"skip_reward_from_rauru",
	"skipped_trials",
	"starting_age",
	"triforce_goal_per_world",
	"warp_songs",
	"zora_fountain",
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
