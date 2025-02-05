package settings

import (
	"fmt"
	"sudonters/zootler/internal/skelly/bitset32"
)

func (this *Zootr) Bitfield(name string) (bitset32.Bitset, error) {
	var val bitset32.Bitset
	switch name {
	case "dungeon_shortcuts":
		panic("not implemented for dungeon_shortcuts")
	case "mq_dungeons_specific":
		panic("not implemented for mq_dungeons_specific")
	case "empty_dungeons_specific":
		panic("not implemented for empty_dungeons_specific")
	case "empty_dungeons_rewards":
		panic("not implemented for empty_dungeons_rewards")
	case "spawn_positions":
		panic("not implemented for spawn_positions")
	case "shuffle_child_trade":
		panic("not implemented for shuffle_child_trade")
	case "adult_trade_start":
		panic("not implemented for adult_trade_start")
	case "key_rings":
		panic("not implemented for key_rings")
	case "silver_rupee_pouches":
		panic("not implemented for silver_rupee_pouches")
	case "disabled_locations":
		panic("not implemented for disabled_locations")
	case "allowed_tricks":
		panic("not implemented for allowed_tricks")
	case "starting_items":
		panic("not implemented for starting_items")
	case "starting_equipment":
		panic("not implemented for starting_equipment")
	case "starting_inventory":
		panic("not implemented for starting_inventory")
	case "misc_hints":
		panic("not implemented for misc_hints")
	case "starting_songs":
		panic("not implemented for starting_songs")
	case "chest_textures_specific":
		panic("not implemented for chest_textures_specific")
	default:
		return val, unknown(name)

	}
	return val, nil
}

func (this *Zootr) String(name string) (string, error) {
	var val string
	switch name {
	case "logic_rules":
		panic("not implemented for logic_rules")
	case "reachable_locations":
		panic("not implemented for reachable_locations")
	case "lacs_condition":
		panic("not implemented for lacs_condition")
	case "bridge":
		panic("not implemented for bridge")
	case "shuffle_ganon_bosskey":
		panic("not implemented for shuffle_ganon_bosskey")
	case "open_forest":
		panic("not implemented for open_forest")
	case "open_kakariko":
		panic("not implemented for open_kakariko")
	case "zora_fountain":
		panic("not implemented for zora_fountain")
	case "gerudo_fortress":
		panic("not implemented for gerudo_fortress")
	case "dungeon_shortcuts_choice":
		panic("not implemented for dungeon_shortcuts_choice")
	case "starting_age":
		panic("not implemented for starting_age")
	case "mq_dungeons_mode":
		panic("not implemented for mq_dungeons_mode")
	case "empty_dungeons_mode":
		panic("not implemented for empty_dungeons_mode")
	case "shuffle_dungeon_entrances":
		panic("not implemented for shuffle_dungeon_entrances")
	case "shuffle_bosses":
		panic("not implemented for shuffle_bosses")
	case "shuffle_song_items":
		panic("not implemented for shuffle_song_items")
	case "shopsanity_prices":
		panic("not implemented for shopsanity_prices")
	case "shopsanity":
		panic("not implemented for shopsanity")
	case "tokensanity":
		panic("not implemented for tokensanity")
	case "shuffle_scrubs":
		panic("not implemented for shuffle_scrubs")
	case "shuffle_freestanding_items":
		panic("not implemented for shuffle_freestanding_items")
	case "shuffle_pots":
		panic("not implemented for shuffle_pots")
	case "shuffle_crates":
		panic("not implemented for shuffle_crates")
	case "shuffle_dungeon_rewards":
		panic("not implemented for shuffle_dungeon_rewards")
	case "shuffle_mapcompass":
		panic("not implemented for shuffle_mapcompass")
	case "shuffle_smallkeys":
		panic("not implemented for shuffle_smallkeys")
	case "shuffle_hideoutkeys":
		panic("not implemented for shuffle_hideoutkeys")
	case "shuffle_interior_entrances":
		panic("not implemented for shuffle_interior_entrances")
	case "shuffle_gerudo_fortress_heart_piece":
		panic("not implemented for shuffle_gerudo_fortress_heart_piece")
	case "shuffle_tcgkeys":
		panic("not implemented for shuffle_tcgkeys")
	case "key_rings_choice":
		panic("not implemented for key_rings_choice")
	case "shuffle_bosskeys":
		panic("not implemented for shuffle_bosskeys")
	case "shuffle_silver_rupees":
		panic("not implemented for shuffle_silver_rupees")
	case "silver_rupee_pouches_choice":
		panic("not implemented for silver_rupee_pouches_choice")
	case "hints":
		panic("not implemented for hints")
	case "hint_dist":
		panic("not implemented for hint_dist")
	case "correct_chest_appearances":
		panic("not implemented for correct_chest_appearances")
	case "ocarina_songs":
		panic("not implemented for ocarina_songs")
	case "damage_multiplier":
		panic("not implemented for damage_multiplier")
	case "deadly_bonks":
		panic("not implemented for deadly_bonks")
	case "starting_tod":
		panic("not implemented for starting_tod")
	case "item_pool_value":
		panic("not implemented for item_pool_value")
	case "junk_ice_traps":
		panic("not implemented for junk_ice_traps")
	case "ice_trap_appearance":
		panic("not implemented for ice_trap_appearance")

	default:
		return val, unknown(name)
	}
}

func (this *Zootr) Float64(name string) (float64, error) {
	var val float64
	switch name {
	case "triforce_count_per_world":
		panic("not implemented for triforce_count_per_world")
	case "triforce_goal_per_world":
		panic("not implemented for triforce_goal_per_world")
	case "lacs_medallions":
		panic("not implemented for lacs_medallions")
	case "lacs_stones":
		panic("not implemented for lacs_stones")
	case "lacs_rewards":
		panic("not implemented for lacs_rewards")
	case "lacs_tokens":
		panic("not implemented for lacs_tokens")
	case "lacs_hearts":
		panic("not implemented for lacs_hearts")
	case "bridge_medallions":
		panic("not implemented for bridge_medallions")
	case "bridge_stones":
		panic("not implemented for bridge_stones")
	case "bridge_rewards":
		panic("not implemented for bridge_rewards")
	case "bridge_tokens":
		panic("not implemented for bridge_tokens")
	case "bridge_hearts":
		panic("not implemented for bridge_hearts")
	case "trials":
		panic("not implemented for trials")
	case "ganon_bosskey_stones":
		panic("not implemented for ganon_bosskey_stones")
	case "ganon_bosskey_rewards":
		panic("not implemented for ganon_bosskey_rewards")
	case "ganon_bosskey_tokens":
		panic("not implemented for ganon_bosskey_tokens")
	case "ganon_bosskey_hearts":
		panic("not implemented for ganon_bosskey_hearts")
	case "mq_dungeons_count":
		panic("not implemented for mq_dungeons_count")
	case "empty_dungeons_count":
		panic("not implemented for empty_dungeons_count")
	case "starting_hearts":
		panic("not implemented for starting_hearts")
	case "fae_torch_count":
		panic("not implemented for fae_torch_count")
	case "chicken_count":
		panic("not implemented for chicken_count")
	case "big_poe_count":
		panic("not implemented for big_poe_count")
	case "custom_ice_trap_percent":
		panic("not implemented for custom_ice_trap_percent")
	case "custom_ice_trap_count":
		panic("not implemented for custom_ice_trap_count")

	default:
		return val, unknown(name)
	}

	return val, nil
}

func (this *Zootr) Bool(name string) (bool, error) {
	var val bool
	switch name {
	case "trials_random":
		panic("not implemented for trials_random")
	case "triforce_hunt":
		panic("not implemented for triforce_hunt")
	case "open_door_of_time":
		panic("not implemented for open_door_of_time")
	case "shuffle_hideout_entrances":
		panic("not implemented for shuffle_hideout_entrances")
	case "shuffle_grotto_entrances":
		panic("not implemented for shuffle_grotto_entrances")
	case "shuffle_ganon_tower":
		panic("not implemented for shuffle_ganon_tower")
	case "shuffle_overworld_entrances":
		panic("not implemented for shuffle_overworld_entrances")
	case "shuffle_gerudo_valley_river_exit":
		panic("not implemented for shuffle_gerudo_valley_river_exit")
	case "owl_drops":
		panic("not implemented for owl_drops")
	case "free_bombchu_drops":
		panic("not implemented for free_bombchu_drops")
	case "warp_songs":
		panic("not implemented for warp_songs")
	case "adult_trade_shuffle":
		panic("not implemented for adult_trade_shuffle")
	case "shuffle_empty_pots":
		panic("not implemented for shuffle_empty_pots")
	case "shuffle_empty_crates":
		panic("not implemented for shuffle_empty_crates")
	case "shuffle_cows":
		panic("not implemented for shuffle_cows")
	case "shuffle_beehives":
		panic("not implemented for shuffle_beehives")
	case "shuffle_wonderitems":
		panic("not implemented for shuffle_wonderitems")
	case "shuffle_kokiri_sword":
		panic("not implemented for shuffle_kokiri_sword")
	case "shuffle_ocarinas":
		panic("not implemented for shuffle_ocarinas")
	case "shuffle_gerudo_card":
		panic("not implemented for shuffle_gerudo_card")
	case "shuffle_beans":
		panic("not implemented for shuffle_beans")
	case "shuffle_expensive_merchants":
		panic("not implemented for shuffle_expensive_merchants")
	case "shuffle_frog_song_rupees":
		panic("not implemented for shuffle_frog_song_rupees")
	case "shuffle_loach_reward":
		panic("not implemented for shuffle_loach_reward")
	case "shuffle_individual_ocarina_notes":
		panic("not implemented for shuffle_individual_ocarina_notes")
	case "keyring_give_bk":
		panic("not implemented for keyring_give_bk")
	case "enhance_map_compass":
		panic("not implemented for enhance_map_compass")
	case "start_with_consumables":
		panic("not implemented for start_with_consumables")
	case "start_with_rupees":
		panic("not implemented for start_with_rupees")
	case "skip_reward_from_rauru":
		panic("not implemented for skip_reward_from_rauru")
	case "no_escape_sequence":
		panic("not implemented for no_escape_sequence")
	case "no_guard_stealth":
		panic("not implemented for no_guard_stealth")
	case "no_epona_race":
		panic("not implemented for no_epona_race")
	case "skip_some_minigame_phases":
		panic("not implemented for skip_some_minigame_phases")
	case "complete_mask_quest":
		panic("not implemented for complete_mask_quest")
	case "useful_cutscenes":
		panic("not implemented for useful_cutscenes")
	case "fast_chests":
		panic("not implemented for fast_chests")
	case "free_scarecrow":
		panic("not implemented for free_scarecrow")
	case "fast_bunny_hood":
		panic("not implemented for fast_bunny_hood")
	case "auto_equip_masks":
		panic("not implemented for auto_equip_masks")
	case "plant_beans":
		panic("not implemented for plant_beans")
	case "easier_fire_arrow_entry":
		panic("not implemented for easier_fire_arrow_entry")
	case "ruto_already_f1_jabu":
		panic("not implemented for ruto_already_f1_jabu")
	case "fast_shadow_boat":
		panic("not implemented for fast_shadow_boat")
	case "chicken_count_random":
		panic("not implemented for chicken_count_random")
	case "clearer_hints":
		panic("not implemented for clearer_hints")
	case "blue_fire_arrows":
		panic("not implemented for blue_fire_arrows")
	case "fix_broken_drops":
		panic("not implemented for fix_broken_drops")
	case "tcg_requires_lens":
		panic("not implemented for tcg_requires_lens")
	case "no_collectible_hearts":
		panic("not implemented for no_collectible_hearts")
	case "one_item_per_dungeon":
		panic("not implemented for one_item_per_dungeon")
	default:
		return val, unknown(name)
	}
}

func unknown(name string) error {
	return fmt.Errorf("%q is not a known setting", name)
}
