package ast

// lexer -> tokens -> parser -> parse tree -> ??? -> ast -> ??? -> ir -> compiler -> bytecode
// goal is read all rules and create optimized IR for them that we store as a
// template when a seed is being generated, the compiler will apply last mile
// optimizations -- which should only be inlining settings and light DCE

type Node interface {
	Type() AstNodeType
}

type AstNodeType uint8

// notice that subscript, tuple, and unary op are gone these are all
// transformed into Call ast -- additionally _SOME_ identifier loads may also
// be transformed into Call ast.
type Comparison struct {
	LHS, RHS Node
	Op       AstCompareOp
}

type BooleanOp struct {
	LHS, RHS Node
	Op       AstBoolOp
}

type Call struct {
	Callee string
	Args   []Node
	Macro  bool
}

type Identifier struct {
	Name string
	Kind AstIdentifierKind
}

type Literal struct {
	Value any
	Kind  AstLiteralKind
}

// nil can be legitmately returned in the case of errors however, collapsing
// unary not into boolean and/or leaves a hole on one side -- instead of having
// special out of band nil knowledge we can plug it with something we can tell
// the IR generator "it's always okay to replace this with a NOP" which also
// means its always okay to not insert _any_ bytecode
type Empty struct{}

func IsEmpty(ast Node) bool {
	return ast.Type() == AST_NODE_EMPTY
}

type AstCompareOp uint8
type AstBoolOp uint8
type AstIdentifierKind uint8
type AstLiteralKind uint8

const (
	AST_CMP_EQ AstCompareOp = iota + 1
	AST_CMP_NQ
	AST_CMP_LT

	AST_BOOL_AND    AstBoolOp = 1
	AST_BOOL_OR               = 2
	AST_BOOL_NEGATE           = 3

	AST_IDENT_UNK AstIdentifierKind = 0
	AST_IDENT_TOK                   = 2
	AST_IDENT_SET                   = 4
	AST_IDENT_TRK                   = 5

	AST_NODE_EMPTY AstNodeType = iota
	AST_NODE_CMP
	AST_NODE_BOOL
	AST_NODE_CALL
	AST_NODE_IDENT
	AST_NODE_LITERAL

	AST_LIT_NUM  AstLiteralKind = 1
	AST_LIT_BOOL                = 2
	AST_LIT_STR                 = 3
)

/*
for quickly stamping out things that touch all nodes
Comparison
BooleanOp
Call
Identifier
Literal
Empty
*/

func (n *Comparison) Type() AstNodeType { return AST_NODE_CMP }
func (n *BooleanOp) Type() AstNodeType  { return AST_NODE_BOOL }
func (n *Call) Type() AstNodeType       { return AST_NODE_CALL }
func (n *Identifier) Type() AstNodeType { return AST_NODE_IDENT }
func (n *Literal) Type() AstNodeType    { return AST_NODE_LITERAL }
func (n *Empty) Type() AstNodeType      { return AST_NODE_EMPTY }
