package main

import (
	"fmt"
	"strings"

	"sudonters/zootler/pkg/rulesparser"

	"github.com/etc-sudonters/substrate/dontio"
)

func newSingleLine() *singleLine {
	return &singleLine{&strings.Builder{}, 0}
}

type singleLine struct {
	b     *strings.Builder
	depth int
}

func (s *singleLine) VisitAttrAccess(a *rulesparser.AttrAccess) {
	a.Target.Visit(s)
	s.b.WriteRune('.')
	a.Attr.Visit(s)
}
func (s *singleLine) VisitBinOp(b *rulesparser.BinOp) {
	s.writeOpenParen()
	s.b.WriteString(keywordColor.Paint(string(b.Op)))
	s.b.WriteRune(' ')
	b.Left.Visit(s)
	s.b.WriteRune(' ')
	b.Right.Visit(s)
	s.writeCloseParen()
}
func (s *singleLine) VisitBoolOp(b *rulesparser.BoolOp) {
	s.writeOpenParen()
	s.b.WriteString(keywordColor.Paint(string(b.Op)))
	s.b.WriteRune(' ')
	b.Left.Visit(s)
	s.b.WriteRune(' ')
	b.Right.Visit(s)
	s.writeCloseParen()
}
func (s *singleLine) VisitBoolean(b *rulesparser.Boolean) {
	fmt.Fprintf(s.b, boolColor.Paint("%t"), b.Value)
}
func (s *singleLine) VisitCall(c *rulesparser.Call) {
	s.writeOpenParen()
	c.Name.Visit(s)
	for _, arg := range c.Args {
		s.b.WriteRune(' ')
		arg.Visit(s)
	}
	s.writeCloseParen()
}
func (s *singleLine) VisitIdentifier(i *rulesparser.Identifier) {
	s.b.WriteString(identColor.Paint(i.Value))
}
func (s *singleLine) VisitNumber(n *rulesparser.Number) {
	fmt.Fprintf(s.b, numColor.Paint("%.0f"), n.Value)
}
func (s *singleLine) VisitString(r *rulesparser.String) {
	s.b.WriteString(strColor.Paint(r.Value))
}
func (s *singleLine) VisitSubscript(r *rulesparser.Subscript) {
	s.writeOpenParen()
	s.b.WriteString("[] ")
	r.Target.Visit(s)
	s.b.WriteRune(' ')
	r.Index.Visit(s)
	s.writeCloseParen()
}
func (s *singleLine) VisitTuple(t *rulesparser.Tuple) {
	s.writeOpenParen()
	t.Elems[0].Visit(s)
	for _, arg := range t.Elems[1:] {
		s.b.WriteRune(' ')
		arg.Visit(s)
	}
	s.writeCloseParen()
}
func (s *singleLine) VisitUnary(u *rulesparser.UnaryOp) {
	s.writeOpenParen()
	s.b.WriteString(keywordColor.Paint(string(u.Op)))
	s.b.WriteRune(' ')
	u.Target.Visit(s)
	s.writeCloseParen()
}

func (s *singleLine) writeOpenParen() {
	s.b.WriteString(s.bracketColor().Paint("("))
	s.depth += 1
}

func (s *singleLine) writeCloseParen() {
	s.depth -= 1
	s.b.WriteString(s.bracketColor().Paint(")"))
}

func (s singleLine) bracketColor() dontio.ForegroundColor {
	return bracketColors[s.depth%len(bracketColors)]
}
