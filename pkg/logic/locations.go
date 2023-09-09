package logic

import (
	"os"

	"github.com/etc-sudonters/zootler/pkg/entity"
	"muzzammil.xyz/jsonc"
)

type (
	Beehive              struct{}
	BottomOfTheWell      struct{}
	Chest                struct{}
	Cow                  struct{}
	Crate                struct{}
	DeathMountainCrater  struct{}
	DeathMountainTrail   struct{}
	DekuScrub            struct{}
	DekuScrubUpgrade     struct{}
	DekuTree             struct{}
	DesertColossus       struct{}
	DodongosCavern       struct{}
	FireTemple           struct{}
	FlyingPot            struct{}
	ForestArea           struct{}
	ForestTemple         struct{}
	Freestanding         struct{}
	GanonsCastle         struct{}
	GanonsTower          struct{}
	GerudosFortress      struct{}
	GerudoTrainingGround struct{}
	GerudoValley         struct{}
	GoldSkulltula        struct{}
	GoronCity            struct{}
	Graveyard            struct{}
	GreatFairie          struct{}
	Grotto               struct{}
	HauntedWasteland     struct{}
	HyruleCastle         struct{}
	HyruleField          struct{}
	IceCavern            struct{}
	JabuJabusBelly       struct{}
	KakarikoVillage      struct{}
	KokiriForest         struct{}
	LakeHylia            struct{}
	LonLonRanch          struct{}
	LostWoods            struct{}
	Market               struct{}
	MasterQuest          struct{}
	Minigame             struct{}
	NeedSpiritualStones  struct{}
	NPC                  struct{}
	OutsideGanonsCastle  struct{}
	Pot                  struct{}
	RupeeTower           struct{}
	SacredForestMeadow   struct{}
	ShadowTemple         struct{}
	Shop                 struct{}
	SkulltulaHouse       struct{}
	SmallCrate           struct{}
	Song                 struct{}
	SpiritTemple         struct{}
	TempleofTime         struct{}
	ThievesHideout       struct{}
	VanillaDungeon       struct{}
	WaterTemple          struct{}
	ZorasDomain          struct{}
	ZorasFountain        struct{}
	ZorasRiver           struct{}
)

func ParseComponentsFromLocationTag(tag string) []entity.Component {
	var comps []entity.Component

	switch tag {
	case "Beehives":
		comps = append(comps, Beehive{})
		break
	case "Bottom of the Well":
		comps = append(comps, BottomOfTheWell{})
		break
	case "Bottom of the Well MQ":
		comps = append(comps, BottomOfTheWell{}, MasterQuest{})
		break
	case "Chests":
		comps = append(comps, Chest{})
		break
	case "Cows":
		comps = append(comps, Cow{})
		break
	case "Crates":
		comps = append(comps, Crate{})
		break
	case "Death Mountain Crater":
		comps = append(comps, DeathMountainCrater{})
		break
	case "Death Mountain Trail":
		comps = append(comps, DeathMountainTrail{})
		break
	case "Deku Scrubs":
		comps = append(comps, DekuScrub{})
		break
	case "Deku Scrub Upgrades":
		comps = append(comps, DekuScrubUpgrade{})
		break
	case "Deku Tree":
		comps = append(comps, DekuTree{})
		break
	case "Deku Tree MQ":
		comps = append(comps, DekuTree{}, MasterQuest{})
		break
	case "Desert Colossus":
		comps = append(comps, DesertColossus{})
		break
	case "Dodongo's Cavern":
		comps = append(comps, DodongosCavern{})
		break
	case "Dodongo's Cavern MQ":
		comps = append(comps, DodongosCavern{}, MasterQuest{})
		break
	case "Fire Temple":
		comps = append(comps, FireTemple{})
		break
	case "Fire Temple MQ":
		comps = append(comps, FireTemple{}, MasterQuest{})
		break
	case "Flying Pots":
		comps = append(comps, FlyingPot{})
		break
	case "Forest Area":
		comps = append(comps, ForestArea{})
		break
	case "Forest Temple":
		comps = append(comps, ForestTemple{})
		break
	case "Forest Temple MQ":
		comps = append(comps, ForestTemple{}, MasterQuest{})
		break
	case "Freestandings":
		comps = append(comps, Freestanding{})
		break
	case "Ganon's Castle":
		comps = append(comps, GanonsCastle{})
		break
	case "Ganon's Castle MQ":
		comps = append(comps, GanonsCastle{}, MasterQuest{})
		break
	case "Ganon's Tower":
		comps = append(comps, GanonsTower{})
		break
	case "Gerudo's Fortress":
		comps = append(comps, GerudosFortress{})
		break
	case "Gerudo Training Ground":
		comps = append(comps, GerudoTrainingGround{})
		break
	case "Gerudo Training Ground MQ":
		comps = append(comps, GerudoTrainingGround{}, MasterQuest{})
		break
	case "Gerudo Valley":
		comps = append(comps, GerudoValley{})
		break
	case "Gold Skulltulas":
		comps = append(comps, GoldSkulltula{})
		break
	case "Goron City":
		comps = append(comps, GoronCity{})
		break
	case "Graveyard":
		comps = append(comps, Graveyard{})
		break
	case "Great Fairies":
		comps = append(comps, GreatFairie{})
		break
	case "Grottos":
		comps = append(comps, Grotto{})
		break
	case "Haunted Wasteland":
		comps = append(comps, HauntedWasteland{})
		break
	case "Hyrule Castle":
		comps = append(comps, HyruleCastle{})
		break
	case "Hyrule Field":
		comps = append(comps, HyruleField{})
		break
	case "Ice Cavern":
		comps = append(comps, IceCavern{})
		break
	case "Ice Cavern MQ":
		comps = append(comps, IceCavern{}, MasterQuest{})
		break
	case "Jabu Jabu's Belly":
		comps = append(comps, JabuJabusBelly{})
		break
	case "Jabu Jabu's Belly MQ":
		comps = append(comps, JabuJabusBelly{}, MasterQuest{})
		break
	case "Kakariko Village":
		comps = append(comps, KakarikoVillage{})
		break
	case "Kokiri Forest":
		comps = append(comps, KokiriForest{})
		break
	case "Lake Hylia":
		comps = append(comps, LakeHylia{})
		break
	case "Lon Lon Ranch":
		comps = append(comps, LonLonRanch{})
		break
	case "Lost Woods":
		comps = append(comps, LostWoods{})
		break
	case "Market":
		comps = append(comps, Market{})
		break
	case "Master Quest":
		comps = append(comps, MasterQuest{})
		break
	case "Minigames":
		comps = append(comps, Minigame{})
		break
	case "Need Spiritual Stones":
		comps = append(comps, NeedSpiritualStones{})
		break
	case "NPCs":
		comps = append(comps, NPC{})
		break
	case "Outside Ganon's Castle":
		comps = append(comps, OutsideGanonsCastle{})
		break
	case "Pots":
		comps = append(comps, Pot{})
		break
	case "Rupee Towers":
		comps = append(comps, RupeeTower{})
		break
	case "Sacred Forest Meadow":
		comps = append(comps, SacredForestMeadow{})
		break
	case "Shadow Temple":
		comps = append(comps, ShadowTemple{})
		break
	case "Shadow Temple MQ":
		comps = append(comps, ShadowTemple{}, MasterQuest{})
		break
	case "Shops":
		comps = append(comps, Shop{})
		break
	case "Skulltula House":
		comps = append(comps, SkulltulaHouse{})
		break
	case "Small Crates":
		comps = append(comps, SmallCrate{})
		break
	case "Songs":
		comps = append(comps, Song{})
		break
	case "Spirit Temple":
		comps = append(comps, SpiritTemple{})
		break
	case "Spirit Temple MQ":
		comps = append(comps, SpiritTemple{}, MasterQuest{})
		break
	case "Temple of Time":
		comps = append(comps, TempleofTime{})
		break
	case "Thieves' Hideout":
		comps = append(comps, ThievesHideout{})
		break
	case "Vanilla Dungeons":
		comps = append(comps, VanillaDungeon{})
		break
	case "Water Temple":
		comps = append(comps, WaterTemple{})
		break
	case "Water Temple MQ":
		comps = append(comps, WaterTemple{}, MasterQuest{})
		break
	case "Zora's Domain":
		comps = append(comps, ZorasDomain{})
		break
	case "Zora's Fountain":
		comps = append(comps, ZorasFountain{})
		break
	case "Zora's River":
		comps = append(comps, ZorasRiver{})
		break
	}

	return comps
}

type PlacementLocation struct {
	Name        string   `json:"name"`
	Type        string   `json:"type"`
	DefaultItem string   `json:"vanilla"`
	Tags        []string `json:"categories"`
}

func ReadLocationFile(fp string) ([]PlacementLocation, error) {
	contents, err := os.ReadFile(fp)
	if err != nil {
		return nil, err
	}

	var locs []PlacementLocation
	if err := jsonc.Unmarshal(contents, &locs); err != nil {
		return nil, err
	}

	return locs, nil
}
