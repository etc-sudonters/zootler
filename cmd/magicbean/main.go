package main

import (
	"fmt"
	"slices"
	"strings"
	"sudonters/zootler/internal/settings"
	"sudonters/zootler/magicbeanvm"
	"sudonters/zootler/magicbeanvm/optimizer"
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

	seedSettings := settings.Default()
	analysis := newanalysis()
	compileEnv := magicbeanvm.NewCompileEnv(
		magicbeanvm.CompilerDefaults(),
		magicbeanvm.CompilerWithTokens(rawTokens),
		magicbeanvm.CompilerWithFunctions(func(env *magicbeanvm.CompileEnv) optimizer.CompilerFunctions {
			return compfuncs{
				constCompileFuncs(true),
				settingCompilerFuncs{
					settings: &seedSettings,
					symbols:  env.Symbols,
				},
				magicbeanvm.ConnectionGeneration(env.Optimize.Context, env.Symbols),
			}

		}),
		func(env *magicbeanvm.CompileEnv) {
			funcBuildErr := env.BuildFunctionTable(ReadHelpers(".data/logic/helpers.json"))
			if funcBuildErr != nil {
				panic(funcBuildErr)
			}
			aliasTokens(env.Symbols, env.Functions, rawTokens)
			analysis.register(env)
		},
	)

	source, _ := SourceRules(locations), FakeSourceRules()
	codeGen := magicbeanvm.Compiler(&compileEnv)

	compiled := make([]magicbeanvm.CompiledSource, len(source))
	var failedCompiles []failedcompile

	for i := range source {
		declaration := &compiled[i]
		declaration.Source = source[i]

		switch declaration.Kind {
		case magicbeanvm.SourceTransit, magicbeanvm.SourceCheck:
			declaration.Source.Destination = fmt.Sprintf(
				"%s -> %s",
				declaration.OriginatingRegion, declaration.Destination,
			)
		}
		symbol := compileEnv.Symbols.Declare(declaration.Destination, declaration.Kind.AsSymbolKind())
		if declaration.Kind == magicbeanvm.SourceEvent {
			compileEnv.Symbols.Alias(symbol, escape(declaration.Destination))
		}
	}

	for i := range source {
		var compileErr error
		compiling := magicbeanvm.CompiledSource{
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

	connections, connectionErr := magicbeanvm.CompileGeneratedConnections(&codeGen)
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
	src *magicbeanvm.Source
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
