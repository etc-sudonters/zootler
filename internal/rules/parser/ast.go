package parser

import (
	"math"
)

type BinOpKind string
type BoolOpKind string
type UnaryOpKind string

type AnyNumeric interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

func Identify(name string) *Identifier {
	return &Identifier{Value: name}
}

func TokenLiteral(t uint64) *Literal {
	return NumberLiteral(t)
}

func BoolLiteral(b bool) *Literal {
	return &Literal{
		Kind:  LiteralBool,
		Value: b,
	}
}

func NumberLiteral[N AnyNumeric](n N) *Literal {
	return &Literal{
		Kind:  LiteralNum,
		Value: float64(n),
	}
}

func StringLiteral(s string) *Literal {
	return &Literal{
		Kind:  LiteralStr,
		Value: s,
	}
}
func MakeCall(callee Expression, args []Expression) *Call {
	return &Call{
		Callee: callee,
		Args:   args,
	}
}

func MakeCallSplat(callee Expression, args ...Expression) *Call {
	return MakeCall(callee, args)
}

var (
	BinOpEq       BinOpKind   = "=="
	BinOpNotEq    BinOpKind   = "!="
	BinOpLt       BinOpKind   = "<"
	BinOpContains BinOpKind   = "in"
	BoolOpAnd     BoolOpKind  = "and"
	BoolOpOr      BoolOpKind  = "or"
	UnaryNot      UnaryOpKind = "not"
)

type ExprType string

const (
	ExprBinOp      = "BinOp"
	ExprBoolOp     = "BoolOp"
	ExprCall       = "Call"
	ExprIdentifier = "Identifier"
	ExprSubscript  = "Subscript"
	ExprTuple      = "Tuple"
	ExprUnaryOp    = "UnaryOp"
	ExprLiteral    = "Literal"
)

type LiteralKind string

const (
	LiteralBool LiteralKind = "Boolean"
	LiteralNum              = "Number"
	LiteralStr              = "String"

	// placeholder
	LiteralToken = LiteralNum
)

type (
	Expression interface {
		Type() ExprType
		exprNode()
	}
)

type FunctionDecl struct {
	Body       Expression
	Identifier string
	Parameters []string
}

type (
	BoolOp struct {
		Left  Expression
		Op    BoolOpKind
		Right Expression
	}

	Literal struct {
		Kind  LiteralKind
		Value any
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

	Tuple struct {
		Elems []Expression
	}

	UnaryOp struct {
		Op     UnaryOpKind
		Target Expression
	}
)

func (b *BinOp) exprNode()      {}
func (b *BoolOp) exprNode()     {}
func (c *Call) exprNode()       {}
func (i *Identifier) exprNode() {}
func (s *Subscript) exprNode()  {}
func (t *Tuple) exprNode()      {}
func (u *UnaryOp) exprNode()    {}
func (l *Literal) exprNode()    {}

func (expr *BinOp) Type() ExprType      { return ExprBinOp }
func (expr *BoolOp) Type() ExprType     { return ExprBoolOp }
func (expr *Call) Type() ExprType       { return ExprCall }
func (expr *Identifier) Type() ExprType { return ExprIdentifier }
func (expr *Subscript) Type() ExprType  { return ExprSubscript }
func (expr *Tuple) Type() ExprType      { return ExprTuple }
func (expr *UnaryOp) Type() ExprType    { return ExprUnaryOp }
func (expr *Literal) Type() ExprType    { return ExprLiteral }

func (expr *Literal) AsBool() (bool, bool) {
	if expr.Kind == LiteralBool {
		return expr.Value.(bool), true
	}
	return false, false
}

func (expr *Literal) AsNumber() (float64, bool) {
	if expr.Kind == LiteralNum {
		return expr.Value.(float64), true
	}

	return math.NaN(), false
}

func (expr *Literal) AsString() (string, bool) {
	if expr.Kind == LiteralStr {
		return expr.Value.(string), true
	}
	return "", false
}
