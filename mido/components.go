package mido

import (
	"sudonters/zootler/mido/ast"
	"sudonters/zootler/mido/compiler"
)

type RawSource string
type ParsedSource ast.Node
type OptimizedSource ast.Node
type Bytecode compiler.Bytecode
