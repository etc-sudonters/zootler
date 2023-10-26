package peruse

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type LexFn func(*StringLexer, any) LexFn

type TokenType int64

const (
	EOF TokenType = -1
	ERR TokenType = -2
)

type Pos int
type Token struct {
	Type    TokenType
	Pos     Pos
	Literal string
}

func (t Token) String() string {
	return fmt.Sprintf("{typ: %d, pos: %d, literal: %.10q...}", t.Type, t.Pos, t.Literal)
}

func (t Token) Is(o TokenType) bool {
	return t.Type == o
}

func NewLexer(input string, initial LexFn, state any) *StringLexer {
	s := new(StringLexer)
	s.Init(input, initial, state)
	return s
}

type StringLexer struct {
	input string
	pos   Pos
	start Pos
	atEOF bool
	token Token
	state any
	fn    LexFn
}

func (l *StringLexer) Init(input string, initial LexFn, state any) {
	l.input = input
	l.pos = Pos(0)
	l.start = Pos(0)
	l.atEOF = false
	l.token = Token{}
	l.state = state
	l.fn = initial
}

func (l *StringLexer) NextToken() Token {
	l.token = Token{Type: EOF, Pos: l.pos}
	fn := l.fn
	for {
		fn = fn(l, l.state)
		if fn == nil {
			return l.token
		}
	}
}

func (l *StringLexer) Word() string {
	return l.input[l.start:l.pos]
}

func (l *StringLexer) Peek() rune {
	r := l.Next()
	l.Prev()
	return r
}

func (l *StringLexer) Next() rune {
	if int(l.pos) >= len(l.input) {
		l.atEOF = true
		return -1
	}

	r, width := utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += Pos(width)
	return r
}

func (l *StringLexer) Prev() {
	if !l.atEOF && l.pos > 0 {
		_, width := utf8.DecodeLastRuneInString(l.input[:l.pos])
		l.pos -= Pos(width)
	}
}

func (l *StringLexer) Emit(typ TokenType) LexFn {
	l.token = Token{
		Type:    typ,
		Pos:     l.start,
		Literal: l.input[l.start:l.pos],
	}

	l.start = l.pos
	return nil
}

func (l *StringLexer) Error(format string, args ...any) LexFn {
	l.token = Token{
		Type:    ERR,
		Pos:     l.pos,
		Literal: fmt.Sprintf(format, args...),
	}
	return nil
}

func (l *StringLexer) Discard() {
	l.start = l.pos
}

func (l *StringLexer) AcceptOneOf(chars string) bool {
	if strings.ContainsRune(chars, l.Next()) {
		return true
	}
	l.Prev()
	return false
}

func (l *StringLexer) AcceptManyOf(chars string) {
	l.AcceptWhile(func(r rune) bool { return strings.ContainsRune(chars, r) })
}

func (l *StringLexer) AcceptWhile(a func(rune) bool) {
	for {
		r := l.Next()
		if r == -1 {
			return
		}
		if !a(r) {
			l.Prev()
			break
		}
	}
}
