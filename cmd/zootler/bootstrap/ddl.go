package bootstrap

import (
	"sudonters/zootler/internal/table"
	"sudonters/zootler/internal/table/columns"
	"sudonters/zootler/magicbean"
	"sudonters/zootler/mido/symbols"
	"sudonters/zootler/zecs"
)

func sizedslice[T zecs.Value](size uint32) zecs.DDL {
	return func() *table.ColumnBuilder {
		return columns.SizedSliceColumn[T](size)
	}
}

func sizedbit[T zecs.Value](size uint32) zecs.DDL {
	return func() *table.ColumnBuilder {
		return columns.SizedBitColumnOf[T](size)
	}
}

func sizedhash[T zecs.Value](capacity uint32) zecs.DDL {
	return func() *table.ColumnBuilder {
		return columns.SizedHashMapColumn[T](capacity)
	}
}

func staticddl() []zecs.DDL {
	return []zecs.DDL{
		sizedslice[magicbean.Name](9000),
		sizedbit[magicbean.Placement](5000),
		sizedhash[symbols.Kind](5000),
		sizedhash[magicbean.Ptr](5000),

		sizedhash[magicbean.RuleParsed](4000),
		sizedhash[magicbean.RuleOptimized](4000),
		sizedhash[magicbean.RuleCompiled](4000),
		sizedhash[magicbean.EdgeKind](4000),
		sizedhash[magicbean.Connection](4000),
		sizedhash[magicbean.RuleSource](4000),
		sizedhash[magicbean.DefaultPlacement](2200),

		columns.HashMapColumn[magicbean.CollectablePriority],
		columns.HashMapColumn[magicbean.HeldAt],
		columns.HashMapColumn[magicbean.HoldsToken],
		columns.HashMapColumn[magicbean.HintRegion],
		columns.HashMapColumn[magicbean.AltHintRegion],
		columns.HashMapColumn[magicbean.DungeonName],
		columns.HashMapColumn[magicbean.Savewarp],
		columns.HashMapColumn[magicbean.Scene],
		columns.HashMapColumn[magicbean.ScriptDecl],
		columns.HashMapColumn[magicbean.ScriptSource],
		columns.HashMapColumn[magicbean.ScriptParsed],
		columns.HashMapColumn[magicbean.AliasingName],
		columns.HashMapColumn[magicbean.HeartPieceCount],
		columns.HashMapColumn[magicbean.OcarinaButton],
		columns.HashMapColumn[magicbean.Song],
		columns.HashMapColumn[magicbean.SmallKeyCount],
		columns.HashMapColumn[magicbean.SmallKeyGroup],
		columns.HashMapColumn[magicbean.SilverRupeeCount],
		columns.HashMapColumn[magicbean.SilverRupeeGroup],

		columns.BitColumnOf[magicbean.Token],
		columns.BitColumnOf[magicbean.Region],
		columns.BitColumnOf[magicbean.IsBossRoom],
		columns.BitColumnOf[magicbean.Empty],
		columns.BitColumnOf[magicbean.Generated],
		columns.BitColumnOf[magicbean.Collectable],
		columns.BitColumnOf[magicbean.Location],
		columns.BitColumnOf[magicbean.TimePassess],
		columns.BitColumnOf[magicbean.BossKey],
		columns.BitColumnOf[magicbean.Compass],
		columns.BitColumnOf[magicbean.Drop],
		columns.BitColumnOf[magicbean.DungeonReward],
		columns.BitColumnOf[magicbean.Event],
		columns.BitColumnOf[magicbean.GanonBossKey],
		columns.BitColumnOf[magicbean.HideoutSmallKey],
		columns.BitColumnOf[magicbean.HideoutSmallKeyRing],
		columns.BitColumnOf[magicbean.Item],
		columns.BitColumnOf[magicbean.Map],
		columns.BitColumnOf[magicbean.Refill],
		columns.BitColumnOf[magicbean.Shop],
		columns.BitColumnOf[magicbean.SilverRupee],
		columns.BitColumnOf[magicbean.SmallKey],
		columns.BitColumnOf[magicbean.SmallKeyRing],
		columns.BitColumnOf[magicbean.TCGSmallKey],
		columns.BitColumnOf[magicbean.TCGSmallKeyRing],
		columns.BitColumnOf[magicbean.GoldSkulltulaToken],
		columns.BitColumnOf[magicbean.Medallion],
		columns.BitColumnOf[magicbean.Stone],
		columns.BitColumnOf[magicbean.Bottle],
	}
}
