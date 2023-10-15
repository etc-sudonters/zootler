package logic

import (
	"io"

	"sudonters/zootler/pkg/entity"
)

/*
Most of OOTR's logic is placed in well structured json files. Rules for moving
through an exit or collecting a token are encoded as strings that are parsed as
Python expressions using Python's std ast module, rewriting and redirecting
calls as necessary, before compiling the ast into function code objects that
can be provided the current state of the world and produces a boolean response
indicating if travel can continue to the desired node. As with any language,
not everything is bootstrapped and these logic expressions are allowed to call
into "built in" methods that cannot be expressed in the logic language or
easily expressed -- "Beware of the Turing tar-pit in which everything is
possible but nothing of interest is easy.", Alan Perlis, Epigrams of
Programming -- we will consider the functionality offered by these "built ins"
to be interfaces and the details to be implementation defined.
*/

// determines if we are able to progress on our path
type Rule interface {
	Fulfill(entity.Queryable) (bool, error)
}

// a rule that is always true or false
type constRule bool

func (c constRule) Fulfill(entity.Queryable) (bool, error) {
	return bool(c), nil
}

// premature optimization is the root of all evil
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

// both embedded rules must be true for this rule to be true
// short circuiting, RHS is never called if LHS is false
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

// either embedded rule must be true for this rule to be true
// short circuits, RHS is never called if LHS is true
type OrRule struct {
	LHS Rule
	RHS Rule
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
