package parser

import (
	"strings"
	"sudonters/zootler/internal/peruse"
	"unicode"
)

const (
	eof        = -1
	spaceChars = " \t\n\r"
	trueWord   = "True"
	falseWord  = "False"
	andWord    = "and"
	orWord     = "or"
	notWord    = "not"

	_ peruse.TokenType = iota
	TokenDot
	TokenOpenParen
	TokenCloseParen
	TokenOpenBracket
	TokenCloseBracket
	TokenIdentifier
	TokenString
	TokenNumber
	TokenAnd
	TokenOr
	TokenTrue
	TokenFalse
	TokenComma
	TokenEq
	TokenNotEq
	TokenLt
	TokenUnary
)

func TokenTypeString(i peruse.TokenType) string {
	switch i {
	case peruse.EOF:
		return "<EOF>"
	case peruse.ERR:
		return "<ERR>"
	case TokenDot:
		return "<DOT>"
	case TokenOpenParen:
		return "<OPENPAREN>"
	case TokenCloseParen:
		return "<CLOSEPAREN>"
	case TokenOpenBracket:
		return "<OPENBRACKET>"
	case TokenCloseBracket:
		return "<CLOSEBRACKET>"
	case TokenIdentifier:
		return "<IDENT>"
	case TokenString:
		return "<STR>"
	case TokenNumber:
		return "<NUMBER>"
	case TokenAnd:
		return "<AND>"
	case TokenOr:
		return "<OR>"
	case TokenTrue:
		return "<TRUE>"
	case TokenFalse:
		return "<FALSE>"
	case TokenComma:
		return "<COMMA>"
	case TokenEq:
		return "<EQ>"
	case TokenNotEq:
		return "<NEQ>"
	case TokenLt:
		return "<LT>"
	case TokenUnary:
		return "<UNARY>"
	default:
		return "<UNKNOWN>"
	}
}

func NewRulesLexer(raw string) *peruse.StringLexer {
	return peruse.NewLexer(raw, lexRule, &ruleLexState{})
}

type ruleLexState struct {
	parenDepth int
	brackDepth int
}

func lexRule(l *peruse.StringLexer, state any) peruse.LexFn {
	s := state.(*ruleLexState)

	switch r := l.Next(); {
	case r == eof:
		if s.parenDepth > 0 || s.brackDepth > 0 {
			return l.Error("unclosed '(' or '['")
		}
		return nil
	case isIdentBegin(r):
		l.Prev()
		return lexIdent
	case isDigit(r):
		l.Prev()
		return lexNumber
	case isWhitespace(r):
		return lexWhitespace
	case r == '(':
		return lexOpenParen
	case r == ')':
		return lexCloseParen
	case r == '[':
		return lexOpenBrack
	case r == ']':
		return lexCloseBrack
	case r == ',':
		return l.Emit(TokenComma)
	case r == '=':
		l.Prev()
		return lexEq
	case r == '!':
		l.Prev()
		return lexNotEq
	case r == '<' || r == '>':
		l.Prev()
		return lexInEq
	case r == '\'':
		return lexStr
	default:
		return l.Error("unrecongized character %#U", r)
	}
}

// scans an identifier
func lexIdent(l *peruse.StringLexer, _ any) peruse.LexFn {
	l.AcceptWhile(isIdentRune)
	switch word := l.Word(); {
	case word == andWord:
		return l.Emit(TokenAnd)
	case word == orWord:
		return l.Emit(TokenOr)
	case word == trueWord:
		return l.Emit(TokenTrue)
	case word == falseWord:
		return l.Emit(TokenFalse)
	case word == notWord:
		return l.Emit(TokenUnary)
	}

	if !atSeparator(l) {
		return unexpected(l.Peek(), l)
	}

	return l.Emit(TokenIdentifier)
}

func lexWhitespace(l *peruse.StringLexer, _ any) peruse.LexFn {
	l.AcceptWhile(isWhitespace)
	l.Discard()
	return lexRule
}

func lexNumber(l *peruse.StringLexer, _ any) peruse.LexFn {
	l.AcceptWhile(isDigit)
	if !atSeparator(l) {
		return unexpected(l.Peek(), l)
	}
	return l.Emit(TokenNumber)
}

// the '(' is already scanned
func lexOpenParen(l *peruse.StringLexer, state any) peruse.LexFn {
	s := state.(*ruleLexState)
	s.parenDepth++
	return l.Emit(TokenOpenParen)
}

// the ')' is already scanned
func lexCloseParen(l *peruse.StringLexer, state any) peruse.LexFn {
	s := state.(*ruleLexState)
	s.parenDepth--
	if s.parenDepth < 0 {
		return unexpected(')', l)
	}
	return l.Emit(TokenCloseParen)
}

// the '[' is already scanned
func lexOpenBrack(l *peruse.StringLexer, state any) peruse.LexFn {
	s := state.(*ruleLexState)
	s.brackDepth++
	return l.Emit(TokenOpenBracket)
}

// the ']' is already scanned
func lexCloseBrack(l *peruse.StringLexer, state any) peruse.LexFn {
	s := state.(*ruleLexState)
	s.brackDepth--
	if s.brackDepth < 0 {
		return unexpected(']', l)
	}
	return l.Emit(TokenCloseBracket)
}

func lexEq(l *peruse.StringLexer, _ any) peruse.LexFn {
	if l.AcceptOneOf("=") && l.AcceptOneOf("=") {
		return l.Emit(TokenEq)
	}

	return unexpected(l.Peek(), l)
}

func lexNotEq(l *peruse.StringLexer, _ any) peruse.LexFn {
	if l.AcceptOneOf("!") && l.AcceptOneOf("=") {
		return l.Emit(TokenNotEq)
	}
	return unexpected(l.Peek(), l)
}

func lexInEq(l *peruse.StringLexer, _ any) peruse.LexFn {
	l.AcceptOneOf("<") // know its one of these
	return l.Emit(TokenLt)
}

// opening ' is already scanned
func lexStr(l *peruse.StringLexer, _ any) peruse.LexFn {
	l.AcceptWhile(func(r rune) bool { return r != '\'' })
	l.Next()
	return l.Emit(TokenString)
}

func isDigit(r rune) bool {
	return unicode.IsDigit(r)
}

func isIdentBegin(r rune) bool {
	return r == '_' || unicode.IsLetter(r)
}

func isIdentRune(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}

func isWhitespace(r rune) bool {
	return strings.ContainsRune(spaceChars, r)
}

func atSeparator(l *peruse.StringLexer) bool {
	r := l.Peek()
	if isWhitespace(r) {
		return true
	}

	switch r {
	case eof, '.', '(', ')', ',', '[', ']':
		return true
	default:
		return false
	}
}

func unexpected(r rune, l *peruse.StringLexer) peruse.LexFn {
	return l.Error("unexpected %q", r)
}
