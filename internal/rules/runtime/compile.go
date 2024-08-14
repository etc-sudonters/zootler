package runtime

import "sudonters/zootler/internal/rules/parser"

func Compile(ast parser.Expression) (*Chunk, error) {
	c := Compiler{new(ChunkBuilder)}
	parser.Visit(c, ast)
	return &c.Chunk, nil
}

type Compiler struct {
	*ChunkBuilder
}

func (c Compiler) VisitBinOp(node *parser.BinOp) error {
	panic("not implemented") // TODO: Implement
}

func (c Compiler) VisitBoolOp(node *parser.BoolOp) error {
	panic("not implemented") // TODO: Implement
}

func (c Compiler) VisitCall(call *parser.Call) error {
	panic("not implemented") // TODO: Implement
}

func (c Compiler) VisitIdentifier(node *parser.Identifier) error {
	panic("not implemented") // TODO: Implement
}

func (c Compiler) VisitSubscript(node *parser.Subscript) error {
	panic("not implemented") // TODO: Implement
}

func (c Compiler) VisitTuple(node *parser.Tuple) error {
	panic("not implemented") // TODO: Implement
}

func (c Compiler) VisitUnary(node *parser.UnaryOp) error {
	panic("not implemented") // TODO: Implement
}

func (c Compiler) VisitLiteral(node *parser.Literal) error {
	panic("not implemented") // TODO: Implement
}
