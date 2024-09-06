package parser

import (
	"github.com/etc-sudonters/substrate/peruse"
	"github.com/etc-sudonters/substrate/skelly/stack"
	"github.com/etc-sudonters/substrate/slipup"
)

type MacroExpander interface {
	FullyExpand(*ExpansionContext) (Expression, error)
}

type MacroCoven struct {
	allMacros      map[string]Macro
	eligibility    map[string]bool
	expanders      map[string]MacroExpander
	defaultXpander MacroExpander
	scopes         stack.S[ExpansionContext]
}

func (mc MacroCoven) IsExpandableMacro(name string) bool {
	return mc.eligibility[name]
}

func (mc *MacroCoven) CreateContext(name string) (*ExpansionContext, func()) {
	var ctx ExpansionContext
	mc.eligibility[name] = false
	ctx.Expanding = mc.allMacros[name]
	stop := func() {
		mc.scopes.Pop()
		mc.eligibility[name] = true
	}

	return &mc.scopes[0], stop
}

func (mc *MacroCoven) Expander(m Macro) MacroExpander {
	xpander := mc.expanders[m.Name]
	if xpander == nil {
		if len(m.Body) != 0 || len(m.Params) > 0 || mc.defaultXpander == nil {
			panic(slipup.Createf("no expander found for macro '%s'", m.Name))
		}
		xpander = mc.defaultXpander
	}

	return xpander
}

type Macro struct {
	Name   string
	Body   []peruse.Token
	Params []string
}

type ExpansionContext struct {
	Args      []MacroArg
	Expanding Macro
	Parser    *RulesParser
}

type MacroArg struct {
	Tokens []peruse.Token
	AST    Expression
}
type CopyPasteExpander struct{}

func (c CopyPasteExpander) FullyExpand(ctx *ExpansionContext) (Expression, error) {
	var body []Token
	switch len(ctx.Args) {
	case 0:
		body = ctx.Expanding.Body
	default:
		body = buildReplacementBody(ctx)
	}
	stream := TokenSlice{tok: body}
	return ctx.Parser.Parse(&stream)
}

type TokenSlice struct {
	idx int
	tok []peruse.Token
}

func (t *TokenSlice) NextToken() peruse.Token {
	if t.idx >= len(t.tok) {
		return peruse.Token{Type: peruse.EOF}
	}

	tok := t.tok[t.idx]
	t.idx += 1
	return tok
}

func buildReplacementBody(ctx *ExpansionContext) []Token {
	var argTokenLen int
	for _, arg := range ctx.Args {
		argTokenLen += len(arg.Tokens)
	}

	body := make([]Token, 0, len(ctx.Expanding.Body)+2*argTokenLen)
	lastSpanStart := -1

	for idx, tok := range ctx.Expanding.Body {
		if tok.Type != TOK_MACRO_ARG_1 && tok.Type != TOK_MACRO_ARG_2 {
			if lastSpanStart == -1 {
				lastSpanStart = idx
			}
			continue
		}
		if lastSpanStart != -1 {
			body = append(body, ctx.Expanding.Body[lastSpanStart:idx]...)
			lastSpanStart = -1
		}

		switch tok.Type {
		case TOK_MACRO_ARG_1:
			body = append(body, ctx.Args[0].Tokens...)
			continue
		case TOK_MACRO_ARG_2:
			body = append(body, ctx.Args[1].Tokens...)
			continue
		}
	}

	return body
}
