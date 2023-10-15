package logic

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/etc-sudonters/zootler/pkg/entity"
)

// jq -Mr 'reduce .[] as $x ([]; . + $x.categories // []) | .[]' data/locations.json
type (
	DefaultItem          string
	RecoveryHeart        struct{}
	ActorOverride        struct{}
	Beehive              struct{}
	Boss                 struct{}
	BossHeart            struct{}
	BottomOfTheWell      struct{}
	Chest                struct{}
	Collectable          struct{}
	Cow                  struct{}
	Crate                struct{}
	Cutscene             struct{}
	DeathMountainCrater  struct{}
	DeathMountainTrail   struct{}
	DekuScrub            struct{}
	DekuScrubUpgrade     struct{}
	DekuTree             struct{}
	DesertColossus       struct{}
	DodongosCavern       struct{}
	FireTemple           struct{}
	Flying               struct{}
	ForestArea           struct{}
	ForestTemple         struct{}
	Freestanding         struct{}
	GanonsCastle         struct{}
	GanonsTower          struct{}
	GerudoTrainingGround struct{}
	GerudoValley         struct{}
	GerudosFortress      struct{}
	GoldSkulltula        struct{}
	GoronCity            struct{}
	Graveyard            struct{}
	GreatFairie          struct{}
	Grotto               struct{}
	GrottoScrub          struct{}
	HauntedWasteland     struct{}
	Hint                 struct{}
	HintStone            struct{}
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
	NPC                  struct{}
	NeedSpiritualStones  struct{}
	OutsideGanonsCastle  struct{}
	Pot                  struct{}
	RupeeTower           struct{}
	SacredForestMeadow   struct{}
	Scrub                struct{}
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

func ParseComponentsFromLocationType(typ string) []entity.Component {
	var comps []entity.Component
	switch typ {
	case "ActorOverride":
		comps = append(comps, ActorOverride{})
		break
	case "Beehive":
		comps = append(comps, Beehive{})
		break
	case "Boss":
		comps = append(comps, Boss{})
		break
	case "BossHeart":
		comps = append(comps, BossHeart{})
		break
	case "Chest":
		comps = append(comps, Chest{})
		break
	case "Collectable":
		comps = append(comps, Collectable{})
		break
	case "Crate":
		comps = append(comps, Crate{})
		break
	case "Cutscene":
		comps = append(comps, Cutscene{})
		break
	case "Drop":
		comps = append(comps, Drop{})
		break
	case "Event":
		comps = append(comps, Event{})
		break
	case "FlyingPot":
		comps = append(comps, Flying{}, Pot{})
		break
	case "Freestanding":
		comps = append(comps, Freestanding{})
		break
	case "GrottoScrub":
		comps = append(comps, GrottoScrub{})
		break
	case "GS Token":
		comps = append(comps, GoldSkulltulaToken{})
		break
	case "Hint":
		comps = append(comps, Hint{})
		break
	case "HintStone":
		comps = append(comps, HintStone{}, Hint{})
		break
	case "NPC":
		comps = append(comps, NPC{})
		break
	case "Pot":
		comps = append(comps, Pot{})
		break
	case "RupeeTower":
		comps = append(comps, RupeeTower{})
		break
	case "Scrub":
		comps = append(comps, Scrub{})
		break
	case "Shop":
		comps = append(comps, Shop{})
		break
	case "SmallCrate":
		comps = append(comps, SmallCrate{})
		break
	case "Song":
		comps = append(comps, Song{})
		break
	}

	return comps
}

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
		comps = append(comps, Flying{}, Pot{})
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

func GetAllLocationComponents(p PlacementLocation) []entity.Component {
	var comps []entity.Component

	comps = append(comps, ParseComponentsFromLocationType(p.Type))
	comps = append(comps, DefaultItem(p.DefaultItem))
	for _, tag := range p.Tags {
		comps = append(comps, ParseComponentsFromLocationTag(tag))
	}

	return comps
}

type PlacementLocation struct {
	Name        string   `json:"name"`
	Type        string   `json:"type"`
	DefaultItem string   `json:"vanilla"`
	Tags        []string `json:"categories"`
}

func (p *PlacementLocation) UnmarshalJSON(data []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	for k, v := range raw {
		if v == nil {
			continue
		}
		switch strings.ToLower(k) {
		case "name":
			p.Name = v.(string)
			break
		case "type":
			p.Type = v.(string)
			break
		case "vanilla":
			p.DefaultItem = v.(string)
			break
		case "categories":
			c := v.([]interface{})
			p.Tags = make([]string, len(c))
			for i, s := range c {
				p.Tags[i] = s.(string)
			}
			break
		}
	}

	return nil
}

func ReadLocations(r io.Reader) ([]PlacementLocation, error) {
	decoder := json.NewDecoder(r)
	var locs []PlacementLocation
	if err := decoder.Decode(&locs); err != nil {
		return nil, fmt.Errorf("while loading locations %w", err)
	}

	return locs, nil
}

func ReadLocationFile(fp string) ([]PlacementLocation, error) {
	fh, err := os.Open(fp)
	if err != nil {
		return nil, fmt.Errorf("when opening %s: %w", fp, err)
	}

	return ReadLocations(fh)
}
