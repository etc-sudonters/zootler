package bootstrap

import (
	"slices"
	"sudonters/zootler/cmd/zootler/z16"
	"sudonters/zootler/internal/settings"
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

func Phase1_InitializeStorage(ddl []zecs.DDL) zecs.Ocm {
	ocm, err := zecs.New()
	PanicWhenErr(err)
	PanicWhenErr(zecs.Apply(&ocm, staticddl()))
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
			PanicWhenErr(loadsymbols(ocm, env.Symbols))
			PanicWhenErr(loadscripts(ocm, env))
			PanicWhenErr(aliassymbols(ocm, env.Symbols))
		},
		installCompilerFunctions(theseSettings),
		installConnectionGenerator(ocm),
		mido.WithBuiltInFunctionDefs(func(*mido.CompileEnv) []objects.BuiltInFunctionDef {
			return magicbean.CreateBuiltInDefs()
		}),
		mido.CompilerWithFastOps(compiler.FastOps{
			"has": compiler.FastHasOp,
		}),
		func(env *mido.CompileEnv) {
			createptrs(ocm, env.Symbols, env.Objects)
		},
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
