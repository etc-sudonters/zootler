package rules

import (
	"fmt"
	"strings"
)

type AstVisitor interface {
	VisitAttrAccess(*AttrAccess)
	VisitBinOp(*BinOp)
	VisitBoolOp(*BoolOp)
	VisitBoolean(*Boolean)
	VisitCall(*Call)
	VisitIdentifier(*Identifier)
	VisitNumber(*Number)
	VisitString(*String)
	VisitSubscript(*Subscript)
	VisitTuple(*Tuple)
}

type TotalRule struct {
	Rule Expression
}

type (
	Node interface {
		Visit(AstVisitor)
	}

	Expression interface {
		Node
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
		Name Expression
		Args []Expression
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
)

type BinOpKind string
type BoolOpKind string

var (
	BinOpEq    BinOpKind  = "=="
	BinOpNotEq BinOpKind  = "!="
	BinOpLt    BinOpKind  = "<"
	BoolOpAnd  BoolOpKind = "and"
	BoolOpOr   BoolOpKind = "or"
)

func BoolOpFromTok(t Item) BoolOpKind {
	switch s := strings.ToLower(t.Value); s {
	case string(BoolOpAnd):
		return BoolOpAnd
	case string(BoolOpOr):
		return BoolOpOr
	default:
		panic(fmt.Errorf("invalid boolop %q", t))
	}
}

func BinOpFromTok(t Item) BinOpKind {
	switch t.Value {
	case string(BinOpLt):
		return BinOpLt
	case string(BinOpEq):
		return BinOpEq
	case string(BinOpNotEq):
		return BinOpNotEq
	default:
		panic(fmt.Errorf("invalid binop %q", t))
	}
}

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

func (expr *AttrAccess) Visit(v AstVisitor) { v.VisitAttrAccess(expr) }
func (expr *BinOp) Visit(v AstVisitor)      { v.VisitBinOp(expr) }
func (expr *BoolOp) Visit(v AstVisitor)     { v.VisitBoolOp(expr) }
func (expr *Boolean) Visit(v AstVisitor)    { v.VisitBoolean(expr) }
func (expr *Call) Visit(v AstVisitor)       { v.VisitCall(expr) }
func (expr *Identifier) Visit(v AstVisitor) { v.VisitIdentifier(expr) }
func (expr *Number) Visit(v AstVisitor)     { v.VisitNumber(expr) }
func (expr *String) Visit(v AstVisitor)     { v.VisitString(expr) }
func (expr *Subscript) Visit(v AstVisitor)  { v.VisitSubscript(expr) }
func (expr *Tuple) Visit(v AstVisitor)      { v.VisitTuple(expr) }
