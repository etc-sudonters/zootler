package parser

import (
	"github.com/etc-sudonters/substrate/peruse"
	"github.com/etc-sudonters/substrate/skelly/stack"
)

func NewRulesParser(macros MacroCoven) RulesParser {
	var r RulesParser
	r.annointedGrammar = annointGrammar(&r)
	r.tokenStreams = &TokenStreamStack{}
	r.tokenStreams.toks = make(stack.S[peruse.TokenStream], 0, 8)
	r.parser = peruse.NewParser(&r.annointedGrammar, r.tokenStreams)
	r.macros = macros
	return r
}

type RulesParser struct {
	parser           *peruse.Parser[Expression]
	tokenStreams     *TokenStreamStack
	annointedGrammar peruse.Grammar[Expression]
	macros           MacroCoven
}

func (r *RulesParser) ParseString(rule string) (Expression, error) {
	return r.ParseTokens(NewRulesLexer(rule))
}

func (r *RulesParser) ParseTokens(stream peruse.TokenStream) (Expression, error) {
	r.PushTokenStream(stream)
	defer r.PopTokenStream()
	r.tokenStreams.Lock()
	r.parser.Consume()
	r.parser.Consume()
	r.tokenStreams.Unlock()
	return r.parser.ParseAt(LOWEST)
}

func (r *RulesParser) PushTokenStream(s peruse.TokenStream) {
	r.tokenStreams.Push(s)
}

func (r *RulesParser) PopTokenStream() peruse.TokenStream {
	stream := r.teardownStream()
	return stream
}

func (r *RulesParser) teardownStream() peruse.TokenStream {
	stream, err := r.tokenStreams.Pop()
	if err != nil {
		panic("could not restore parser state")
	}

	return stream
}

func (r *RulesParser) mustOurParser(p *peruse.Parser[Expression]) {
	if r.parser != p {
		panic("unknown parser")
	}
}

func (r *RulesParser) parseCall(p *peruse.Parser[Expression], left Expression, bp peruse.Precedence) (Expression, error) {
	r.mustOurParser(p)
	expansion, didExpand, err := r.tryExpandCall(left, bp)
	if err != nil {
		return nil, err
	}
	if didExpand {
		return expansion, nil
	}

	return parseCall(r.parser, left, bp)
}

func (r *RulesParser) parseIdentifier(p *peruse.Parser[Expression]) (Expression, error) {
	r.mustOurParser(p)
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

func (r *RulesParser) tryExpandCall(left Expression, bp peruse.Precedence) (Expression, bool, error) {
	return nil, false, nil
}

func (r *RulesParser) tryExpandIdent(ident *Identifier) (Expression, bool, error) {
	return nil, false, nil
}

func annointGrammar(r *RulesParser) peruse.Grammar[Expression] {
	g := NewRulesGrammar()
	g.Parse(TokenIdentifier, r.parseIdentifier)
	g.Infix(PARENS, r.parseCall, TokenOpenParen)
	return g
}
