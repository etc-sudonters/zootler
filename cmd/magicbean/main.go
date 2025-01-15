package main

import (
	"fmt"
	"slices"
	"strings"
	"sudonters/zootler/internal/settings"
	"sudonters/zootler/mido"
	"sudonters/zootler/mido/ast"
	"sudonters/zootler/mido/compiler"
	"sudonters/zootler/mido/objects"
	"sudonters/zootler/mido/optimizer"
	"sudonters/zootler/mido/symbols"
)

func constBuiltInFunc([]objects.Object) (objects.Object, error) {
	return objects.True, nil
}

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

	seedSettings := settings.Default()
	_ = seedSettings
	analysis := newanalysis()
	compileEnv := mido.NewCompileEnv(
		mido.CompilerWithConnectionGeneration(func(env *mido.CompileEnv) func(*symbols.Sym) {
			return func(s *symbols.Sym) {
				env.Objects.AssociateSymbol(s, objects.PackTaggedPtr32(objects.PtrToken, 0xdeadbeef))
			}
		}),
		mido.CompilerDefaults(),
		mido.CompilerWithTokens(rawTokens),
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
		func(env *mido.CompileEnv) {
			funcBuildErr := env.BuildFunctionTable(ReadHelpers(".data/logic/helpers.json"))
			if funcBuildErr != nil {
				panic(funcBuildErr)
			}
			aliasTokens(env.Symbols, env.Functions, rawTokens)
			analysis.register(env)
			for i := range rawTokens {
				symbol := env.Symbols.LookUpByName(rawTokens[i])
				env.Objects.AssociateSymbol(symbol, objects.PackPtr32(uint32(i)))
			}
			for i, name := range settings.Names() {
				symbol := env.Symbols.LookUpByName(name)
				env.Objects.AssociateSymbol(symbol, objects.PackPtr32(uint32(i)))
			}
		},
	)

	realSource, fakeSource := SourceRules(locations), FakeSourceRules()
	_, _ = realSource, fakeSource
	source := realSource
	codeGen := mido.Compiler(&compileEnv)

	compiled := make([]mido.CompiledSource, len(source))
	var failedCompiles []failedcompile

	for i := range source {
		declaration := &compiled[i]
		declaration.Source = source[i]

		switch declaration.Kind {
		case mido.SourceTransit:
			declaration.Source.Destination = fmt.Sprintf(
				"%s -> %s",
				declaration.OriginatingRegion, declaration.Destination,
			)
		}
		symbol := compileEnv.Symbols.Declare(declaration.Destination, declaration.Kind.AsSymbolKind())
		if declaration.Kind == mido.SourceEvent {
			compileEnv.Symbols.Alias(symbol, escape(declaration.Destination))
			compileEnv.Objects.AssociateSymbol(symbol, objects.PackTaggedPtr32(objects.PtrToken, 0xdeadbeef))
		}
	}

	for i := range source {
		var compileErr error
		compiling := mido.CompiledSource{
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

	connections, connectionErr := mido.CompileGeneratedConnections(&codeGen)
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
	src *mido.Source
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
