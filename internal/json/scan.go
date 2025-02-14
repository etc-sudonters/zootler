package json

import (
	"bufio"
	"errors"
	"fmt"
	"io"
)

func NewScanner(r io.Reader) *Scanner {
	this := &Scanner{
		inner: bufio.NewScanner(r),
	}
	this.inner.Split(this.split)
	return this
}

type Scanner struct {
	inner   *bufio.Scanner
	scanned scanned
}

type lexeme struct {
	scanned scanned
	body    []byte
}

func (this *Scanner) Next() (lexeme, error) {
	var l lexeme
	if !this.Scan() {
		l.scanned = scanned_eof
		err := this.Err()
		if err != nil {
			l.scanned = scanned_err
		}
		return l, err
	}
	l.scanned = this.scanned
	l.body = this.inner.Bytes()
	return l, nil
}

func (this *Scanner) Scan() bool {
	return this.inner.Scan()
}

func (this *Scanner) Lexeme() (scanned, []byte) {
	return this.scanned, this.inner.Bytes()
}

func (this *Scanner) Err() error {
	return this.inner.Err()
}

func (this *Scanner) split(buffer []byte, atEof bool) (int, []byte, error) {
	this.scanned = scanned_eof
	if len(buffer) == 0 && atEof {
		return 0, nil, bufio.ErrFinalToken
	}

	char, n := buffer[0], 1

	var err error
	if isWhitespace(char) {
		char, err = this.scanWhitespace(buffer, atEof, &n)
		if err != nil {
			return n, nil, err
		}
	}

	if char == eof {
		return n, nil, bufio.ErrFinalToken
	}

	switch char {
	case '{':
		this.scanned = scanned_obj_open
		return n, []byte{'{'}, nil
	case '}':
		this.scanned = scanned_obj_close
		return n, []byte{'}'}, nil
	case '[':
		this.scanned = scanned_arr_open
		return n, []byte{'['}, nil
	case ']':
		this.scanned = scanned_arr_close
		return n, []byte{']'}, nil
	case ':':
		this.scanned = scanned_colon
		return n, []byte{':'}, nil
	case ',':
		this.scanned = scanned_comma
		return n, []byte{','}, nil
	case '"':
		token, err := this.scanString(buffer, atEof, &n)
		return n, token, err
	case 't':
		token, err := this.scanLiteral(buffer, atEof, "true", scanned_true, &n)
		return n, token, err
	case 'f':
		token, err := this.scanLiteral(buffer, atEof, "false", scanned_false, &n)
		return n, token, err
	case 'n':
		token, err := this.scanLiteral(buffer, atEof, "null", scanned_null, &n)
		return n, token, err
	case '#', '/':
		token, err := this.scanComment(buffer, atEof, &n)
		return n, token, err
	}

	if canBeginNumber(char) {
		token, err := this.scanNumber(buffer, atEof, &n)
		return n, token, err
	}

	return n, nil, this.invalidToken(string(char))
}

func isWhitespace(char byte) bool {
	return char == ' ' || char == '\n' || char == '\t' || char == '\r'
}

func canBeginNumber(char byte) bool {
	return ('0' <= char && char <= '9') || char == '.' || char == '-'
}

func (this *Scanner) scanWhitespace(buffer []byte, _ bool, n *int) (byte, error) {
	startAt := *n
	found := eof
	for _, char := range buffer[startAt:] {
		(*n)++
		if !isWhitespace(char) {
			found = char
			break
		}
	}

	return found, nil
}

func (this *Scanner) scanString(buffer []byte, atEof bool, n *int) ([]byte, error) {
	size := len(buffer) - (*n)
	token := make([]byte, size)
	startAt := *n
	i := 0

	for _, char := range buffer[startAt:] {
		token[i] = char
		if char == '"' {
			if token[i-1] != '\\' {
				break
			}
			i--
			// overwrite \
			token[i] = char
		}
		i++
		(*n)++
	}

	if (*n) >= len(buffer) || token[i] != '"' {
		if !atEof {
			(*n) = startAt - 1
			return nil, nil
		}
		return nil, errors.New("unterminated string")
	}

	(*n)++
	this.scanned = scanned_string
	return token[:i], nil
}

func (this *Scanner) scanComment(buffer []byte, atEof bool, n *int) ([]byte, error) {
	startAt := *n
	if buffer[startAt-1] == '/' {
		if buffer[startAt] != '/' {
			return nil, this.invalidToken(string(buffer[(*n)-1]))
		}
		(*n)++
	}
	beginOfComment := *n
	last := buffer[beginOfComment]

	for _, char := range buffer[beginOfComment:] {
		last = char
		if char == '\n' || char == '\r' {
			break
		}
		(*n)++
	}

	if last != '\n' && last != '\r' {
		if !atEof {
			(*n) = startAt - 1
			return nil, nil
		}
	}

	this.scanned = scanned_comment
	token := buffer[beginOfComment:(*n)]
	if !atEof {
		(*n)++
	}
	return token, nil
}

func (this *Scanner) scanNumber(buffer []byte, atEof bool, n *int) ([]byte, error) {
	startAt := *n
	scannedSep := buffer[(*n)-1] == '.'
	scannedSign := buffer[(*n)-1] == '-'
	foundEnd := false

	for _, char := range buffer[startAt:] {
		if '0' <= char && char <= '9' {
			(*n)++
			continue
		}

		if char == '.' {
			if scannedSep {
				(*n) = startAt - 1
				return nil, this.invalidToken(string(buffer[*n]))
			}
			scannedSep = true
			(*n)++
			continue
		}

		if char == '-' {
			if scannedSign {
				(*n) = startAt - 1
				return nil, this.invalidToken(string(buffer[*n]))
			}
			scannedSign = true
			(*n)++
			continue
		}

		foundEnd = true
		break
	}

	if !foundEnd && !atEof {
		(*n) = startAt - 1
		return nil, nil
	}

	this.scanned = scanned_number
	token := buffer[startAt-1 : *n]
	return token, nil
}

func (this *Scanner) scanLiteral(buffer []byte, atEof bool, literal string, scanned scanned, n *int) ([]byte, error) {
	bytes := []byte(literal)
	size := len(bytes)
	(*n)--
	if size > len(buffer) {
		if atEof {
			return nil, this.invalidToken(string(buffer))
		}
		return nil, nil
	}

	segment := buffer[*n : (*n)+size]
	str := string(segment)

	if literal == str {
		(*n) += size
		this.scanned = scanned
		return bytes, nil
	}

	return nil, this.invalidToken(str)
}

func (this *Scanner) invalidToken(invalid string) error {
	this.scanned = scanned_err
	return fmt.Errorf("invalid token %q", invalid)
}

type scanned uint8

const (
	_ scanned = iota
	scanned_obj_open
	scanned_obj_close
	scanned_arr_open
	scanned_arr_close
	scanned_string
	scanned_comment
	scanned_number
	scanned_true
	scanned_false
	scanned_null
	scanned_comma
	scanned_colon
	scanned_eof scanned = 0xF0
	scanned_err scanned = 0xFE

	eof byte = 0
)

func (this scanned) String() string {
	switch this {
	case scanned_obj_open:
		return "OBJ_OPEN"
	case scanned_obj_close:
		return "OBJ_CLOSE"
	case scanned_arr_open:
		return "ARR_OPEN"
	case scanned_arr_close:
		return "ARR_CLOSE"
	case scanned_string:
		return "STRING"
	case scanned_comment:
		return "COMMENT"
	case scanned_number:
		return "NUMBER"
	case scanned_true:
		return "TRUE"
	case scanned_false:
		return "FALSE"
	case scanned_null:
		return "NULL"
	case scanned_comma:
		return "COMMA"
	case scanned_colon:
		return "COLON"
	case scanned_eof:
		return "EOF"
	case scanned_err:
		return "ERR"
	default:
		return "UNKNOWN"
	}
}
