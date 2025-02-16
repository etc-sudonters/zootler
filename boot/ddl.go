package boot

import (
	"sudonters/libzootr/components"
	"sudonters/libzootr/internal/table"
	"sudonters/libzootr/internal/table/columns"
	"sudonters/libzootr/mido/symbols"
	"sudonters/libzootr/zecs"
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
		sizedslice[components.Name](9000),
		sizedbit[components.PlacementLocationMarker](5000),
		sizedhash[symbols.Kind](5000),
		sizedhash[components.Ptr](5000),

		sizedhash[components.RuleParsed](4000),
		sizedhash[components.RuleOptimized](4000),
		sizedhash[components.RuleCompiled](4000),
		sizedhash[components.EdgeKind](4000),
		sizedhash[components.Connection](4000),
		sizedhash[components.RuleSource](4000),
		sizedhash[components.DefaultPlacement](2200),

		columns.HashMapColumn[components.CollectablePriority],
		columns.HashMapColumn[components.HoldsToken],
		columns.HashMapColumn[components.HintRegion],
		columns.HashMapColumn[components.AltHintRegion],
		columns.HashMapColumn[components.DungeonName],
		columns.HashMapColumn[components.Savewarp],
		columns.HashMapColumn[components.Scene],
		columns.HashMapColumn[components.ScriptDecl],
		columns.HashMapColumn[components.ScriptSource],
		columns.HashMapColumn[components.ScriptParsed],
		columns.HashMapColumn[components.AliasingName],
		columns.HashMapColumn[components.OcarinaNote],
		columns.HashMapColumn[components.SongNotes],
		columns.HashMapColumn[components.DungeonGroup],
		columns.HashMapColumn[components.SilverRupeePuzzle],
		columns.HashMapColumn[components.Song],
		columns.HashMapColumn[components.Price],

		columns.BitColumnOf[components.Skipped],
		columns.BitColumnOf[components.Collected],
		columns.BitColumnOf[components.TokenMarker],
		columns.BitColumnOf[components.RegionMarker],
		columns.BitColumnOf[components.IsBossRoom],
		columns.BitColumnOf[components.Empty],
		columns.BitColumnOf[components.Generated],
		columns.BitColumnOf[components.Collectable],
		columns.BitColumnOf[components.LocationMarker],
		columns.BitColumnOf[components.TimePassess],
		columns.BitColumnOf[components.Compass],
		columns.BitColumnOf[components.Drop],
		columns.BitColumnOf[components.DungeonReward],
		columns.BitColumnOf[components.Event],
		columns.BitColumnOf[components.Item],
		columns.BitColumnOf[components.Map],
		columns.BitColumnOf[components.Refill],
		columns.BitColumnOf[components.Shop],
		columns.BitColumnOf[components.SilverRupee],
		columns.BitColumnOf[components.SilverRupeePouch],
		columns.BitColumnOf[components.SmallKey],
		columns.BitColumnOf[components.BossKey],
		columns.BitColumnOf[components.DungeonKeyRing],
		columns.BitColumnOf[components.GoldSkulltulaToken],
		columns.BitColumnOf[components.Medallion],
		columns.BitColumnOf[components.Stone],
		columns.BitColumnOf[components.Bottle],
		columns.BitColumnOf[components.WorldGraphRoot],
	}
}
