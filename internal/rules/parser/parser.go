package parser

import (
	"fmt"
	"strconv"
	"strings"
	"github.com/etc-sudonters/substrate/slipup"

	"github.com/etc-sudonters/substrate/peruse"
)

const (
	_ peruse.Precedence = iota
	LOWEST
	OR
	AND
	NOT
	EQ
	INDEX
	PARENS
)

func ParseFunctionDecl(decl, body string) (FunctionDecl, error) {
	var f FunctionDecl

	funcDecl, funcDeclErr := Parse(decl)
	if funcDeclErr != nil {
		return f, slipup.Describe(funcDeclErr, "while parsing function decl")
	}

	switch d := funcDecl.(type) {
	case *Identifier:
		f.Identifier = d.Value
		break
	case *Call:
		ident, wasIdent := d.Callee.(*Identifier)
		if !wasIdent {
			return f, slipup.Createf("unsupported function decl identifier: %v", d)
		}

		f.Identifier = ident.Value
		f.Parameters = make([]string, len(d.Args))

		for i := range d.Args {
			switch a := d.Args[i].(type) {
			case *Identifier:
				f.Parameters[i] = a.Value
				break
			default:
				return f, slipup.Createf("unsupported function parameter identifier: %v", d)
			}
		}

		break
	default:
		return f, slipup.Createf("unsupported function decl identifier: %v", d)
	}

	funcBody, funcBodyErr := Parse(body)
	if funcBodyErr != nil {
		return f, slipup.Describe(funcBodyErr, "while parsing function body")
	}

	f.Body = funcBody
	return f, nil
}

func Parse(raw string) (Expression, error) {
	l := NewRulesLexer(raw)
	p := NewRulesParser(l)
	return p.Parse()
}

func NewRulesParser(l *peruse.StringLexer) *peruse.Parser[Expression] {
	return peruse.NewParser(NewRulesGrammar(), l)
}

func NewRulesGrammar() peruse.Grammar[Expression] {
	g := peruse.NewGrammar[Expression]()

	g.Parse(TokenTrue, parseBool)
	g.Parse(TokenFalse, parseBool)
	g.Parse(TokenIdentifier, parseIdentifierExpr)
	g.Parse(TokenNumber, parseNumber)
	g.Parse(TokenOpenParen, parseParenExpr)
	g.Parse(TokenString, parseString)
	g.Parse(TokenUnaryNot, parsePrefixNot)

	g.Infix(OR, parseBoolOpExpr, TokenOr)
	g.Infix(AND, parseBoolOpExpr, TokenAnd)
	g.Infix(EQ, parseBinOp, TokenEq, TokenNotEq, TokenLt, TokenContains)
	g.Infix(INDEX, parseSubscript, TokenOpenBracket)
	g.Infix(PARENS, parseCall, TokenOpenParen)

	return g
}

func parseParenExpr(p *peruse.Parser[Expression]) (Expression, error) {
	p.Consume()
	e, err := p.ParseAt(LOWEST)
	if err != nil {
		return nil, err
	}

	if p.Next.Is(TokenComma) {
		e, err = parseTuple(p, e)
		if err != nil {
			return nil, err
		}
	}

	if !p.Expect(TokenCloseParen) {
		return nil, fmt.Errorf(
			"PARENEXPR: expected %q but got %q",
			TokenTypeString(TokenCloseParen),
			p.Next,
		)
	}

	return e, nil
}

func parseTuple(p *peruse.Parser[Expression], left Expression) (Expression, error) {
	elems := []Expression{left}

	for p.Expect(TokenComma) {
		p.Consume()
		elem, err := p.ParseAt(LOWEST)
		if err != nil {
			return nil, err
		}

		elems = append(elems, elem)
	}

	return &Tuple{Elems: elems}, nil
}

func parseIdentifierExpr(p *peruse.Parser[Expression]) (Expression, error) {
	return &Identifier{p.Cur.Literal}, nil
}

func parseBoolOpExpr(p *peruse.Parser[Expression], left Expression, parentPrecedence peruse.Precedence) (Expression, error) {
	thisTok := p.Cur
	p.Consume()
	right, err := p.ParseAt(parentPrecedence)
	if err != nil {
		return nil, err
	}
	b := BoolOp{
		Left:  left,
		Op:    BoolOpFromTok(thisTok),
		Right: right,
	}

	return &b, nil
}

func parseBinOp(p *peruse.Parser[Expression], left Expression, bp peruse.Precedence) (Expression, error) {
	thisTok := p.Cur
	p.Consume()
	right, err := p.ParseAt(bp)
	if err != nil {
		return nil, err
	}

	b := BinOp{
		Left:  left,
		Op:    BinOpFromTok(thisTok),
		Right: right,
	}

	return &b, nil
}

func parseCall(p *peruse.Parser[Expression], left Expression, bp peruse.Precedence) (Expression, error) {
	if p.Expect(TokenCloseParen) { // fn()
		return &Call{Callee: left}, nil
	}

	var args []Expression
	for {
		p.Consume()
		arg, err := p.ParseAt(LOWEST)
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
		if !p.Expect(TokenComma) {
			break
		}
	}

	if err := p.ExpectOrError(TokenCloseParen); err != nil {
		return nil, peruse.UnexpectedAt("FNCALL", err)
	}

	c := Call{Callee: left, Args: args}
	return &c, nil
}

func parseSubscript(p *peruse.Parser[Expression], left Expression, bp peruse.Precedence) (Expression, error) {
	p.Consume()
	index, err := p.ParseAt(LOWEST)
	if err != nil {
		return nil, err
	}

	if err = p.ExpectOrError(TokenCloseBracket); err != nil {
		return nil, peruse.UnexpectedAt("SUBSCRIPT", err)
	}

	s := Subscript{Target: left, Index: index}
	return &s, nil
}

func parseString(p *peruse.Parser[Expression]) (Expression, error) {
	s := &Literal{Value: p.Cur.Literal, Kind: LiteralStr}
	return s, nil
}

func parseNumber(p *peruse.Parser[Expression]) (Expression, error) {
	n, err := strconv.ParseFloat(p.Cur.Literal, 64)
	if err != nil {
		return nil, fmt.Errorf("cannot parse %q as number", p.Cur)
	}
	return &Literal{Value: n, Kind: LiteralNum}, nil
}

func parseBool(p *peruse.Parser[Expression]) (Expression, error) {
	return &Literal{Value: p.Cur.Literal == trueWord, Kind: LiteralBool}, nil
}

func parsePrefixNot(p *peruse.Parser[Expression]) (Expression, error) {
	thisTok := p.Cur
	p.Consume()
	target, err := p.ParseAt(NOT)
	if err != nil {
		return nil, err
	}
	switch thisTok.Literal {
	case notWord:
		u := UnaryOp{
			Op:     UnaryNot,
			Target: target,
		}
		return &u, nil
	default:
		return nil, fmt.Errorf("unexpected unary op %q", thisTok)
	}
}

func UnaryOpFromTok(t peruse.Token) UnaryOpKind {
	switch t.Literal {
	case string(UnaryNot):
		return UnaryNot
	default:
		panic(fmt.Errorf("invalid unaryop %q", t))
	}
}

func BoolOpFromTok(t peruse.Token) BoolOpKind {
	switch s := strings.ToLower(t.Literal); s {
	case string(BoolOpAnd):
		return BoolOpAnd
	case string(BoolOpOr):
		return BoolOpOr
	default:
		panic(fmt.Errorf("invalid boolop %q", t))
	}
}

func BinOpFromTok(t peruse.Token) BinOpKind {
	switch t.Literal {
	case string(BinOpLt):
		return BinOpLt
	case string(BinOpEq):
		return BinOpEq
	case string(BinOpNotEq):
		return BinOpNotEq
	case string(BinOpContains):
		return BinOpContains
	default:
		panic(fmt.Errorf("invalid binop %q", t))
	}
}
