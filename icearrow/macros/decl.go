package macros

import "strings"

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
