package ast

import (
	"fmt"
	"hash"
	"hash/fnv"
)

func Hash(node Node) uint64 {
	h := fnv.New64()
	if err := Hash64(node, h); err != nil {
		panic(err)
	}
	return h.Sum64()
}

func Hash64(node Node, h64 hash.Hash64) error {
	hasher := hash64{h64}

	visitor := Visitor{
		AnyOf:      hasher.AnyOf,
		Boolean:    hasher.Boolean,
		Compare:    hasher.Compare,
		Every:      hasher.Every,
		Identifier: hasher.Identifier,
		Invert:     hasher.Invert,
		Invoke:     hasher.Invoke,
		Number:     hasher.Number,
		String:     hasher.String,
	}

	return visitor.Visit(node)
}

type hash64 struct {
	hash.Hash64
}

func (this hash64) AnyOf(node AnyOf, visit Visiting) error {
	fmt.Fprint(this, node.Kind())
	fmt.Fprintf(this, "%d", len(node))
	return visit.All(node)
}

func (this hash64) Boolean(node Boolean, visit Visiting) error {
	fmt.Fprint(this, node.Kind())
	fmt.Fprintf(this, "%t", node)
	return nil
}

func (this hash64) Compare(node Compare, visit Visiting) error {
	fmt.Fprint(this, node.Kind())
	fmt.Fprint(this, node.Op)
	return visit.All([]Node{node.LHS, node.RHS})
}

func (this hash64) Every(node Every, visit Visiting) error {
	fmt.Fprint(this, node.Kind())
	fmt.Fprintf(this, "%d", len(node))
	return visit.All(node)
}

func (this hash64) Identifier(node Identifier, visit Visiting) error {
	fmt.Fprint(this, node.Kind())
	fmt.Fprintf(this, "%x", node)
	return nil
}

func (this hash64) Invert(node Invert, visit Visiting) error {
	fmt.Fprint(this, node.Kind())
	return visit(node.Inner)
}

func (this hash64) Invoke(node Invoke, visit Visiting) error {
	fmt.Fprint(this, node.Kind())
	if err := visit(node.Target); err != nil {
		return err
	}

	fmt.Fprintf(this, "%d", len(node.Args))
	return visit.All(node.Args)
}

func (this hash64) Number(node Number, visit Visiting) error {
	fmt.Fprint(this, node.Kind())
	fmt.Fprintf(this, "%f", node)
	return nil
}

func (this hash64) String(node String, visit Visiting) error {
	fmt.Fprint(this, node.Kind())
	fmt.Fprint(this, node)
	return nil
}
