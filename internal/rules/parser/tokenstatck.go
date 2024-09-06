package parser

import (
	"errors"

	"github.com/etc-sudonters/substrate/peruse"
	"github.com/etc-sudonters/substrate/skelly/stack"
)

type TokenStream = peruse.TokenStream
type Token = peruse.Token

var (
	ErrLockedStack = errors.New("token stream stack is locked")
	ErrEmptyStack  = errors.New("no token streams available")

	EOFToken = Token{Type: peruse.EOF}
)

type stackflag uint64

func (s stackflag) Has(f stackflag) bool {
	return s&f == f
}

func (s stackflag) Add(f stackflag) stackflag {
	return s | f
}

func (s stackflag) Remove(f stackflag) stackflag {
	return s & (^f)
}

type TokenStreamStack struct {
	// a single bool will get padded to hell and back so here's the full 64
	// bits the upper 32 bits are defined here, the lower 32 bits are free
	// space for callers to manipulate as needed
	flags stackflag
	toks  stack.S[TokenStream]
	buf   []Token
}

func (tss *TokenStreamStack) NextToken() Token {
	if tss.flags.Has(END_OF_STREAM) {
		if !tss.flags.Has(LOCKED) {
			return tss.nextTokenSlow()
		} else {
			return EOFToken
		}
	}

	top, err := tss.toks.Top()
	if err != nil { // empty stack
		tss.flags = tss.flags.Add(END_OF_STREAM)
		return EOFToken
	}

	tok := tss.maybeBuffer((*top).NextToken())
	if tok.Type == peruse.EOF {
		tss.flags = tss.flags.Add(END_OF_STREAM)
	}

	return tok
}

func (tss *TokenStreamStack) Lock() {
	tss.flags = tss.flags.Add(LOCKED)
}

func (tss *TokenStreamStack) Unlock() {
	tss.flags = tss.flags.Remove(LOCKED)
}

func (tss *TokenStreamStack) Flag(u uint32) {
	tss.flags = tss.flags.Add(stackflag(u))
}

func (tss *TokenStreamStack) Unflag(u uint32) {
	tss.flags = tss.flags.Remove(stackflag(u))
}

func (tss *TokenStreamStack) HasFlag(u uint32) bool {
	return tss.flags.Has(stackflag(u))
}

func (tss *TokenStreamStack) Push(t TokenStream) error {
	if tss.flags.Has(LOCKED) {
		return ErrLockedStack
	}

	tss.flags = tss.flags.Remove(END_OF_STREAM)
	tss.toks.Push(t)
	return nil
}

func (tss *TokenStreamStack) Pop() (TokenStream, error) {
	if tss.flags.Has(LOCKED) {
		return nil, ErrLockedStack
	}
	if len(tss.toks) == 0 {
		return nil, ErrEmptyStack
	}
	return tss.toks.Pop()
}

func (tss *TokenStreamStack) nextTokenSlow() Token {
	tok := EOFToken
	for tok.Type == peruse.EOF {
		tss.toks.Pop()
		ts, noMoreStreams := tss.toks.Top()
		if noMoreStreams != nil {
			return EOFToken
		}
		tok := (*ts).NextToken()
		if tok.Type != EOFToken.Type {
			tss.flags = tss.flags.Remove(END_OF_STREAM)
			return tok
		}
	}

	panic("unreachable")
}

func (tss *TokenStreamStack) maybeBuffer(t Token) Token {
	if tss.flags.Has(BUFFERING) && t.Type != peruse.EOF && t.Type != peruse.ERR {
		tss.buf = append(tss.buf, t)
	}
	return t
}

func (tss *TokenStreamStack) StartBuffer() {
	tss.flags = tss.flags.Add(BUFFERING)
}

func (tss *TokenStreamStack) StopBuffer() []Token {
	tss.flags = tss.flags.Remove(BUFFERING)
	buf := tss.buf
	tss.buf = nil
	return buf
}

const (
	// STACKFLAGS
	_      stackflag = 0
	LOCKED           = 1 << (iota + 33)
	END_OF_STREAM
	BUFFERING
)
