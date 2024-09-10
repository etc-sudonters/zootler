package macros

import (
	"slices"
	"strings"
	"sudonters/zootler/icearrow/parser"

	"github.com/etc-sudonters/substrate/peruse"
)

func CreateScriptedMacro(xpndr *Expansions, decl, body string) error {
	name, args := quickanddirtyDeclParse(decl)
	tokens := slices.Collect(peruse.AllTokens(parser.NewRulesLexer(body)))
	return xpndr.Declare(name, args, tokens, DefaultExpander)
}

// this is similar/the same as upstream OOTR
func quickanddirtyDeclParse(decl string) (string, []string) {
	if !strings.Contains(decl, "(") {
		return decl, nil
	}

	decl = strings.TrimSuffix(decl, ")")
	splitDecl := strings.Split(decl, "(")
	args := strings.Split(splitDecl[1], ",")
	return splitDecl[0], args
}
