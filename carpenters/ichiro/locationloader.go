package ichiro

import (
	"sudonters/zootler/internal"
	"sudonters/zootler/internal/components"
	"sudonters/zootler/internal/table"

	"github.com/etc-sudonters/substrate/slipup"
)

type LocationComponents struct {
	Name       string   `json:"name"`
	Type       string   `json:"type"`
	Default    string   `json:"vanilla"`
	Categories []string `json:"categories"`
}

func (loc LocationComponents) EntityName() components.Name {
	return components.Name(loc.Name)
}

func (loc LocationComponents) AsComponents() table.Values {
	vt := table.Values{components.Location{}, loc.kind()}
	if loc.Default != "" {
		vt = append(vt, components.DefaultItemName(loc.Default))
	}
	return loc.categories(vt)
}

func (loc LocationComponents) kind() table.Value {
	switch internal.Normalize(loc.Type) {
	case "beehive":
		return components.Beehive{}
	case "boss":
		return components.Boss{}
	case "bossheart":
		return components.BossHeart{}
	case "chest":
		return components.Chest{}
	case "collectable":
		return components.Collectable{}
	case "crate":
		return components.Crate{}
	case "cutscene":
		return components.Cutscene{}
	case "drop":
		return components.Drop{}
	case "event":
		return components.Event{}
	case "flyingpot":
		return components.FlyingPot{}
	case "freestanding":
		return components.Freestanding{}
	case "grottoscrub":
		return components.GrottoScrub{}
	case "gstoken":
		return components.GoldSkulltulaToken{}
	case "hint":
		return components.Hint{}
	case "hintstone":
		return components.HintStone{}
	case "maskshop":
		return components.MaskShop{}
	case "npc":
		return components.NPC{}
	case "pot":
		return components.Pot{}
	case "rupeetower":
		return components.RupeeTower{}
	case "scrub":
		return components.Scrub{}
	case "shop":
		return components.Shop{}
	case "silverrupee":
		return components.SilverRupee{}
	case "smallcrate":
		return components.SmallCrate{}
	case "song":
		return components.Song{}
	case "wonderitem":
		return components.Wonderitem{}
	}

	panic(slipup.Createf("unknown location type '%s'", loc.Type))
}

func (loc LocationComponents) categories(vt table.Values) table.Values {
	assign := func(vs ...table.Value) {
		vt = append(vt, vs...)
	}

	for _, category := range loc.Categories {
		switch internal.Normalize(category) {
		case "beehives":
			assign(components.Beehive{})
			continue
		case "bottomofthewell":
			assign(components.BottomoftheWell{})
			continue
		case "bottomofthewellmq":
			assign(components.BottomoftheWell{}, components.MasterQuest{})
			continue
		case "chests":
			assign(components.Chest{})
			continue
		case "cows":
			assign(components.Cows{})
			continue
		case "crates":
			assign(components.Crate{})
			continue
		case "deathmountain":
			assign(components.DeathMountain{})
			continue
		case "deathmountaincrater":
			assign(components.DeathMountainCrater{})
			continue
		case "deathmountaintrail":
			assign(components.DeathMountainTrail{})
			continue
		case "dekuscrubs":
			assign(components.DekuScrubs{})
			continue
		case "dekuscrubupgrades":
			assign(components.DekuScrubs{}, components.DekuScrubUpgrades{})
			continue
		case "dekutree":
			assign(components.DekuTree{})
			continue
		case "dekutreemq":
			assign(components.DekuTreeMQ{}, components.MasterQuest{})
			continue
		case "desertcolossus":
			assign(components.DesertColossus{})
			continue
		case "dodongoscavern":
			assign(components.DodongosCavern{})
			continue
		case "dodongoscavernmq":
			assign(components.DodongosCavernMQ{}, components.MasterQuest{})
			continue
		case "dungeonrewards":
			assign(components.DungeonReward{})
			continue
		case "firetemple":
			assign(components.FireTemple{})
			continue
		case "firetemplemq":
			assign(components.FireTempleMQ{}, components.MasterQuest{})
			continue
		case "flyingpots":
			assign(components.FlyingPot{})
			continue
		case "forest":
			assign(components.Forest{})
			continue
		case "forestarea":
			assign(components.ForestArea{})
			continue
		case "foresttemple":
			assign(components.ForestTemple{})
			continue
		case "foresttemplemq":
			assign(components.ForestTempleMQ{}, components.MasterQuest{})
			continue
		case "freestandings":
			assign(components.Freestandings{})
			continue
		case "ganonscastle":
			assign(components.GanonsCastle{})
			continue
		case "ganonscastlemq":
			assign(components.GanonsCastleMQ{}, components.MasterQuest{})
			continue
		case "ganonstower":
			assign(components.GanonsTower{})
			continue
		case "gerudo":
			assign(components.Gerudo{})
			continue
		case "gerudosfortress":
			assign(components.GerudosFortress{})
			continue
		case "gerudotrainingground":
			assign(components.GerudoTrainingGround{})
			continue
		case "gerudotraininggroundmq":
			assign(components.GerudoTrainingGroundMQ{}, components.MasterQuest{})
			continue
		case "gerudovalley":
			assign(components.GerudoValley{})
			continue
		case "goldskulltulas":
			assign(components.GoldSkulltulaToken{})
			continue
		case "goroncity":
			assign(components.GoronCity{})
			continue
		case "graveyard":
			assign(components.Graveyard{})
			continue
		case "greatfairies":
			assign(components.GreatFairies{})
			continue
		case "grottos":
			assign(components.Grottos{})
			continue
		case "hauntedwasteland":
			assign(components.HauntedWasteland{})
			continue
		case "hyrulecastle":
			assign(components.HyruleCastle{})
			continue
		case "hyrulefield":
			assign(components.HyruleField{})
			continue
		case "icecavern":
			assign(components.IceCavern{})
			continue
		case "icecavernmq":
			assign(components.IceCavernMQ{}, components.MasterQuest{})
			continue
		case "jabujabusbelly":
			assign(components.JabuJabusBelly{})
			continue
		case "jabujabusbellymq":
			assign(components.JabuJabusBellyMQ{}, components.MasterQuest{})
			continue
		case "kakariko":
			assign(components.Kakariko{})
			continue
		case "kakarikovillage":
			assign(components.KakarikoVillage{})
			continue
		case "kokiriforest":
			assign(components.KokiriForest{})
			continue
		case "lakehylia":
			assign(components.LakeHylia{})
			continue
		case "lonlonranch":
			assign(components.LonLonRanch{})
			continue
		case "lostwoods":
			assign(components.LostWoods{})
			continue
		case "market":
			assign(components.Market{})
			continue
		case "masterquest":
			assign(components.MasterQuest{})
			continue
		case "minigames":
			assign(components.Minigames{})
			continue
		case "needspiritualstones":
			assign(components.NeedSpiritualStones{})
			continue
		case "npcs":
			assign(components.NPC{})
			continue
		case "outsideganonscastle":
			assign(components.OutsideGanonsCastle{})
			continue
		case "pots":
			assign(components.Pot{})
			continue
		case "rupeetowers":
			assign(components.RupeeTower{})
			continue
		case "sacredforestmeadow":
			assign(components.SacredForestMeadow{})
			continue
		case "shadowtemple":
			assign(components.ShadowTemple{})
			continue
		case "shadowtemplemq":
			assign(components.ShadowTempleMQ{}, components.MasterQuest{})
			continue
		case "shops":
			assign(components.Shop{})
			continue
		case "silverrupees":
			assign(components.SilverRupee{})
			continue
		case "skulltulahouse":
			assign(components.SkulltulaHouse{})
			continue
		case "smallcrates":
			assign(components.SmallCrate{})
			continue
		case "songs":
			assign(components.OcarinaSong{})
			continue
		case "spirittemple":
			assign(components.SpiritTemple{})
			continue
		case "spirittemplemq":
			assign(components.SpiritTempleMQ{}, components.MasterQuest{})
			continue
		case "templeoftime":
			assign(components.TempleofTime{})
			continue
		case "thieveshideout":
			assign(components.ThievesHideout{})
			continue
		case "vanilladungeons":
			assign(components.VanillaDungeons{})
			continue
		case "watertemple":
			assign(components.WaterTemple{})
			continue
		case "watertemplemq":
			assign(components.WaterTempleMQ{}, components.MasterQuest{})
			continue
		case "wonderitem":
			assign(components.Wonderitem{})
			continue
		case "zorasdomain":
			assign(components.ZorasDomain{})
			continue
		case "zorasfountain":
			assign(components.ZorasFountain{})
			continue
		case "zorasriver":
			assign(components.ZorasRiver{})
			continue
		default:
			panic(slipup.Createf("unknown category '%s'", category))
		}
	}
	return vt
}
