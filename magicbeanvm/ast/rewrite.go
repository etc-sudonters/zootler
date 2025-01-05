package ast

import (
	"errors"
	"fmt"
)

type CouldNotRewrite struct {
	Node  Node
	Cause error
}

func (err CouldNotRewrite) Error() string {
	if err.Cause != nil {
		return fmt.Sprintf("could not rewrite %T: %s", err.Node, err.Cause)
	}

	return fmt.Sprintf("could not rewrite: %#v", err.Node)
}

func DontRewrite[N Node]() RewriteFunc[N] {
	return func(n N, _ Rewriting) (Node, error) {
		return n, nil
	}
}

func RewriteWithEvery(original Node, rw []Rewriter) (Node, error) {
	var err error
	node := original
	for i := range rw {
		node, err = rw[i].Rewrite(node)
		if node == nil {
			break
		}
		if err != nil {
			break
		}
	}

	return node, err
}

type Rewriting func(Node) (Node, error)
type RewriteFunc[T Node] func(T, Rewriting) (Node, error)

func (r Rewriting) All(ast []Node) ([]Node, error) {
	var err error
	rewritten := make([]Node, len(ast))
	for i := range ast {
		rewritten[i], err = r(ast[i])
		if err != nil {
			err = errors.Join(err)
		}
	}
	return rewritten, err
}

type Rewriter struct {
	AnyOf      RewriteFunc[AnyOf]
	Bool       RewriteFunc[Bool]
	Compare    RewriteFunc[Compare]
	Every      RewriteFunc[Every]
	Identifier RewriteFunc[Identifier]
	Invert     RewriteFunc[Invert]
	Invoke     RewriteFunc[Invoke]
	Number     RewriteFunc[Number]
	String     RewriteFunc[String]
}

func (v Rewriter) Rewrite(ast Node) (Node, error) {
	if ast == nil {
		panic("rewriting nil node")
	}
	switch ast := ast.(type) {
	case AnyOf:
		if v.AnyOf == nil {
			return v.anyof(ast, v.Rewrite)
		}
		return v.AnyOf(ast, v.Rewrite)
	case Bool:
		if v.Bool == nil {
			return v.boolean(ast)
		}
		return v.Bool(ast, v.Rewrite)
	case Compare:
		if v.Compare == nil {
			return v.compare(ast, v.Rewrite)
		}
		return v.Compare(ast, v.Rewrite)
	case Every:
		if v.Every == nil {
			return v.every(ast, v.Rewrite)
		}
		return v.Every(ast, v.Rewrite)
	case Identifier:
		if v.Identifier == nil {
			return v.identifier(ast)
		}
		return v.Identifier(ast, v.Rewrite)
	case Invert:
		if v.Invert == nil {
			return v.invert(ast)
		}
		return v.Invert(ast, v.Rewrite)
	case Invoke:
		if v.Invoke == nil {
			return v.invoke(ast)
		}
		return v.Invoke(ast, v.Rewrite)
	case Number:
		if v.Number == nil {
			return v.number(ast)
		}
		return v.Number(ast, v.Rewrite)
	case String:
		if v.String == nil {
			return v.str(ast)
		}
		return v.String(ast, v.Rewrite)
	default:
		return nil, CouldNotRewrite{ast, UnknownNode}
	}
}

func (v Rewriter) anyof(anyof AnyOf, rewrite Rewriting) (Node, error) {
	items, err := rewrite.All(anyof)
	return AnyOf(items), err
}

func (v Rewriter) boolean(b Bool) (Node, error) {
	return b, nil
}

func (v Rewriter) compare(compare Compare, rewrite Rewriting) (Node, error) {
	operands, err := rewrite.All([]Node{compare.LHS, compare.RHS})
	return Compare{
		LHS: operands[0],
		RHS: operands[1],
		Op:  compare.Op,
	}, err
}

func (v Rewriter) every(every Every, rewrite Rewriting) (Node, error) {
	items, err := rewrite.All(every)
	return Every(items), err
}

func (v Rewriter) identifier(i Identifier) (Node, error) {
	return i, nil
}

func (v Rewriter) invert(invert Invert) (Node, error) {
	return v.Rewrite(invert.Inner)
}

func (v Rewriter) invoke(invoke Invoke) (Node, error) {
	var rewritten Invoke
	var err error
	rewritten.Target, err = v.Rewrite(invoke.Target)
	rewritten.Args = make([]Node, len(invoke.Args))

	for i := range invoke.Args {
		var argErr error
		rewritten.Args[i], argErr = v.Rewrite(invoke.Args[i])
		if argErr != nil {
			err = errors.Join(argErr)
		}
	}
	return rewritten, err
}

func (v Rewriter) number(n Number) (Node, error) {
	return n, nil
}

func (v Rewriter) str(s String) (Node, error) {
	return s, nil
}
