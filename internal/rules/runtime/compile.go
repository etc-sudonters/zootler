package runtime

import (
	"sudonters/zootler/internal/rules/parser"
	"sudonters/zootler/internal/slipup"
)

func Compile(ast parser.Expression) (*Chunk, error) {
	c := Compiler{new(ChunkBuilder)}
	if err := parser.Visit(c, ast); err != nil {
		return nil, err
	}
	c.chunk.SetReturn()
	c.chunk.Return()
	return &c.chunk.Chunk, nil
}

type Compiler struct {
	chunk *ChunkBuilder
	funcs *FuncNamespace
}

func (c Compiler) VisitBinOp(node *parser.BinOp) error {
	if err := parser.Visit(c, node.Left); err != nil {
		return slipup.Describef(err, "while writing binop %+v", node)
	}
	if err := parser.Visit(c, node.Right); err != nil {
		return slipup.Describef(err, "while writing binop %+v", node)
	}
	switch node.Op {
	case parser.BinOpEq:
		c.chunk.Equal()
		return nil
	case parser.BinOpNotEq:
		c.chunk.NotEqual()
		return nil
	case parser.BinOpLt:
		c.chunk.LessThan()
		return nil
	default:
		return slipup.Createf("unsupported BinOpKind '%s'", node.Op)
	}
}

func (c Compiler) VisitBoolOp(node *parser.BoolOp) error {
	if err := parser.Visit(c, node.Left); err != nil {
		return slipup.Describef(err, "while writing boolop %+v", node)
	}
	if err := parser.Visit(c, node.Right); err != nil {
		return slipup.Describef(err, "while writing boolop %+v", node)
	}

	switch node.Op {
	case parser.BoolOpAnd:
		c.chunk.And()
		return nil
	case parser.BoolOpOr:
		c.chunk.Or()
		return nil
	}

	return slipup.Createf("unsupported boolopkind '%s'", node.Op)
}

func (c Compiler) VisitCall(call *parser.Call) error {
	callee, wasIdent := call.Callee.(*parser.Identifier)
	if !wasIdent {
		return slipup.Createf("expected identifier, received: '%+v'", callee)
	}

	size := len(call.Args)

	for _, arg := range call.Args {
		if err := parser.Visit(c, arg); err != nil {
			return err
		}
	}

	c.chunk.Call(callee.Value, size)
	return nil
}

func (c Compiler) VisitIdentifier(node *parser.Identifier) error {
	if c.funcs.IsFunc(node.Value) {
		c.chunk.Call(node.Value, 0)
		return nil
	}

	c.chunk.LoadIdentifier(node.Value)
	return nil
}

func (c Compiler) VisitSubscript(node *parser.Subscript) error {
	return slipup.NotImplemented("subscripts: settings, etc should be arranged before compile")
}

func (c Compiler) VisitTuple(node *parser.Tuple) error {
	if len(node.Elems) != 2 {
		return slipup.Createf("expected 2 arguments for has, got %d", len(node.Elems))
	}

	parser.Visit(c, node.Elems[0])
	parser.Visit(c, node.Elems[1])
	c.chunk.Call("has", 2)
	return nil
}

func (c Compiler) VisitUnary(node *parser.UnaryOp) error {
	switch node.Op {
	case parser.UnaryNot:
		c.chunk.Not()
		return nil
	default:
		return slipup.Createf("unsuported unary operator: '%+v", node)
	}
}

func (c Compiler) VisitLiteral(node *parser.Literal) error {
	v, err := ValueFrom(node.Value)
	if err != nil {
		return slipup.Describef(err, "%s", node.Kind)
	}
	c.chunk.LoadConst(v)
	return nil
}
