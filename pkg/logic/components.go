package logic

import (
	"sudonters/zootler/pkg/rules/ast"
)

type (
	RawRule    string
	ParsedRule struct {
		R ast.Expression
	}
)
