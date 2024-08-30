package main

import (
	"context"
	"fmt"
	"github.com/etc-sudonters/substrate/slipup"
	"sudonters/zootler/internal"
	"sudonters/zootler/internal/components"
	"sudonters/zootler/internal/query"
	"sudonters/zootler/internal/table"
)

type FileLocation struct {
	Name       string   `json:"name"`
	Type       string   `json:"type"`
	Default    string   `json:"vanilla"`
	Categories []string `json:"categories"`
}

func (item FileLocation) GetName() components.Name {
	return components.Name(item.Name)
}

func (item FileLocation) AddComponents(rid table.RowId, storage query.Engine) error {
	if err := storage.SetValues(rid, table.Values{components.Location{}}); err != nil {
		return err
	}

	if err := item.kind(rid, storage); err != nil {
		return err
	}

	if err := item.categories(rid, storage); err != nil {
		return err
	}

	return nil
}

func (item FileLocation) kind(rid table.RowId, storage query.Engine) error {
	switch internal.Normalize(item.Type) {
	case "beehive":
		return storage.SetValues(rid, table.Values{components.Beehive{}})
	case "boss":
		return storage.SetValues(rid, table.Values{components.Boss{}})
	case "bossheart":
		return storage.SetValues(rid, table.Values{components.BossHeart{}})
	case "chest":
		return storage.SetValues(rid, table.Values{components.Chest{}})
	case "collectable":
		return storage.SetValues(rid, table.Values{components.Collectable{}})
	case "crate":
		return storage.SetValues(rid, table.Values{components.Crate{}})
	case "cutscene":
		return storage.SetValues(rid, table.Values{components.Cutscene{}})
	case "drop":
		return storage.SetValues(rid, table.Values{components.Drop{}})
	case "event":
		return storage.SetValues(rid, table.Values{components.Event{}})
	case "flyingpot":
		return storage.SetValues(rid, table.Values{components.FlyingPot{}})
	case "freestanding":
		return storage.SetValues(rid, table.Values{components.Freestanding{}})
	case "grottoscrub":
		return storage.SetValues(rid, table.Values{components.GrottoScrub{}})
	case "gstoken":
		return storage.SetValues(rid, table.Values{components.GoldSkulltulaToken{}})
	case "hint":
		return storage.SetValues(rid, table.Values{components.Hint{}})
	case "hintstone":
		return storage.SetValues(rid, table.Values{components.HintStone{}})
	case "maskshop":
		return storage.SetValues(rid, table.Values{components.MaskShop{}})
	case "npc":
		return storage.SetValues(rid, table.Values{components.NPC{}})
	case "pot":
		return storage.SetValues(rid, table.Values{components.Pot{}})
	case "rupeetower":
		return storage.SetValues(rid, table.Values{components.RupeeTower{}})
	case "scrub":
		return storage.SetValues(rid, table.Values{components.Scrub{}})
	case "shop":
		return storage.SetValues(rid, table.Values{components.Shop{}})
	case "silverrupee":
		return storage.SetValues(rid, table.Values{components.SilverRupee{}})
	case "smallcrate":
		return storage.SetValues(rid, table.Values{components.SmallCrate{}})
	case "song":
		return storage.SetValues(rid, table.Values{components.Song{}})
	case "wonderitem":
		return storage.SetValues(rid, table.Values{components.Wonderitem{}})
	}

	return fmt.Errorf("unknown location type '%s'", item.Type)
}

func (item FileLocation) categories(rid table.RowId, storage query.Engine) error {
	for _, category := range item.Categories {
		switch internal.Normalize(category) {
		case "beehives":
			if err := storage.SetValues(rid, table.Values{components.Beehive{}}); err != nil {
				return err
			}
			break
		case "bottomofthewell":
			if err := storage.SetValues(rid, table.Values{components.BottomoftheWell{}}); err != nil {
				return err
			}
			break
		case "bottomofthewellmq":
			if err := storage.SetValues(rid, table.Values{components.BottomoftheWellMQ{}}); err != nil {
				return err
			}
			break
		case "chests":
			if err := storage.SetValues(rid, table.Values{components.Chest{}}); err != nil {
				return err
			}
			break
		case "cows":
			if err := storage.SetValues(rid, table.Values{components.Cows{}}); err != nil {
				return err
			}
			break
		case "crates":
			if err := storage.SetValues(rid, table.Values{components.Crate{}}); err != nil {
				return err
			}
			break
		case "deathmountain":
			if err := storage.SetValues(rid, table.Values{components.DeathMountain{}}); err != nil {
				return err
			}
			break
		case "deathmountaincrater":
			if err := storage.SetValues(rid, table.Values{components.DeathMountainCrater{}}); err != nil {
				return err
			}
			break
		case "deathmountaintrail":
			if err := storage.SetValues(rid, table.Values{components.DeathMountainTrail{}}); err != nil {
				return err
			}
			break
		case "dekuscrubs":
			if err := storage.SetValues(rid, table.Values{components.DekuScrubs{}}); err != nil {
				return err
			}
			break
		case "dekuscrubupgrades":
			if err := storage.SetValues(rid, table.Values{components.DekuScrubUpgrades{}}); err != nil {
				return err
			}
			break
		case "dekutree":
			if err := storage.SetValues(rid, table.Values{components.DekuTree{}}); err != nil {
				return err
			}
			break
		case "dekutreemq":
			if err := storage.SetValues(rid, table.Values{components.DekuTreeMQ{}}); err != nil {
				return err
			}
			break
		case "desertcolossus":
			if err := storage.SetValues(rid, table.Values{components.DesertColossus{}}); err != nil {
				return err
			}
			break
		case "dodongoscavern":
			if err := storage.SetValues(rid, table.Values{components.DodongosCavern{}}); err != nil {
				return err
			}
			break
		case "dodongoscavernmq":
			if err := storage.SetValues(rid, table.Values{components.DodongosCavernMQ{}}); err != nil {
				return err
			}
			break
		case "dungeonrewards":
			if err := storage.SetValues(rid, table.Values{components.DungeonReward{}}); err != nil {
				return err
			}
			break
		case "firetemple":
			if err := storage.SetValues(rid, table.Values{components.FireTemple{}}); err != nil {
				return err
			}
			break
		case "firetemplemq":
			if err := storage.SetValues(rid, table.Values{components.FireTempleMQ{}}); err != nil {
				return err
			}
			break
		case "flyingpots":
			if err := storage.SetValues(rid, table.Values{components.FlyingPot{}}); err != nil {
				return err
			}
			break
		case "forest":
			if err := storage.SetValues(rid, table.Values{components.Forest{}}); err != nil {
				return err
			}
			break
		case "forestarea":
			if err := storage.SetValues(rid, table.Values{components.ForestArea{}}); err != nil {
				return err
			}
			break
		case "foresttemple":
			if err := storage.SetValues(rid, table.Values{components.ForestTemple{}}); err != nil {
				return err
			}
			break
		case "foresttemplemq":
			if err := storage.SetValues(rid, table.Values{components.ForestTempleMQ{}}); err != nil {
				return err
			}
			break
		case "freestandings":
			if err := storage.SetValues(rid, table.Values{components.Freestandings{}}); err != nil {
				return err
			}
			break
		case "ganonscastle":
			if err := storage.SetValues(rid, table.Values{components.GanonsCastle{}}); err != nil {
				return err
			}
			break
		case "ganonscastlemq":
			if err := storage.SetValues(rid, table.Values{components.GanonsCastleMQ{}}); err != nil {
				return err
			}
			break
		case "ganonstower":
			if err := storage.SetValues(rid, table.Values{components.GanonsTower{}}); err != nil {
				return err
			}
			break
		case "gerudo":
			if err := storage.SetValues(rid, table.Values{components.Gerudo{}}); err != nil {
				return err
			}
			break
		case "gerudosfortress":
			if err := storage.SetValues(rid, table.Values{components.GerudosFortress{}}); err != nil {
				return err
			}
			break
		case "gerudotrainingground":
			if err := storage.SetValues(rid, table.Values{components.GerudoTrainingGround{}}); err != nil {
				return err
			}
			break
		case "gerudotraininggroundmq":
			if err := storage.SetValues(rid, table.Values{components.GerudoTrainingGroundMQ{}}); err != nil {
				return err
			}
			break
		case "gerudovalley":
			if err := storage.SetValues(rid, table.Values{components.GerudoValley{}}); err != nil {
				return err
			}
			break
		case "goldskulltulas":
			if err := storage.SetValues(rid, table.Values{components.GoldSkulltulas{}}); err != nil {
				return err
			}
			break
		case "goroncity":
			if err := storage.SetValues(rid, table.Values{components.GoronCity{}}); err != nil {
				return err
			}
			break
		case "graveyard":
			if err := storage.SetValues(rid, table.Values{components.Graveyard{}}); err != nil {
				return err
			}
			break
		case "greatfairies":
			if err := storage.SetValues(rid, table.Values{components.GreatFairies{}}); err != nil {
				return err
			}
			break
		case "grottos":
			if err := storage.SetValues(rid, table.Values{components.Grottos{}}); err != nil {
				return err
			}
			break
		case "hauntedwasteland":
			if err := storage.SetValues(rid, table.Values{components.HauntedWasteland{}}); err != nil {
				return err
			}
			break
		case "hyrulecastle":
			if err := storage.SetValues(rid, table.Values{components.HyruleCastle{}}); err != nil {
				return err
			}
			break
		case "hyrulefield":
			if err := storage.SetValues(rid, table.Values{components.HyruleField{}}); err != nil {
				return err
			}
			break
		case "icecavern":
			if err := storage.SetValues(rid, table.Values{components.IceCavern{}}); err != nil {
				return err
			}
			break
		case "icecavernmq":
			if err := storage.SetValues(rid, table.Values{components.IceCavernMQ{}}); err != nil {
				return err
			}
			break
		case "jabujabusbelly":
			if err := storage.SetValues(rid, table.Values{components.JabuJabusBelly{}}); err != nil {
				return err
			}
			break
		case "jabujabusbellymq":
			if err := storage.SetValues(rid, table.Values{components.JabuJabusBellyMQ{}}); err != nil {
				return err
			}
			break
		case "kakariko":
			if err := storage.SetValues(rid, table.Values{components.Kakariko{}}); err != nil {
				return err
			}
			break
		case "kakarikovillage":
			if err := storage.SetValues(rid, table.Values{components.KakarikoVillage{}}); err != nil {
				return err
			}
			break
		case "kokiriforest":
			if err := storage.SetValues(rid, table.Values{components.KokiriForest{}}); err != nil {
				return err
			}
			break
		case "lakehylia":
			if err := storage.SetValues(rid, table.Values{components.LakeHylia{}}); err != nil {
				return err
			}
			break
		case "lonlonranch":
			if err := storage.SetValues(rid, table.Values{components.LonLonRanch{}}); err != nil {
				return err
			}
			break
		case "lostwoods":
			if err := storage.SetValues(rid, table.Values{components.LostWoods{}}); err != nil {
				return err
			}
			break
		case "market":
			if err := storage.SetValues(rid, table.Values{components.Market{}}); err != nil {
				return err
			}
			break
		case "masterquest":
			if err := storage.SetValues(rid, table.Values{components.MasterQuest{}}); err != nil {
				return err
			}
			break
		case "minigames":
			if err := storage.SetValues(rid, table.Values{components.Minigames{}}); err != nil {
				return err
			}
			break
		case "needspiritualstones":
			if err := storage.SetValues(rid, table.Values{components.NeedSpiritualStones{}}); err != nil {
				return err
			}
			break
		case "npcs":
			if err := storage.SetValues(rid, table.Values{components.NPC{}}); err != nil {
				return err
			}
			break
		case "outsideganonscastle":
			if err := storage.SetValues(rid, table.Values{components.OutsideGanonsCastle{}}); err != nil {
				return err
			}
			break
		case "pots":
			if err := storage.SetValues(rid, table.Values{components.Pot{}}); err != nil {
				return err
			}
			break
		case "rupeetowers":
			if err := storage.SetValues(rid, table.Values{components.RupeeTower{}}); err != nil {
				return err
			}
			break
		case "sacredforestmeadow":
			if err := storage.SetValues(rid, table.Values{components.SacredForestMeadow{}}); err != nil {
				return err
			}
			break
		case "shadowtemple":
			if err := storage.SetValues(rid, table.Values{components.ShadowTemple{}}); err != nil {
				return err
			}
			break
		case "shadowtemplemq":
			if err := storage.SetValues(rid, table.Values{components.ShadowTempleMQ{}}); err != nil {
				return err
			}
			break
		case "shops":
			if err := storage.SetValues(rid, table.Values{components.Shop{}}); err != nil {
				return err
			}
			break
		case "silverrupees":
			if err := storage.SetValues(rid, table.Values{components.SilverRupee{}}); err != nil {
				return err
			}
			break
		case "skulltulahouse":
			if err := storage.SetValues(rid, table.Values{components.SkulltulaHouse{}}); err != nil {
				return err
			}
			break
		case "smallcrates":
			if err := storage.SetValues(rid, table.Values{components.SmallCrate{}}); err != nil {
				return err
			}
			break
		case "songs":
			if err := storage.SetValues(rid, table.Values{components.OcarinaSong{}}); err != nil {
				return err
			}
			break
		case "spirittemple":
			if err := storage.SetValues(rid, table.Values{components.SpiritTemple{}}); err != nil {
				return err
			}
			break
		case "spirittemplemq":
			if err := storage.SetValues(rid, table.Values{components.SpiritTempleMQ{}}); err != nil {
				return err
			}
			break
		case "templeoftime":
			if err := storage.SetValues(rid, table.Values{components.TempleofTime{}}); err != nil {
				return err
			}
			break
		case "thieveshideout":
			if err := storage.SetValues(rid, table.Values{components.ThievesHideout{}}); err != nil {
				return err
			}
			break
		case "vanilladungeons":
			if err := storage.SetValues(rid, table.Values{components.VanillaDungeons{}}); err != nil {
				return err
			}
			break
		case "watertemple":
			if err := storage.SetValues(rid, table.Values{components.WaterTemple{}}); err != nil {
				return err
			}
			break
		case "watertemplemq":
			if err := storage.SetValues(rid, table.Values{components.WaterTempleMQ{}}); err != nil {
				return err
			}
			break
		case "wonderitem":
			if err := storage.SetValues(rid, table.Values{components.Wonderitem{}}); err != nil {
				return err
			}
			break
		case "zorasdomain":
			if err := storage.SetValues(rid, table.Values{components.ZorasDomain{}}); err != nil {
				return err
			}
			break
		case "zorasfountain":
			if err := storage.SetValues(rid, table.Values{components.ZorasFountain{}}); err != nil {
				return err
			}
			break
		case "zorasriver":
			if err := storage.SetValues(rid, table.Values{components.ZorasRiver{}}); err != nil {
				return err
			}
			break
		default:
			return fmt.Errorf("unknown category '%s'", category)
		}
	}
	return nil
}

type AttachDefaultItem struct {
	items map[string]table.RowId
}

func (a *AttachDefaultItem) Init(_ context.Context, storage query.Engine) error {
	if a.items == nil {
		a.items = make(map[string]table.RowId, 256)
	}
	q := storage.CreateQuery()
	q.Load(query.MustAsColumnId[components.Name](storage))
	q.Exists(query.MustAsColumnId[components.CollectableGameToken](storage))
	names, err := storage.Retrieve(q)
	if err != nil {
		return slipup.Describe(err, "while building item name map")
	}

	if names.Len() == 0 {
		return slipup.Createf("no items load")
	}

	for id, tup := range names.All {
		name, castErr := internal.TypeAssert[components.Name](tup.Values[0])
		if castErr != nil {
			return slipup.Describef(castErr, "gathering from row %d", id)
		}
		a.items[string(name)] = id
	}

	return nil
}

func (a *AttachDefaultItem) Components(_ context.Context, id table.RowId, l FileLocation, e query.Engine) error {
	if l.Default == "" {
		return nil
	}

	itemId, exists := a.items[l.Default]
	if !exists {
		return slipup.Createf("item '%s' was not found in the loaded items list", l.Default)
	}

	return e.SetValues(id, table.Values{components.DefaultItem(itemId)})
}
