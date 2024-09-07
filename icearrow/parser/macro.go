package parser

import (
	"errors"
	"slices"
	"strings"

	"github.com/etc-sudonters/substrate/peruse"
	"github.com/etc-sudonters/substrate/skelly/stack"
	"github.com/etc-sudonters/substrate/slipup"
)

var (
	ErrMacroExists = errors.New("macro already declared")
)

type MacroExpander interface {
	FullyExpand(*ExpansionContext) (Expression, error)
}

func InitiateCoven(mc *MacroCoven) MacroBuilder {
	return MacroBuilder{mc}
}

type MacroBuilder struct {
	mc *MacroCoven
}

func (mb MacroBuilder) AddScriptedMacro(decl, body string) error {
	name, params := quickanddirtyDeclParse(decl)
	if _, exists := mb.mc.allMacros[name]; exists {
		return ErrMacroExists
	}

	m := Macro{
		Name:   name,
		Params: params,
		Body:   mb.macroBodyTokens(body, params),
	}

	mb.mc.allMacros[name] = m
	mb.mc.eligibility[name] = true
	return nil
}

func (mb MacroBuilder) AddBuiltInMacro(name string, params []string, expander MacroExpander) error {
	if _, exists := mb.mc.allMacros[name]; exists {
		return ErrMacroExists
	}

	m := Macro{
		Name:   name,
		Params: params,
	}

	mb.mc.allMacros[name] = m
	mb.mc.eligibility[name] = true
	mb.mc.expanders[name] = expander
	return nil
}

func (mb MacroBuilder) macroBodyTokens(body string, params []string) []Token {
	tokens := slices.Collect(peruse.AllTokens(NewRulesLexer(body)))
	if len(params) == 0 {
		return tokens
	}

	repl := map[string]peruse.TokenType{
		params[0]: TOK_MACRO_ARG_0,
	}
	if len(params) > 1 {
		repl[params[1]] = TOK_MACRO_ARG_1
	}
	if len(params) > 2 {
		panic("macro passed too many params")
	}

	for i, tok := range tokens {
		if tok.Type == TokenIdentifier {
			replacement, shouldReplace := repl[tok.Literal]
			if shouldReplace {
				// we can keep the same literal
				tok.Type = replacement
				tokens[i] = tok
			}
		}
	}

	return tokens
}

type Macro struct {
	Name   string
	Body   []peruse.Token // this _will_ have EOF at end
	Params []string
}

func NewCoven(defaultXpander MacroExpander) MacroCoven {
	var mc MacroCoven
	mc.allMacros = make(map[string]Macro, 256)
	mc.eligibility = make(map[string]bool, 256)
	mc.expanders = make(map[string]MacroExpander, 8) // very unique expanders
	mc.defaultXpander = defaultXpander
	mc.scopes = stack.Make[ExpansionContext](0, 8)
	return mc
}

func DefaultCoven() MacroCoven {
	return NewCoven(CopyPasteExpander{})
}

type MacroCoven struct {
	allMacros      map[string]Macro
	eligibility    map[string]bool
	expanders      map[string]MacroExpander
	defaultXpander MacroExpander
	scopes         *stack.S[ExpansionContext]
}

func (mc MacroCoven) IsExpandableMacro(name string) bool {
	return mc.eligibility[name]
}

func (mc *MacroCoven) CreateContext(name string) (*ExpansionContext, func()) {
	var ctx ExpansionContext
	mc.eligibility[name] = false
	ctx.Expanding = mc.allMacros[name]
	mc.scopes.Push(ctx)
	stop := func() {
		mc.scopes.Pop()
		mc.eligibility[name] = true
	}

	// borrow through top to enforce stack's ownership
	top, _ := mc.scopes.Top()
	return top, stop
}

func (mc *MacroCoven) Expander(m Macro) MacroExpander {
	xpander := mc.expanders[m.Name]
	if xpander == nil {
		if len(m.Body) != 0 || mc.defaultXpander == nil {
			panic(slipup.Createf("no expander found for macro '%s'", m.Name))
		}
		xpander = mc.defaultXpander
	}

	return xpander
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
	return ctx.Parser.ParseTokens(&stream)
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
	body := []Token{}
	for _, tok := range ctx.Expanding.Body {
		if tok.Type != TOK_MACRO_ARG_0 && tok.Type != TOK_MACRO_ARG_1 {
			body = append(body, tok)
			continue
		}

		switch tok.Type {
		case TOK_MACRO_ARG_0:
			body = append(body, ctx.Args[0].Tokens...)
			continue
		case TOK_MACRO_ARG_1:
			body = append(body, ctx.Args[1].Tokens...)
			continue
		}
	}

	return body
}

// this is similar/the same as upstream OOTR
func quickanddirtyDeclParse(decl string) (string, []string) {
	if !strings.Contains(decl, "(") {
		return decl, nil
	}

	decl = strings.TrimSuffix(decl, ")")
	splitDecl := strings.Split(decl, "(")
	args := strings.Split(splitDecl[1], ",")
	return splitDecl[0], args
}
