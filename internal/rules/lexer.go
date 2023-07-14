package rules

import (
	"fmt"
	"strings"
	"unicode"
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

func (i item) fungible(o item) bool {
	return i.typ == o.typ && i.pos == o.pos && i.val == o.val
}

const (
	trueWord  = "True"
	falseWord = "False"
	noneWord  = "None"
	andWord   = "and"
	orWord    = "or"
)

const (
	itemEof itemType = iota
	itemErr
	itemDot
	itemQuote
	itemOpenParen
	itemCloseParen
	itemOpenBrack
	itemCloseBrack
	itemIdent
	itemStr
	itemNumber
	itemBoolOp
	itemBool
	itemNone
	itemComma
	itemCmp
)

func (i itemType) String() string {
	switch i {
	case itemEof:
		return "<EOF>"
	case itemErr:
		return "<ERR>"
	case itemDot:
		return "<DOT>"
	case itemQuote:
		return "<SINGLEQUOTE>"
	case itemOpenParen:
		return "<OPENPAREN>"
	case itemCloseParen:
		return "<CLOSEPAREN>"
	case itemIdent:
		return "<IDENT>"
	case itemStr:
		return "<STR>"
	case itemNumber:
		return "<NUMBER>"
	case itemBoolOp:
		return "<BOOLOP>"
	case itemBool:
		return "<BOOL>"
	case itemNone:
		return "<NONE>"
	case itemComma:
		return "<COMMA>"
	case itemCmp:
		return "<CMP>"
	default:
		return "<UNKNOWN>"
	}
}

const (
	eof        = -1
	spaceChars = " \t\n\r"
)

func (i item) String() string {
	repr := &strings.Builder{}
	repr.WriteString("{typ: ")
	repr.WriteString(i.typ.String())
	repr.WriteString(fmt.Sprintf(", pos: %d, val:", i.pos))

	if len(i.val) > 10 && i.typ != itemErr {
		fmt.Fprintf(repr, "%.10q...}", i.val)
		return repr.String()
	}

	fmt.Fprintf(repr, "%q}", i.val)
	return repr.String()
}

func lex(name, rule string) *lexer {
	return &lexer{
		name:  name,
		input: rule,
	}
}

type lexer struct {
	name       string
	input      string
	pos        pos
	start      pos
	atEOF      bool
	item       item
	parenDepth int
	brackDepth int
}

func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

func (l *lexer) next() rune {
	if int(l.pos) >= len(l.input) {
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

func (l *lexer) acceptOne(valid string) bool {
	if strings.ContainsRune(valid, l.next()) {
		return true
	}
	l.backup()
	return false
}

func (l *lexer) acceptRun(valid string) {
	l.acceptFn(func(r rune) bool {
		return strings.ContainsRune(valid, r)
	})
}

func (l *lexer) acceptFn(accept func(rune) bool) {
	for {
		r := l.next()
		if r == eof {
			return
		}
		if !accept(r) {
			l.backup()
			break
		}
	}
}

func (l *lexer) errorf(format string, args ...any) stateFn {
	l.item = item{
		typ: itemErr,
		pos: l.pos,
		val: fmt.Sprintf(format, args...),
	}
	return nil
}

func (l *lexer) nextItem() item {
	l.item = item{typ: itemEof, pos: l.pos}
	state := lexRule
	for {
		state = state(l)
		if state == nil {
			return l.item
		}
	}
}

func lexRule(l *lexer) stateFn {
	switch r := l.next(); {
	case r == eof:
		if l.parenDepth > 0 || l.brackDepth > 0 {
			return l.errorf("unclosed '(' or '['")
		}
		return nil
	case isIdentBegin(r):
		l.backup()
		return lexIdent
	case isDigit(r):
		l.backup()
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
		return l.emit(itemComma)
	case r == '=':
		l.backup()
		return lexEq
	case r == '!':
		l.backup()
		return lexNotEq
	case r == '<' || r == '>':
		l.backup()
		return lexInEq
	case r == '\'':
		return lexStr
	default:
		return l.errorf("unrecongized character %#U", r)
	}
}

// scans an identifier
func lexIdent(l *lexer) stateFn {
	l.acceptFn(isIdentRune)
	switch word := l.input[l.start:l.pos]; {
	case word == andWord || word == orWord:
		return l.emit(itemBoolOp)
	}

	if !atSeparator(l) {
		return unexpected(l.peek(), l)
	}

	return l.emit(itemIdent)
}

func lexWhitespace(l *lexer) stateFn {
	l.acceptFn(isWhitespace)
	l.ignore()
	return lexRule
}

func lexNumber(l *lexer) stateFn {
	l.acceptFn(isDigit)
	if !atSeparator(l) {
		return unexpected(l.peek(), l)
	}
	return l.emit(itemNumber)
}

// the '(' is already scanned
func lexOpenParen(l *lexer) stateFn {
	l.parenDepth++
	return l.emit(itemOpenParen)
}

// the ')' is already scanned
func lexCloseParen(l *lexer) stateFn {
	l.parenDepth--
	if l.parenDepth < 0 {
		return unexpected(')', l)
	}
	return l.emit(itemCloseParen)
}

// the '[' is already scanned
func lexOpenBrack(l *lexer) stateFn {
	l.brackDepth++
	return l.emit(itemOpenBrack)
}

// the ']' is already scanned
func lexCloseBrack(l *lexer) stateFn {
	l.brackDepth--
	if l.brackDepth < 0 {
		return unexpected(')', l)
	}
	return l.emit(itemCloseBrack)
}

func lexEq(l *lexer) stateFn {
	if l.acceptOne("=") && l.acceptOne("=") {
		return l.emit(itemCmp)
	}

	return unexpected(l.peek(), l)
}

func lexNotEq(l *lexer) stateFn {
	if l.acceptOne("!") && l.acceptOne("=") {
		return l.emit(itemCmp)
	}
	return unexpected(l.peek(), l)
}

func lexInEq(l *lexer) stateFn {
	l.acceptOne("<>") // know its one of these
	l.acceptOne("=")  // might be hanging out too
	return l.emit(itemCmp)
}

// opening ' is already scanned
func lexStr(l *lexer) stateFn {
	l.acceptFn(func(r rune) bool { return r != '\'' })
	l.next()
	return l.emit(itemStr)
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

func atSeparator(l *lexer) bool {
	r := l.peek()
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

func unexpected(r rune, l *lexer) stateFn {
	return l.errorf("unexpected %q", r)
}
