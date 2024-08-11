package compiler

import (
	"sudonters/zootler/internal/rules/bytecode"
	"sudonters/zootler/internal/rules/parser"
)

func Compile(ast parser.Expression) (bytecode.Chunk, error) {
	var b bytecode.Chunk
	return b, nil
}
