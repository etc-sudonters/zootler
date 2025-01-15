package compile

import (
	"slices"
	"sudonters/zootler/mido"
	"sudonters/zootler/mido/ast"
	"sudonters/zootler/mido/compiler"
	"sudonters/zootler/mido/objects"
	"sudonters/zootler/mido/optimizer"
)

func constCompileFunc([]ast.Node, ast.Rewriting) (ast.Node, error) {
	return ast.Boolean(true), nil
}

func CreateCompileEnv(options ...mido.ConfigureCompiler) mido.CompileEnv {
	defaults := []mido.ConfigureCompiler{
		mido.CompilerDefaults(),
		mido.WithCompilerFunctions(func(*mido.CompileEnv) optimizer.CompilerFunctionTable {
			return optimizer.CompilerFunctionTable{
				"load_setting":           constCompileFunc,
				"load_setting_2":         constCompileFunc,
				"compare_setting":        constCompileFunc,
				"region_has_shortcuts":   constCompileFunc,
				"is_trick_enabled":       constCompileFunc,
				"had_night_start":        constCompileFunc,
				"has_all_notes_for_song": constCompileFunc,
				"at_dampe_time":          constCompileFunc,
				"at_day":                 constCompileFunc,
				"at_night":               constCompileFunc,
			}
		}),
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

	/*
		mido.CompilerWithTokens(rawTokens),
		func(env *mido.CompileEnv) {
			funcBuildErr := env.BuildFunctionTable(ReadHelpers(".data/logic/helpers.json"))
			if funcBuildErr != nil {
				panic(funcBuildErr)
			}
			aliasTokens(env.Symbols, env.Functions, rawTokens)
			analysis.register(env)
			for i := range rawTokens {
				env.Objects.AddPointer(rawTokens[i], objects.Pointer(objects.OpaquePointer(i), objects.PtrToken))
			}
			for i, name := range settings.Names() {
				env.Objects.AddPointer(name, objects.Pointer(objects.OpaquePointer(i), objects.PtrSetting))
			}
		},
	*/
	configure := slices.Concat(defaults, options)
	return mido.NewCompileEnv(configure...)
}
