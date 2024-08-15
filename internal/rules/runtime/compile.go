package runtime

import (
	"sudonters/zootler/internal/rules/parser"
	"sudonters/zootler/internal/slipup"
)

type Compiler struct {
	funcs   *FuncNamespace
	globals *ExecutionEnvironment
}

type compiling struct {
	*ChunkBuilder
	funcs *FuncNamespace
	env   *ExecutionEnvironment
}

func CompilerUsing(funcs *FuncNamespace, globals *ExecutionEnvironment) Compiler {
	return Compiler{funcs, globals}
}

func NewCompiler() Compiler {
	return Compiler{
		funcs:   NewFuncNamespace(),
		globals: NewEnv(),
	}
}

func (c Compiler) CompileFunctionDecl(decl parser.FunctionDecl) error {
	var f CompiledFunc
	c.funcs.DeclFunction(decl.Identifier)
	f.arity = len(decl.Parameters)
	f.env = c.globals.ChildScope()
	f.chunk = new(Chunk)
	builder := &ChunkBuilder{*f.chunk}

	for _, name := range decl.Parameters {
		builder.DeclareIdentifier(name)
	}

	compiling := compiling{builder, c.funcs, f.env}
	if err := c.compileUnit(decl.Body, compiling); err != nil {
		return err
	}

	c.funcs.AddFunc(decl.Identifier, &f)
	return nil
}

func (c Compiler) CompileEdgeRule(ast parser.Expression) (*Chunk, error) {
	code := compiling{new(ChunkBuilder), c.funcs, c.globals}
	if err := c.compileUnit(ast, code); err != nil {
		return nil, err
	}
	return &code.Chunk, nil
}

func (c Compiler) compileUnit(ast parser.Expression, compiling compiling) error {
	compiling.Preamble()
	if err := parser.Visit(compiling, ast); err != nil {
		return err
	}
	compiling.Epilogue()
	return nil

}

func (c compiling) Preamble() {}

func (c compiling) Epilogue() {
	c.SetReturn()
	c.Return()
}

func (c compiling) VisitBinOp(node *parser.BinOp) error {
	if err := parser.Visit(c, node.Left); err != nil {
		return slipup.Describef(err, "while writing binop %+v", node)
	}
	if err := parser.Visit(c, node.Right); err != nil {
		return slipup.Describef(err, "while writing binop %+v", node)
	}
	switch node.Op {
	case parser.BinOpEq:
		c.Equal()
		return nil
	case parser.BinOpNotEq:
		c.NotEqual()
		return nil
	case parser.BinOpLt:
		c.LessThan()
		return nil
	default:
		return slipup.Createf("unsupported BinOpKind '%s'", node.Op)
	}
}

func (c compiling) VisitBoolOp(node *parser.BoolOp) error {
	if err := parser.Visit(c, node.Left); err != nil {
		return slipup.Describef(err, "while writing boolop %+v", node)
	}
	if err := parser.Visit(c, node.Right); err != nil {
		return slipup.Describef(err, "while writing boolop %+v", node)
	}

	switch node.Op {
	case parser.BoolOpAnd:
		c.And()
		return nil
	case parser.BoolOpOr:
		c.Or()
		return nil
	}

	return slipup.Createf("unsupported boolopkind '%s'", node.Op)
}

func (c compiling) VisitCall(call *parser.Call) error {
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

	c.Call(callee.Value, size)
	return nil
}

func (c compiling) VisitIdentifier(node *parser.Identifier) error {
	if c.funcs.IsFunc(node.Value) {
		c.Call(node.Value, 0)
		return nil
	}

	c.LoadIdentifier(node.Value)
	return nil
}

func (c compiling) VisitSubscript(node *parser.Subscript) error {
	return slipup.NotImplemented("subscripts: settings, etc should be arranged before compile")
}

func (c compiling) VisitTuple(node *parser.Tuple) error {
	if len(node.Elems) != 2 {
		return slipup.Createf("expected 2 arguments for has, got %d", len(node.Elems))
	}

	parser.Visit(c, node.Elems[0])
	parser.Visit(c, node.Elems[1])
	c.Call("has", 2)
	return nil
}

func (c compiling) VisitUnary(node *parser.UnaryOp) error {
	switch node.Op {
	case parser.UnaryNot:
		c.Not()
		return nil
	default:
		return slipup.Createf("unsuported unary operator: '%+v", node)
	}
}

func (c compiling) VisitLiteral(node *parser.Literal) error {
	v, err := ValueFrom(node.Value)
	if err != nil {
		return slipup.Describef(err, "%s", node.Kind)
	}
	c.LoadConst(v)
	return nil
}
