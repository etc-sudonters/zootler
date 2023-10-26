package peruse

import (
	"fmt"
	"sudonters/zootler/internal/mirrors"
)

const (
	INVALID_PRECEDENCE Precedence = iota
	LOWEST
)

func UnexpectedAt(where string, unexpected error) error {
	return fmt.Errorf("parsing %s: %w", where, unexpected)
}

type UnexpectedToken struct {
	Have Token
}

func (u UnexpectedToken) Error() string {
	return fmt.Sprintf("unexpected token %q", u.Have)
}

type InvalidToken struct {
	Have   Token
	Wanted TokenType
}

func (e InvalidToken) Error() string {
	return fmt.Sprintf("expected %q but found %q", e.Wanted, e.Have)
}

type Parselet[T any] func(*Parser[T]) (T, error)
type InflixParselet[T any] func(*Parser[T], T, Precedence) (T, error)
type Precedence uint

type Grammar[T any] struct {
	parselets  map[TokenType]Parselet[T]
	infixes    map[TokenType]InflixParselet[T]
	precedence map[TokenType]Precedence
}

func NewGrammar[T any]() Grammar[T] {
	var g Grammar[T]
	g.parselets = make(map[TokenType]Parselet[T])
	g.infixes = make(map[TokenType]InflixParselet[T])
	g.precedence = make(map[TokenType]Precedence)
	return g
}

func (g Grammar[T]) Parse(t TokenType, p Parselet[T]) {
	g.parselets[t] = p
}

func (g Grammar[T]) Infix(p Precedence, i InflixParselet[T], ts ...TokenType) {
	for _, t := range ts {
		g.infixes[t] = i
		g.precedence[t] = p
	}
}

func (g Grammar[T]) Precedence(t TokenType) Precedence {
	return g.precedence[t]
}

func NewParser[T any](g Grammar[T], l *StringLexer) *Parser[T] {
	p := new(Parser[T])
	p.g = g
	p.l = l
	// cycle first two tokens into place
	p.Consume()
	p.Consume()
	return p
}

type Parser[T any] struct {
	g         Grammar[T]
	l         *StringLexer
	Cur, Next Token
	empty     T
}

func (p *Parser[T]) HasMore() bool {
	return !(p.Next.Is(EOF) || p.Next.Is(ERR))
}

func (p *Parser[T]) Parse() (T, error) {
	return p.ParseAt(LOWEST)
}

func (p *Parser[T]) ParseAt(prd Precedence) (T, error) {
	parser, ok := p.g.parselets[p.Cur.Type]

	if !ok {
		return mirrors.Empty[T](), UnexpectedToken{p.Cur}
	}

	left, err := parser(p)

	if err != nil {
		return mirrors.Empty[T](), err
	}

	for thisPrd := p.NextPrecedence(); prd < thisPrd; thisPrd = p.NextPrecedence() {
		p.Consume()

		parselet, exists := p.g.infixes[p.Cur.Type]
		if !exists {
			break
		}

		left, err = parselet(p, left, thisPrd)
		if err != nil {
			return left, err
		}
	}

	return left, nil
}

func (p *Parser[T]) NextPrecedence() Precedence {
	return p.g.precedence[p.Next.Type]
}

func (p *Parser[T]) Consume() {
	p.Cur = p.Next
	p.Next = p.l.NextToken()
}

func (p *Parser[T]) Expect(n TokenType) bool {
	if p.Next.Is(n) {
		p.Consume()
		return true
	}
	return false
}

func (p *Parser[T]) ExpectOrError(n TokenType) error {
	if !p.Expect(n) {
		return InvalidToken{Wanted: n, Have: p.Next}
	}

	return nil
}
