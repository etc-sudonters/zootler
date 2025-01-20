package main

import (
	"sudonters/zootler/cmd/magicbean/bootstrap"
	"sudonters/zootler/internal/settings"
	"sudonters/zootler/mido"
	"sudonters/zootler/mido/objects"
)

func main() {
	paths := bootstrap.LoadPaths{
		Tokens:     ".data/data/items.json",
		Placements: ".data/data/locations.json",
		Scripts:    ".data/logic/helpers.json",
		Relations:  ".data/logic/glitchless/",
	}
	theseSettings := settings.Default()
	ocm := bootstrap.Phase1_InitializeStorage(nil)
	bootstrap.PanicWhenErr(bootstrap.Phase2_ImportFromFiles(&ocm, paths))

	compileEnv := bootstrap.Phase3_ConfigureCompiler(
		&ocm, &theseSettings,
	)

	codegen := mido.Compiler(&compileEnv)

	bootstrap.PanicWhenErr(bootstrap.Phase4_Compile(
		&ocm, &codegen,
	))

	world := bootstrap.Phase5_CreateWorld(&ocm, &theseSettings, objects.TableFrom(compileEnv.Objects))
	_ = world
}
