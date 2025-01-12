package main

import (
	"fmt"
	"slices"
	"strings"
	"sudonters/zootler/internal/settings"
	"sudonters/zootler/midologic"
	"sudonters/zootler/midologic/ast"
	"sudonters/zootler/midologic/compiler"
	"sudonters/zootler/midologic/objects"
	"sudonters/zootler/midologic/optimizer"
	"sudonters/zootler/midologic/symbols"
)

func main() {
	rawTokens, tokenErr := loadTokensNames(".data/data/items.json")
	if tokenErr != nil {
		panic(tokenErr)
	}

	locations, locationErr := readLogicFiles(".data/logic/glitchless")
	if locationErr != nil {
		panic(locationErr)
	}

	var constCompileFunc optimizer.CompilerFunction = func([]ast.Node, ast.Rewriting) (ast.Node, error) {
		return ast.Boolean(true), nil
	}

	var constBuiltInFunc objects.BuiltInFn = func([]objects.Object) (objects.Object, error) {
		return objects.Boolean(true), nil
	}

	seedSettings := settings.Default()
	_ = seedSettings
	analysis := newanalysis()
	compileEnv := midologic.NewCompileEnv(
		midologic.CompilerWithConnectionGeneration(func(env *midologic.CompileEnv) func(*symbols.Sym) {
			return func(s *symbols.Sym) {
				env.Objects.AddPointer(s.Name, objects.Pointer(objects.OpaquePointer(0xdead), objects.PtrToken))
			}
		}),
		midologic.CompilerDefaults(),
		midologic.CompilerWithTokens(rawTokens),
		midologic.WithCompilerFunctions(func(*midologic.CompileEnv) optimizer.CompilerFunctionTable {
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
		midologic.WithBuiltInFunctions(func(*midologic.CompileEnv) objects.BuiltInFunctions {
			return objects.BuiltInFunctions{
				{Name: "has", Params: 2, Fn: constBuiltInFunc},
				{Name: "has_anyof", Params: -1, Fn: constBuiltInFunc},
				{Name: "has_every", Params: -1, Fn: constBuiltInFunc},
				{Name: "is_adult", Params: 0, Fn: constBuiltInFunc},
				{Name: "is_child", Params: 0, Fn: constBuiltInFunc},
				{Name: "has_bottle", Params: 0, Fn: constBuiltInFunc},
				{Name: "has_dungeon_rewards", Params: 1, Fn: constBuiltInFunc},
				{Name: "has_hearts", Params: 1, Fn: constBuiltInFunc},
				{Name: "has_medallions", Params: 1, Fn: constBuiltInFunc},
				{Name: "has_stones", Params: 1, Fn: constBuiltInFunc},
				{Name: "is_starting_age", Params: 0, Fn: constBuiltInFunc},
			}
		}),
		midologic.CompilerWithFastOps(compiler.FastOps{
			"has": compiler.FastHasOp,
		}),
		func(env *midologic.CompileEnv) {
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
	)

	realSource, fakeSource := SourceRules(locations), FakeSourceRules()
	_, _ = realSource, fakeSource
	source := realSource
	codeGen := midologic.Compiler(&compileEnv)

	compiled := make([]midologic.CompiledSource, len(source))
	var failedCompiles []failedcompile

	for i := range source {
		declaration := &compiled[i]
		declaration.Source = source[i]

		switch declaration.Kind {
		case midologic.SourceTransit:
			declaration.Source.Destination = fmt.Sprintf(
				"%s -> %s",
				declaration.OriginatingRegion, declaration.Destination,
			)
		}
		symbol := compileEnv.Symbols.Declare(declaration.Destination, declaration.Kind.AsSymbolKind())
		if declaration.Kind == midologic.SourceEvent {
			compileEnv.Symbols.Alias(symbol, escape(declaration.Destination))
		}
		compileEnv.Objects.AddPointer(symbol.Name, objects.Pointer(objects.OpaquePointer(0xdead), objects.PtrToken))
	}

	for i := range source {
		var compileErr error
		compiling := midologic.CompiledSource{
			Source: source[i],
		}
		compiling.ByteCode, compileErr = codeGen.CompileSource(&compiling.Source)
		compiled[i] = compiling
		if compileErr != nil {
			failedCompiles = append(failedCompiles, failedcompile{
				err: compileErr,
				src: &compiled[i].Source,
			})
		}
	}

	connections, connectionErr := midologic.CompileGeneratedConnections(&codeGen)
	if connectionErr != nil {
		panic(connectionErr)
	}

	compiled = slices.Concat(compiled, connections)

	ExectuteAll(&compileEnv, compiled)
	DisassembleAll(compiled)
	analysis.Report()
	SymbolReport(compileEnv.Symbols)

	if len(failedCompiles) > 0 {
		fmt.Printf("%04d FAILED COMPILATIONS\n", len(failedCompiles))
		for _, failure := range failedCompiles {
			fmt.Println(failure)
		}
	}
}

type failedcompile struct {
	err error
	src *midologic.Source
}

func (this failedcompile) String() string {
	var str strings.Builder

	fmt.Fprintf(&str, "%q %s -> %s\n", this.src.Kind, this.src.OriginatingRegion, this.src.Destination)
	fmt.Fprintln(&str, this.err.Error())
	if this.src.String != "" {
		fmt.Fprintln(&str, this.src.String)
	}

	return str.String()
}
