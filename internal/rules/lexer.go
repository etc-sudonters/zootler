package rules

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

type stateFn func(*lexer) stateFn

type ItemType int
type Pos int
type Item struct {
	Type  ItemType
	Pos   Pos
	Value string
}

const (
	trueWord  = "True"
	falseWord = "False"
	andWord   = "and"
	orWord    = "or"
)

const (
	ItemEof ItemType = iota
	ItemErr
	ItemDot
	ItemQuote
	ItemOpenParen
	ItemCloseParen
	ItemOpenBracket
	ItemCloseBracket
	ItemIdentifier
	ItemString
	ItemNumber
	ItemAnd
	ItemOr
	ItemBool
	ItemComma
	ItemCompare
)

func (i ItemType) String() string {
	switch i {
	case ItemEof:
		return "<EOF>"
	case ItemErr:
		return "<ERR>"
	case ItemDot:
		return "<DOT>"
	case ItemQuote:
		return "<SINGLEQUOTE>"
	case ItemOpenParen:
		return "<OPENPAREN>"
	case ItemCloseParen:
		return "<CLOSEPAREN>"
	case ItemIdentifier:
		return "<IDENT>"
	case ItemString:
		return "<STR>"
	case ItemNumber:
		return "<NUMBER>"
	case ItemAnd:
		return "<ANDOP>"
	case ItemOr:
		return "<OROP>"
	case ItemBool:
		return "<BOOL>"
	case ItemComma:
		return "<COMMA>"
	case ItemCompare:
		return "<CMP>"
	default:
		return "<UNKNOWN>"
	}
}

const (
	eof        = -1
	spaceChars = " \t\n\r"
)

func (i Item) String() string {
	repr := &strings.Builder{}
	repr.WriteString("{typ: ")
	repr.WriteString(i.Type.String())
	repr.WriteString(fmt.Sprintf(", pos: %d, val:", i.Pos))

	if len(i.Value) > 10 && i.Type != ItemErr {
		fmt.Fprintf(repr, "%.10q...}", i.Value)
		return repr.String()
	}

	fmt.Fprintf(repr, "%q}", i.Value)
	return repr.String()
}

func (i Item) Is(t ItemType) bool {
	return i.Type == t
}

func NewLexer(name, rule string) *lexer {
	return &lexer{
		name:  name,
		input: rule,
	}
}

type lexer struct {
	name       string
	input      string
	pos        Pos
	start      Pos
	atEOF      bool
	item       Item
	parenDepth int
	brackDepth int
}

func (l *lexer) nextItem() Item {
	l.item = Item{Type: ItemEof, Pos: l.pos}
	state := lexRule
	for {
		state = state(l)
		if state == nil {
			return l.item
		}
	}
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
	l.pos += Pos(width)
	return r
}

func (l *lexer) backup() {
	if !l.atEOF && l.pos > 0 {
		_, width := utf8.DecodeLastRuneInString(l.input[:l.pos])
		l.pos -= Pos(width)
	}
}

func (l *lexer) thisItem(typ ItemType) Item {
	i := Item{
		Type:  typ,
		Pos:   l.start,
		Value: l.input[l.start:l.pos],
	}
	l.start = l.pos
	return i
}

func (l *lexer) emit(typ ItemType) stateFn {
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
	l.item = Item{
		Type:  ItemErr,
		Pos:   l.pos,
		Value: fmt.Sprintf(format, args...),
	}
	return nil
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
		return l.emit(ItemComma)
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
	case word == andWord:
		return l.emit(ItemAnd)
	case word == orWord:
		return l.emit(ItemOr)
	case word == trueWord || word == falseWord:
		return l.emit(ItemBool)
	}

	if !atSeparator(l) {
		return unexpected(l.peek(), l)
	}

	return l.emit(ItemIdentifier)
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
	return l.emit(ItemNumber)
}

// the '(' is already scanned
func lexOpenParen(l *lexer) stateFn {
	l.parenDepth++
	return l.emit(ItemOpenParen)
}

// the ')' is already scanned
func lexCloseParen(l *lexer) stateFn {
	l.parenDepth--
	if l.parenDepth < 0 {
		return unexpected(')', l)
	}
	return l.emit(ItemCloseParen)
}

// the '[' is already scanned
func lexOpenBrack(l *lexer) stateFn {
	l.brackDepth++
	return l.emit(ItemOpenBracket)
}

// the ']' is already scanned
func lexCloseBrack(l *lexer) stateFn {
	l.brackDepth--
	if l.brackDepth < 0 {
		return unexpected(']', l)
	}
	return l.emit(ItemCloseBracket)
}

func lexEq(l *lexer) stateFn {
	if l.acceptOne("=") && l.acceptOne("=") {
		return l.emit(ItemCompare)
	}

	return unexpected(l.peek(), l)
}

func lexNotEq(l *lexer) stateFn {
	if l.acceptOne("!") && l.acceptOne("=") {
		return l.emit(ItemCompare)
	}
	return unexpected(l.peek(), l)
}

func lexInEq(l *lexer) stateFn {
	l.acceptOne("<>") // know its one of these
	l.acceptOne("=")  // might be hanging out too
	return l.emit(ItemCompare)
}

// opening ' is already scanned
func lexStr(l *lexer) stateFn {
	l.acceptFn(func(r rune) bool { return r != '\'' })
	l.next()
	return l.emit(ItemString)
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
