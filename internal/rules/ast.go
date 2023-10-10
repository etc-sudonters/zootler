package rules

import (
	"fmt"
	"strings"
)

type (
	Node interface{}

	Expression interface {
		Node
		exprNode()
	}
)

type (
	TotalRule struct {
		Rule Expression
	}

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
		Members []Expression
	}

	String struct {
		Value string
	}
)

type BinOpKind string
type BoolOpKind string

var (
	BinOpLt   BinOpKind  = "<"
	BoolOpAnd BoolOpKind = "and"
	BoolOpOr  BoolOpKind = "or"
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
