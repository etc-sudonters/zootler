package saburo

import (
	"sudonters/zootler/icearrow/parser"
	"sudonters/zootler/internal"

	"github.com/etc-sudonters/substrate/slipup"
)

func LoadScriptedMacros(mb parser.MacroBuilder, path string) error {
	decls, macroScriptReadErr := internal.ReadJsonFileStringMap(path)
	if macroScriptReadErr != nil {
		return slipup.Describe(macroScriptReadErr, "while reading macro file")
	}

	for decl, body := range decls {
		if declErr := mb.AddScriptedMacro(decl, body); declErr != nil {
			return slipup.Describef(declErr, "while adding macro '%s'", decl)
		}
	}

	return nil
}
