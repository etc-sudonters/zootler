package main

import (
	"os"
	"sudonters/zootler/internal/query"
	"sudonters/zootler/internal/table"
	"sudonters/zootler/internal/table/columns"
	"sudonters/zootler/pkg/world/components"

	"muzzammil.xyz/jsonc"
)

type IntoTableValues interface {
	TableValues() (table.Values, error)
}

type DataFileLoader[T IntoTableValues] string

func (l DataFileLoader[T]) Configure(storage query.Engine) error {
	raw, err := os.ReadFile(string(l))
	if err != nil {
		return err
	}

	var items []T

	if err := jsonc.Unmarshal(raw, &items); err != nil {
		return err
	}

	for _, item := range items {
		values, err := item.TableValues()
		if err != nil {
			return err
		}
		_, err = storage.InsertRow(values...)
		if err != nil {
			return err
		}
	}

	return nil

}

type CreateStorage struct{}

func (cs CreateStorage) Configure(storage query.Engine) error {
	storage.CreateColumn(table.BuildColumnOf[components.Alias](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.Beehive](columns.NewBit(components.Beehive{})))
	storage.CreateColumn(table.BuildColumnOf[components.BossHeart](columns.NewBit(components.BossHeart{})))
	storage.CreateColumn(table.BuildColumnOf[components.BossKey](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.Boss](columns.NewBit(components.Boss{})))
	storage.CreateColumn(table.BuildColumnOf[components.Bottle](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.BottomoftheWellMQ](columns.NewBit(components.BottomoftheWellMQ{})))
	storage.CreateColumn(table.BuildColumnOf[components.BottomoftheWell](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.Chest](columns.NewBit(components.Chest{})))
	storage.CreateColumn(table.BuildColumnOf[components.Collectable](columns.NewBit(components.Collectable{})))
	storage.CreateColumn(table.BuildColumnOf[components.Collected](columns.NewBit(components.Collected{})))
	storage.CreateColumn(table.BuildColumnOf[components.Compass](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.Count](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.Cows](columns.NewBit(components.Cows{})))
	storage.CreateColumn(table.BuildColumnOf[components.Crate](columns.NewBit(components.Crate{})))
	storage.CreateColumn(table.BuildColumnOf[components.Cutscene](columns.NewBit(components.Cutscene{})))
	storage.CreateColumn(table.BuildColumnOf[components.DeathMountainCrater](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.DeathMountainTrail](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.DeathMountain](columns.NewBit(components.DeathMountain{})))
	storage.CreateColumn(table.BuildColumnOf[components.DefaultItem](columns.NewSlice()))
	storage.CreateColumn(table.BuildColumnOf[components.DekuScrubUpgrades](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.DekuScrubs](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.DekuTreeMQ](columns.NewBit(components.DekuTreeMQ{})))
	storage.CreateColumn(table.BuildColumnOf[components.DekuTree](columns.NewBit(components.DekuTree{})))
	storage.CreateColumn(table.BuildColumnOf[components.DesertColossus](columns.NewBit(components.DesertColossus{})))
	storage.CreateColumn(table.BuildColumnOf[components.DodongosCavernMQ](columns.NewBit(components.DodongosCavernMQ{})))
	storage.CreateColumn(table.BuildColumnOf[components.DodongosCavern](columns.NewBit(components.DodongosCavern{})))
	storage.CreateColumn(table.BuildColumnOf[components.Drop](columns.NewBit(components.Drop{})))
	storage.CreateColumn(table.BuildColumnOf[components.DungeonReward](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.Event](columns.NewBit(components.Event{})))
	storage.CreateColumn(table.BuildColumnOf[components.FireTempleMQ](columns.NewBit(components.FireTempleMQ{})))
	storage.CreateColumn(table.BuildColumnOf[components.FireTemple](columns.NewBit(components.FireTemple{})))
	storage.CreateColumn(table.BuildColumnOf[components.FlyingPot](columns.NewBit(components.FlyingPot{})))
	storage.CreateColumn(table.BuildColumnOf[components.ForestArea](columns.NewBit(components.ForestArea{})))
	storage.CreateColumn(table.BuildColumnOf[components.ForestTempleMQ](columns.NewBit(components.ForestTempleMQ{})))
	storage.CreateColumn(table.BuildColumnOf[components.ForestTemple](columns.NewBit(components.ForestTemple{})))
	storage.CreateColumn(table.BuildColumnOf[components.Forest](columns.NewBit(components.Forest{})))
	storage.CreateColumn(table.BuildColumnOf[components.Freestanding](columns.NewBit(components.Freestanding{})))
	storage.CreateColumn(table.BuildColumnOf[components.GanonBossKey](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.GanonsCastleMQ](columns.NewBit(components.GanonsCastleMQ{})))
	storage.CreateColumn(table.BuildColumnOf[components.GanonsCastle](columns.NewBit(components.GanonsCastle{})))
	storage.CreateColumn(table.BuildColumnOf[components.GanonsTower](columns.NewBit(components.GanonsTower{})))
	storage.CreateColumn(table.BuildColumnOf[components.GerudoTrainingGroundMQ](columns.NewBit(components.GerudoTrainingGroundMQ{})))
	storage.CreateColumn(table.BuildColumnOf[components.GerudoTrainingGround](columns.NewBit(components.GerudoTrainingGround{})))
	storage.CreateColumn(table.BuildColumnOf[components.GerudoValley](columns.NewBit(components.GerudoValley{})))
	storage.CreateColumn(table.BuildColumnOf[components.Gerudo](columns.NewBit(components.Gerudo{})))
	storage.CreateColumn(table.BuildColumnOf[components.GerudosFortress](columns.NewBit(components.GerudosFortress{})))
	storage.CreateColumn(table.BuildColumnOf[components.GoldSkulltulaToken](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.GoldSkulltulas](columns.NewBit(components.GoldSkulltulas{})))
	storage.CreateColumn(table.BuildColumnOf[components.GoronCity](columns.NewBit(components.GoronCity{})))
	storage.CreateColumn(table.BuildColumnOf[components.Graveyard](columns.NewBit(components.Graveyard{})))
	storage.CreateColumn(table.BuildColumnOf[components.GreatFairies](columns.NewBit(components.GreatFairies{})))
	storage.CreateColumn(table.BuildColumnOf[components.GrottoScrub](columns.NewBit(components.GrottoScrub{})))
	storage.CreateColumn(table.BuildColumnOf[components.Grottos](columns.NewBit(components.Grottos{})))
	storage.CreateColumn(table.BuildColumnOf[components.HauntedWasteland](columns.NewBit(components.HauntedWasteland{})))
	storage.CreateColumn(table.BuildColumnOf[components.HideoutSmallKey](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.HintStone](columns.NewBit(components.HintStone{})))
	storage.CreateColumn(table.BuildColumnOf[components.Hint](columns.NewBit(components.Hint{})))
	storage.CreateColumn(table.BuildColumnOf[components.HyruleCastle](columns.NewBit(components.HyruleCastle{})))
	storage.CreateColumn(table.BuildColumnOf[components.HyruleField](columns.NewBit(components.HyruleField{})))
	storage.CreateColumn(table.BuildColumnOf[components.IceCavernMQ](columns.NewBit(components.IceCavernMQ{})))
	storage.CreateColumn(table.BuildColumnOf[components.IceCavern](columns.NewBit(components.IceCavern{})))
	storage.CreateColumn(table.BuildColumnOf[components.Inhabited](columns.NewSlice()))
	storage.CreateColumn(table.BuildColumnOf[components.Inhabits](columns.NewSlice()))
	storage.CreateColumn(table.BuildColumnOf[components.Item](columns.NewBit(components.Item{})))
	storage.CreateColumn(table.BuildColumnOf[components.JabuJabusBellyMQ](columns.NewBit(components.JabuJabusBellyMQ{})))
	storage.CreateColumn(table.BuildColumnOf[components.JabuJabusBelly](columns.NewBit(components.JabuJabusBelly{})))
	storage.CreateColumn(table.BuildColumnOf[components.Junk](columns.NewBit(components.Junk{})))
	storage.CreateColumn(table.BuildColumnOf[components.KakarikoVillage](columns.NewBit(components.KakarikoVillage{})))
	storage.CreateColumn(table.BuildColumnOf[components.Kakariko](columns.NewBit(components.Kakariko{})))
	storage.CreateColumn(table.BuildColumnOf[components.KokiriForest](columns.NewBit(components.KokiriForest{})))
	storage.CreateColumn(table.BuildColumnOf[components.LakeHylia](columns.NewBit(components.LakeHylia{})))
	storage.CreateColumn(table.BuildColumnOf[components.Location](columns.NewBit(components.Location{})))
	storage.CreateColumn(table.BuildColumnOf[components.Locked](columns.NewBit(components.Locked{})))
	storage.CreateColumn(table.BuildColumnOf[components.LonLonRanch](columns.NewBit(components.LonLonRanch{})))
	storage.CreateColumn(table.BuildColumnOf[components.LostWoods](columns.NewBit(components.LostWoods{})))
	storage.CreateColumn(table.BuildColumnOf[components.Map](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.Market](columns.NewBit(components.Market{})))
	storage.CreateColumn(table.BuildColumnOf[components.MaskShop](columns.NewBit(components.MaskShop{})))
	storage.CreateColumn(table.BuildColumnOf[components.MasterQuest](columns.NewBit(components.MasterQuest{})))
	storage.CreateColumn(table.BuildColumnOf[components.Medallion](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.Minigames](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.NPC](columns.NewBit(components.NPC{})))
	storage.CreateColumn(table.BuildColumnOf[components.Name](columns.NewSlice()))
	storage.CreateColumn(table.BuildColumnOf[components.NeedSpiritualStones](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.OcarinaNote](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.OutsideGanonsCastle](columns.NewBit(components.OutsideGanonsCastle{})))
	storage.CreateColumn(table.BuildColumnOf[components.Placeable](columns.NewBit(components.Placeable{})))
	storage.CreateColumn(table.BuildColumnOf[components.Pot](columns.NewBit(components.Pot{})))
	storage.CreateColumn(table.BuildColumnOf[components.Price](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.Refill](columns.NewBit(components.Refill{})))
	storage.CreateColumn(table.BuildColumnOf[components.RupeeTower](columns.NewBit(components.RupeeTower{})))
	storage.CreateColumn(table.BuildColumnOf[components.SacredForestMeadow](columns.NewBit(components.SacredForestMeadow{})))
	storage.CreateColumn(table.BuildColumnOf[components.Scrub](columns.NewBit(components.Scrub{})))
	storage.CreateColumn(table.BuildColumnOf[components.ShadowTempleMQ](columns.NewBit(components.ShadowTempleMQ{})))
	storage.CreateColumn(table.BuildColumnOf[components.ShadowTemple](columns.NewBit(components.ShadowTemple{})))
	storage.CreateColumn(table.BuildColumnOf[components.ShopObject](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.Shop](columns.NewBit(components.Shop{})))
	storage.CreateColumn(table.BuildColumnOf[components.SilverRupee](columns.NewBit(components.SilverRupee{})))
	storage.CreateColumn(table.BuildColumnOf[components.SkulltulaHouse](columns.NewBit(components.SkulltulaHouse{})))
	storage.CreateColumn(table.BuildColumnOf[components.SmallCrate](columns.NewBit(components.SmallCrate{})))
	storage.CreateColumn(table.BuildColumnOf[components.SmallKey](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.Song](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.Spawn](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.SpiritTempleMQ](columns.NewBit(components.SpiritTempleMQ{})))
	storage.CreateColumn(table.BuildColumnOf[components.SpiritTemple](columns.NewBit(components.SpiritTemple{})))
	storage.CreateColumn(table.BuildColumnOf[components.SpiritualStone](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.TempleofTime](columns.NewBit(components.TempleofTime{})))
	storage.CreateColumn(table.BuildColumnOf[components.ThievesHideout](columns.NewBit(components.ThievesHideout{})))
	storage.CreateColumn(table.BuildColumnOf[components.Token](columns.NewBit(components.Token{})))
	storage.CreateColumn(table.BuildColumnOf[components.Trade](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.VanillaDungeons](columns.NewBit(components.VanillaDungeons{})))
	storage.CreateColumn(table.BuildColumnOf[components.WaterTempleMQ](columns.NewBit(components.WaterTempleMQ{})))
	storage.CreateColumn(table.BuildColumnOf[components.WaterTemple](columns.NewBit(components.WaterTemple{})))
	storage.CreateColumn(table.BuildColumnOf[components.Wonderitem](columns.NewBit(components.Wonderitem{})))
	storage.CreateColumn(table.BuildColumnOf[components.ZorasDomain](columns.NewBit(components.ZorasDomain{})))
	storage.CreateColumn(table.BuildColumnOf[components.ZorasFountain](columns.NewBit(components.ZorasFountain{})))
	storage.CreateColumn(table.BuildColumnOf[components.ZorasRiver](columns.NewBit(components.ZorasRiver{})))

	return nil
}
