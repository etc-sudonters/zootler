package json

import (
	"fmt"
	"io"
	"strconv"
	"sudonters/libzootr/internal"
)

type Token struct {
	Kind Kind
	Body []byte
}

func unexpectedToken(t *Token) error {
	return fmt.Errorf("unexpected token %s: %q", scanned(t.Kind), string(t.Body))
}

type Kind uint8

const (
	OBJ_OPEN  = Kind(scanned_obj_open)
	OBJ_CLOSE = Kind(scanned_obj_close)
	ARR_OPEN  = Kind(scanned_arr_open)
	ARR_CLOSE = Kind(scanned_arr_close)
	STRING    = Kind(scanned_string)
	COMMENT   = Kind(scanned_comment)
	NUMBER    = Kind(scanned_number)
	TRUE      = Kind(scanned_true)
	FALSE     = Kind(scanned_false)
	NULL      = Kind(scanned_null)
	comma     = Kind(scanned_comma)
	colon     = Kind(scanned_colon)
	EOF       = Kind(scanned_eof)
	ERR       = Kind(scanned_err)
)

type Parser struct {
	scanner    *Scanner
	curr, peek Token
}

type composite struct {
	obj bool
	n   int
}

func NewParser(scanner *Scanner) *Parser {
	p := new(Parser)
	p.scanner = scanner
	p.Next()
	p.Next()
	return p
}

func (this *Parser) Discard() error {
	switch this.curr.Kind {
	case ARR_OPEN:
		if arr, err := this.ReadArray(); err != nil {
			return err
		} else {
			return arr.DiscardRemaining()
		}
	case OBJ_OPEN:
		if obj, err := this.ReadObject(); err != nil {
			return err
		} else {
			return obj.DiscardRemaining()
		}
	default:
		this.Next()
		return nil
	}
}

func (this *Parser) Next() bool {
	if this.peek.Kind == EOF {
		return false
	}

	for {
		lexeme, err := this.scanner.Next()
		internal.PanicOnError(err)

		if lexeme.scanned == scanned_comment {
			continue
		}

		this.curr = this.peek
		this.peek.Kind = Kind(lexeme.scanned)
		this.peek.Body = make([]byte, len(lexeme.body))
		copy(this.peek.Body, lexeme.body)
		return true
	}
}

func (this *Parser) Current() Token {
	return this.curr
}

func (this *Parser) Peek() Token {
	return this.peek
}

func (this *Parser) ReadObject() (*ObjectParser, error) {
	_, err := this.expect(OBJ_OPEN)
	if err != nil {
		return nil, err
	}

	this.Next()
	return &ObjectParser{this}, nil
}

func (this *Parser) ReadArray() (*ArrayParser, error) {
	_, err := this.expect(ARR_OPEN)
	if err != nil {
		return nil, err
	}
	this.Next()
	return &ArrayParser{this}, nil
}

func (this *Parser) ReadString() (string, error) {
	token, err := this.expect(STRING)
	if err != nil {
		return "", err
	}

	this.Next()
	return string(token.Body), nil
}

func (this *Parser) ReadInt() (int, error) {
	token, err := this.expect(NUMBER)
	if err != nil {
		return 0, err
	}

	number, err := strconv.Atoi(string(token.Body))
	if err != nil {
		return 0, fmt.Errorf("failed to parse number: %w", err)
	}

	this.Next()
	return number, nil
}

func (this *Parser) ReadFloat() (float64, error) {
	token, err := this.expect(NUMBER)
	if err != nil {
		return 0, err
	}

	number, err := strconv.ParseFloat(string(token.Body), 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse number: %w", err)
	}
	this.Next()
	return number, nil
}

func (this *Parser) ReadBool() (bool, error) {
	if this.curr.Kind == TRUE {
		this.Next()
		return true, nil
	} else if this.curr.Kind == FALSE {
		this.Next()
		return false, nil
	} else {
		return false, unexpectedToken(&this.curr)
	}
}

func (this *Parser) expect(expected Kind) (Token, error) {
	if this.curr.Kind == EOF {
		return this.curr, io.EOF
	}

	if this.curr.Kind != expected {
		return Token{Kind: ERR}, unexpectedToken(&this.curr)
	}

	return this.curr, nil
}

type ObjectParser struct {
	p *Parser
}

func (this *ObjectParser) ReadPropertyName() (string, error) {
	str, err := this.p.ReadString()
	if err != nil {
		return "", err
	}
	if this.p.curr.Kind != colon {
		return str, unexpectedToken(&this.p.curr)
	}
	this.p.Next()
	return str, nil
}

func (this *ObjectParser) ReadInt() (int, error) {
	number, err := this.p.ReadInt()
	if err != nil {
		return 0, err
	}
	maybeReadComma(this.p)
	return number, nil
}

func (this *ObjectParser) ReadFloat() (float64, error) {
	number, err := this.p.ReadFloat()
	if err != nil {
		return 0, err
	}
	maybeReadComma(this.p)
	return number, nil
}

func (this *ObjectParser) ReadString() (string, error) {
	str, err := this.p.ReadString()
	if err != nil {
		return "", err
	}
	maybeReadComma(this.p)
	return str, nil
}

func (this *ObjectParser) ReadBool() (bool, error) {
	boolean, err := this.p.ReadBool()
	if err != nil {
		return false, err
	}
	maybeReadComma(this.p)
	return boolean, nil
}

func (this *ObjectParser) DiscardValue() error {
	if err := this.p.Discard(); err != nil {
		return fmt.Errorf("unexpected end of object: %w", err)
	}
	maybeReadComma(this.p)
	return nil
}

func (this *ObjectParser) ReadEnd() error {
	_, err := this.p.expect(OBJ_CLOSE)
	if err != nil {
		return err
	}
	this.p.Next()
	maybeReadComma(this.p)
	return nil
}

func (this *ObjectParser) More() bool {
	return this.p.curr.Kind != OBJ_CLOSE
}

func (this *ObjectParser) ReadObject() (*ObjectParser, error) {
	return this.p.ReadObject()
}

func (this *ObjectParser) ReadArray() (*ArrayParser, error) {
	return this.p.ReadArray()
}

func (this *ObjectParser) DiscardRemaining() error {
	for this.More() {
		if _, err := this.ReadPropertyName(); err != nil {
			return err
		}
		if err := this.DiscardValue(); err != nil {
			return err
		}
	}
	this.ReadEnd()
	return nil
}

type ArrayParser struct {
	p *Parser
}

func (this *ArrayParser) ReadInt() (int, error) {
	number, err := this.p.ReadInt()
	if err != nil {
		return 0, err
	}
	maybeReadComma(this.p)
	return number, nil
}

func (this *ArrayParser) ReadFloat() (float64, error) {
	number, err := this.p.ReadFloat()
	if err != nil {
		return 0, err
	}
	maybeReadComma(this.p)
	return number, nil
}

func (this *ArrayParser) ReadString() (string, error) {
	str, err := this.p.ReadString()
	if err != nil {
		return "", err
	}
	maybeReadComma(this.p)
	return str, nil
}

func (this *ArrayParser) ReadBool() (bool, error) {
	boolean, err := this.p.ReadBool()
	if err != nil {
		return false, err
	}
	maybeReadComma(this.p)
	return boolean, nil
}

func (this *ArrayParser) ReadEnd() error {
	_, err := this.p.expect(ARR_CLOSE)
	if err != nil {
		return err
	}
	this.p.Next()
	maybeReadComma(this.p)
	return nil
}

func (this *ArrayParser) DiscardValue() error {
	if err := this.p.Discard(); err != nil {
		return fmt.Errorf("unexpected end of array: %w", err)
	}
	maybeReadComma(this.p)
	return nil
}

func (this *ArrayParser) More() bool {
	return this.p.curr.Kind != ARR_CLOSE
}

func (this *ArrayParser) ReadArray() (*ArrayParser, error) {
	return this.p.ReadArray()
}

func (this *ArrayParser) ReadObject() (*ObjectParser, error) {
	return this.p.ReadObject()
}

func (this *ArrayParser) DiscardRemaining() error {
	for this.More() {
		if err := this.DiscardValue(); err != nil {
			return err
		}
	}
	this.ReadEnd()
	return nil
}

func maybeReadComma(parser *Parser) {
	if parser.curr.Kind == comma {
		parser.Next()
	}
}
