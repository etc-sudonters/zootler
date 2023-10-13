package rules

import (
	"fmt"
	"strconv"
)

func expectedFrom(where string, unexpected error) error {
	return fmt.Errorf("parsing %s: %w", where, unexpected)
}

type UnexpectedToken struct {
	Have Item
}

func (u UnexpectedToken) Error() string {
	return fmt.Sprintf("unexpected token %q", u.Have)
}

type InvalidToken struct {
	Have   Item
	Wanted ItemType
}

func (e InvalidToken) Error() string {
	return fmt.Sprintf("expected %q but found %q", e.Wanted, e.Have)
}

type prefixParselet func() (Expression, error)
type infixParselet func(Expression, bindingPower) (Expression, error)
type bindingPower int

const (
	_ bindingPower = iota
	LOWEST
	AND
	ACCESS
	EQ
	LT
	PREFIX
	FUNC
)

type parser struct {
	lexer           *lexer
	curTok          Item
	nextTok         Item
	prefixParselets map[ItemType]prefixParselet
	infixParselets  map[ItemType]infixParselet
	infixBinding    map[ItemType]bindingPower
}

func NewParser(l *lexer) *parser {
	p := &parser{
		lexer:           l,
		prefixParselets: make(map[ItemType]prefixParselet),
		infixParselets:  make(map[ItemType]infixParselet),
		infixBinding:    make(map[ItemType]bindingPower),
	}

	addPrefixParselet := func(parselet prefixParselet, ts ...ItemType) {
		for _, t := range ts {
			p.prefixParselets[t] = parselet
		}
	}

	addInfixParselet := func(bp bindingPower, parselet infixParselet, ts ...ItemType) {
		for _, t := range ts {
			p.infixParselets[t] = parselet
			p.infixBinding[t] = bp
		}
	}

	addInfixParselet(FUNC, p.parseSubscript, ItemOpenBracket)
	addInfixParselet(AND, p.parseBoolOpExpr, ItemAnd, ItemOr)
	addInfixParselet(LOWEST, p.parseAttrAccess, ItemDot)
	addInfixParselet(EQ, p.parseBinOp, ItemEq, ItemNotEq)
	addInfixParselet(LT, p.parseBinOp, ItemLt)
	addInfixParselet(FUNC, p.parseCall, ItemOpenParen)

	addPrefixParselet(p.parseBool, ItemTrue, ItemFalse)
	addPrefixParselet(p.parseIdentifierExpr, ItemIdentifier)
	addPrefixParselet(p.parseNumber, ItemNumber)
	addPrefixParselet(p.parseParenExpr, ItemOpenParen)
	addPrefixParselet(p.parseString, ItemString)
	addPrefixParselet(p.parsePrefixUnaryOp, ItemUnary)

	// rotate first token into current
	p.consume()
	p.consume()

	return p
}

func (p *parser) ParseTotalRule() (*TotalRule, error) {
	rule, err := p.parseRule(LOWEST)
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

func (p *parser) parseRule(bp bindingPower) (Expression, error) {
	prefix, ok := p.prefixParselets[p.curTok.Type]
	if !ok {
		return nil, UnexpectedToken{p.curTok}
	}
	left, err := prefix()
	if err != nil {
		return nil, err
	}
	for bp < p.nextPrecedence() {
		thisBp := p.nextPrecedence()
		p.consume()
		infix := p.infixParselets[p.curTok.Type]
		left, err = infix(left, thisBp)
		if err != nil {
			return nil, err
		}
	}

	return left, nil
}

func (p *parser) nextPrecedence() bindingPower {
	return p.infixBinding[p.nextTok.Type]
}

func (p *parser) expectOrErr(next ItemType) error {
	if !p.expect(next) {
		return InvalidToken{Wanted: next, Have: p.nextTok}
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

func (p *parser) parseParenExpr() (Expression, error) {
	p.consume()
	e, err := p.parseRule(LOWEST)
	if err != nil {
		return nil, err
	}

	if p.nextTok.Is(ItemComma) {
		e, err = p.parseTuple(e)
		if err != nil {
			return nil, err
		}
	}

	if err := p.expectOrErr(ItemCloseParen); err != nil {
		return nil, expectedFrom("PARENEXPR", err)
	}

	return e, nil
}

func (p *parser) parseTuple(left Expression) (Expression, error) {
	elems := []Expression{left}

	for p.expect(ItemComma) {
		p.consume()
		elem, err := p.parseRule(LOWEST)
		if err != nil {
			return nil, err
		}

		elems = append(elems, elem)
	}

	return &Tuple{elems}, nil
}

func (p *parser) parseIdentifierExpr() (Expression, error) {
	return &Identifier{p.curTok.Value}, nil
}

func (p *parser) parseBoolOpExpr(left Expression, bp bindingPower) (Expression, error) {
	thisTok := p.curTok
	p.consume()
	right, err := p.parseRule(bp)
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

func (p *parser) parseBinOp(left Expression, bp bindingPower) (Expression, error) {
	thisTok := p.curTok
	p.consume()
	right, err := p.parseRule(bp)
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

func (p *parser) parseCall(left Expression, bp bindingPower) (Expression, error) {
	if p.expect(ItemCloseParen) { // fn()
		return &Call{Name: left}, nil
	}

	var args []Expression
	for {
		p.consume()
		arg, err := p.parseRule(LOWEST)
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
		if !p.expect(ItemComma) {
			break
		}
	}

	if err := p.expectOrErr(ItemCloseParen); err != nil {
		return nil, expectedFrom("FNCALL", err)
	}

	c := Call{Name: left, Args: args}
	return &c, nil
}

func (p *parser) parseSubscript(left Expression, bp bindingPower) (Expression, error) {
	p.consume()
	index, err := p.parseRule(LOWEST)
	if err != nil {
		return nil, err
	}

	if err = p.expectOrErr(ItemCloseBracket); err != nil {
		return nil, expectedFrom("SUBSCRIPT", err)
	}

	s := Subscript{Target: left, Index: index}
	return &s, nil
}

func (p *parser) parseString() (Expression, error) {
	s := &String{Value: p.curTok.Value}
	return s, nil
}

func (p *parser) parseAttrAccess(target Expression, bp bindingPower) (Expression, error) {
	p.consume()
	attr, err := p.parseRule(bp)
	if err != nil {
		return nil, err
	}
	return &AttrAccess{
		Target: target,
		Attr:   attr,
	}, nil
}

func (p *parser) parseNumber() (Expression, error) {
	n, err := strconv.ParseFloat(p.curTok.Value, 64)
	if err != nil {
		return nil, fmt.Errorf("cannot parse %q as number", p.curTok)
	}
	return &Number{n}, nil
}

func (p *parser) parseBool() (Expression, error) {
	switch p.curTok.Value {
	case trueWord:
		return &Boolean{Value: true}, nil
	case falseWord:
		return &Boolean{Value: false}, nil
	default:
		return nil, fmt.Errorf("unexpected boolean value %s", p.curTok)
	}
}

func (p *parser) parsePrefixUnaryOp() (Expression, error) {
	thisTok := p.curTok
	p.consume()
	target, err := p.parseRule(PREFIX)
	if err != nil {
		return nil, err
	}
	switch thisTok.Value {
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
