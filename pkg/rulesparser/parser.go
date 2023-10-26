package rulesparser

import (
	"fmt"
	"strconv"
	"sudonters/zootler/internal/peruse"
)

const (
	_ peruse.Precedence = iota
	LOWEST
	AND
	ACCESS
	EQ
	LT
	PREFIX
	FUNC
)

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
	g.Parse(TokenUnary, parsePrefixUnaryOp)

	g.Infix(AND, parseBoolOpExpr, TokenAnd, TokenOr)
	g.Infix(LOWEST, parseAttrAccess, TokenDot)
	g.Infix(EQ, parseBinOp, TokenEq, TokenNotEq)
	g.Infix(FUNC, parseCall, TokenOpenParen)
	g.Infix(FUNC, parseSubscript, TokenOpenBracket)

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

	if err := p.ExpectOrError(TokenCloseParen); err != nil {
		return nil, peruse.UnexpectedAt("PARENEXPR", err)
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

	return &Tuple{elems}, nil
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
		return &Call{Name: left}, nil
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

	c := Call{Name: left, Args: args}
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
	s := &String{Value: p.Cur.Literal}
	return s, nil
}

func parseAttrAccess(p *peruse.Parser[Expression], target Expression, bp peruse.Precedence) (Expression, error) {
	p.Consume()
	attr, err := p.ParseAt(bp)
	if err != nil {
		return nil, err
	}
	return &AttrAccess{
		Target: target,
		Attr:   attr,
	}, nil
}

func parseNumber(p *peruse.Parser[Expression]) (Expression, error) {
	n, err := strconv.ParseFloat(p.Cur.Literal, 64)
	if err != nil {
		return nil, fmt.Errorf("cannot parse %q as number", p.Cur)
	}
	return &Number{n}, nil
}

func parseBool(p *peruse.Parser[Expression]) (Expression, error) {
	switch p.Cur.Literal {
	case trueWord:
		return &Boolean{Value: true}, nil
	case falseWord:
		return &Boolean{Value: false}, nil
	default:
		return nil, fmt.Errorf("unexpected boolean value %s", p.Cur)
	}
}

func parsePrefixUnaryOp(p *peruse.Parser[Expression]) (Expression, error) {
	thisTok := p.Cur
	p.Consume()
	target, err := p.ParseAt(PREFIX)
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
