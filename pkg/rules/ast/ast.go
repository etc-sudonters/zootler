package ast

type BinOpKind string
type BoolOpKind string
type UnaryOpKind string

var (
	BinOpEq    BinOpKind   = "=="
	BinOpNotEq BinOpKind   = "!="
	BinOpLt    BinOpKind   = "<"
	BoolOpAnd  BoolOpKind  = "and"
	BoolOpOr   BoolOpKind  = "or"
	UnaryNot   UnaryOpKind = "not"
)

type ExprType string

const (
	ExprAttrAccess ExprType = "AttrAccess"
	ExprBinOp               = "BinOp"
	ExprBoolOp              = "BoolOp"
	ExprBoolean             = "Boolean"
	ExprCall                = "Call"
	ExprIdentifier          = "Identifier"
	ExprNumber              = "Number"
	ExprString              = "String"
	ExprSubscript           = "Subscript"
	ExprTuple               = "Tuple"
	ExprUnaryOp             = "UnaryOp"
)

type (
	Expression interface {
		Type() ExprType
		exprNode()
	}
)

type (
	Boolean struct {
		Value bool
	}

	BoolOp struct {
		Left  Expression
		Op    BoolOpKind
		Right Expression
	}

	Number struct {
		Value float64
	}

	Identifier struct {
		Value string
	}

	BinOp struct {
		Left  Expression
		Op    BinOpKind
		Right Expression
	}

	Call struct {
		Callee Expression
		Args   []Expression
	}

	Subscript struct {
		Target Expression
		Index  Expression
	}

	AttrAccess struct {
		Target Expression
		Attr   Expression
	}

	Tuple struct {
		Elems []Expression
	}

	String struct {
		Value string
	}

	UnaryOp struct {
		Op     UnaryOpKind
		Target Expression
	}
)

func (a *AttrAccess) exprNode() {}
func (b *BinOp) exprNode()      {}
func (b *BoolOp) exprNode()     {}
func (c *Boolean) exprNode()    {}
func (c *Call) exprNode()       {}
func (i *Identifier) exprNode() {}
func (n *Number) exprNode()     {}
func (s *String) exprNode()     {}
func (s *Subscript) exprNode()  {}
func (t *Tuple) exprNode()      {}
func (u *UnaryOp) exprNode()    {}

func (expr *AttrAccess) Visit(v Visitor) { v.VisitAttrAccess(expr) }
func (expr *BinOp) Visit(v Visitor)      { v.VisitBinOp(expr) }
func (expr *BoolOp) Visit(v Visitor)     { v.VisitBoolOp(expr) }
func (expr *Boolean) Visit(v Visitor)    { v.VisitBoolean(expr) }
func (expr *Call) Visit(v Visitor)       { v.VisitCall(expr) }
func (expr *Identifier) Visit(v Visitor) { v.VisitIdentifier(expr) }
func (expr *Number) Visit(v Visitor)     { v.VisitNumber(expr) }
func (expr *String) Visit(v Visitor)     { v.VisitString(expr) }
func (expr *Subscript) Visit(v Visitor)  { v.VisitSubscript(expr) }
func (expr *Tuple) Visit(v Visitor)      { v.VisitTuple(expr) }
func (expr *UnaryOp) Visit(v Visitor)    { v.VisitUnary(expr) }

func (expr *AttrAccess) Type() ExprType { return ExprAttrAccess }
func (expr *BinOp) Type() ExprType      { return ExprBinOp }
func (expr *BoolOp) Type() ExprType     { return ExprBoolOp }
func (expr *Boolean) Type() ExprType    { return ExprBoolean }
func (expr *Call) Type() ExprType       { return ExprCall }
func (expr *Identifier) Type() ExprType { return ExprIdentifier }
func (expr *Number) Type() ExprType     { return ExprNumber }
func (expr *String) Type() ExprType     { return ExprString }
func (expr *Subscript) Type() ExprType  { return ExprSubscript }
func (expr *Tuple) Type() ExprType      { return ExprTuple }
func (expr *UnaryOp) Type() ExprType    { return ExprUnaryOp }
