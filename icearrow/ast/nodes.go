package ast

// lexer -> tokens -> parser -> parse tree -> ??? -> ast -> ??? -> ir -> compiler -> bytecode
// goal is read all rules and create optimized IR for them that we store as a
// template when a seed is being generated, the compiler will apply last mile
// optimizations -- which should only be inlining settings and light DCE

type Node interface {
	Type() AstNodeType
	String() string
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
	AST_CMP_EQ AstCompareOp = 1
	AST_CMP_NQ              = 2
	AST_CMP_LT              = 3

	AST_BOOL_AND    AstBoolOp = 4
	AST_BOOL_OR               = 5
	AST_BOOL_NEGATE           = 6

	AST_IDENT_UNK AstIdentifierKind = 0x00
	AST_IDENT_EXP                   = 0x01
	AST_IDENT_TOK                   = 0x02
	AST_IDENT_VAR                   = 0x03
	AST_IDENT_SET                   = 0x04
	AST_IDENT_TRK                   = 0x05
	AST_IDENT_BIF                   = 0x06
	AST_IDENT_EVT                   = 0x07
	AST_IDENT_SYM                   = 0x08
	AST_IDENT_UNP                   = 0xFF

	AST_NODE_EMPTY AstNodeType = iota
	AST_NODE_CMP
	AST_NODE_BOOL
	AST_NODE_CALL
	AST_NODE_IDENT
	AST_NODE_LITERAL

	AST_LIT_NUM  AstLiteralKind = 0x01
	AST_LIT_BOOL                = 0x02
	AST_LIT_STR                 = 0x03
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

func (n *Comparison) String() string {
	var r AstRender
	Visit(&r, n)
	return r.String()
}

func (n *BooleanOp) String() string {
	var r AstRender
	Visit(&r, n)
	return r.String()
}

func (n *Call) String() string {
	var r AstRender
	Visit(&r, n)
	return r.String()
}

func (n *Identifier) String() string {
	var r AstRender
	Visit(&r, n)
	return r.String()
}

func (n *Literal) String() string {
	var r AstRender
	Visit(&r, n)
	return r.String()
}

func (n *Empty) String() string {
	var r AstRender
	Visit(&r, n)
	return r.String()
}

func LiteralBool(b bool) *Literal {
	l := new(Literal)
	l.Kind = AST_LIT_BOOL
	l.Value = b
	return l
}

func LiteralNumber(f64 float64) *Literal {
	l := new(Literal)
	l.Kind = AST_LIT_NUM
	l.Value = f64
	return l
}
