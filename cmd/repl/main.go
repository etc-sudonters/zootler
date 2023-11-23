package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sudonters/zootler/internal/entity"
	"sudonters/zootler/internal/entity/bitpool"
	"sudonters/zootler/internal/entity/componenttable"
	"sudonters/zootler/pkg/logic"
	"sudonters/zootler/pkg/logic/interpreter"
	"sudonters/zootler/pkg/rules/ast"
	"sudonters/zootler/pkg/rules/parser"
	"sudonters/zootler/pkg/world"
	"sudonters/zootler/pkg/world/components"

	"github.com/etc-sudonters/substrate/mirrors"
	"github.com/etc-sudonters/substrate/skelly/hashset"
	"muzzammil.xyz/jsonc"
)

func main() {
	tbl := componenttable.New(30000)
	pool := bitpool.FromTable(tbl, 400)
	b := world.NewBuilder(pool, tbl)

	env := interpreter.NewEnv()
	rewriter := interpreter.NewInliner(env)
	rewriter.Settings = make(map[string]any, 0)
	rewriter.SkippedTrials = make(map[string]bool, 0)
	rewriter.Tricks = defaultTricks()
	rewriter.Builder = b
	loadHelpers("inputs/logic", env, rewriter)

	for name, value := range defaultSettings() {
		env.Set(name, interpreter.Box(value))
	}

	I := interpreter.New(env)

	tokens, err := b.Pool.Query(entity.FilterBuilder{}.
		With(mirrors.TypeOf[components.Token]()).
		With(mirrors.TypeOf[components.Name]()).
		Build())

	if err != nil {
		panic(err)
	}

	var itemName components.Name
	typedStrings := mirrors.NewTypedStrings()

	for _, t := range tokens {
		t.Get(&itemName)
		literal := logic.EscapeName(string(itemName))
		component := typedStrings.InstanceOf(literal)
		t.Add(component)
		env.Set(literal, interpreter.Token{
			Component: typedStrings.Typed(literal),
			Literal:   literal,
		})
	}
	tbl.RowOf(mirrors.TypeOf[components.Collected]())

	rules, err := b.Pool.Query(entity.FilterBuilder{}.
		With(mirrors.TypeOf[logic.RawRule]()).
		With(mirrors.TypeOf[world.Edge]()).
		With(mirrors.TypeOf[world.FromName]()).
		With(mirrors.TypeOf[components.Name]()).
		Build())

	if err != nil {
		panic(err)
	}

	var rule logic.RawRule
	var region world.FromName

	env.SetBuiltIn("at_day", 0, interpreter.AtDay)
	env.SetBuiltIn("at_night", 0, interpreter.AtNigt)
	env.SetBuiltIn("at_dampe_time", 0, interpreter.AtDampe)

	env.SetBuiltIn("has", 2, interpreter.Zoot_HasQuantityOf{
		Entities: b.Pool,
		Strs:     b.TypedStrs,
	})

	env.SetBuiltIn("has_medallions", 1, interpreter.Zoot_HasMedallions{
		Has: interpreter.Zoot_HasQuantityOf{
			Entities: b.Pool,
			Strs:     b.TypedStrs,
		},
	})

	env.SetBuiltIn("region_has_shortcuts", 1, interpreter.Zoot_RegionHasShortcuts{
		// TODO: add to settings and peel off from there
		RegionalShortcuts: hashset.New[string](),
	})

	env.SetBuiltIn("has_bottle", 0, interpreter.Zoot_HasBottle{
		Has: interpreter.Zoot_HasQuantityOf{
			Entities: b.Pool,
			Strs:     b.TypedStrs,
		},
	})

	// argument to rule
	env.SetBool("spot", false) // TODO absolutely not the correct thing
	env.SetNumber("tod", 0)
	env.SetString("age", "child")

	// runtime needs to calculate these properties based on zootr's logic
	env.SetBool("skip_child_zelda", true)

	// wat, these are all for projectile check
	env.SetString("child", "child")
	env.SetString("adult", "adult")
	env.SetString("both", "both")
	env.SetString("either", "either")

	var results []interpreter.Value

	for _, bearer := range rules {
		bearer.Get(&rule)
		bearer.Get(&region)
		rewriter.RegionName = string(region)
		p, err := parser.Parse(string(rule))
		if err != nil {
			panic(err)
		}
		p = rewriter.Rewrite(p, env)
		if r, rewrote := rewriter.Make0ArityFnCall(p, env); rewrote {
			p = rewriter.Rewrite(r, env)
		}
		val := I.Evaluate(p, env)
		results = append(results, val)
	}

}

func loadHelpers(logicDir string, env interpreter.Environment, rewriter *interpreter.Inliner) {
	contents, err := os.ReadFile(filepath.Join(logicDir, "LogicHelpers.json"))
	if err != nil {
		panic(err)
	}

	var helpers map[string]string

	if err := jsonc.Unmarshal(contents, &helpers); err != nil {
		panic(err)
	}

	passed := true

	for raw, helper := range helpers {
		helper = compressWhiteSpace(helper)
		decl, err := parser.Parse(raw)
		if err != nil {
			passed = false
			fmt.Fprintf(os.Stdout, "Name:\t%s\n", raw)
			fmt.Fprintf(os.Stdout, "FAILED TO PARSE: %s\n", helper)
			fmt.Fprintf(os.Stdout, "ERROR: %s\n", err.Error())
			continue
		}
		rule, err := parser.Parse(helper)
		if err != nil {
			passed = false
			fmt.Fprintf(os.Stdout, "Name:\t%s\n", raw)
			fmt.Fprintf(os.Stdout, "FAILED TO PARSE: %s\n", helper)
			fmt.Fprintf(os.Stdout, "ERROR: %s\n", err.Error())
			continue
		}

		interpreter.FunctionDecl(decl, rule, env)
	}

	if !passed {
		os.Exit(99)
	}
}

func toZootCallable(rawDecl string, body ast.Expression) interpreter.Fn {
	decl, err := parser.Parse(rawDecl)
	if err != nil {
		panic(err)
	}

	switch decl := decl.(type) {
	// 0 arity function decls look like identifiers
	case *ast.Identifier:
		return interpreter.Fn{
			Name:   decl,
			Body:   body,
			Params: nil,
		}
	// N arity decls look like calls
	case *ast.Call:
		name := decl.Callee.(*ast.Identifier)
		params := make([]string, len(decl.Args))
		for i := range params {
			params[i] = decl.Args[i].(*ast.Identifier).Value
		}
		return interpreter.Fn{
			Name:   name,
			Body:   body,
			Params: params,
		}
	default:
		panic(fmt.Errorf("expected Ident or Call for decl, got %T", decl))
	}
}

func compressWhiteSpace[S ~string](s S) S {
	return S(strings.Join(strings.Fields(string(s)), " "))
}

func defaultSettings() map[string]any {
	return map[string]any{
		"show_seed_info":                          true,
		"user_message":                            "",
		"world_count":                             1,
		"create_spoiler":                          true,
		"randomize_settings":                      false,
		"logic_rules":                             "glitchless",
		"reachable_locations":                     "all",
		"triforce_hunt":                           false,
		"lacs_condition":                          "vanilla",
		"bridge":                                  "medallions",
		"bridge_medallions":                       6,
		"trials_random":                           false,
		"trials":                                  0,
		"shuffle_ganon_bosskey":                   "remove",
		"shuffle_bosskeys":                        "dungeon",
		"shuffle_smallkeys":                       "dungeon",
		"shuffle_hideoutkeys":                     "vanilla",
		"shuffle_tcgkeys":                         "vanilla",
		"key_rings_choice":                        "off",
		"shuffle_silver_rupees":                   "vanilla",
		"shuffle_mapcompass":                      "startwith",
		"enhance_map_compass":                     false,
		"open_forest":                             "closed_deku",
		"open_kakariko":                           "open",
		"open_door_of_time":                       true,
		"zora_fountain":                           "open",
		"gerudo_fortress":                         "fast",
		"dungeon_shortcuts_choice":                "off",
		"starting_age":                            "adult",
		"mq_dungeons_mode":                        "vanilla",
		"empty_dungeons_mode":                     "none",
		"shuffle_interior_entrances":              "off",
		"shuffle_grotto_entrances":                false,
		"shuffle_dungeon_entrances":               "off",
		"shuffle_bosses":                          "off",
		"shuffle_overworld_entrances":             false,
		"shuffle_gerudo_valley_river_exit":        false,
		"owl_drops":                               true,
		"warp_songs":                              false,
		"free_bombchu_drops":                      false,
		"one_item_per_dungeon":                    false,
		"shuffle_song_items":                      "song",
		"shopsanity":                              "off",
		"tokensanity":                             "off",
		"shuffle_scrubs":                          "off",
		"shuffle_freestanding_items":              "off",
		"shuffle_pots":                            "off",
		"shuffle_crates":                          "off",
		"shuffle_cows":                            false,
		"shuffle_beehives":                        false,
		"shuffle_kokiri_sword":                    true,
		"shuffle_ocarinas":                        false,
		"shuffle_gerudo_card":                     false,
		"shuffle_beans":                           false,
		"shuffle_expensive_merchants":             false,
		"shuffle_frog_song_rupees":                false,
		"shuffle_individual_ocarina_notes":        true,
		"shuffle_loach_reward":                    "off",
		"logic_no_night_tokens_without_suns_song": false,
		"start_with_consumables":                  true,
		"start_with_rupees":                       false,
		"starting_hearts":                         3,
		"no_escape_sequence":                      true,
		"no_guard_stealth":                        true,
		"no_epona_race":                           true,
		"skip_some_minigame_phases":               true,
		"complete_mask_quest":                     false,
		"useful_cutscenes":                        false,
		"fast_chests":                             true,
		"free_scarecrow":                          false,
		"fast_bunny_hood":                         true,
		"auto_equip_masks":                        false,
		"plant_beans":                             false,
		"chicken_count_random":                    false,
		"chicken_count":                           7,
		"big_poe_count_random":                    false,
		"big_poe_count":                           1,
		"easier_fire_arrow_entry":                 false,
		"ruto_already_f1_jabu":                    false,
		"ocarina_songs":                           "off",
		"correct_chest_appearances":               "both",
		"minor_items_as_major_chest":              false,
		"invisible_chests":                        false,
		"correct_potcrate_appearances":            "textures_content",
		"key_appearance_match_dungeon":            false,
		"clearer_hints":                           true,
		"hints":                                   "always",
		"hint_dist":                               "tournament",
		"text_shuffle":                            "none",
		"damage_multiplier":                       "normal",
		"deadly_bonks":                            "none",
		"no_collectible_hearts":                   false,
		"starting_tod":                            "default",
		"blue_fire_arrows":                        false,
		"fix_broken_drops":                        false,
		"item_pool_value":                         "balanced",
		"junk_ice_traps":                          "off",
		"ice_trap_appearance":                     "junk_only",
		"adult_trade_shuffle":                     false,
		"bridge_tokens":                           100,
		"ganon_bosskey_tokens":                    999,
		"lacs_tokens":                             999,
	}
}

func defaultTricks() map[string]bool {
	return map[string]bool{
		"visible_collisions":          true,
		"grottos_without_agony":       true,
		"fewer_tunic_requirements":    true,
		"rusted_switches":             true,
		"man_on_roof":                 true,
		"windmill_poh":                true,
		"crater_bean_poh_with_hovers": true,
		"dc_jump":                     true,
		"lens_botw":                   true,
		"child_deadhand":              true,
		"forest_vines":                true,
		"lens_shadow":                 true,
		"lens_shadow_platform":        true,
		"lens_bongo":                  true,
		"lens_spirit":                 true,
		"lens_gtg":                    true,
		"lens_castle":                 true,
	}
}
