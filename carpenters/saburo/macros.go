package saburo

import (
	"sudonters/zootler/icearrow/parser"
	"sudonters/zootler/internal"
)

func LoadScriptedMacros(mb parser.MacroBuilder, path string) error {
	decls, macroScriptReadErr := internal.ReadJsonFileStringMap(path)
	paniconerr(macroScriptReadErr)

	for decl, body := range decls {
		if declareErr := mb.AddScriptedMacro(decl, body); declareErr != nil {
			paniconerr(declareErr)
		}
	}

	return nil
}
