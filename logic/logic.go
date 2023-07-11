package logic

import (
	"io"

	"github.com/etc-sudonters/zootler/entity"
	"github.com/etc-sudonters/zootler/ioutils"
)

/*
Most of OOTR's logic is placed in well structured json files. Rules for moving
through an exit or collecting a token are encoded as strings that are parsed as
Python expressions using Python's std ast module, rewriting and redirecting
calls as necessary, before compiling the ast into function code objects that
can be provided the current state of the world and produces a boolean response
indicating if travel can continue to the desired node. As with any language,
not everything is bootstrapped and these logic expressions are allowed to call
into "built in" methods that cannot be expressed in the logic language (or
easily expressed -- "Beware of the Turing tar-pit in which everything is
possible but nothing of interest is easy.", Alan Perlis, Epigrams of
Programming -- we will consider the functionality offered by these "built ins"
to be interfaces and the details to be implementation defined.
*/

// we should be able to recover human understandable representation of this rule
type Rule interface {
	Fulfill(entity.Queryable) (bool, error)
	WriteTo(io.Writer) (int64, error)
}

type constRule bool

func (c constRule) Fulfill(entity.Queryable) (bool, error) {
	return bool(c), nil
}

func (c constRule) WriteTo(w io.Writer) (n int64, err error) {
	var m int
	if c {
		m, err = w.Write([]byte("true"))
	} else {
		m, err = w.Write([]byte("false"))
	}

	n = int64(m)
	return
}

var TrueRule Rule = constRule(true)
var FalseRule Rule = constRule(false)

type AndRule struct {
	LHS Rule
	RHS Rule
}

func (r AndRule) Fulfill(q entity.Queryable) (bool, error) {
	ok, err := r.LHS.Fulfill(q)
	if err != nil {
		return false, err
	}

	if !ok {
		return false, nil
	}

	return r.RHS.Fulfill(q)
}

func (r AndRule) WriteTo(w io.Writer) (int64, error) {
	e := &ioutils.ErrorCarryingWriter{W: w}
	c := &ioutils.CountingWriter{W: e}
	c.Write([]byte("("))
	r.LHS.WriteTo(c)
	c.Write([]byte(" and "))
	r.RHS.WriteTo(c)
	c.Write([]byte(")"))
	return c.N, e.Err
}

type OrRule struct {
	LHS Rule
	RHS Rule
}

func (r OrRule) WriteTo(w io.Writer) (int64, error) {
	e := &ioutils.ErrorCarryingWriter{W: w}
	c := &ioutils.CountingWriter{W: e}
	c.Write([]byte("("))
	r.LHS.WriteTo(c)
	c.Write([]byte(" or "))
	r.RHS.WriteTo(c)
	c.Write([]byte(")"))
	return c.N, e.Err
}

func (r OrRule) Fulfill(q entity.Queryable) (bool, error) {
	ok, err := r.LHS.Fulfill(q)
	if err != nil {
		return false, err
	}

	if ok {
		return true, nil
	}

	return r.RHS.Fulfill(q)
}
