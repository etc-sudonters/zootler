package main

import (
	"io/fs"
	"path/filepath"
	"sudonters/zootler/internal"

	"github.com/etc-sudonters/substrate/slipup"
)

var globals = []string{
	"age",
	"Forest",
	"Fire",
	"Water",
	"Shadow",
	"Spirit",
	"Light",
}

var tokens = []string{
	"Dins_Fire",
	"Nayrus_Love",
	"Farores_Wind",
	"Magic_Meter",
	"Bow",
	"Megaton_Hammer",
	"Iron_Boots",
	"Hover_Boots",
	"Mirror_Shield",
	"Slingshot",
	"Boomerang",
	"Kokiri_Sword",
	"Buy_Goron_Tunic",
}

var compTime = []string{
	"load_setting",
	"load_setting_2",
}

var builtIns = []string{
	"has_dungeon_shortcuts",
	"is_trial_skipped",
	"at_dampe_time",
	"at_day",
	"at_night",
	"had_night_start",
	"has_bottle",
	"has_hearts",
	"has_stones",
	"is_adult",
	"is_child",
	"is_starting_age",
}

var settings = []string{
	"logic_rules",
	"adult_trade_shuffle",
	"big_poe_count",
	"bridge",
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
	"gerudo_fortress",
	"lacs_condition",
	"lacs_hearts",
	"lacs_medallions",
	"open_door_of_time",
	"open_forest",
	"open_kakariko",
	"plant_beans",
	"chicken_count",
	"selected_adult_trade_item",
	"shuffle_dungeon_entrances",
	"shuffle_empty_pots",
	"shuffle_expensive_merchants",
	"shuffle_ganon_bosskey",
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
	"triforce_goal_per_world",
	"warp_songs",
	"zora_fountain",

	"bridge_hearts",
	"bridge_medallions",
	"bridge_rewards",
	"bridge_stones",
	"bridge_tokens",
	"ganon_bosskey_tokens_hearts",
	"ganon_bosskey_tokens_medallions",
	"ganon_bosskey_tokens_rewards",
	"ganon_bosskey_tokens_stones",
	"ganon_bosskey_tokens_tokens",
	"lacs_hearts",
	"lacs_medallions",
	"lacs_rewards",
	"lacs_stones",
	"lacs_tokens",
	"starting_age",
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
			namedRules := []map[string]string{
				node.Locations,
				node.Events,
				node.Exits,
			}

			for _, pairs := range namedRules {
				for name, body := range pairs {
					rules = append(rules, loadingRule{
						node.RegionName, name, body,
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
