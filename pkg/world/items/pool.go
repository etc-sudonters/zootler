package items

import "math"
import "sudonters/zootler/pkg/world/settings"

func BuildItemPool(settings.SeedSettings) map[string]int {
	return DefaultItemPool()
}

func DefaultItemPool() map[string]int {
	return map[string]int{
		"Arrows (10)":                          8,
		"Arrows (30)":                          6,
		"Arrows (5)":                           3,
		"Biggoron Sword":                       1,
		"Bolero of Fire":                       1,
		"Bomb Bag":                             3,
		"Bombchus (10)":                        3,
		"Bombchus (20)":                        1,
		"Bombchus (5)":                         1,
		"Bombs (10)":                           2,
		"Bombs (20)":                           2,
		"Bombs (5)":                            8,
		"Boomerang":                            1,
		"Bottle with Blue Potion":              1,
		"Bottle with Red Potion":               2,
		"Bow":                                  3,
		"Claim Check":                          1,
		"Deku Nut Capacity":                    2,
		"Deku Nuts (10)":                       1,
		"Deku Nuts (5)":                        4,
		"Deku Seeds (30)":                      4,
		"Deku Shield":                          4,
		"Deku Stick (1)":                       3,
		"Deku Stick Capacity":                  2,
		"Dins Fire":                            1,
		"Double Defense":                       1,
		"Eponas Song":                          1,
		"Farores Wind":                         1,
		"Fire Arrows":                          1,
		"Goron Tunic":                          1,
		"Heart Container":                      8,
		"Hover Boots":                          1,
		"Hylian Shield":                        2,
		"Ice Arrows":                           1,
		"Iron Boots":                           1,
		"Kokiri Sword":                         1,
		"Lens of Truth":                        1,
		"Light Arrows":                         1,
		"Magic Meter":                          2,
		"Megaton Hammer":                       1,
		"Minuet of Forest":                     1,
		"Mirror Shield":                        1,
		"Nayrus Love":                          1,
		"Nocturne of Shadow":                   1,
		"Piece of Heart (Treasure Chest Game)": 1,
		"Piece of Heart":                       35,
		"Prelude of Light":                     1,
		"Progressive Hookshot":                 2,
		"Progressive Scale":                    2,
		"Progressive Strength Upgrade":         3,
		"Progressive Wallet":                   2,
		"Recovery Heart":                       11,
		"Requiem of Spirit":                    1,
		"Rupee (1)":                            1,
		"Rupees (20)":                          6,
		"Rupees (200)":                         6,
		"Rupees (5)":                           23,
		"Rupees (50)":                          7,
		"Rutos Letter":                         1,
		"Sarias Song":                          1,
		"Serenade of Water":                    1,
		"Slingshot":                            3,
		"Song of Storms":                       1,
		"Song of Time":                         1,
		"Stone of Agony":                       1,
		"Suns Song":                            1,
		"Zeldas Lullaby":                       1,
		"Zora Tunic":                           1,
	}
}

type GetItemId int

const (
	GI_MISSING                                                GetItemId = -1
	GI_NONE                                                             = 0x0000
	GI_BOMBS_5                                                          = 0x0001
	GI_DEKU_NUTS_5                                                      = 0x0002
	GI_BOMBCHUS_10                                                      = 0x0003
	GI_BOW                                                              = 0x0004
	GI_SLINGSHOT                                                        = 0x0005
	GI_BOOMERANG                                                        = 0x0006
	GI_DEKU_STICKS_1                                                    = 0x0007
	GI_HOOKSHOT                                                         = 0x0008
	GI_LONGSHOT                                                         = 0x0009
	GI_LENS_OF_TRUTH                                                    = 0x000A
	GI_ZELDAS_LETTER                                                    = 0x000B
	GI_OCARINA_OF_TIME                                                  = 0x000C
	GI_HAMMER                                                           = 0x000D
	GI_COJIRO                                                           = 0x000E
	GI_BOTTLE_EMPTY                                                     = 0x000F
	GI_BOTTLE_POTION_RED                                                = 0x0010
	GI_BOTTLE_POTION_GREEN                                              = 0x0011
	GI_BOTTLE_POTION_BLUE                                               = 0x0012
	GI_BOTTLE_FAIRY                                                     = 0x0013
	GI_BOTTLE_MILK_FULL                                                 = 0x0014
	GI_BOTTLE_RUTOS_LETTER                                              = 0x0015
	GI_MAGIC_BEAN                                                       = 0x0016
	GI_MASK_SKULL                                                       = 0x0017
	GI_MASK_SPOOKY                                                      = 0x0018
	GI_CHICKEN                                                          = 0x0019
	GI_MASK_KEATON                                                      = 0x001A
	GI_MASK_BUNNY_HOOD                                                  = 0x001B
	GI_MASK_TRUTH                                                       = 0x001C
	GI_POCKET_EGG                                                       = 0x001D
	GI_POCKET_CUCCO                                                     = 0x001E
	GI_ODD_MUSHROOM                                                     = 0x001F
	GI_ODD_POTION                                                       = 0x0020
	GI_POACHERS_SAW                                                     = 0x0021
	GI_BROKEN_GORONS_SWORD                                              = 0x0022
	GI_PRESCRIPTION                                                     = 0x0023
	GI_EYEBALL_FROG                                                     = 0x0024
	GI_EYE_DROPS                                                        = 0x0025
	GI_CLAIM_CHECK                                                      = 0x0026
	GI_SWORD_KOKIRI                                                     = 0x0027
	GI_SWORD_KNIFE                                                      = 0x0028
	GI_SHIELD_DEKU                                                      = 0x0029
	GI_SHIELD_HYLIAN                                                    = 0x002A
	GI_SHIELD_MIRROR                                                    = 0x002B
	GI_TUNIC_GORON                                                      = 0x002C
	GI_TUNIC_ZORA                                                       = 0x002D
	GI_BOOTS_IRON                                                       = 0x002E
	GI_BOOTS_HOVER                                                      = 0x002F
	GI_QUIVER_40                                                        = 0x0030
	GI_QUIVER_50                                                        = 0x0031
	GI_BOMB_BAG_20                                                      = 0x0032
	GI_BOMB_BAG_30                                                      = 0x0033
	GI_BOMB_BAG_40                                                      = 0x0034
	GI_SILVER_GAUNTLETS                                                 = 0x0035
	GI_GOLD_GAUNTLETS                                                   = 0x0036
	GI_SCALE_SILVER                                                     = 0x0037
	GI_SCALE_GOLDEN                                                     = 0x0038
	GI_STONE_OF_AGONY                                                   = 0x0039
	GI_GERUDOS_CARD                                                     = 0x003A
	GI_OCARINA_FAIRY                                                    = 0x003B
	GI_DEKU_SEEDS_5                                                     = 0x003C
	GI_HEART_CONTAINER                                                  = 0x003D
	GI_HEART_PIECE                                                      = 0x003E
	GI_BOSS_KEY                                                         = 0x003F
	GI_COMPASS                                                          = 0x0040
	GI_DUNGEON_MAP                                                      = 0x0041
	GI_SMALL_KEY                                                        = 0x0042
	GI_MAGIC_JAR_SMALL                                                  = 0x0043
	GI_MAGIC_JAR_LARGE                                                  = 0x0044
	GI_WALLET_ADULT                                                     = 0x0045
	GI_WALLET_GIANT                                                     = 0x0046
	GI_WEIRD_EGG                                                        = 0x0047
	GI_RECOVERY_HEART                                                   = 0x0048
	GI_ARROWS_5                                                         = 0x0049
	GI_ARROWS_10                                                        = 0x004A
	GI_ARROWS_30                                                        = 0x004B
	GI_RUPEE_GREEN                                                      = 0x004C
	GI_RUPEE_BLUE                                                       = 0x004D
	GI_RUPEE_RED                                                        = 0x004E
	GI_HEART_CONTAINER_2                                                = 0x004F
	GI_MILK                                                             = 0x0050
	GI_MASK_GORON                                                       = 0x0051
	GI_MASK_ZORA                                                        = 0x0052
	GI_MASK_GERUDO                                                      = 0x0053
	GI_GORONS_BRACELET                                                  = 0x0054
	GI_RUPEE_PURPLE                                                     = 0x0055
	GI_RUPEE_GOLD                                                       = 0x0056
	GI_SWORD_BIGGORON                                                   = 0x0057
	GI_ARROW_FIRE                                                       = 0x0058
	GI_ARROW_ICE                                                        = 0x0059
	GI_ARROW_LIGHT                                                      = 0x005A
	GI_SKULL_TOKEN                                                      = 0x005B
	GI_DINS_FIRE                                                        = 0x005C
	GI_FARORES_WIND                                                     = 0x005D
	GI_NAYRUS_LOVE                                                      = 0x005E
	GI_BULLET_BAG_30                                                    = 0x005F
	GI_BULLET_BAG_40                                                    = 0x0060
	GI_DEKU_STICKS_5                                                    = 0x0061
	GI_DEKU_STICKS_10                                                   = 0x0062
	GI_DEKU_NUTS_5_2                                                    = 0x0063
	GI_DEKU_NUTS_10                                                     = 0x0064
	GI_BOMBS_1                                                          = 0x0065
	GI_BOMBS_10                                                         = 0x0066
	GI_BOMBS_20                                                         = 0x0067
	GI_BOMBS_30                                                         = 0x0068
	GI_DEKU_SEEDS_30                                                    = 0x0069
	GI_BOMBCHUS_5                                                       = 0x006A
	GI_BOMBCHUS_20                                                      = 0x006B
	GI_BOTTLE_FISH                                                      = 0x006C
	GI_BOTTLE_BUGS                                                      = 0x006D
	GI_BOTTLE_BLUE_FIRE                                                 = 0x006E
	GI_BOTTLE_POE                                                       = 0x006F
	GI_BOTTLE_BIG_POE                                                   = 0x0070
	GI_DOOR_KEY                                                         = 0x0071
	GI_RUPEE_GREEN_LOSE                                                 = 0x0072
	GI_RUPEE_BLUE_LOSE                                                  = 0x0073
	GI_RUPEE_RED_LOSE                                                   = 0x0074
	GI_RUPEE_PURPLE_LOSE                                                = 0x0075
	GI_HEART_PIECE_WIN                                                  = 0x0076
	GI_DEKU_STICK_UPGRADE_20                                            = 0x0077
	GI_DEKU_STICK_UPGRADE_30                                            = 0x0078
	GI_DEKU_NUT_UPGRADE_30                                              = 0x0079
	GI_DEKU_NUT_UPGRADE_40                                              = 0x007A
	GI_BULLET_BAG_50                                                    = 0x007B
	GI_ICE_TRAP                                                         = 0x007C
	GI_TEXT_0                                                           = 0x007D
	GI_CAPPED_PIECE_OF_HEART                                            = 0x007D
	GI_VANILLA_MAX                                                      = 0x007E
	GI_CAPPED_HEART_CONTAINER                                           = 0x007E
	GI_CAPPED_PIECE_OF_HEART_CHESTGAME                                  = 0x007F
	GI_PROGRESSIVE_HOOKSHOT                                             = 0x0080
	GI_PROGRESSIVE_STRENGTH                                             = 0x0081
	GI_PROGRESSIVE_BOMB_BAG                                             = 0x0082
	GI_PROGRESSIVE_BOW                                                  = 0x0083
	GI_PROGRESSIVE_SLINGSHOT                                            = 0x0084
	GI_PROGRESSIVE_WALLET                                               = 0x0085
	GI_PROGRESSIVE_SCALE                                                = 0x0086
	GI_PROGRESSIVE_NUT_CAPACITY                                         = 0x0087
	GI_PROGRESSIVE_STICK_CAPACITY                                       = 0x0088
	GI_PROGRESSIVE_BOMBCHUS                                             = 0x0089
	GI_PROGRESSIVE_MAGIC_METER                                          = 0x008A
	GI_PROGRESSIVE_OCARINA                                              = 0x008B
	GI_BOTTLE_WITH_RED_POTION                                           = 0x008C
	GI_BOTTLE_WITH_GREEN_POTION                                         = 0x008D
	GI_BOTTLE_WITH_BLUE_POTION                                          = 0x008E
	GI_BOTTLE_WITH_FAIRY                                                = 0x008F
	GI_BOTTLE_WITH_FISH                                                 = 0x0090
	GI_BOTTLE_WITH_BLUE_FIRE                                            = 0x0091
	GI_BOTTLE_WITH_BUGS                                                 = 0x0092
	GI_BOTTLE_WITH_BIG_POE                                              = 0x0093
	GI_BOTTLE_WITH_POE                                                  = 0x0094
	GI_BOSS_KEY_FOREST_TEMPLE                                           = 0x0095
	GI_BOSS_KEY_FIRE_TEMPLE                                             = 0x0096
	GI_BOSS_KEY_WATER_TEMPLE                                            = 0x0097
	GI_BOSS_KEY_SPIRIT_TEMPLE                                           = 0x0098
	GI_BOSS_KEY_SHADOW_TEMPLE                                           = 0x0099
	GI_BOSS_KEY_GANONS_CASTLE                                           = 0x009A
	GI_COMPASS_DEKU_TREE                                                = 0x009B
	GI_COMPASS_DODONGOS_CAVERN                                          = 0x009C
	GI_COMPASS_JABU_JABU                                                = 0x009D
	GI_COMPASS_FOREST_TEMPLE                                            = 0x009E
	GI_COMPASS_FIRE_TEMPLE                                              = 0x009F
	GI_COMPASS_WATER_TEMPLE                                             = 0x00A0
	GI_COMPASS_SPIRIT_TEMPLE                                            = 0x00A1
	GI_COMPASS_SHADOW_TEMPLE                                            = 0x00A2
	GI_COMPASS_BOTTOM_OF_THE_WELL                                       = 0x00A3
	GI_COMPASS_ICE_CAVERN                                               = 0x00A4
	GI_MAP_DEKU_TREE                                                    = 0x00A5
	GI_MAP_DODONGOS_CAVERN                                              = 0x00A6
	GI_MAP_JABU_JABU                                                    = 0x00A7
	GI_MAP_FOREST_TEMPLE                                                = 0x00A8
	GI_MAP_FIRE_TEMPLE                                                  = 0x00A9
	GI_MAP_WATER_TEMPLE                                                 = 0x00AA
	GI_MAP_SPIRIT_TEMPLE                                                = 0x00AB
	GI_MAP_SHADOW_TEMPLE                                                = 0x00AC
	GI_MAP_BOTTOM_OF_THE_WELL                                           = 0x00AD
	GI_MAP_ICE_CAVERN                                                   = 0x00AE
	GI_SMALL_KEY_FOREST_TEMPLE                                          = 0x00AF
	GI_SMALL_KEY_FIRE_TEMPLE                                            = 0x00B0
	GI_SMALL_KEY_WATER_TEMPLE                                           = 0x00B1
	GI_SMALL_KEY_SPIRIT_TEMPLE                                          = 0x00B2
	GI_SMALL_KEY_SHADOW_TEMPLE                                          = 0x00B3
	GI_SMALL_KEY_BOTTOM_OF_THE_WELL                                     = 0x00B4
	GI_SMALL_KEY_GERUDO_TRAINING                                        = 0x00B5
	GI_SMALL_KEY_THIEVES_HIDEOUT                                        = 0x00B6
	GI_SMALL_KEY_GANONS_CASTLE                                          = 0x00B7
	GI_DOUBLE_DEFENSE                                                   = 0x00B8
	GI_MAGIC_METER                                                      = 0x00B9
	GI_DOUBLE_MAGIC                                                     = 0x00BA
	GI_MINUET_OF_FOREST                                                 = 0x00BB
	GI_BOLERO_OF_FIRE                                                   = 0x00BC
	GI_SERENADE_OF_WATER                                                = 0x00BD
	GI_REQUIEM_OF_SPIRIT                                                = 0x00BE
	GI_NOCTURNE_OF_SHADOW                                               = 0x00BF
	GI_PRELUDE_OF_LIGHT                                                 = 0x00C0
	GI_ZELDAS_LULLABY                                                   = 0x00C1
	GI_EPONAS_SONG                                                      = 0x00C2
	GI_SARIAS_SONG                                                      = 0x00C3
	GI_SUNS_SONG                                                        = 0x00C4
	GI_SONG_OF_TIME                                                     = 0x00C5
	GI_SONG_OF_STORMS                                                   = 0x00C6
	GI_TYCOONS_WALLET                                                   = 0x00C7
	GI_REDUNDANT_LETTER_BOTTLE                                          = 0x00C8
	GI_MAGIC_BEAN_PACK                                                  = 0x00C9
	GI_TRIFORCE_PIECE                                                   = 0x00CA
	GI_SMALL_KEY_RING_FOREST_TEMPLE                                     = 0x00CB
	GI_SMALL_KEY_RING_FIRE_TEMPLE                                       = 0x00CC
	GI_SMALL_KEY_RING_WATER_TEMPLE                                      = 0x00CD
	GI_SMALL_KEY_RING_SPIRIT_TEMPLE                                     = 0x00CE
	GI_SMALL_KEY_RING_SHADOW_TEMPLE                                     = 0x00CF
	GI_SMALL_KEY_RING_BOTTOM_OF_THE_WELL                                = 0x00D0
	GI_SMALL_KEY_RING_GERUDO_TRAINING                                   = 0x00D1
	GI_SMALL_KEY_RING_THIEVES_HIDEOUT                                   = 0x00D2
	GI_SMALL_KEY_RING_GANONS_CASTLE                                     = 0x00D3
	GI_BOMBCHU_BAG_20                                                   = 0x00D4
	GI_BOMBCHU_BAG_10                                                   = 0x00D5
	GI_BOMBCHU_BAG_5                                                    = 0x00D6
	GI_SMALL_KEY_RING_TREASURE_CHEST_GAME                               = 0x00D7
	GI_SILVER_RUPEE_DODONGOS_CAVERN_STAIRCASE                           = 0x00D8
	GI_SILVER_RUPEE_ICE_CAVERN_SPINNING_SCYTHE                          = 0x00D9
	GI_SILVER_RUPEE_ICE_CAVERN_PUSH_BLOCK                               = 0x00DA
	GI_SILVER_RUPEE_BOTTOM_OF_THE_WELL_BASEMENT                         = 0x00DB
	GI_SILVER_RUPEE_SHADOW_TEMPLE_SCYTHE_SHORTCUT                       = 0x00DC
	GI_SILVER_RUPEE_SHADOW_TEMPLE_INVISIBLE_BLADES                      = 0x00DD
	GI_SILVER_RUPEE_SHADOW_TEMPLE_HUGE_PIT                              = 0x00DE
	GI_SILVER_RUPEE_SHADOW_TEMPLE_INVISIBLE_SPIKES                      = 0x00DF
	GI_SILVER_RUPEE_GERUDO_TRAINING_GROUND_SLOPES                       = 0x00E0
	GI_SILVER_RUPEE_GERUDO_TRAINING_GROUND_LAVA                         = 0x00E1
	GI_SILVER_RUPEE_GERUDO_TRAINING_GROUND_WATER                        = 0x00E2
	GI_SILVER_RUPEE_SPIRIT_TEMPLE_CHILD_EARLY_TORCHES                   = 0x00E3
	GI_SILVER_RUPEE_SPIRIT_TEMPLE_ADULT_BOULDERS                        = 0x00E4
	GI_SILVER_RUPEE_SPIRIT_TEMPLE_LOBBY_AND_LOWER_ADULT                 = 0x00E5
	GI_SILVER_RUPEE_SPIRIT_TEMPLE_SUN_BLOCK                             = 0x00E6
	GI_SILVER_RUPEE_SPIRIT_TEMPLE_ADULT_CLIMB                           = 0x00E7
	GI_SILVER_RUPEE_GANONS_CASTLE_SPIRIT_TRIAL                          = 0x00E8
	GI_SILVER_RUPEE_GANONS_CASTLE_LIGHT_TRIAL                           = 0x00E9
	GI_SILVER_RUPEE_GANONS_CASTLE_FIRE_TRIAL                            = 0x00EA
	GI_SILVER_RUPEE_GANONS_CASTLE_SHADOW_TRIAL                          = 0x00EB
	GI_SILVER_RUPEE_GANONS_CASTLE_WATER_TRIAL                           = 0x00EC
	GI_SILVER_RUPEE_GANONS_CASTLE_FOREST_TRIAL                          = 0x00ED
	GI_SILVER_RUPEE_POUCH_DODONGOS_CAVERN_STAIRCASE                     = 0x00EE
	GI_SILVER_RUPEE_POUCH_ICE_CAVERN_SPINNING_SCYTHE                    = 0x00EF
	GI_SILVER_RUPEE_POUCH_ICE_CAVERN_PUSH_BLOCK                         = 0x00F0
	GI_SILVER_RUPEE_POUCH_BOTTOM_OF_THE_WELL_BASEMENT                   = 0x00F1
	GI_SILVER_RUPEE_POUCH_SHADOW_TEMPLE_SCYTHE_SHORTCUT                 = 0x00F2
	GI_SILVER_RUPEE_POUCH_SHADOW_TEMPLE_INVISIBLE_BLADES                = 0x00F3
	GI_SILVER_RUPEE_POUCH_SHADOW_TEMPLE_HUGE_PIT                        = 0x00F4
	GI_SILVER_RUPEE_POUCH_SHADOW_TEMPLE_INVISIBLE_SPIKES                = 0x00F5
	GI_SILVER_RUPEE_POUCH_GERUDO_TRAINING_GROUND_SLOPES                 = 0x00F6
	GI_SILVER_RUPEE_POUCH_GERUDO_TRAINING_GROUND_LAVA                   = 0x00F7
	GI_SILVER_RUPEE_POUCH_GERUDO_TRAINING_GROUND_WATER                  = 0x00F8
	GI_SILVER_RUPEE_POUCH_SPIRIT_TEMPLE_CHILD_EARLY_TORCHES             = 0x00F9
	GI_SILVER_RUPEE_POUCH_SPIRIT_TEMPLE_ADULT_BOULDERS                  = 0x00FA
	GI_SILVER_RUPEE_POUCH_SPIRIT_TEMPLE_LOBBY_AND_LOWER_ADULT           = 0x00FB
	GI_SILVER_RUPEE_POUCH_SPIRIT_TEMPLE_SUN_BLOCK                       = 0x00FC
	GI_SILVER_RUPEE_POUCH_SPIRIT_TEMPLE_ADULT_CLIMB                     = 0x00FD
	GI_SILVER_RUPEE_POUCH_GANONS_CASTLE_SPIRIT_TRIAL                    = 0x00FE
	GI_SILVER_RUPEE_POUCH_GANONS_CASTLE_LIGHT_TRIAL                     = 0x00FF
	GI_SILVER_RUPEE_POUCH_GANONS_CASTLE_FIRE_TRIAL                      = 0x0100
	GI_SILVER_RUPEE_POUCH_GANONS_CASTLE_SHADOW_TRIAL                    = 0x0101
	GI_SILVER_RUPEE_POUCH_GANONS_CASTLE_WATER_TRIAL                     = 0x0102
	GI_SILVER_RUPEE_POUCH_GANONS_CASTLE_FOREST_TRIAL                    = 0x0103
	GI_OCARINA_BUTTON_A                                                 = 0x0104
	GI_OCARINA_BUTTON_C_UP                                              = 0x0105
	GI_OCARINA_BUTTON_C_DOWN                                            = 0x0106
	GI_OCARINA_BUTTON_C_LEFT                                            = 0x0107
	GI_OCARINA_BUTTON_C_RIGHT                                           = 0x0108
	GI_BOSS_KEY_MODEL_FOREST_TEMPLE                                     = 0x0109
	GI_BOSS_KEY_MODEL_FIRE_TEMPLE                                       = 0x010A
	GI_BOSS_KEY_MODEL_WATER_TEMPLE                                      = 0x010B
	GI_BOSS_KEY_MODEL_SPIRIT_TEMPLE                                     = 0x010C
	GI_BOSS_KEY_MODEL_SHADOW_TEMPLE                                     = 0x010D
	GI_BOSS_KEY_MODEL_GANONS_CASTLE                                     = 0x010E
	GI_SMALL_KEY_MODEL_FOREST_TEMPLE                                    = 0x010F
	GI_SMALL_KEY_MODEL_FIRE_TEMPLE                                      = 0x0110
	GI_SMALL_KEY_MODEL_WATER_TEMPLE                                     = 0x0111
	GI_SMALL_KEY_MODEL_SPIRIT_TEMPLE                                    = 0x0112
	GI_SMALL_KEY_MODEL_SHADOW_TEMPLE                                    = 0x0113
	GI_SMALL_KEY_MODEL_BOTTOM_OF_THE_WELL                               = 0x0114
	GI_SMALL_KEY_MODEL_GERUDO_TRAINING                                  = 0x0115
	GI_SMALL_KEY_MODEL_THIEVES_HIDEOUT                                  = 0x0116
	GI_SMALL_KEY_MODEL_GANONS_CASTLE                                    = 0x0117
	GI_SMALL_KEY_MODEL_CHEST_GAME                                       = 0x0118
	GI_RANDO_MAX                                                        = 0x0119
)

type Priority int

const (
	PriorityNormal Priority = iota
	PriorityMajor
	PriorityAdvancement
)

type rawtoken struct {
	name     string
	priority Priority
	itemId   GetItemId
	special  map[string]any
}

var item_table = map[string]any{
	"Bombs (5)":                            rawtoken{"Item", PriorityNormal, GI_BOMBS_5, map[string]any{"junk": 8}},
	"Deku Nuts (5)":                        rawtoken{"Item", PriorityNormal, GI_DEKU_NUTS_5, map[string]any{"junk": 5}},
	"Bombchus (10)":                        rawtoken{"Item", PriorityAdvancement, GI_BOMBCHUS_10, nil},
	"Boomerang":                            rawtoken{"Item", PriorityAdvancement, GI_BOOMERANG, nil},
	"Deku Stick (1)":                       rawtoken{"Item", PriorityNormal, GI_DEKU_STICKS_1, map[string]any{"junk": 5}},
	"Lens of Truth":                        rawtoken{"Item", PriorityAdvancement, GI_LENS_OF_TRUTH, nil},
	"Megaton Hammer":                       rawtoken{"Item", PriorityAdvancement, GI_HAMMER, nil},
	"Cojiro":                               rawtoken{"Item", PriorityAdvancement, GI_COJIRO, map[string]any{"trade": true}},
	"Bottle":                               rawtoken{"Item", PriorityAdvancement, GI_BOTTLE_EMPTY, map[string]any{"bottle": math.Inf(1)}},
	"Blue Potion":                          rawtoken{"Item", PriorityAdvancement, GI_BOTTLE_POTION_BLUE, nil}, // distinct from shop item
	"Bottle with Milk":                     rawtoken{"Item", PriorityAdvancement, GI_BOTTLE_MILK_FULL, map[string]any{"bottle": math.Inf(1)}},
	"Rutos Letter":                         rawtoken{"Item", PriorityAdvancement, GI_BOTTLE_RUTOS_LETTER, nil},
	"Deliver Letter":                       rawtoken{"Item", PriorityAdvancement, GI_MISSING, map[string]any{"bottle": math.Inf(1)}},
	"Sell Big Poe":                         rawtoken{"Item", PriorityAdvancement, GI_MISSING, map[string]any{"bottle": math.Inf(1)}},
	"Magic Bean":                           rawtoken{"Item", PriorityAdvancement, GI_MAGIC_BEAN, map[string]any{"progressive": 10}},
	"Skull Mask":                           rawtoken{"Item", PriorityAdvancement, GI_MASK_SKULL, map[string]any{"trade": true, "object": 0x0136}},
	"Spooky Mask":                          rawtoken{"Item", PriorityAdvancement, GI_MASK_SPOOKY, map[string]any{"trade": true, "object": 0x0135}},
	"Chicken":                              rawtoken{"Item", PriorityAdvancement, GI_CHICKEN, map[string]any{"trade": true}},
	"Keaton Mask":                          rawtoken{"Item", PriorityAdvancement, GI_MASK_KEATON, map[string]any{"trade": true, "object": 0x0134}},
	"Bunny Hood":                           rawtoken{"Item", PriorityAdvancement, GI_MASK_BUNNY_HOOD, map[string]any{"trade": true, "object": 0x0137}},
	"Mask of Truth":                        rawtoken{"Item", PriorityAdvancement, GI_MASK_TRUTH, map[string]any{"trade": true, "object": 0x0138}},
	"Pocket Egg":                           rawtoken{"Item", PriorityAdvancement, GI_POCKET_EGG, map[string]any{"trade": true}},
	"Pocket Cucco":                         rawtoken{"Item", PriorityAdvancement, GI_POCKET_CUCCO, map[string]any{"trade": true}},
	"Odd Mushroom":                         rawtoken{"Item", PriorityAdvancement, GI_ODD_MUSHROOM, map[string]any{"trade": true}},
	"Odd Potion":                           rawtoken{"Item", PriorityAdvancement, GI_ODD_POTION, map[string]any{"trade": true}},
	"Poachers Saw":                         rawtoken{"Item", PriorityAdvancement, GI_POACHERS_SAW, map[string]any{"trade": true}},
	"Broken Sword":                         rawtoken{"Item", PriorityAdvancement, GI_BROKEN_GORONS_SWORD, map[string]any{"trade": true}},
	"Prescription":                         rawtoken{"Item", PriorityAdvancement, GI_PRESCRIPTION, map[string]any{"trade": true}},
	"Eyeball Frog":                         rawtoken{"Item", PriorityAdvancement, GI_EYEBALL_FROG, map[string]any{"trade": true}},
	"Eyedrops":                             rawtoken{"Item", PriorityAdvancement, GI_EYE_DROPS, map[string]any{"trade": true}},
	"Claim Check":                          rawtoken{"Item", PriorityAdvancement, GI_CLAIM_CHECK, map[string]any{"trade": true}},
	"Kokiri Sword":                         rawtoken{"Item", PriorityAdvancement, GI_SWORD_KOKIRI, nil},
	"Giants Knife":                         rawtoken{"Item", PriorityNormal, GI_SWORD_KNIFE, nil},
	"Deku Shield":                          rawtoken{"Item", PriorityNormal, GI_SHIELD_DEKU, nil},
	"Hylian Shield":                        rawtoken{"Item", PriorityNormal, GI_SHIELD_HYLIAN, nil},
	"Mirror Shield":                        rawtoken{"Item", PriorityAdvancement, GI_SHIELD_MIRROR, nil},
	"Goron Tunic":                          rawtoken{"Item", PriorityAdvancement, GI_TUNIC_GORON, nil},
	"Zora Tunic":                           rawtoken{"Item", PriorityAdvancement, GI_TUNIC_ZORA, nil},
	"Iron Boots":                           rawtoken{"Item", PriorityAdvancement, GI_BOOTS_IRON, nil},
	"Hover Boots":                          rawtoken{"Item", PriorityAdvancement, GI_BOOTS_HOVER, nil},
	"Stone of Agony":                       rawtoken{"Item", PriorityAdvancement, GI_STONE_OF_AGONY, nil},
	"Gerudo Membership Card":               rawtoken{"Item", PriorityAdvancement, GI_GERUDOS_CARD, nil},
	"Heart Container":                      rawtoken{"Item", PriorityAdvancement, GI_HEART_CONTAINER, map[string]any{"alias": map[string]int{"Piece of Heart": 4}, "progressive": math.Inf(1)}},
	"Piece of Heart":                       rawtoken{"Item", PriorityAdvancement, GI_HEART_PIECE, map[string]any{"progressive": math.Inf(1)}},
	"Boss Key":                             rawtoken{"BossKey", PriorityAdvancement, GI_BOSS_KEY, nil},
	"Compass":                              rawtoken{"Compass", PriorityNormal, GI_COMPASS, nil},
	"Map":                                  rawtoken{"Map", PriorityNormal, GI_DUNGEON_MAP, nil},
	"Small Key":                            rawtoken{"SmallKey", PriorityAdvancement, GI_SMALL_KEY, map[string]any{"progressive": math.Inf(1)}},
	"Weird Egg":                            rawtoken{"Item", PriorityAdvancement, GI_WEIRD_EGG, map[string]any{"trade": true}},
	"Recovery Heart":                       rawtoken{"Item", PriorityNormal, GI_RECOVERY_HEART, map[string]any{"junk": 0}},
	"Arrows (5)":                           rawtoken{"Item", PriorityNormal, GI_ARROWS_5, map[string]any{"junk": 8}},
	"Arrows (10)":                          rawtoken{"Item", PriorityNormal, GI_ARROWS_10, map[string]any{"junk": 2}},
	"Arrows (30)":                          rawtoken{"Item", PriorityNormal, GI_ARROWS_30, map[string]any{"junk": 0}},
	"Rupee (1)":                            rawtoken{"Item", PriorityNormal, GI_RUPEE_GREEN, map[string]any{"junk": -1}},
	"Rupees (5)":                           rawtoken{"Item", PriorityNormal, GI_RUPEE_BLUE, map[string]any{"junk": 10}},
	"Rupees (20)":                          rawtoken{"Item", PriorityNormal, GI_RUPEE_RED, map[string]any{"junk": 4}},
	"Milk":                                 rawtoken{"Item", PriorityNormal, GI_MILK, nil},
	"Goron Mask":                           rawtoken{"Item", PriorityNormal, GI_MASK_GORON, map[string]any{"trade": true, "object": 0x0150}},
	"Zora Mask":                            rawtoken{"Item", PriorityNormal, GI_MASK_ZORA, map[string]any{"trade": true, "object": 0x0151}},
	"Gerudo Mask":                          rawtoken{"Item", PriorityNormal, GI_MASK_GERUDO, map[string]any{"trade": true, "object": 0x0152}},
	"Rupees (50)":                          rawtoken{"Item", PriorityNormal, GI_RUPEE_PURPLE, map[string]any{"junk": 1}},
	"Rupees (200)":                         rawtoken{"Item", PriorityNormal, GI_RUPEE_GOLD, map[string]any{"junk": 0}},
	"Biggoron Sword":                       rawtoken{"Item", PriorityNormal, GI_SWORD_BIGGORON, nil},
	"Fire Arrows":                          rawtoken{"Item", PriorityAdvancement, GI_ARROW_FIRE, nil},
	"Ice Arrows":                           rawtoken{"Item", PriorityAdvancement, GI_ARROW_ICE, nil},
	"Blue Fire Arrows":                     rawtoken{"Item", PriorityAdvancement, GI_ARROW_ICE, nil},
	"Light Arrows":                         rawtoken{"Item", PriorityAdvancement, GI_ARROW_LIGHT, nil},
	"Gold Skulltula Token":                 rawtoken{"Token", PriorityAdvancement, GI_SKULL_TOKEN, map[string]any{"progressive": math.Inf(1)}},
	"Dins Fire":                            rawtoken{"Item", PriorityAdvancement, GI_DINS_FIRE, nil},
	"Farores Wind":                         rawtoken{"Item", PriorityAdvancement, GI_FARORES_WIND, nil},
	"Nayrus Love":                          rawtoken{"Item", PriorityAdvancement, GI_NAYRUS_LOVE, nil},
	"Deku Nuts (10)":                       rawtoken{"Item", PriorityNormal, GI_DEKU_NUTS_10, map[string]any{"junk": 0}},
	"Bomb (1)":                             rawtoken{"Item", PriorityNormal, GI_BOMBS_1, map[string]any{"junk": -1}},
	"Bombs (10)":                           rawtoken{"Item", PriorityNormal, GI_BOMBS_10, map[string]any{"junk": 2}},
	"Bombs (20)":                           rawtoken{"Item", PriorityNormal, GI_BOMBS_20, map[string]any{"junk": 0}},
	"Deku Seeds (30)":                      rawtoken{"Item", PriorityNormal, GI_DEKU_SEEDS_30, map[string]any{"junk": 5}},
	"Bombchus (5)":                         rawtoken{"Item", PriorityAdvancement, GI_BOMBCHUS_5, nil},
	"Bombchus (20)":                        rawtoken{"Item", PriorityAdvancement, GI_BOMBCHUS_20, nil},
	"Small Key (Treasure Chest Game)":      rawtoken{"TCGSmallKey", PriorityAdvancement, GI_DOOR_KEY, map[string]any{"progressive": math.Inf(1)}},
	"Rupee (Treasure Chest Game) (1)":      rawtoken{"Item", PriorityNormal, GI_RUPEE_GREEN_LOSE, nil},
	"Rupees (Treasure Chest Game) (5)":     rawtoken{"Item", PriorityNormal, GI_RUPEE_BLUE_LOSE, nil},
	"Rupees (Treasure Chest Game) (20)":    rawtoken{"Item", PriorityNormal, GI_RUPEE_RED_LOSE, nil},
	"Rupees (Treasure Chest Game) (50)":    rawtoken{"Item", PriorityNormal, GI_RUPEE_PURPLE_LOSE, nil},
	"Piece of Heart (Treasure Chest Game)": rawtoken{"Item", PriorityAdvancement, GI_HEART_PIECE_WIN, map[string]any{"alias": map[string]int{"Piece of Heart": 1}, "progressive": math.Inf(1)}},
	"Ice Trap":                             rawtoken{"Item", PriorityNormal, GI_ICE_TRAP, map[string]any{"junk": 0}},
	"Progressive Hookshot":                 rawtoken{"Item", PriorityAdvancement, GI_PROGRESSIVE_HOOKSHOT, map[string]any{"progressive": 2}},
	"Progressive Strength Upgrade":         rawtoken{"Item", PriorityAdvancement, GI_PROGRESSIVE_STRENGTH, map[string]any{"progressive": 3}},
	"Bomb Bag":                             rawtoken{"Item", PriorityAdvancement, GI_PROGRESSIVE_BOMB_BAG, nil},
	"Bow":                                  rawtoken{"Item", PriorityAdvancement, GI_PROGRESSIVE_BOW, nil},
	"Slingshot":                            rawtoken{"Item", PriorityAdvancement, GI_PROGRESSIVE_SLINGSHOT, nil},
	"Progressive Wallet":                   rawtoken{"Item", PriorityAdvancement, GI_PROGRESSIVE_WALLET, map[string]any{"progressive": 3}},
	"Progressive Scale":                    rawtoken{"Item", PriorityAdvancement, GI_PROGRESSIVE_SCALE, map[string]any{"progressive": 2}},
	"Deku Nut Capacity":                    rawtoken{"Item", PriorityNormal, GI_PROGRESSIVE_NUT_CAPACITY, nil},
	"Deku Stick Capacity":                  rawtoken{"Item", PriorityNormal, GI_PROGRESSIVE_STICK_CAPACITY, nil},
	"Bombchus":                             rawtoken{"Item", PriorityAdvancement, GI_PROGRESSIVE_BOMBCHUS, nil},
	"Magic Meter":                          rawtoken{"Item", PriorityAdvancement, GI_PROGRESSIVE_MAGIC_METER, nil},
	"Ocarina":                              rawtoken{"Item", PriorityAdvancement, GI_PROGRESSIVE_OCARINA, nil},
	"Bottle with Red Potion":               rawtoken{"Item", PriorityAdvancement, GI_BOTTLE_WITH_RED_POTION, map[string]any{"bottle": true, "shop_object": 0x0F}},
	"Bottle with Green Potion":             rawtoken{"Item", PriorityAdvancement, GI_BOTTLE_WITH_GREEN_POTION, map[string]any{"bottle": true, "shop_object": 0x0F}},
	"Bottle with Blue Potion":              rawtoken{"Item", PriorityAdvancement, GI_BOTTLE_WITH_BLUE_POTION, map[string]any{"bottle": true, "shop_object": 0x0F}},
	"Bottle with Fairy":                    rawtoken{"Item", PriorityAdvancement, GI_BOTTLE_WITH_FAIRY, map[string]any{"bottle": true, "shop_object": 0x0F}},
	"Bottle with Fish":                     rawtoken{"Item", PriorityAdvancement, GI_BOTTLE_WITH_FISH, map[string]any{"bottle": true, "shop_object": 0x0F}},
	"Bottle with Blue Fire":                rawtoken{"Item", PriorityAdvancement, GI_BOTTLE_WITH_BLUE_FIRE, map[string]any{"bottle": true, "shop_object": 0x0F}},
	"Bottle with Bugs":                     rawtoken{"Item", PriorityAdvancement, GI_BOTTLE_WITH_BUGS, map[string]any{"bottle": true, "shop_object": 0x0F}},
	"Bottle with Big Poe":                  rawtoken{"Item", PriorityAdvancement, GI_BOTTLE_WITH_BIG_POE, map[string]any{"shop_object": 0x0F}},
	"Bottle with Poe":                      rawtoken{"Item", PriorityAdvancement, GI_BOTTLE_WITH_POE, map[string]any{"bottle": true, "shop_object": 0x0F}},
	"Boss Key (Forest Temple)":             rawtoken{"BossKey", PriorityAdvancement, GI_BOSS_KEY_FOREST_TEMPLE, nil},
	"Boss Key (Fire Temple)":               rawtoken{"BossKey", PriorityAdvancement, GI_BOSS_KEY_FIRE_TEMPLE, nil},
	"Boss Key (Water Temple)":              rawtoken{"BossKey", PriorityAdvancement, GI_BOSS_KEY_WATER_TEMPLE, nil},
	"Boss Key (Spirit Temple)":             rawtoken{"BossKey", PriorityAdvancement, GI_BOSS_KEY_SPIRIT_TEMPLE, nil},
	"Boss Key (Shadow Temple)":             rawtoken{"BossKey", PriorityAdvancement, GI_BOSS_KEY_SHADOW_TEMPLE, nil},
	"Boss Key (Ganons Castle)":             rawtoken{"GanonBossKey", PriorityAdvancement, GI_BOSS_KEY_GANONS_CASTLE, nil},
	"Compass (Deku Tree)":                  rawtoken{"Compass", PriorityMajor, GI_COMPASS_DEKU_TREE, nil},
	"Compass (Dodongos Cavern)":            rawtoken{"Compass", PriorityMajor, GI_COMPASS_DODONGOS_CAVERN, nil},
	"Compass (Jabu Jabus Belly)":           rawtoken{"Compass", PriorityMajor, GI_COMPASS_JABU_JABU, nil},
	"Compass (Forest Temple)":              rawtoken{"Compass", PriorityMajor, GI_COMPASS_FOREST_TEMPLE, nil},
	"Compass (Fire Temple)":                rawtoken{"Compass", PriorityMajor, GI_COMPASS_FIRE_TEMPLE, nil},
	"Compass (Water Temple)":               rawtoken{"Compass", PriorityMajor, GI_COMPASS_WATER_TEMPLE, nil},
	"Compass (Spirit Temple)":              rawtoken{"Compass", PriorityMajor, GI_COMPASS_SPIRIT_TEMPLE, nil},
	"Compass (Shadow Temple)":              rawtoken{"Compass", PriorityMajor, GI_COMPASS_SHADOW_TEMPLE, nil},
	"Compass (Bottom of the Well)":         rawtoken{"Compass", PriorityMajor, GI_COMPASS_BOTTOM_OF_THE_WELL, nil},
	"Compass (Ice Cavern)":                 rawtoken{"Compass", PriorityMajor, GI_COMPASS_ICE_CAVERN, nil},
	"Map (Deku Tree)":                      rawtoken{"Map", PriorityMajor, GI_MAP_DEKU_TREE, nil},
	"Map (Dodongos Cavern)":                rawtoken{"Map", PriorityMajor, GI_MAP_DODONGOS_CAVERN, nil},
	"Map (Jabu Jabus Belly)":               rawtoken{"Map", PriorityMajor, GI_MAP_JABU_JABU, nil},
	"Map (Forest Temple)":                  rawtoken{"Map", PriorityMajor, GI_MAP_FOREST_TEMPLE, nil},
	"Map (Fire Temple)":                    rawtoken{"Map", PriorityMajor, GI_MAP_FIRE_TEMPLE, nil},
	"Map (Water Temple)":                   rawtoken{"Map", PriorityMajor, GI_MAP_WATER_TEMPLE, nil},
	"Map (Spirit Temple)":                  rawtoken{"Map", PriorityMajor, GI_MAP_SPIRIT_TEMPLE, nil},
	"Map (Shadow Temple)":                  rawtoken{"Map", PriorityMajor, GI_MAP_SHADOW_TEMPLE, nil},
	"Map (Bottom of the Well)":             rawtoken{"Map", PriorityMajor, GI_MAP_BOTTOM_OF_THE_WELL, nil},
	"Map (Ice Cavern)":                     rawtoken{"Map", PriorityMajor, GI_MAP_ICE_CAVERN, nil},
	"Small Key (Forest Temple)":            rawtoken{"SmallKey", PriorityAdvancement, GI_SMALL_KEY_FOREST_TEMPLE, map[string]any{"progressive": math.Inf(1)}},
	"Small Key (Fire Temple)":              rawtoken{"SmallKey", PriorityAdvancement, GI_SMALL_KEY_FIRE_TEMPLE, map[string]any{"progressive": math.Inf(1)}},
	"Small Key (Water Temple)":             rawtoken{"SmallKey", PriorityAdvancement, GI_SMALL_KEY_WATER_TEMPLE, map[string]any{"progressive": math.Inf(1)}},
	"Small Key (Spirit Temple)":            rawtoken{"SmallKey", PriorityAdvancement, GI_SMALL_KEY_SPIRIT_TEMPLE, map[string]any{"progressive": math.Inf(1)}},
	"Small Key (Shadow Temple)":            rawtoken{"SmallKey", PriorityAdvancement, GI_SMALL_KEY_SHADOW_TEMPLE, map[string]any{"progressive": math.Inf(1)}},
	"Small Key (Bottom of the Well)":       rawtoken{"SmallKey", PriorityAdvancement, GI_SMALL_KEY_BOTTOM_OF_THE_WELL, map[string]any{"progressive": math.Inf(1)}},
	"Small Key (Gerudo Training Ground)":   rawtoken{"SmallKey", PriorityAdvancement, GI_SMALL_KEY_GERUDO_TRAINING, map[string]any{"progressive": math.Inf(1)}},
	"Small Key (Thieves Hideout)":          rawtoken{"HideoutSmallKey", PriorityAdvancement, GI_SMALL_KEY_THIEVES_HIDEOUT, map[string]any{"progressive": math.Inf(1)}},
	"Small Key (Ganons Castle)":            rawtoken{"SmallKey", PriorityAdvancement, GI_SMALL_KEY_GANONS_CASTLE, map[string]any{"progressive": math.Inf(1)}},
	"Double Defense":                       rawtoken{"Item", PriorityNormal, GI_DOUBLE_DEFENSE, nil},
	"Buy Magic Bean":                       rawtoken{"Item", PriorityAdvancement, GI_MAGIC_BEAN, map[string]any{"alias": map[string]int{"Magic Bean": 10}, "progressive": 10}},
	"Magic Bean Pack":                      rawtoken{"Item", PriorityAdvancement, GI_MAGIC_BEAN_PACK, map[string]any{"alias": map[string]int{"Magic Bean": 10}, "progressive": 10}},
	"Triforce Piece":                       rawtoken{"Item", PriorityAdvancement, GI_TRIFORCE_PIECE, map[string]any{"progressive": math.Inf(1)}},
	"Zeldas Letter":                        rawtoken{"Item", PriorityAdvancement, GI_ZELDAS_LETTER, map[string]any{"trade": true}},
	"Time Travel":                          rawtoken{"Event", PriorityAdvancement, GI_MISSING, nil},
	"Scarecrow Song":                       rawtoken{"Event", PriorityAdvancement, GI_MISSING, nil},
	"Triforce":                             rawtoken{"Event", PriorityAdvancement, GI_MISSING, nil},

	"Small Key Ring (Forest Temple)":          rawtoken{"SmallKey", PriorityAdvancement, GI_SMALL_KEY_RING_FOREST_TEMPLE, map[string]any{"alias": map[string]int{"Small Key (Forest Temple)": 10}, "progressive": math.Inf(1)}},
	"Small Key Ring (Fire Temple)":            rawtoken{"SmallKey", PriorityAdvancement, GI_SMALL_KEY_RING_FIRE_TEMPLE, map[string]any{"alias": map[string]int{"Small Key (Fire Temple)": 10}, "progressive": math.Inf(1)}},
	"Small Key Ring (Water Temple)":           rawtoken{"SmallKey", PriorityAdvancement, GI_SMALL_KEY_RING_WATER_TEMPLE, map[string]any{"alias": map[string]int{"Small Key (Water Temple)": 10}, "progressive": math.Inf(1)}},
	"Small Key Ring (Spirit Temple)":          rawtoken{"SmallKey", PriorityAdvancement, GI_SMALL_KEY_RING_SPIRIT_TEMPLE, map[string]any{"alias": map[string]int{"Small Key (Spirit Temple)": 10}, "progressive": math.Inf(1)}},
	"Small Key Ring (Shadow Temple)":          rawtoken{"SmallKey", PriorityAdvancement, GI_SMALL_KEY_RING_SHADOW_TEMPLE, map[string]any{"alias": map[string]int{"Small Key (Shadow Temple)": 10}, "progressive": math.Inf(1)}},
	"Small Key Ring (Bottom of the Well)":     rawtoken{"SmallKey", PriorityAdvancement, GI_SMALL_KEY_RING_BOTTOM_OF_THE_WELL, map[string]any{"alias": map[string]int{"Small Key (Bottom of the Well)": 10}, "progressive": math.Inf(1)}},
	"Small Key Ring (Gerudo Training Ground)": rawtoken{"SmallKey", PriorityAdvancement, GI_SMALL_KEY_RING_GERUDO_TRAINING, map[string]any{"alias": map[string]int{"Small Key (Gerudo Training Ground)": 10}, "progressive": math.Inf(1)}},
	"Small Key Ring (Thieves Hideout)":        rawtoken{"HideoutSmallKey", PriorityAdvancement, GI_SMALL_KEY_RING_THIEVES_HIDEOUT, map[string]any{"alias": map[string]int{"Small Key (Thieves Hideout)": 10}, "progressive": math.Inf(1)}},
	"Small Key Ring (Ganons Castle)":          rawtoken{"SmallKey", PriorityAdvancement, GI_SMALL_KEY_RING_GANONS_CASTLE, map[string]any{"alias": map[string]int{"Small Key (Ganons Castle)": 10}, "progressive": math.Inf(1)}},
	"Small Key Ring (Treasure Chest Game)":    rawtoken{"TCGSmallKey", PriorityAdvancement, GI_SMALL_KEY_RING_TREASURE_CHEST_GAME, map[string]any{"alias": map[string]int{"Small Key (Treasure Chest Game)": 10}, "progressive": math.Inf(1)}},

	"Silver Rupee (Dodongos Cavern Staircase)":           rawtoken{"SilverRupee", PriorityAdvancement, GI_SILVER_RUPEE_DODONGOS_CAVERN_STAIRCASE, map[string]any{"progressive": 5}},
	"Silver Rupee (Ice Cavern Spinning Scythe)":          rawtoken{"SilverRupee", PriorityAdvancement, GI_SILVER_RUPEE_ICE_CAVERN_SPINNING_SCYTHE, map[string]any{"progressive": 5}},
	"Silver Rupee (Ice Cavern Push Block)":               rawtoken{"SilverRupee", PriorityAdvancement, GI_SILVER_RUPEE_ICE_CAVERN_PUSH_BLOCK, map[string]any{"progressive": 5}},
	"Silver Rupee (Bottom of the Well Basement)":         rawtoken{"SilverRupee", PriorityAdvancement, GI_SILVER_RUPEE_BOTTOM_OF_THE_WELL_BASEMENT, map[string]any{"progressive": 5}},
	"Silver Rupee (Shadow Temple Scythe Shortcut)":       rawtoken{"SilverRupee", PriorityAdvancement, GI_SILVER_RUPEE_SHADOW_TEMPLE_SCYTHE_SHORTCUT, map[string]any{"progressive": 5}},
	"Silver Rupee (Shadow Temple Invisible Blades)":      rawtoken{"SilverRupee", PriorityAdvancement, GI_SILVER_RUPEE_SHADOW_TEMPLE_INVISIBLE_BLADES, map[string]any{"progressive": 10}},
	"Silver Rupee (Shadow Temple Huge Pit)":              rawtoken{"SilverRupee", PriorityAdvancement, GI_SILVER_RUPEE_SHADOW_TEMPLE_HUGE_PIT, map[string]any{"progressive": 5}},
	"Silver Rupee (Shadow Temple Invisible Spikes)":      rawtoken{"SilverRupee", PriorityAdvancement, GI_SILVER_RUPEE_SHADOW_TEMPLE_INVISIBLE_SPIKES, map[string]any{"progressive": 10}},
	"Silver Rupee (Gerudo Training Ground Slopes)":       rawtoken{"SilverRupee", PriorityAdvancement, GI_SILVER_RUPEE_GERUDO_TRAINING_GROUND_SLOPES, map[string]any{"progressive": 5}},
	"Silver Rupee (Gerudo Training Ground Lava)":         rawtoken{"SilverRupee", PriorityAdvancement, GI_SILVER_RUPEE_GERUDO_TRAINING_GROUND_LAVA, map[string]any{"progressive": 6}},
	"Silver Rupee (Gerudo Training Ground Water)":        rawtoken{"SilverRupee", PriorityAdvancement, GI_SILVER_RUPEE_GERUDO_TRAINING_GROUND_WATER, map[string]any{"progressive": 5}},
	"Silver Rupee (Spirit Temple Child Early Torches)":   rawtoken{"SilverRupee", PriorityAdvancement, GI_SILVER_RUPEE_SPIRIT_TEMPLE_CHILD_EARLY_TORCHES, map[string]any{"progressive": 5}},
	"Silver Rupee (Spirit Temple Adult Boulders)":        rawtoken{"SilverRupee", PriorityAdvancement, GI_SILVER_RUPEE_SPIRIT_TEMPLE_ADULT_BOULDERS, map[string]any{"progressive": 5}},
	"Silver Rupee (Spirit Temple Lobby and Lower Adult)": rawtoken{"SilverRupee", PriorityAdvancement, GI_SILVER_RUPEE_SPIRIT_TEMPLE_LOBBY_AND_LOWER_ADULT, map[string]any{"progressive": 5}},
	"Silver Rupee (Spirit Temple Sun Block)":             rawtoken{"SilverRupee", PriorityAdvancement, GI_SILVER_RUPEE_SPIRIT_TEMPLE_SUN_BLOCK, map[string]any{"progressive": 5}},
	"Silver Rupee (Spirit Temple Adult Climb)":           rawtoken{"SilverRupee", PriorityAdvancement, GI_SILVER_RUPEE_SPIRIT_TEMPLE_ADULT_CLIMB, map[string]any{"progressive": 5}},
	"Silver Rupee (Ganons Castle Spirit Trial)":          rawtoken{"SilverRupee", PriorityAdvancement, GI_SILVER_RUPEE_GANONS_CASTLE_SPIRIT_TRIAL, map[string]any{"progressive": 5}},
	"Silver Rupee (Ganons Castle Light Trial)":           rawtoken{"SilverRupee", PriorityAdvancement, GI_SILVER_RUPEE_GANONS_CASTLE_LIGHT_TRIAL, map[string]any{"progressive": 5}},
	"Silver Rupee (Ganons Castle Fire Trial)":            rawtoken{"SilverRupee", PriorityAdvancement, GI_SILVER_RUPEE_GANONS_CASTLE_FIRE_TRIAL, map[string]any{"progressive": 5}},
	"Silver Rupee (Ganons Castle Shadow Trial)":          rawtoken{"SilverRupee", PriorityAdvancement, GI_SILVER_RUPEE_GANONS_CASTLE_SHADOW_TRIAL, map[string]any{"progressive": 5}},
	"Silver Rupee (Ganons Castle Water Trial)":           rawtoken{"SilverRupee", PriorityAdvancement, GI_SILVER_RUPEE_GANONS_CASTLE_WATER_TRIAL, map[string]any{"progressive": 5}},
	"Silver Rupee (Ganons Castle Forest Trial)":          rawtoken{"SilverRupee", PriorityAdvancement, GI_SILVER_RUPEE_GANONS_CASTLE_FOREST_TRIAL, map[string]any{"progressive": 5}},

	"Silver Rupee Pouch (Dodongos Cavern Staircase)":           rawtoken{"SilverRupee", PriorityAdvancement, GI_SILVER_RUPEE_POUCH_DODONGOS_CAVERN_STAIRCASE, map[string]any{"alias": map[string]int{"Silver Rupee (Dodongos Cavern Staircase)": 10}, "progressive": 1}},
	"Silver Rupee Pouch (Ice Cavern Spinning Scythe)":          rawtoken{"SilverRupee", PriorityAdvancement, GI_SILVER_RUPEE_POUCH_ICE_CAVERN_SPINNING_SCYTHE, map[string]any{"alias": map[string]int{"Silver Rupee (Ice Cavern Spinning Scythe)": 10}, "progressive": 1}},
	"Silver Rupee Pouch (Ice Cavern Push Block)":               rawtoken{"SilverRupee", PriorityAdvancement, GI_SILVER_RUPEE_POUCH_ICE_CAVERN_PUSH_BLOCK, map[string]any{"alias": map[string]int{"Silver Rupee (Ice Cavern Push Block)": 10}, "progressive": 1}},
	"Silver Rupee Pouch (Bottom of the Well Basement)":         rawtoken{"SilverRupee", PriorityAdvancement, GI_SILVER_RUPEE_POUCH_BOTTOM_OF_THE_WELL_BASEMENT, map[string]any{"alias": map[string]int{"Silver Rupee (Bottom of the Well Basement)": 10}, "progressive": 1}},
	"Silver Rupee Pouch (Shadow Temple Scythe Shortcut)":       rawtoken{"SilverRupee", PriorityAdvancement, GI_SILVER_RUPEE_POUCH_SHADOW_TEMPLE_SCYTHE_SHORTCUT, map[string]any{"alias": map[string]int{"Silver Rupee (Shadow Temple Scythe Shortcut)": 10}, "progressive": 1}},
	"Silver Rupee Pouch (Shadow Temple Invisible Blades)":      rawtoken{"SilverRupee", PriorityAdvancement, GI_SILVER_RUPEE_POUCH_SHADOW_TEMPLE_INVISIBLE_BLADES, map[string]any{"alias": map[string]int{"Silver Rupee (Shadow Temple Invisible Blades)": 10}, "progressive": 1}},
	"Silver Rupee Pouch (Shadow Temple Huge Pit)":              rawtoken{"SilverRupee", PriorityAdvancement, GI_SILVER_RUPEE_POUCH_SHADOW_TEMPLE_HUGE_PIT, map[string]any{"alias": map[string]int{"Silver Rupee (Shadow Temple Huge Pit)": 10}, "progressive": 1}},
	"Silver Rupee Pouch (Shadow Temple Invisible Spikes)":      rawtoken{"SilverRupee", PriorityAdvancement, GI_SILVER_RUPEE_POUCH_SHADOW_TEMPLE_INVISIBLE_SPIKES, map[string]any{"alias": map[string]int{"Silver Rupee (Shadow Temple Invisible Spikes)": 10}, "progressive": 1}},
	"Silver Rupee Pouch (Gerudo Training Ground Slopes)":       rawtoken{"SilverRupee", PriorityAdvancement, GI_SILVER_RUPEE_POUCH_GERUDO_TRAINING_GROUND_SLOPES, map[string]any{"alias": map[string]int{"Silver Rupee (Gerudo Training Ground Slopes)": 10}, "progressive": 1}},
	"Silver Rupee Pouch (Gerudo Training Ground Lava)":         rawtoken{"SilverRupee", PriorityAdvancement, GI_SILVER_RUPEE_POUCH_GERUDO_TRAINING_GROUND_LAVA, map[string]any{"alias": map[string]int{"Silver Rupee (Gerudo Training Ground Lava)": 10}, "progressive": 1}},
	"Silver Rupee Pouch (Gerudo Training Ground Water)":        rawtoken{"SilverRupee", PriorityAdvancement, GI_SILVER_RUPEE_POUCH_GERUDO_TRAINING_GROUND_WATER, map[string]any{"alias": map[string]int{"Silver Rupee (Gerudo Training Ground Water)": 10}, "progressive": 1}},
	"Silver Rupee Pouch (Spirit Temple Child Early Torches)":   rawtoken{"SilverRupee", PriorityAdvancement, GI_SILVER_RUPEE_POUCH_SPIRIT_TEMPLE_CHILD_EARLY_TORCHES, map[string]any{"alias": map[string]int{"Silver Rupee (Spirit Temple Child Early Torches)": 10}, "progressive": 1}},
	"Silver Rupee Pouch (Spirit Temple Adult Boulders)":        rawtoken{"SilverRupee", PriorityAdvancement, GI_SILVER_RUPEE_POUCH_SPIRIT_TEMPLE_ADULT_BOULDERS, map[string]any{"alias": map[string]int{"Silver Rupee (Spirit Temple Adult Boulders)": 10}, "progressive": 1}},
	"Silver Rupee Pouch (Spirit Temple Lobby and Lower Adult)": rawtoken{"SilverRupee", PriorityAdvancement, GI_SILVER_RUPEE_POUCH_SPIRIT_TEMPLE_LOBBY_AND_LOWER_ADULT, map[string]any{"alias": map[string]int{"Silver Rupee (Spirit Temple Lobby and Lower Adult)": 10}, "progressive": 1}},
	"Silver Rupee Pouch (Spirit Temple Sun Block)":             rawtoken{"SilverRupee", PriorityAdvancement, GI_SILVER_RUPEE_POUCH_SPIRIT_TEMPLE_SUN_BLOCK, map[string]any{"alias": map[string]int{"Silver Rupee (Spirit Temple Sun Block)": 10}, "progressive": 1}},
	"Silver Rupee Pouch (Spirit Temple Adult Climb)":           rawtoken{"SilverRupee", PriorityAdvancement, GI_SILVER_RUPEE_POUCH_SPIRIT_TEMPLE_ADULT_CLIMB, map[string]any{"alias": map[string]int{"Silver Rupee (Spirit Temple Adult Climb)": 10}, "progressive": 1}},
	"Silver Rupee Pouch (Ganons Castle Spirit Trial)":          rawtoken{"SilverRupee", PriorityAdvancement, GI_SILVER_RUPEE_POUCH_GANONS_CASTLE_SPIRIT_TRIAL, map[string]any{"alias": map[string]int{"Silver Rupee (Ganons Castle Spirit Trial)": 10}, "progressive": 1}},
	"Silver Rupee Pouch (Ganons Castle Light Trial)":           rawtoken{"SilverRupee", PriorityAdvancement, GI_SILVER_RUPEE_POUCH_GANONS_CASTLE_LIGHT_TRIAL, map[string]any{"alias": map[string]int{"Silver Rupee (Ganons Castle Light Trial)": 10}, "progressive": 1}},
	"Silver Rupee Pouch (Ganons Castle Fire Trial)":            rawtoken{"SilverRupee", PriorityAdvancement, GI_SILVER_RUPEE_POUCH_GANONS_CASTLE_FIRE_TRIAL, map[string]any{"alias": map[string]int{"Silver Rupee (Ganons Castle Fire Trial)": 10}, "progressive": 1}},
	"Silver Rupee Pouch (Ganons Castle Shadow Trial)":          rawtoken{"SilverRupee", PriorityAdvancement, GI_SILVER_RUPEE_POUCH_GANONS_CASTLE_SHADOW_TRIAL, map[string]any{"alias": map[string]int{"Silver Rupee (Ganons Castle Shadow Trial)": 10}, "progressive": 1}},
	"Silver Rupee Pouch (Ganons Castle Water Trial)":           rawtoken{"SilverRupee", PriorityAdvancement, GI_SILVER_RUPEE_POUCH_GANONS_CASTLE_WATER_TRIAL, map[string]any{"alias": map[string]int{"Silver Rupee (Ganons Castle Water Trial)": 10}, "progressive": 1}},
	"Silver Rupee Pouch (Ganons Castle Forest Trial)":          rawtoken{"SilverRupee", PriorityAdvancement, GI_SILVER_RUPEE_POUCH_GANONS_CASTLE_FOREST_TRIAL, map[string]any{"alias": map[string]int{"Silver Rupee (Ganons Castle Forest Trial)": 10}, "progressive": 1}},

	"Ocarina A Button":       rawtoken{"Item", PriorityAdvancement, GI_OCARINA_BUTTON_A, map[string]any{"ocarina_button": true}},
	"Ocarina C up Button":    rawtoken{"Item", PriorityAdvancement, GI_OCARINA_BUTTON_C_UP, map[string]any{"ocarina_button": true}},
	"Ocarina C down Button":  rawtoken{"Item", PriorityAdvancement, GI_OCARINA_BUTTON_C_DOWN, map[string]any{"ocarina_button": true}},
	"Ocarina C left Button":  rawtoken{"Item", PriorityAdvancement, GI_OCARINA_BUTTON_C_LEFT, map[string]any{"ocarina_button": true}},
	"Ocarina C right Button": rawtoken{"Item", PriorityAdvancement, GI_OCARINA_BUTTON_C_RIGHT, map[string]any{"ocarina_button": true}},

	// Event items otherwise generated by generic event logic
	// can be defined here to enforce their appearance in playthroughs.
	"Water Temple Clear": rawtoken{"Event", PriorityAdvancement, GI_MISSING, nil},
	"Forest Trial Clear": rawtoken{"Event", PriorityAdvancement, GI_MISSING, nil},
	"Fire Trial Clear":   rawtoken{"Event", PriorityAdvancement, GI_MISSING, nil},
	"Water Trial Clear":  rawtoken{"Event", PriorityAdvancement, GI_MISSING, nil},
	"Shadow Trial Clear": rawtoken{"Event", PriorityAdvancement, GI_MISSING, nil},
	"Spirit Trial Clear": rawtoken{"Event", PriorityAdvancement, GI_MISSING, nil},
	"Light Trial Clear":  rawtoken{"Event", PriorityAdvancement, GI_MISSING, nil},
	"Epona":              rawtoken{"Event", PriorityAdvancement, GI_MISSING, nil},

	"Deku Stick Drop":  rawtoken{"Drop", PriorityAdvancement, GI_MISSING, nil},
	"Deku Nut Drop":    rawtoken{"Drop", PriorityAdvancement, GI_MISSING, nil},
	"Blue Fire":        rawtoken{"Drop", PriorityAdvancement, GI_MISSING, nil},
	"Fairy":            rawtoken{"Drop", PriorityAdvancement, GI_MISSING, nil},
	"Fish":             rawtoken{"Drop", PriorityAdvancement, GI_MISSING, nil},
	"Bugs":             rawtoken{"Drop", PriorityAdvancement, GI_MISSING, nil},
	"Big Poe":          rawtoken{"Drop", PriorityAdvancement, GI_MISSING, nil},
	"Bombchu Drop":     rawtoken{"Drop", PriorityAdvancement, GI_MISSING, nil},
	"Deku Shield Drop": rawtoken{"Drop", PriorityAdvancement, GI_MISSING, nil},

	// Consumable refills defined mostly to placate "starting with" options
	"Arrows":      rawtoken{"Refill", PriorityNormal, GI_MISSING, nil},
	"Bombs":       rawtoken{"Refill", PriorityNormal, GI_MISSING, nil},
	"Deku Seeds":  rawtoken{"Refill", PriorityNormal, GI_MISSING, nil},
	"Deku Sticks": rawtoken{"Refill", PriorityNormal, GI_MISSING, nil},
	"Deku Nuts":   rawtoken{"Refill", PriorityNormal, GI_MISSING, nil},
	"Rupees":      rawtoken{"Refill", PriorityNormal, GI_MISSING, nil},

	"Minuet of Forest":   rawtoken{"Song", PriorityAdvancement, GI_MINUET_OF_FOREST, map[string]any{"text_id": 0x73, "song_id": 0x02, "item_id": 0x5A}},
	"Bolero of Fire":     rawtoken{"Song", PriorityAdvancement, GI_BOLERO_OF_FIRE, map[string]any{"text_id": 0x74, "song_id": 0x03, "item_id": 0x5B}},
	"Serenade of Water":  rawtoken{"Song", PriorityAdvancement, GI_SERENADE_OF_WATER, map[string]any{"text_id": 0x75, "song_id": 0x04, "item_id": 0x5C}},
	"Requiem of Spirit":  rawtoken{"Song", PriorityAdvancement, GI_REQUIEM_OF_SPIRIT, map[string]any{"text_id": 0x76, "song_id": 0x05, "item_id": 0x5D}},
	"Nocturne of Shadow": rawtoken{"Song", PriorityAdvancement, GI_NOCTURNE_OF_SHADOW, map[string]any{"text_id": 0x77, "song_id": 0x06, "item_id": 0x5E}},
	"Prelude of Light":   rawtoken{"Song", PriorityAdvancement, GI_PRELUDE_OF_LIGHT, map[string]any{"text_id": 0x78, "song_id": 0x07, "item_id": 0x5F}},
	"Zeldas Lullaby":     rawtoken{"Song", PriorityAdvancement, GI_ZELDAS_LULLABY, map[string]any{"text_id": 0xD4, "song_id": 0x0A, "item_id": 0x60}},
	"Eponas Song":        rawtoken{"Song", PriorityAdvancement, GI_EPONAS_SONG, map[string]any{"text_id": 0xD2, "song_id": 0x09, "item_id": 0x61}},
	"Sarias Song":        rawtoken{"Song", PriorityAdvancement, GI_SARIAS_SONG, map[string]any{"text_id": 0xD1, "song_id": 0x08, "item_id": 0x62}},
	"Suns Song":          rawtoken{"Song", PriorityAdvancement, GI_SUNS_SONG, map[string]any{"text_id": 0xD3, "song_id": 0x0B, "item_id": 0x63}},
	"Song of Time":       rawtoken{"Song", PriorityAdvancement, GI_SONG_OF_TIME, map[string]any{"text_id": 0xD5, "song_id": 0x0C, "item_id": 0x64}},
	"Song of Storms":     rawtoken{"Song", PriorityAdvancement, GI_SONG_OF_STORMS, map[string]any{"text_id": 0xD6, "song_id": 0x0D, "item_id": 0x65}},

	// shop is weird and has its own rules for these values
	"Buy Deku Nut (5)":             rawtoken{"Shop", PriorityAdvancement, 0x00, map[string]any{"object": 0x00BB, "price": 15}},
	"Buy Arrows (30)":              rawtoken{"Shop", PriorityMajor, 0x01, map[string]any{"object": 0x00D8, "price": 60}},
	"Buy Arrows (50)":              rawtoken{"Shop", PriorityMajor, 0x02, map[string]any{"object": 0x00D8, "price": 90}},
	"Buy Bombs (5) for 25 Rupees":  rawtoken{"Shop", PriorityMajor, 0x03, map[string]any{"object": 0x00CE, "price": 25}},
	"Buy Deku Nut (10)":            rawtoken{"Shop", PriorityAdvancement, 0x04, map[string]any{"object": 0x00BB, "price": 30}},
	"Buy Deku Stick (1)":           rawtoken{"Shop", PriorityAdvancement, 0x05, map[string]any{"object": 0x00C7, "price": 10}},
	"Buy Bombs (10)":               rawtoken{"Shop", PriorityMajor, 0x06, map[string]any{"object": 0x00CE, "price": 50}},
	"Buy Fish":                     rawtoken{"Shop", PriorityAdvancement, 0x07, map[string]any{"object": 0x00F4, "price": 200}},
	"Buy Red Potion for 30 Rupees": rawtoken{"Shop", PriorityMajor, 0x08, map[string]any{"object": 0x00EB, "price": 30}},
	"Buy Green Potion":             rawtoken{"Shop", PriorityMajor, 0x09, map[string]any{"object": 0x00EB, "price": 30}},
	"Buy Blue Potion":              rawtoken{"Shop", PriorityMajor, 0x0A, map[string]any{"object": 0x00EB, "price": 100}},
	"Buy Hylian Shield":            rawtoken{"Shop", PriorityAdvancement, 0x0C, map[string]any{"object": 0x00DC, "price": 80}},
	"Buy Deku Shield":              rawtoken{"Shop", PriorityAdvancement, 0x0D, map[string]any{"object": 0x00CB, "price": 40}},
	"Buy Goron Tunic":              rawtoken{"Shop", PriorityAdvancement, 0x0E, map[string]any{"object": 0x00F2, "price": 200}},
	"Buy Zora Tunic":               rawtoken{"Shop", PriorityAdvancement, 0x0F, map[string]any{"object": 0x00F2, "price": 300}},
	"Buy Heart":                    rawtoken{"Shop", PriorityMajor, 0x10, map[string]any{"object": 0x00B7, "price": 10}},
	"Buy Bombchu (10)":             rawtoken{"Shop", PriorityAdvancement, 0x15, map[string]any{"object": 0x00D9, "price": 99}},
	"Buy Bombchu (20)":             rawtoken{"Shop", PriorityAdvancement, 0x16, map[string]any{"object": 0x00D9, "price": 180}},
	"Buy Bombchu (5)":              rawtoken{"Shop", PriorityAdvancement, 0x18, map[string]any{"object": 0x00D9, "price": 60}},
	"Buy Deku Seeds (30)":          rawtoken{"Shop", PriorityMajor, 0x1D, map[string]any{"object": 0x0119, "price": 30}},
	"Sold Out":                     rawtoken{"Shop", PriorityMajor, 0x26, map[string]any{"object": 0x0148}},
	"Buy Blue Fire":                rawtoken{"Shop", PriorityAdvancement, 0x27, map[string]any{"object": 0x0173, "price": 300}},
	"Buy Bottle Bug":               rawtoken{"Shop", PriorityAdvancement, 0x28, map[string]any{"object": 0x0174, "price": 50}},
	"Buy Poe":                      rawtoken{"Shop", PriorityMajor, 0x2A, map[string]any{"object": 0x0176, "price": 30}},
	"Buy Fairy\"s Spirit":          rawtoken{"Shop", PriorityAdvancement, 0x2B, map[string]any{"object": 0x0177, "price": 50}},
	"Buy Arrows (10)":              rawtoken{"Shop", PriorityMajor, 0x2C, map[string]any{"object": 0x00D8, "price": 20}},
	"Buy Bombs (20)":               rawtoken{"Shop", PriorityMajor, 0x2D, map[string]any{"object": 0x00CE, "price": 80}},
	"Buy Bombs (30)":               rawtoken{"Shop", PriorityMajor, 0x2E, map[string]any{"object": 0x00CE, "price": 120}},
	"Buy Bombs (5) for 35 Rupees":  rawtoken{"Shop", PriorityMajor, 0x2F, map[string]any{"object": 0x00CE, "price": 35}},
	"Buy Red Potion for 40 Rupees": rawtoken{"Shop", PriorityMajor, 0x30, map[string]any{"object": 0x00EB, "price": 40}},
	"Buy Red Potion for 50 Rupees": rawtoken{"Shop", PriorityMajor, 0x31, map[string]any{"object": 0x00EB, "price": 50}},

	"Kokiri Emerald":   rawtoken{"DungeonReward", PriorityAdvancement, GI_MISSING, map[string]any{"stone": true, "addr2_data": 0x80, "bit_mask": 0x00040000, "item_id": 0x6C, "actor_type": 0x13, "object_id": 0x00AD}},
	"Goron Ruby":       rawtoken{"DungeonReward", PriorityAdvancement, GI_MISSING, map[string]any{"stone": true, "addr2_data": 0x81, "bit_mask": 0x00080000, "item_id": 0x6D, "actor_type": 0x14, "object_id": 0x00AD}},
	"Zora Sapphire":    rawtoken{"DungeonReward", PriorityAdvancement, GI_MISSING, map[string]any{"stone": true, "addr2_data": 0x82, "bit_mask": 0x00100000, "item_id": 0x6E, "actor_type": 0x15, "object_id": 0x00AD}},
	"Forest Medallion": rawtoken{"DungeonReward", PriorityAdvancement, GI_MISSING, map[string]any{"medallion": true, "addr2_data": 0x3E, "bit_mask": 0x00000001, "item_id": 0x66, "actor_type": 0x0B, "object_id": 0x00BA}},
	"Fire Medallion":   rawtoken{"DungeonReward", PriorityAdvancement, GI_MISSING, map[string]any{"medallion": true, "addr2_data": 0x3C, "bit_mask": 0x00000002, "item_id": 0x67, "actor_type": 0x09, "object_id": 0x00BA}},
	"Water Medallion":  rawtoken{"DungeonReward", PriorityAdvancement, GI_MISSING, map[string]any{"medallion": true, "addr2_data": 0x3D, "bit_mask": 0x00000004, "item_id": 0x68, "actor_type": 0x0A, "object_id": 0x00BA}},
	"Spirit Medallion": rawtoken{"DungeonReward", PriorityAdvancement, GI_MISSING, map[string]any{"medallion": true, "addr2_data": 0x3F, "bit_mask": 0x00000008, "item_id": 0x69, "actor_type": 0x0C, "object_id": 0x00BA}},
	"Shadow Medallion": rawtoken{"DungeonReward", PriorityAdvancement, GI_MISSING, map[string]any{"medallion": true, "addr2_data": 0x41, "bit_mask": 0x00000010, "item_id": 0x6A, "actor_type": 0x0D, "object_id": 0x00BA}},
	"Light Medallion":  rawtoken{"DungeonReward", PriorityAdvancement, GI_MISSING, map[string]any{"medallion": true, "addr2_data": 0x40, "bit_mask": 0x00000020, "item_id": 0x6B, "actor_type": 0x0E, "object_id": 0x00BA}},
}
