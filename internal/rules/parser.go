package rules

import "fmt"

type unexpectedToken struct {
	have Item
}

func (u unexpectedToken) Error() string {
	return fmt.Sprintf("unexpected token %q", u.have)
}

type expectedToken struct {
	have   Item
	wanted ItemType
}

func (e expectedToken) Error() string {
	return fmt.Sprintf("expected %q but found %q", e.wanted, e.have)
}

type prefixParselet func(precedence) (Expression, error)
type infixParselet func(Expression, precedence) (Expression, error)
type precedence int

type parser struct {
	lexer            *lexer
	curTok           Item
	nextTok          Item
	prefixParselets  map[ItemType]prefixParselet
	infixParselets   map[ItemType]infixParselet
	prefixPrecedence map[ItemType]precedence
	infixPrecedence  map[ItemType]precedence
}

func NewParser(l *lexer) *parser {
	p := &parser{
		lexer:            l,
		prefixParselets:  make(map[ItemType]prefixParselet),
		infixParselets:   make(map[ItemType]infixParselet),
		prefixPrecedence: make(map[ItemType]precedence),
		infixPrecedence:  make(map[ItemType]precedence),
	}

	addPrefixParselet := func(t ItemType, pr precedence, parselet prefixParselet) {
		p.prefixParselets[t] = parselet
		p.prefixPrecedence[t] = pr
	}

	addInfixParselet := func(t ItemType, pr precedence, parselet infixParselet) {
		p.infixParselets[t] = parselet
		p.infixPrecedence[t] = pr
	}

	addPrefixParselet(ItemOpenParen, 0, p.parseParenExpr)
	addPrefixParselet(ItemIdentifier, 0, p.parseIdentifierExpr)
	addInfixParselet(ItemAnd, 1, p.parseBoolOpExpr)
	addInfixParselet(ItemOr, 2, p.parseBoolOpExpr)
	addInfixParselet(ItemDot, 3, p.parseAttrAccess)
	addInfixParselet(ItemOpenBracket, 4, p.parseSubscript)
	addInfixParselet(ItemCompare, 5, p.parseBinOp)
	addInfixParselet(ItemOpenParen, 8, p.parseCall)
	addInfixParselet(ItemString, 9, p.parseString)

	// rotate first token into current
	p.consume()
	p.consume()

	return p
}

func (p *parser) ParseTotalRule() (*TotalRule, error) {
	rule, err := p.parseRule(precedence(0))
	if err != nil {
		return nil, err
	}
	t := &TotalRule{rule}
	return t, nil
}

func (p *parser) consume() {
	p.curTok = p.nextTok
	p.nextTok = p.lexer.nextItem()
}

func (p *parser) parseRule(pr precedence) (Expression, error) {
	prefix, ok := p.prefixParselets[p.curTok.Type]
	if !ok {
		return nil, unexpectedToken{p.curTok}
	}
	left, err := prefix(pr)
	if err != nil {
		return nil, err
	}
	for pr < p.nextPrecedence() {
		p.consume()
		infix := p.infixParselets[p.curTok.Type]
		left, err = infix(left, pr)
		if err != nil {
			return nil, err
		}
	}

	return left, nil
}

func (p *parser) nextPrecedence() precedence {
	return p.infixPrecedence[p.nextTok.Type]
}

func (p *parser) expectOrErr(next ItemType) error {
	if !p.expect(next) {
		return expectedToken{wanted: next, have: p.nextTok}
	}
	return nil
}

func (p *parser) expect(next ItemType) bool {
	if p.nextTok.Is(next) {
		p.consume()
		return true
	}
	return false
}

func (p *parser) parseParenExpr(precedence) (Expression, error) {
	e, err := p.parseRule(0)
	if err != nil {
		return nil, err
	}

	if p.expect(ItemComma) {
		next, err := p.parseRule(0)
		if err != nil {
			return nil, err
		}
		args := []Expression{e, next}

		for p.expect(ItemComma) {
			arg, err := p.parseRule(0)
			if err != nil {
				return nil, err
			}

			args = append(args, arg)
		}

		e = &Tuple{args}
	}

	if err := p.expectOrErr(ItemCloseParen); err != nil {
		return nil, err
	}

	return e, nil
}

func (p *parser) parseIdentifierExpr(precedence) (Expression, error) {
	return &Identifier{p.curTok.Value}, nil
}

func (p *parser) parseBoolOpExpr(left Expression, pr precedence) (Expression, error) {
	thisTok := p.curTok
	right, err := p.parseRule(pr)
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

func (p *parser) parseBinOp(left Expression, pr precedence) (Expression, error) {
	thisTok := p.curTok
	right, err := p.parseRule(pr)
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

func (p *parser) parseCall(left Expression, pr precedence) (Expression, error) {
	if p.expect(ItemCloseParen) {
		return &Call{Name: left}, nil
	}

	var args []Expression
	for {
		arg, err := p.parseRule(0)
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
		if !p.nextTok.Is(ItemCloseParen) {
			break
		}
		p.consume()
	}

	if err := p.expectOrErr(ItemCloseParen); err != nil {
		return nil, err
	}

	return &Call{
		Name: left,
		Args: args,
	}, nil
}

func (p *parser) parseSubscript(left Expression, pr precedence) (Expression, error) {
	index, err := p.parseRule(0)
	if err != nil {
		return nil, err
	}

	if err = p.expectOrErr(ItemCloseBracket); err != nil {
		return nil, err
	}

	s := Subscript{Target: left, Index: index}
	return &s, nil
}

func (p *parser) parseString(precedence) (Expression, error) {
	return &String{Value: p.curTok.Value}, nil
}

func (p *parser) parseAttrAccess(Expression, precedence) (Expression, error) {
	panic("not impled")
}
