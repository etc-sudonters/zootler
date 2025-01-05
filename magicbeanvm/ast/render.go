package ast

import (
	"fmt"
	"strings"
)

type renderer struct {
	*strings.Builder
}

func Render(node Node) string {
	var sb strings.Builder
	render := renderer{&sb}
	v := Visitor{
		AnyOf:      (&render).AnyOf,
		Bool:       (&render).Bool,
		Compare:    (&render).Compare,
		Every:      (&render).Every,
		Identifier: (&render).Identifier,
		Invert:     (&render).Invert,
		Invoke:     (&render).Invoke,
		Number:     (&render).Number,
		String:     (&render).String,
	}
	v.Visit(node)
	return sb.String()
}

func (r *renderer) AnyOf(node AnyOf, visit Visiting) error {
	r.WriteString("(any-of ")
	r.joinWith(visit, node, " ")
	r.WriteRune(')')
	return nil
}

func (r *renderer) Bool(node Bool, visit Visiting) error {
	switch node {
	case true:
		r.WriteString("True")
	case false:
		r.WriteString("False")
	}
	return nil
}
func (r *renderer) Compare(node Compare, visit Visiting) error {
	r.WriteRune('(')
	switch node.Op {
	case CompareEq:
		r.WriteString("== ")
	case CompareNq:
		r.WriteString("!= ")
	case CompareLt:
		r.WriteString("< ")
	}
	r.joinWith(visit, []Node{node.LHS, node.RHS}, " ")
	r.WriteRune(')')
	return nil
}

func (r *renderer) Every(node Every, visit Visiting) error {
	r.WriteString("(every ")
	r.joinWith(visit, node, " ")
	r.WriteRune(')')
	return nil
}

func (r *renderer) Identifier(node Identifier, visit Visiting) error {
	fmt.Fprintf(r, "($%04X %q)", node.Symbol.Index, node.Symbol.Name)
	return nil
}

func (r *renderer) Invert(node Invert, visit Visiting) error {
	r.WriteString("(not ")
	visit(node.Inner)
	r.WriteRune(')')
	return nil
}
func (r *renderer) Invoke(node Invoke, visit Visiting) error {
	r.WriteString("(invoke ")
	inner := append([]Node{node.Target}, node.Args...)
	r.joinWith(visit, inner, " ")
	r.WriteRune(')')
	return nil
}
func (r *renderer) Number(node Number, visit Visiting) error {
	fmt.Fprintf(r, "%f", node)
	return nil
}
func (r *renderer) String(node String, visit Visiting) error {
	fmt.Fprintf(r, "s%q", node)
	return nil
}

func (r *renderer) joinWith(visit Visiting, nodes []Node, join string) error {
	length := len(nodes)
	for i := range nodes {
		visit(nodes[i])
		if i != length-1 {
			r.WriteString(join)
		}
	}

	return nil
}
