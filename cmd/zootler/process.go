package main

import (
	"fmt"
	"regexp"
	"strings"
	"sudonters/zootler/internal/table"
	"sudonters/zootler/pkg/world/components"
)

var alphaOnly = regexp.MustCompile("[^a-z]+")

func normalize(s string) string {
	return alphaOnly.ReplaceAllString(strings.ToLower(s), "")
}

type FileItem struct {
	Name        string                 `json:"name"`
	Type        string                 `json:"type"`
	Advancement bool                   `json:"advancement"`
	Priority    bool                   `json:"priority"`
	Special     map[string]interface{} `json:"special"`
}

type FileLocation struct {
	Name       string   `json:"name"`
	Type       string   `json:"type"`
	Default    string   `json:"vanilla"`
	Categories []string `json:"categories"`
}

func (item FileItem) TableValues() (table.Values, error) {
	values := table.Values{
		components.Name(item.Name),
		components.Token{},
	}

	if kind, err := item.kind(); err != nil {
		return nil, err
	} else {
		values = append(values, kind)
	}

	if item.Advancement {
		values = append(values, components.Advancement{})
	}

	if item.Priority {
		values = append(values, components.Priority{})
	}

	if special, err := item.special(); err != nil {
		return nil, err
	} else {
		values = append(values, special)
	}

	return values, nil
}

func (item FileItem) kind() (table.Value, error) {
	switch normalize(item.Type) {
	case "bosskey":
		return components.BossKey{}, nil
	case "compass":
		return components.Compass{}, nil
	case "drop":
		return components.Drop{}, nil
	case "dungeonreward":
		return components.DungeonReward{}, nil
	case "event":
		return components.Event{}, nil
	case "ganonbosskey":
		return components.GanonBossKey{}, nil
	case "hideoutsmallkey":
		return components.HideoutSmallKey{}, nil
	case "item":
		return components.Item{}, nil
	case "map":
		return components.Map{}, nil
	case "refill":
		return components.Refill{}, nil
	case "shop":
		return components.Shop{}, nil
	case "silverrupee":
		return components.SilverRupee{}, nil
	case "smallkey":
		return components.SmallKey{}, nil
	case "song":
		return components.Song{}, nil
	case "tcgsmallkey":
		return components.TCGSmallKey{}, nil
	case "token":
		return components.GoldSkulltulaToken{}, nil
	}
	return nil, fmt.Errorf("unknown item type '%s'", item.Type)
}

func (item FileItem) special() (table.Values, error) {
	var values table.Values

	if price, ok := item.Special["price"]; ok {
		if price, ok := price.(float64); ok {
			values = append(values, components.Price(price))
		}
	}

	if _, ok := item.Special["bottle"]; ok {
		values = append(values, components.Bottle{})
	}

	if _, ok := item.Special["ocarina_button"]; ok {
		values = append(values, components.OcarinaButton{})
	}

	if _, ok := item.Special["junk"]; ok {
		values = append(values, components.Junk{})
	}

	if _, ok := item.Special["medallion"]; ok {
		values = append(values, components.Medallion{})
	}

	if _, ok := item.Special["stone"]; ok {
		values = append(values, components.SpiritualStone{})
	}

	if _, ok := item.Special["trade"]; ok {
		values = append(values, components.Trade{})
	}

	return values, nil
}

func (item FileLocation) TableValues() (table.Values, error) {
	values := table.Values{components.Name(item.Name), components.Location{}}

	if kind, err := item.kind(); err != nil {
		return nil, err
	} else {
		values = append(values, kind)
	}

	if categories, err := item.categories(); err != nil {
		return nil, err
	} else {
		values = append(values, categories...)
	}

	return values, nil
}

func (item FileLocation) kind() (table.Value, error) {
	switch normalize(item.Type) {
	case "beehive":
		return components.Beehive{}, nil
	case "boss":
		return components.Boss{}, nil
	case "bossheart":
		return components.BossHeart{}, nil
	case "chest":
		return components.Chest{}, nil
	case "collectable":
		return components.Collectable{}, nil
	case "crate":
		return components.Crate{}, nil
	case "cutscene":
		return components.Cutscene{}, nil
	case "drop":
		return components.Drop{}, nil
	case "event":
		return components.Event{}, nil
	case "flyingpot":
		return components.FlyingPot{}, nil
	case "freestanding":
		return components.Freestanding{}, nil
	case "grottoscrub":
		return components.GrottoScrub{}, nil
	case "gstoken":
		return components.GoldSkulltulaToken{}, nil
	case "hint":
		return components.Hint{}, nil
	case "hintstone":
		return components.HintStone{}, nil
	case "maskshop":
		return components.MaskShop{}, nil
	case "npc":
		return components.NPC{}, nil
	case "pot":
		return components.Pot{}, nil
	case "rupeetower":
		return components.RupeeTower{}, nil
	case "scrub":
		return components.Scrub{}, nil
	case "shop":
		return components.Shop{}, nil
	case "silverrupee":
		return components.SilverRupee{}, nil
	case "smallcrate":
		return components.SmallCrate{}, nil
	case "song":
		return components.Song{}, nil
	case "wonderitem":
		return components.Wonderitem{}, nil
	}

	return nil, fmt.Errorf("unknown location type '%s'", item.Type)
}

func (item FileLocation) categories() (values table.Values, err error) {
	for _, category := range item.Categories {
		switch normalize(category) {
		case "beehives":
			values = append(values, components.Beehive{})
			break
		case "bottomofthewell":
			values = append(values, components.BottomoftheWell{})
			break
		case "bottomofthewellmq":
			values = append(values, components.BottomoftheWellMQ{})
			break
		case "chests":
			values = append(values, components.Chest{})
			break
		case "cows":
			values = append(values, components.Cows{})
			break
		case "crates":
			values = append(values, components.Crate{})
			break
		case "deathmountain":
			values = append(values, components.DeathMountain{})
			break
		case "deathmountaincrater":
			values = append(values, components.DeathMountainCrater{})
			break
		case "deathmountaintrail":
			values = append(values, components.DeathMountainTrail{})
			break
		case "dekuscrubs":
			values = append(values, components.DekuScrubs{})
			break
		case "dekuscrubupgrades":
			values = append(values, components.DekuScrubUpgrades{})
			break
		case "dekutree":
			values = append(values, components.DekuTree{})
			break
		case "dekutreemq":
			values = append(values, components.DekuTreeMQ{})
			break
		case "desertcolossus":
			values = append(values, components.DesertColossus{})
			break
		case "dodongoscavern":
			values = append(values, components.DodongosCavern{})
			break
		case "dodongoscavernmq":
			values = append(values, components.DodongosCavernMQ{})
			break
		case "dungeonrewards":
			values = append(values, components.DungeonReward{})
			break
		case "firetemple":
			values = append(values, components.FireTemple{})
			break
		case "firetemplemq":
			values = append(values, components.FireTempleMQ{})
			break
		case "flyingpots":
			values = append(values, components.FlyingPot{})
			break
		case "forest":
			values = append(values, components.Forest{})
			break
		case "forestarea":
			values = append(values, components.ForestArea{})
			break
		case "foresttemple":
			values = append(values, components.ForestTemple{})
			break
		case "foresttemplemq":
			values = append(values, components.ForestTempleMQ{})
			break
		case "freestandings":
			values = append(values, components.Freestanding{})
			break
		case "ganonscastle":
			values = append(values, components.GanonsCastle{})
			break
		case "ganonscastlemq":
			values = append(values, components.GanonsCastleMQ{})
			break
		case "ganonstower":
			values = append(values, components.GanonsTower{})
			break
		case "gerudo":
			values = append(values, components.Gerudo{})
			break
		case "gerudosfortress":
			values = append(values, components.GerudosFortress{})
			break
		case "gerudotrainingground":
			values = append(values, components.GerudoTrainingGround{})
			break
		case "gerudotraininggroundmq":
			values = append(values, components.GerudoTrainingGroundMQ{})
			break
		case "gerudovalley":
			values = append(values, components.GerudoValley{})
			break
		case "goldskulltulas":
			values = append(values, components.GoldSkulltulas{})
			break
		case "goroncity":
			values = append(values, components.GoronCity{})
			break
		case "graveyard":
			values = append(values, components.Graveyard{})
			break
		case "greatfairies":
			values = append(values, components.GreatFairies{})
			break
		case "grottos":
			values = append(values, components.Grottos{})
			break
		case "hauntedwasteland":
			values = append(values, components.HauntedWasteland{})
			break
		case "hyrulecastle":
			values = append(values, components.HyruleCastle{})
			break
		case "hyrulefield":
			values = append(values, components.HyruleField{})
			break
		case "icecavern":
			values = append(values, components.IceCavern{})
			break
		case "icecavernmq":
			values = append(values, components.IceCavernMQ{})
			break
		case "jabujabusbelly":
			values = append(values, components.JabuJabusBelly{})
			break
		case "jabujabusbellymq":
			values = append(values, components.JabuJabusBellyMQ{})
			break
		case "kakariko":
			values = append(values, components.Kakariko{})
			break
		case "kakarikovillage":
			values = append(values, components.KakarikoVillage{})
			break
		case "kokiriforest":
			values = append(values, components.KokiriForest{})
			break
		case "lakehylia":
			values = append(values, components.LakeHylia{})
			break
		case "lonlonranch":
			values = append(values, components.LonLonRanch{})
			break
		case "lostwoods":
			values = append(values, components.LostWoods{})
			break
		case "market":
			values = append(values, components.Market{})
			break
		case "masterquest":
			values = append(values, components.MasterQuest{})
			break
		case "minigames":
			values = append(values, components.Minigames{})
			break
		case "needspiritualstones":
			values = append(values, components.NeedSpiritualStones{})
			break
		case "npcs":
			values = append(values, components.NPC{})
			break
		case "outsideganonscastle":
			values = append(values, components.OutsideGanonsCastle{})
			break
		case "pots":
			values = append(values, components.Pot{})
			break
		case "rupeetowers":
			values = append(values, components.RupeeTower{})
			break
		case "sacredforestmeadow":
			values = append(values, components.SacredForestMeadow{})
			break
		case "shadowtemple":
			values = append(values, components.ShadowTemple{})
			break
		case "shadowtemplemq":
			values = append(values, components.ShadowTempleMQ{})
			break
		case "shops":
			values = append(values, components.Shop{})
			break
		case "silverrupees":
			values = append(values, components.SilverRupee{})
			break
		case "skulltulahouse":
			values = append(values, components.SkulltulaHouse{})
			break
		case "smallcrates":
			values = append(values, components.SmallCrate{})
			break
		case "songs":
			values = append(values, components.Song{})
			break
		case "spirittemple":
			values = append(values, components.SpiritTemple{})
			break
		case "spirittemplemq":
			values = append(values, components.SpiritTempleMQ{})
			break
		case "templeoftime":
			values = append(values, components.TempleofTime{})
			break
		case "thieveshideout":
			values = append(values, components.ThievesHideout{})
			break
		case "vanilladungeons":
			values = append(values, components.VanillaDungeons{})
			break
		case "watertemple":
			values = append(values, components.WaterTemple{})
			break
		case "watertemplemq":
			values = append(values, components.WaterTempleMQ{})
			break
		case "wonderitem":
			values = append(values, components.Wonderitem{})
			break
		case "zorasdomain":
			values = append(values, components.ZorasDomain{})
			break
		case "zorasfountain":
			values = append(values, components.ZorasFountain{})
			break
		case "zorasriver":
			values = append(values, components.ZorasRiver{})
			break
		default:
			err = fmt.Errorf("unknown category '%s'", category)
			return
		}
	}

	return
}
