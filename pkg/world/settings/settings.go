package settings

type SeedSettings struct {
	Logic    LogicRuleSet
	ItemPool ItemPool
	LogicSettings
	ShuffleSettings
}

type LogicSettings struct {
	KokriForest     KokiriForest
	KakGate         KakarikoGate
	DoorOfTime      DoorOfTime
	Fountain        ZorasFountain
	Bridge          BridgeRequirement
	TowerTrials     TowerTrialCount
	StartingAge     StartingAge
	ChildTradeQuest ChildTradeQuest
	AdultTradeItems AdultTradeItems
}

type ShuffleSettings struct {
	ShuffleSongs            SongShuffle
	ShuffleShops            ShopShuffle
	ShuffleTokens           GoldTokenShuffle
	ShuffleScrubs           ScrubShuffle
	ShufflePots             PotShuffle
	ShuffleCrate            CrateShuffle
	ShuffleCows             CowShuffle
	ShuffleBeehinves        BeehiveShuffle
	ShuffleKokriSword       KokriSwordShuffle
	ShuffleOcarinas         OcarinaShuffle
	ShuffleGerudoCard       GerudoCardShuffle
	ShuffleMagicBeans       MagicBeanShuffle
	ShuffleRepeatMerchants  RepeatMerchantShuffle
	ShuffleFrogRupees       FrogRupeeShuffle
	ShuffleMapsAndCompasses MapsAndCompassesShuffle
	ShuffleSmallKeys        SmallKeyShuffle
	ShuffleBossKeys         BossKeyShuffle
	ShuffleTowerBossKey     TowerBossKeyShuffle
	ShuffleChestGameKeys    ChestGameKeyShuffle
}

func DefaultSettings() map[string]any {
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
