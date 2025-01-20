package bootstrap

import (
	"slices"
	"sudonters/zootler/cmd/magicbean/z16"
	"sudonters/zootler/internal/settings"
	"sudonters/zootler/internal/table"
	"sudonters/zootler/internal/table/columns"
	"sudonters/zootler/magicbean"
	"sudonters/zootler/mido"
	"sudonters/zootler/mido/compiler"
	"sudonters/zootler/mido/objects"
	"sudonters/zootler/zecs"
)

func PanicWhenErr(err error) {
	if err != nil {
		panic(err)
	}
}

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

func Phase1_InitializeStorage(ddl []zecs.DDL) zecs.Ocm {
	ocm, err := zecs.New()
	PanicWhenErr(err)
	PanicWhenErr(zecs.Apply(&ocm, []zecs.DDL{
		sizedslice[magicbean.Name](9000),
		sizedbit[magicbean.Placement](5000),
		sizedhash[magicbean.Connection](3300),
		sizedhash[magicbean.RuleSource](2500),
		sizedhash[magicbean.RuleParsed](2500),
		sizedhash[magicbean.RuleOptimized](2500),
		sizedhash[magicbean.RuleCompiled](2500),
		sizedhash[magicbean.DefaultPlacement](2200),

		columns.HashMapColumn[magicbean.Ptr],
		columns.HashMapColumn[magicbean.EdgeKind],
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
		columns.BitColumnOf[magicbean.Song],
		columns.BitColumnOf[magicbean.TCGSmallKey],
		columns.BitColumnOf[magicbean.TCGSmallKeyRing],
		columns.BitColumnOf[magicbean.GoldSkulltulaToken],
	}))
	return ocm
}

func Phase2_ImportFromFiles(ocm *zecs.Ocm, paths LoadPaths) error {
	tokens := z16.NewTokens(ocm)
	nodes := z16.NewNodes(ocm)
	PanicWhenErr(storeScripts(ocm, paths))
	PanicWhenErr(storeTokens(tokens, paths))
	PanicWhenErr(storePlacements(nodes, tokens, paths))
	PanicWhenErr(storeRelations(nodes, tokens, paths))
	return nil
}

func Phase3_ConfigureCompiler(ocm *zecs.Ocm, theseSettings *settings.Zootr, options ...mido.ConfigureCompiler) mido.CompileEnv {
	defaults := []mido.ConfigureCompiler{
		mido.CompilerDefaults(),
		func(env *mido.CompileEnv) {
			PanicWhenErr(loadsymbols(ocm, env.Symbols, env.Objects))
			PanicWhenErr(loadscripts(ocm, env))
			PanicWhenErr(aliassymbols(ocm, env.Symbols))
		},
		installCompilerFunctions(theseSettings),
		installConnectionGenerator(ocm),
		mido.WithBuiltInFunctionDefs(func(*mido.CompileEnv) []objects.BuiltInFunctionDef {
			return []objects.BuiltInFunctionDef{
				{Name: "has", Params: 2},
				{Name: "has_anyof", Params: -1},
				{Name: "has_every", Params: -1},
				{Name: "is_adult", Params: 0},
				{Name: "is_child", Params: 0},
				{Name: "has_bottle", Params: 0},
				{Name: "has_dungeon_rewards", Params: 1},
				{Name: "has_hearts", Params: 1},
				{Name: "has_medallions", Params: 1},
				{Name: "has_stones", Params: 1},
				{Name: "is_starting_age", Params: 0},
			}
		}),
		mido.CompilerWithFastOps(compiler.FastOps{
			"has": compiler.FastHasOp,
		}),
	}
	defaults = slices.Concat(defaults, options)
	return mido.NewCompileEnv(defaults...)
}

func Phase4_Compile(ocm *zecs.Ocm, compiler *mido.CodeGen) error {
	PanicWhenErr(parseall(ocm, compiler))
	PanicWhenErr(optimizeall(ocm, compiler))
	PanicWhenErr(compileall(ocm, compiler))
	return nil
}

func Phase5_CreateWorld(ocm *zecs.Ocm, settings *settings.Zootr, objects objects.Table) magicbean.ExplorableWorld {
	xplore := explorableworldfrom(ocm)
	return xplore
}
