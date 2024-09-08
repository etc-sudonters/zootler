package parser

import (
	"errors"

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

func (ps *ParserStack) BufferTokens(capacityHint int) (func() []peruse.Token, error) {
	if ps.activeTokens == nil {
		return nil, errors.New("no active token stream")
	}
	active := ps.activeTokens
	buffer := new(bufferingtokens)
	buffer.stream = active
	buffer.buffer = make([]peruse.Token, 0, capacityHint)
	ps.activeTokens = buffer

	return func() []peruse.Token {
		if ps.activeTokens != buffer {
			panic("attempting to end buffer with different active tokens")
		}
		ps.activeTokens = active
		return buffer.buffer
	}, nil
}

func (ps *ParserStack) expect(typ peruse.TokenType) bool {
	return ps.parser.Expect(typ)
}

func (ps *ParserStack) consume() {
	ps.parser.Consume()
}

func (ps *ParserStack) parseAt(bp peruse.Precedence) (Expression, error) {
	return ps.parser.ParseAt(bp)
}

type MacroEnableFlags uint8

const (
	MACROS_DISABLE      MacroEnableFlags = MACROS_DISABLE_OBJ | MACROS_DISABLE_FUNC
	MACROS_DISABLE_OBJ                   = 2
	MACROS_DISABLE_FUNC                  = 4
)

func NewParserStack(macros MacroCoven, macroFlags MacroEnableFlags) *ParserStack {
	ps := new(ParserStack)
	ps.macros = macros
	ps.grammar = annoint(ps, macroFlags)
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
	if ident, didNotCast := AssertAs[*Identifier](left); didNotCast == nil {
		expansion, didExpand, err := r.tryExpandCall(ident, bp)
		if err != nil {
			return nil, err
		}
		if didExpand {
			return expansion, nil
		}
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

func (r *ParserStack) tryExpandCall(ident *Identifier, bp peruse.Precedence) (Expression, bool, error) {
	if !r.macros.IsExpandableMacro(ident.Value) {
		return nil, false, nil
	}
	expansionCtx, closeCtx := r.macros.CreateContext(ident.Value, r)
	defer closeCtx()
	expander := r.macros.Expander(expansionCtx.Expanding)
	fragment, expansionErr := expander.FullyExpand(expansionCtx)
	return fragment, fragment != nil && expansionErr == nil, expansionErr
}

func (r *ParserStack) tryExpandIdent(ident *Identifier) (Expression, bool, error) {
	if !r.macros.IsExpandableMacro(ident.Value) {
		return nil, false, nil
	}
	expandCtx, closeCtx := r.macros.CreateContext(ident.Value, r)
	defer closeCtx()
	if len(expandCtx.Expanding.Params) != 0 {
		return nil, false, nil
	}
	expander := r.macros.Expander(expandCtx.Expanding)
	fragment, expansionErr := expander.FullyExpand(expandCtx)
	return fragment, fragment != nil && expansionErr == nil, expansionErr

}

func annoint(r *ParserStack, flags MacroEnableFlags) peruse.Grammar[Expression] {
	g := NewRulesGrammar()
	if flags&MACROS_DISABLE_OBJ != MACROS_DISABLE_OBJ {
		g.Parse(TokenIdentifier, r.parseIdentifier)
	}

	if flags&MACROS_DISABLE_FUNC != MACROS_DISABLE_FUNC {
		g.Infix(PARENS, r.parseCall, TokenOpenParen)
	}
	return g
}

type bufferingtokens struct {
	buffer []peruse.Token
	stream peruse.TokenStream
}

func (b *bufferingtokens) NextToken() peruse.Token {
	tok := b.stream.NextToken()
	b.buffer = append(b.buffer, tok)
	return tok
}
