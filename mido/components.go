package mido

import (
	"sudonters/libzootr/mido/ast"
	"sudonters/libzootr/mido/compiler"
)

type RawSource string
type ParsedSource ast.Node
type OptimizedSource ast.Node
type Bytecode compiler.Bytecode
