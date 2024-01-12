package logic

import "sudonters/zootler/pkg/rules/parser"

type (
	RawRule    string
	ParsedRule struct {
		R parser.Expression
	}
)
