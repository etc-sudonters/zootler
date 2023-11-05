package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sudonters/zootler/internal/astrender"
	"sudonters/zootler/internal/entity"
	"sudonters/zootler/pkg/logic"
	"sudonters/zootler/pkg/logic/interpreter"
	"sudonters/zootler/pkg/rules/ast"
	"sudonters/zootler/pkg/rules/parser"
	"sudonters/zootler/pkg/world"
	"sudonters/zootler/pkg/worldloader"

	"muzzammil.xyz/jsonc"
)

func main() {
	env := interpreter.NewEnv()
	rewriter := interpreter.NewRewriter(defaultTricks())
	loadHelpers("inputs/logic", env, rewriter)

	for name, value := range defaultSettings() {
		env.Set(name, interpreter.Box(value))
	}
	b := world.NewBuilder()
	loader := worldloader.FileSystemLoader{
		LogicDirectory: "inputs/logic",
		DataDirectory:  "inputs/data",
	}

	if err := loader.LoadInto(context.TODO(), b); err != nil {
		panic(err)
	}

	w := b.Build()

	tokens, err := w.Entities.Query([]entity.Selector{
		entity.With[logic.Token]{},
		entity.With[world.Name]{},
	})

	if err != nil {
		panic(err)
	}

	var itemName world.Name

	for _, t := range tokens {
		t.Get(&itemName)
		env.Set(worldloader.EscapeName(string(itemName)), interpreter.Box(t.Model()))
	}

	rules, err := w.Entities.Query([]entity.Selector{
		entity.With[logic.RawRule]{},
		entity.With[world.Edge]{},
		entity.With[world.FromName]{},
		entity.With[world.Name]{},
	})

	if err != nil {
		panic(err)
	}

	var rule logic.RawRule
	var region world.FromName
	var edgeName world.Name
	s := astrender.NewSexpr(astrender.DontTheme())
	identGetter := func(ident *ast.Identifier) ast.Expression {
		v, ok := env.Get(ident.Value)
		if !ok {
			return nil
		}

		if interpreter.CanLiteralfy(v) {
			return interpreter.Literalify(v)
		}

		if fn, ok := v.(interpreter.Fn); ok {
			return fn.Body
		}

		if fn, ok := v.(interpreter.PartiallyEvaluatedFn); ok {
			return fn.Body
		}

		return nil
	}

	for _, bearer := range rules {
		s.Clear()
		bearer.Get(&rule)
		bearer.Get(&region)
		bearer.Get(&edgeName)
		rewriter.SetRegion(string(region))

		p, err := parser.Parse(string(rule))
		if err != nil {
			panic(err)
		}
		s.SetIdentifier(nil)
		ast.Visit(s, p)
		fmt.Printf("Name: %s\nRaw s-expr: %s\n\n", edgeName, s.String())
		s.Clear()

		s.SetIdentifier(identGetter)
		p = rewriter.Rewrite(p, env)
		ast.Visit(s, p)
		fmt.Printf("Rewritten: %s\n\n\n\n", s.String())
		fmt.Scan()
	}
}

func loadHelpers(logicDir string, env interpreter.Environment, rewriter *interpreter.Rewriter) {
	contents, err := os.ReadFile(filepath.Join(logicDir, "LogicHelpers.json"))
	if err != nil {
		panic(err)
	}

	var helpers map[string]string

	if err := jsonc.Unmarshal(contents, &helpers); err != nil {
		panic(err)
	}

	passed := true

	for decl, helper := range helpers {
		helper = compressWhiteSpace(helper)
		rule, err := parser.Parse(helper)
		if err != nil {
			passed = false
			fmt.Fprintf(os.Stdout, "Name:\t%s\n", decl)
			fmt.Fprintf(os.Stdout, "FAILED TO PARSE: %s\n", helper)
			fmt.Fprintf(os.Stdout, "ERROR: %s\n", err.Error())
			continue
		}

		fn := toZootCallable(decl, rule)
		fn.Body = rewriter.Rewrite(fn.Body, env)
		if interpreter.CanReifyLiteral(fn.Body) && fn.Arity() == 0 {
			env.Set(fn.Name, interpreter.ReifyLiteral(fn.Body))
		} else {
			env.Set(fn.Name, fn)
		}
	}

	if !passed {
		os.Exit(99)
	}
}

func toZootCallable(decl string, body ast.Expression) interpreter.Fn {
	var err error
	if decl, declErr := parser.Parse(decl); declErr == nil {
		switch decl := decl.(type) {
		case *ast.Identifier:
			return interpreter.Fn{
				Name:   decl.Value,
				Body:   body,
				Params: nil,
			}
		case *ast.Call:
			name := decl.Callee.(*ast.Identifier).Value
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
			err = fmt.Errorf("expected Ident or Call got %T", decl)
		}
	} else {
		err = declErr
	}

	panic(err)
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
		"starting_age":                            "random",
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
		"shuffle_individual_ocarina_notes":        false,
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
