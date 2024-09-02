package main

import (
	"context"
	"github.com/etc-sudonters/substrate/dontio"
	"github.com/etc-sudonters/substrate/slipup"
	"strings"
	"sudonters/zootler/internal"
	"sudonters/zootler/internal/app"
	"sudonters/zootler/internal/components"
	"sudonters/zootler/internal/query"
	"sudonters/zootler/internal/table"
	"sudonters/zootler/internal/table/columns"
	"sudonters/zootler/internal/table/indexes"
)

type IntoComponents interface {
	GetName() components.Name
	AddComponents(table.RowId, query.Engine) error
}

type AdditionalComponents[T any] interface {
	Init(context.Context, query.Engine) error
	Components(context.Context, table.RowId, T, query.Engine) error
}

type DataFileLoader[T IntoComponents] struct {
	Path      string
	IncludeMQ bool
	Add       AdditionalComponents[T]
}

func (l DataFileLoader[T]) Setup(z *app.Zootlr) error {
	ctx := z.Ctx()
	storage := z.Engine()
	if l.Add != nil {
		if initErr := l.Add.Init(ctx, storage); initErr != nil {
			return slipup.Describef(initErr, "while preparing to read file: '%s'", l.Path)
		}
	}

	dontio.WriteLineOut(ctx, "reading file '%s'", l.Path)
	ts, err := internal.ReadJsonFileAs[[]T](l.Path)
	if err != nil {
		return slipup.Describe(err, l.Path)
	}

	for _, t := range ts {
		name := t.GetName()
		if !l.IncludeMQ && strings.Contains(strings.ToLower(string(name)), "mq") {
			continue
		}

		row, insertErr := storage.InsertRow(name)
		if insertErr != nil {
			return slipup.Describef(insertErr, "while creating row for %+v", t)
		}

		if valuesErr := t.AddComponents(row, storage); valuesErr != nil {
			return slipup.Describef(valuesErr, "while adding core components to %+v", t)
		}

		if l.Add != nil {
			if additonalErr := l.Add.Components(ctx, row, t, storage); additonalErr != nil {
				return slipup.Describef(additonalErr, "while adding additional components to %+v", t)
			}
		}
	}

	return nil
}

type CreateScheme struct {
	DDL []DDL
}
type DDL func() *table.ColumnBuilder

func indexed(ddl DDL, i table.Index) DDL {
	return func() *table.ColumnBuilder {
		return ddl().Index(i)
	}
}

func (cs CreateScheme) Setup(z *app.Zootlr) error {
	dontio.WriteLineOut(z.Ctx(), "running DDL")
	storage := z.Engine()
	for _, ddl := range cs.DDL {
		if _, err := storage.CreateColumn(ddl()); err != nil {
			return err
		}
	}

	return nil
}

func NormalizedNameIndex[T ~string](c T) (internal.NormalizedStr, bool) {
	return internal.Normalize(string(c)), true
}

func MakeDDL() []DDL {
	return []DDL{
		indexed(
			columns.SliceColumn[components.Name],
			indexes.CreateUniqueHashIndex(NormalizedNameIndex[components.Name])),
		indexed(
			columns.HashMapColumn[components.Dungeon],
			indexes.CreateHashIndex(NormalizedNameIndex[components.Dungeon]),
		),
		indexed(
			columns.HashMapColumn[components.HintRegion],
			indexes.CreateHashIndex(func(s components.HintRegion) (internal.NormalizedStr, bool) {
				return internal.Normalize(s.Name), true
			}),
		),

		columns.HashMapColumn[components.CompiledRule],
		columns.HashMapColumn[components.Count],
		columns.HashMapColumn[components.Price],
		columns.HashMapColumn[components.ShopObject],
		columns.HashMapColumn[components.OcarinaButton],
		columns.HashMapColumn[components.OcarinaNote],
		columns.HashMapColumn[components.OcarinaSong],
		columns.HashMapColumn[components.Song],
		columns.HashMapColumn[components.RawLogic],
		columns.HashMapColumn[components.Edge],
		columns.HashMapColumn[components.SavewarpName],
		columns.HashMapColumn[components.DefaultItem],
		columns.HashMapColumn[components.Collected],

		// bit columns only track singletons
		columns.BitColumnOf[components.ExitEdge],
		columns.BitColumnOf[components.CheckEdge],
		columns.BitColumnOf[components.EventEdge],
		columns.BitColumnOf[components.TimePasses],
		columns.BitColumnOf[components.Helper],
		columns.BitColumnOf[components.Advancement],
		columns.BitColumnOf[components.Beehive],
		columns.BitColumnOf[components.BossHeart],
		columns.BitColumnOf[components.BossKey],
		columns.BitColumnOf[components.Boss],
		columns.BitColumnOf[components.BossRoom],
		columns.BitColumnOf[components.Bottle],
		columns.BitColumnOf[components.BottomoftheWellMQ],
		columns.BitColumnOf[components.BottomoftheWell],
		columns.BitColumnOf[components.Chest],
		columns.BitColumnOf[components.CollectableGameToken],
		columns.BitColumnOf[components.Collectable],
		columns.BitColumnOf[components.Compass],
		columns.BitColumnOf[components.Cows],
		columns.BitColumnOf[components.Crate],
		columns.BitColumnOf[components.Cutscene],
		columns.BitColumnOf[components.DeathMountainCrater],
		columns.BitColumnOf[components.DeathMountainTrail],
		columns.BitColumnOf[components.DeathMountain],
		columns.BitColumnOf[components.DekuScrubUpgrades],
		columns.BitColumnOf[components.DekuScrubs],
		columns.BitColumnOf[components.DekuTreeMQ],
		columns.BitColumnOf[components.DekuTree],
		columns.BitColumnOf[components.DesertColossus],
		columns.BitColumnOf[components.DodongosCavernMQ],
		columns.BitColumnOf[components.DodongosCavern],
		columns.BitColumnOf[components.Drop],
		columns.BitColumnOf[components.DungeonReward],
		columns.BitColumnOf[components.Event],
		columns.BitColumnOf[components.FireTempleMQ],
		columns.BitColumnOf[components.FireTemple],
		columns.BitColumnOf[components.FlyingPot],
		columns.BitColumnOf[components.ForestArea],
		columns.BitColumnOf[components.ForestTempleMQ],
		columns.BitColumnOf[components.ForestTemple],
		columns.BitColumnOf[components.Forest],
		columns.BitColumnOf[components.Freestandings],
		columns.BitColumnOf[components.Freestanding],
		columns.BitColumnOf[components.GanonBossKey],
		columns.BitColumnOf[components.GanonsCastleMQ],
		columns.BitColumnOf[components.GanonsCastle],
		columns.BitColumnOf[components.GanonsTower],
		columns.BitColumnOf[components.GerudoTrainingGroundMQ],
		columns.BitColumnOf[components.GerudoTrainingGround],
		columns.BitColumnOf[components.GerudoValley],
		columns.BitColumnOf[components.Gerudo],
		columns.BitColumnOf[components.GerudosFortress],
		columns.BitColumnOf[components.GoldSkulltulaToken],
		columns.BitColumnOf[components.GoldSkulltulas],
		columns.BitColumnOf[components.GoronCity],
		columns.BitColumnOf[components.Graveyard],
		columns.BitColumnOf[components.GreatFairies],
		columns.BitColumnOf[components.GrottoScrub],
		columns.BitColumnOf[components.Grottos],
		columns.BitColumnOf[components.HauntedWasteland],
		columns.BitColumnOf[components.HideoutSmallKey],
		columns.BitColumnOf[components.HintStone],
		columns.BitColumnOf[components.Hint],
		columns.BitColumnOf[components.HyruleCastle],
		columns.BitColumnOf[components.HyruleField],
		columns.BitColumnOf[components.IceCavernMQ],
		columns.BitColumnOf[components.IceCavern],
		columns.BitColumnOf[components.Item],
		columns.BitColumnOf[components.JabuJabusBellyMQ],
		columns.BitColumnOf[components.JabuJabusBelly],
		columns.BitColumnOf[components.Junk],
		columns.BitColumnOf[components.KakarikoVillage],
		columns.BitColumnOf[components.Kakariko],
		columns.BitColumnOf[components.KokiriForest],
		columns.BitColumnOf[components.LakeHylia],
		columns.BitColumnOf[components.Location],
		columns.BitColumnOf[components.LonLonRanch],
		columns.BitColumnOf[components.LostWoods],
		columns.BitColumnOf[components.Map],
		columns.BitColumnOf[components.Market],
		columns.BitColumnOf[components.MaskShop],
		columns.BitColumnOf[components.MasterQuest],
		columns.BitColumnOf[components.Medallion],
		columns.BitColumnOf[components.Minigames],
		columns.BitColumnOf[components.NPC],
		columns.BitColumnOf[components.NeedSpiritualStones],
		columns.BitColumnOf[components.OutsideGanonsCastle],
		columns.BitColumnOf[components.Pot],
		columns.BitColumnOf[components.Priority],
		columns.BitColumnOf[components.Refill],
		columns.BitColumnOf[components.RupeeTower],
		columns.BitColumnOf[components.SacredForestMeadow],
		columns.BitColumnOf[components.Scrub],
		columns.BitColumnOf[components.ShadowTempleMQ],
		columns.BitColumnOf[components.ShadowTemple],
		columns.BitColumnOf[components.Shop],
		columns.BitColumnOf[components.SilverRupee],
		columns.BitColumnOf[components.SkulltulaHouse],
		columns.BitColumnOf[components.SmallCrate],
		columns.BitColumnOf[components.SmallKey],
		columns.BitColumnOf[components.SpiritTempleMQ],
		columns.BitColumnOf[components.SpiritTemple],
		columns.BitColumnOf[components.SpiritualStone],
		columns.BitColumnOf[components.TCGSmallKey],
		columns.BitColumnOf[components.TempleofTime],
		columns.BitColumnOf[components.ThievesHideout],
		columns.BitColumnOf[components.Trade],
		columns.BitColumnOf[components.VanillaDungeons],
		columns.BitColumnOf[components.WaterTempleMQ],
		columns.BitColumnOf[components.WaterTemple],
		columns.BitColumnOf[components.Wonderitem],
		columns.BitColumnOf[components.ZorasDomain],
		columns.BitColumnOf[components.ZorasFountain],
		columns.BitColumnOf[components.ZorasRiver],
	}
}
