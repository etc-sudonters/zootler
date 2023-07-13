package python

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type stateFn func(*lexer) stateFn
type itemType int
type pos int
type item struct {
	typ itemType
	pos pos
	val string
}

const (
	itemEof itemType = iota
	itemErr
	itemName
	itemStr
	itemFuncOpen
	itemFuncClose
	itemOpenParen
	itemCloseParen
	itemNumber
	itemDot
	itemQuote
	itemOr
	itemAnd
)

const (
	eof        = -1
	spaceChars = " \t\n\r"
)

func (i item) String() string {
	switch i.typ {
	case itemEof:
		return "<EOF>"
	case itemErr:
		return i.val
	case itemName:
		return fmt.Sprintf("<%s>", i.val)
	}

	if len(i.val) > 10 {
		return fmt.Sprintf("%.10q...", i.val)
	}

	return fmt.Sprintf("%q", i.val)
}

type lexer struct {
	name       string
	input      string
	pos        pos
	start      pos
	atEOF      bool
	parenDepth int
	item       item
}

func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

func (l *lexer) next() rune {
	if int(l.pos) > len(l.input) {
		l.atEOF = true
		return eof
	}

	r, width := utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += pos(width)
	return r
}

func (l *lexer) backup() {
	if !l.atEOF && l.pos > 0 {
		_, width := utf8.DecodeLastRuneInString(l.input[:l.pos])
		l.pos -= pos(width)
	}
}

func (l *lexer) thisItem(typ itemType) item {
	i := item{
		typ: typ,
		pos: l.start,
		val: l.input[l.start:l.pos],
	}
	l.start = l.pos
	return i
}

func (l *lexer) emit(typ itemType) stateFn {
	l.item = l.thisItem(typ)
	return nil
}

func (l *lexer) ignore() {
	l.start = l.pos
}

func (l *lexer) accept(valid string) bool {
	if strings.ContainsRune(valid, l.next()) {
		return true
	}
	l.backup()
	return false
}

func (l *lexer) acceptRun(valid string) {
	for strings.ContainsRune(valid, l.next()) {
	}
	l.backup()
}

func (l *lexer) errorf(format string, args ...any) stateFn {
	l.item = item{
		typ: itemErr,
		pos: l.pos,
		val: fmt.Sprintf(format, args...),
	}
	l.pos = 0
	l.start = 0
	l.input = l.input[:0]
	return nil
}

func (l *lexer) nextItem() item {
	l.item = item{}
	state := lexRule
	for {
		state = state(l)
		if state == nil {
			return l.item
		}
	}
}

func lexRule(l *lexer) stateFn {
	return nil
}
