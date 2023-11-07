package parser

import (
	"fmt"
	"strconv"
	"strings"
	"sudonters/zootler/internal/peruse"
	"sudonters/zootler/pkg/rules/ast"
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

func Parse(raw string) (ast.Expression, error) {
	l := NewRulesLexer(raw)
	p := NewRulesParser(l)
	return p.Parse()
}

func NewRulesParser(l *peruse.StringLexer) *peruse.Parser[ast.Expression] {
	return peruse.NewParser(NewRulesGrammar(), l)
}

func NewRulesGrammar() peruse.Grammar[ast.Expression] {
	g := peruse.NewGrammar[ast.Expression]()

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

func parseParenExpr(p *peruse.Parser[ast.Expression]) (ast.Expression, error) {
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

func parseTuple(p *peruse.Parser[ast.Expression], left ast.Expression) (ast.Expression, error) {
	elems := []ast.Expression{left}

	for p.Expect(TokenComma) {
		p.Consume()
		elem, err := p.ParseAt(LOWEST)
		if err != nil {
			return nil, err
		}

		elems = append(elems, elem)
	}

	return &ast.Tuple{Elems: elems}, nil
}

func parseIdentifierExpr(p *peruse.Parser[ast.Expression]) (ast.Expression, error) {
	return &ast.Identifier{p.Cur.Literal}, nil
}

func parseBoolOpExpr(p *peruse.Parser[ast.Expression], left ast.Expression, parentPrecedence peruse.Precedence) (ast.Expression, error) {
	thisTok := p.Cur
	p.Consume()
	right, err := p.ParseAt(parentPrecedence)
	if err != nil {
		return nil, err
	}
	b := ast.BoolOp{
		Left:  left,
		Op:    BoolOpFromTok(thisTok),
		Right: right,
	}

	return &b, nil
}

func parseBinOp(p *peruse.Parser[ast.Expression], left ast.Expression, bp peruse.Precedence) (ast.Expression, error) {
	thisTok := p.Cur
	p.Consume()
	right, err := p.ParseAt(bp)
	if err != nil {
		return nil, err
	}

	b := ast.BinOp{
		Left:  left,
		Op:    BinOpFromTok(thisTok),
		Right: right,
	}

	return &b, nil
}

func parseCall(p *peruse.Parser[ast.Expression], left ast.Expression, bp peruse.Precedence) (ast.Expression, error) {
	if p.Expect(TokenCloseParen) { // fn()
		return &ast.Call{Callee: left}, nil
	}

	var args []ast.Expression
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

	c := ast.Call{Callee: left, Args: args}
	return &c, nil
}

func parseSubscript(p *peruse.Parser[ast.Expression], left ast.Expression, bp peruse.Precedence) (ast.Expression, error) {
	p.Consume()
	index, err := p.ParseAt(LOWEST)
	if err != nil {
		return nil, err
	}

	if err = p.ExpectOrError(TokenCloseBracket); err != nil {
		return nil, peruse.UnexpectedAt("SUBSCRIPT", err)
	}

	s := ast.Subscript{Target: left, Index: index}
	return &s, nil
}

func parseString(p *peruse.Parser[ast.Expression]) (ast.Expression, error) {
	s := &ast.Literal{Value: p.Cur.Literal, Kind: ast.LiteralStr}
	return s, nil
}

func parseNumber(p *peruse.Parser[ast.Expression]) (ast.Expression, error) {
	n, err := strconv.ParseFloat(p.Cur.Literal, 64)
	if err != nil {
		return nil, fmt.Errorf("cannot parse %q as number", p.Cur)
	}
	return &ast.Literal{Value: n, Kind: ast.LiteralNum}, nil
}

func parseBool(p *peruse.Parser[ast.Expression]) (ast.Expression, error) {
	return &ast.Literal{Value: p.Cur.Literal == trueWord, Kind: ast.LiteralBool}, nil
}

func parsePrefixNot(p *peruse.Parser[ast.Expression]) (ast.Expression, error) {
	thisTok := p.Cur
	p.Consume()
	target, err := p.ParseAt(NOT)
	if err != nil {
		return nil, err
	}
	switch thisTok.Literal {
	case notWord:
		u := ast.UnaryOp{
			Op:     ast.UnaryNot,
			Target: target,
		}
		return &u, nil
	default:
		return nil, fmt.Errorf("unexpected unary op %q", thisTok)
	}
}

func UnaryOpFromTok(t peruse.Token) ast.UnaryOpKind {
	switch t.Literal {
	case string(ast.UnaryNot):
		return ast.UnaryNot
	default:
		panic(fmt.Errorf("invalid unaryop %q", t))
	}
}

func BoolOpFromTok(t peruse.Token) ast.BoolOpKind {
	switch s := strings.ToLower(t.Literal); s {
	case string(ast.BoolOpAnd):
		return ast.BoolOpAnd
	case string(ast.BoolOpOr):
		return ast.BoolOpOr
	default:
		panic(fmt.Errorf("invalid boolop %q", t))
	}
}

func BinOpFromTok(t peruse.Token) ast.BinOpKind {
	switch t.Literal {
	case string(ast.BinOpLt):
		return ast.BinOpLt
	case string(ast.BinOpEq):
		return ast.BinOpEq
	case string(ast.BinOpNotEq):
		return ast.BinOpNotEq
	case string(ast.BinOpContains):
		return ast.BinOpContains
	default:
		panic(fmt.Errorf("invalid binop %q", t))
	}
}
