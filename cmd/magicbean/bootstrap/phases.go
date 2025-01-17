package bootstrap

import (
	"slices"
	"sudonters/zootler/cmd/magicbean/z16"
	"sudonters/zootler/internal/settings"
	"sudonters/zootler/mido"
	"sudonters/zootler/mido/compiler"
	"sudonters/zootler/mido/objects"
	"sudonters/zootler/zecs"
)

func panicWhenErr(err error) {
	if err != nil {
		panic(err)
	}
}

func Phase1_InitializeStorage(ddl ...zecs.DDL) zecs.Ocm {
	ocm, err := zecs.New()
	panicWhenErr(err)
	panicWhenErr(zecs.Apply(&ocm, []zecs.DDL{
		nil,
	}))
	panicWhenErr(zecs.Apply(&ocm, ddl))
	return ocm
}

func Phase2_ImportFromFiles(ocm *zecs.Ocm, paths LoadPaths) error {
	tokens := z16.NewTokens(ocm)
	nodes := z16.NewRegions(ocm)
	panicWhenErr(storeScripts(ocm, paths))
	panicWhenErr(storeTokens(tokens, paths))
	panicWhenErr(storeplacements(nodes, tokens, paths))
	panicWhenErr(storeRelations(nodes, tokens, paths))
	return nil
}

func Phase3_ConfigureCompiler(ocm *zecs.Ocm, settings *settings.Zootr, options ...mido.ConfigureCompiler) mido.CompileEnv {
	defaults := []mido.ConfigureCompiler{
		mido.CompilerDefaults(),
		func(env *mido.CompileEnv) {
			panicWhenErr(loadsymbols(ocm, env.Symbols))
			panicWhenErr(loadptrs(ocm, env.Objects))
			panicWhenErr(loadscripts(ocm, env))
			panicWhenErr(aliassymbols(ocm, env.Symbols, env.Functions))
		},
		installSettings(settings),
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
	panic("not implemented")
}

func Phase5_CreateWorld(ocm *zecs.Ocm, settings *settings.Zootr) any {
	g := buildgraph(ocm)
	edges := buildedgetable(ocm)
	_, _ = g, edges
	panic("not implemented")
}
