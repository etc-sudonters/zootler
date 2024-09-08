package parser

import (
	"github.com/etc-sudonters/substrate/peruse"
	"github.com/etc-sudonters/substrate/skelly/stack"
)

type ParserStack struct {
	grammar peruse.Grammar[Expression]
	macros  MacroCoven
	parser  *peruse.Parser[Expression]
	states  *stack.S[state]

	activeTokens peruse.TokenStream
}

func NewParserStack(macros MacroCoven) *ParserStack {
	ps := new(ParserStack)
	ps.macros = macros
	ps.grammar = annoint(ps)
	ps.states = stack.Make[state](0, 8)
	ps.activeTokens = new(TokenSlice)
	ps.parser = peruse.NewParser(&ps.grammar, stacktokens{ps})
	return ps
}

func (ps *ParserStack) ParseString(rule string) (Expression, error) {
	return ps.ParseTokens(NewRulesLexer(rule))
}

func (ps *ParserStack) ParseTokens(ts peruse.TokenStream) (Expression, error) {
	restore := ps.stashParserState()
	defer restore()
	ps.activeTokens = ts
	ps.primeParser()
	return ps.parser.ParseAt(LOWEST)
}

func (ps *ParserStack) stashParserState() func() {
	var state state
	state.cur = ps.parser.Cur
	state.next = ps.parser.Next
	state.tokens = ps.activeTokens
	ps.states.Push(state)

	return func() {
		state, err := ps.states.Pop()
		if err != nil {
			panic(err)
		}
		ps.activeTokens = state.tokens
		ps.parser.Cur = state.cur
		ps.parser.Next = state.next
	}
}

func (ps *ParserStack) primeParser() {
	ps.parser.Consume()
	ps.parser.Consume()
}

type stacktokens struct {
	ps *ParserStack
}

func (st stacktokens) NextToken() peruse.Token {
	return st.ps.activeTokens.NextToken()
}

type state struct {
	cur, next peruse.Token
	tokens    peruse.TokenStream
}

func (r *ParserStack) parseCall(p *peruse.Parser[Expression], left Expression, bp peruse.Precedence) (Expression, error) {
	expansion, didExpand, err := r.tryExpandCall(left, bp)
	if err != nil {
		return nil, err
	}
	if didExpand {
		return expansion, nil
	}

	return parseCall(r.parser, left, bp)
}

func (r *ParserStack) parseIdentifier(p *peruse.Parser[Expression]) (Expression, error) {
	ident, err := parseIdentifierExpr(p)
	if err != nil {
		return nil, err
	}
	expansion, didExpand, err := r.tryExpandIdent(ident.(*Identifier))
	if err != nil {
		return nil, err
	}
	if didExpand {
		return expansion, nil
	}

	return ident, nil
}

func (r *ParserStack) tryExpandCall(left Expression, bp peruse.Precedence) (Expression, bool, error) {
	return nil, false, nil
}

func (r *ParserStack) tryExpandIdent(ident *Identifier) (Expression, bool, error) {
	return nil, false, nil
}

func annoint(r *ParserStack) peruse.Grammar[Expression] {
	g := NewRulesGrammar()
	g.Parse(TokenIdentifier, r.parseIdentifier)
	g.Infix(PARENS, r.parseCall, TokenOpenParen)
	return g
}
