package bootstrap

import (
	"sudonters/zootler/internal/settings"
	"sudonters/zootler/mido"
	"sudonters/zootler/mido/ast"
	"sudonters/zootler/mido/objects"
	"sudonters/zootler/mido/symbols"
	"sudonters/zootler/zecs"
)

func loadsymbols(ocm *zecs.Ocm, syms *symbols.Table) error {
	panic("not implemented")
}

func loadptrs(ocm *zecs.Ocm, objs *objects.Builder) error {
	panic("not implemented")
}

func loadscripts(ocm *zecs.Ocm, env *mido.CompileEnv) error {
	panic("not implemented")
}

func aliassymbols(ocm *zecs.Ocm, syms *symbols.Table, scripts *ast.PartialFunctionTable) error {
	panic("not implemented")
}

func installSettings(_ *settings.Zootr) mido.ConfigureCompiler {
	return func(_ *mido.CompileEnv) {}
}

func installConnectionGenerator(_ *zecs.Ocm) mido.ConfigureCompiler {
	return func(_ *mido.CompileEnv) {}
}
