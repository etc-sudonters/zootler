package bootstrap

import (
	"context"
	"os"
	"slices"
	"sudonters/libzootr/magicbean"
	"sudonters/libzootr/magicbean/tracking"
	"sudonters/libzootr/mido"
	"sudonters/libzootr/mido/ast"
	"sudonters/libzootr/mido/objects"
	"sudonters/libzootr/mido/optimizer"
	"sudonters/libzootr/settings"
	"sudonters/libzootr/zecs"
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

func Phase2_ImportFromFiles(ctx context.Context, settings *settings.Model, ocm *zecs.Ocm, set *tracking.Set, paths LoadPaths) error {
	PanicWhenErr(storeScripts(ocm, paths))
	PanicWhenErr(storeTokens(set.Tokens, paths))
	PanicWhenErr(storePlacements(set.Nodes, set.Tokens, paths))
	PanicWhenErr(storeRelations(set.Nodes, set.Tokens, paths))

	if paths.Spoiler != "" {
		fh, err := os.Open(paths.Spoiler)
		PanicWhenErr(err)
		defer func() {
			fh.Close()
		}()
		PanicWhenErr(LoadSpoilerData(ctx, fh, settings, &set.Nodes, &set.Tokens))
	}
	return nil
}

func Phase3_ConfigureCompiler(ocm *zecs.Ocm, theseSettings *settings.Model, options ...mido.ConfigureCompiler) mido.CompileEnv {
	defaults := []mido.ConfigureCompiler{
		mido.CompilerDefaults(),
		mido.CompilerWithGenerationSettings(settings.Names()),
		func(env *mido.CompileEnv) {
			env.Optimize.AddOptimizer(func(env *mido.CompileEnv) ast.Rewriter {
				reader := settings.Reader{Model: theseSettings}
				return optimizer.InlineSettings(reader, env.Symbols)
			})
			PanicWhenErr(loadsymbols(ocm, env.Symbols))
			PanicWhenErr(loadscripts(ocm, env))
			PanicWhenErr(aliassymbols(ocm, env.Symbols))
		},
		installCompilerFunctions(theseSettings),
		installConnectionGenerator(ocm),
		mido.WithBuiltInFunctionDefs(func(*mido.CompileEnv) []objects.BuiltInFunctionDef {
			return magicbean.CreateBuiltInDefs()
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
