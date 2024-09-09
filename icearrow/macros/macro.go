package macros

import (
	"errors"
	"slices"
	"sudonters/zootler/icearrow/parser"

	"github.com/etc-sudonters/substrate/peruse"
)

func ExpandWith(macros *Expansions, tokens peruse.TokenStream) *tokenslicestreamer {
	expanded := macros.Expand(tokens)
	return &tokenslicestreamer{expanded, 0}
}

const (
	TOK_MACRO_ARG_0 peruse.TokenType = 0xFFFA0000
	TOK_MACRO_ARG_1                  = 0xFFFA0001
)

type Macro struct {
	Name   string
	params []string
	body   []peruse.Token
}

func (m Macro) tombstoneBody() {
	params := map[string]peruse.TokenType{
		m.params[0]: TOK_MACRO_ARG_0,
	}
	if len(m.params) > 1 {
		params[m.params[1]] = TOK_MACRO_ARG_1
	}

	for idx, tok := range m.body {
		if tok.Is(parser.TokenIdentifier) {
			if newType, replace := params[tok.Literal]; replace {
				tok.Type = newType
				m.body[idx] = tok
			}
		}
	}
}

func (m Macro) untombstoneBody(arguments []ExpandedArg) []peruse.Token {
	var body []peruse.Token
	for _, tok := range m.body {
		switch tok.Type {
		case TOK_MACRO_ARG_0:
			body = slices.Concat(body, arguments[0][:])
			break
		case TOK_MACRO_ARG_1:
			body = slices.Concat(body, arguments[1][:])
			break
		default:
			body = append(body, tok)
			break
		}
	}

	return body
}

func untombstoner(ctx *ExpansionContext) ExpandedTokens {
	return ctx.Expanding.untombstoneBody(ctx.Args)
}

var DefaultExpander FuncExpander = untombstoner

type FuncExpander func(*ExpansionContext) ExpandedTokens

func (f FuncExpander) Expand(ctx *ExpansionContext) ExpandedTokens {
	return f(ctx)
}

type Expander interface {
	Expand(*ExpansionContext) ExpandedTokens
}

type Expansions struct {
	macros      map[string]Macro
	eligibility map[string]bool
	expanders   map[string]Expander
}

type ExpandedArg []peruse.Token
type ExpandedTokens []peruse.Token

type ExpansionContext struct {
	Expanding Macro
	Args      []ExpandedArg
}

func NewMacroExpansions() Expansions {
	var e Expansions
	e.macros = make(map[string]Macro)
	e.eligibility = make(map[string]bool)
	e.expanders = make(map[string]Expander)
	return e
}

func (e *Expansions) Declare(name string, params []string, body []peruse.Token, xpndr Expander) error {
	if _, exists := e.macros[name]; exists {
		return errors.New("macro already declared")
	}

	m := Macro{name, params, body}
	if len(m.params) != 0 {
		m.tombstoneBody()
	}
	e.macros[name] = m
	e.eligibility[name] = true
	e.expanders[name] = xpndr

	return nil
}

func (e *Expansions) Expand(stream peruse.TokenStream) ExpandedTokens {
	var body ExpandedTokens
loop:
	for {
		tok := stream.NextToken()
		switch tok.Type {
		case peruse.EOF:
			break loop
		case peruse.ERR:
			body = append(body, tok)
			break loop
		case parser.TokenIdentifier:
			body = e.tryExpand(tok, stream, body)
			continue loop
		default:
			body = append(body, tok)
			continue loop
		}
	}
	return body
}

func (e *Expansions) tryExpand(tok peruse.Token, stream peruse.TokenStream, body ExpandedTokens) ExpandedTokens {
	name := tok.Literal
	if !e.eligibility[name] {
		body = append(body, tok)
		return body
	}
	e.eligibility[name] = false
	defer func() { e.eligibility[name] = true }()
	macro := e.macros[name]
	xpndr := e.expanders[name]
	var ctx ExpansionContext
	ctx.Expanding = macro
	if len(macro.params) != 0 {
		ctx.collectArguments(stream)
	}

	expandedTokens := xpndr.Expand(&ctx)
	expandedStream := tokenslicestreamer{expandedTokens, 0}
	reexpanded := e.Expand(&expandedStream)
	return slices.Concat(body, reexpanded)
}

func (e *ExpansionContext) collectArguments(stream peruse.TokenStream) {
	e.Args = make([]ExpandedArg, len(e.Expanding.params))
	p := miniparse{stream, peruse.Token{Type: parser.TokenOpenParen}}
	for _, arg := range e.Args {
		p.consume()
		for !p.expect(parser.TokenComma) && !p.expect(parser.TokenCloseParen) {
			arg = append(arg, p.cur)
		}
	}

	if !p.expect(parser.TokenCloseParen) {
		panic("expected to find TokenCloseParen at end of macro invocation")
	}
	p.consume()
}

type miniparse struct {
	tokens peruse.TokenStream
	cur    peruse.Token
}

func (m miniparse) expect(tt peruse.TokenType) bool {
	return m.cur.Is(tt)
}

func (m miniparse) consume() {
	m.cur = m.tokens.NextToken()
}

type tokenslicestreamer struct {
	tokens []peruse.Token
	idx    int
}

func (t *tokenslicestreamer) NextToken() peruse.Token {
	if t.idx >= len(t.tokens) {
		return peruse.Token{Type: peruse.EOF}
	}

	tok := t.tokens[t.idx]
	t.idx++
	return tok
}
