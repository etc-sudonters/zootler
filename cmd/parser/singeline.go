package main

import (
	"fmt"
	"strings"

	"sudonters/zootler/internal/ioutil"
	"sudonters/zootler/internal/rules"
)

func newSingleLine() *singleLine {
	return &singleLine{&strings.Builder{}, 0}
}

type singleLine struct {
	b     *strings.Builder
	depth int
}

func (s *singleLine) VisitAttrAccess(a *rules.AttrAccess) {
	a.Target.Visit(s)
	s.b.WriteRune('.')
	a.Attr.Visit(s)
}
func (s *singleLine) VisitBinOp(b *rules.BinOp) {
	s.writeOpenParen()
	s.b.WriteString(keywordColor.Paint(string(b.Op)))
	s.b.WriteRune(' ')
	b.Left.Visit(s)
	s.b.WriteRune(' ')
	b.Right.Visit(s)
	s.writeCloseParen()
}
func (s *singleLine) VisitBoolOp(b *rules.BoolOp) {
	s.writeOpenParen()
	s.b.WriteString(keywordColor.Paint(string(b.Op)))
	s.b.WriteRune(' ')
	b.Left.Visit(s)
	s.b.WriteRune(' ')
	b.Right.Visit(s)
	s.writeCloseParen()
}
func (s *singleLine) VisitBoolean(b *rules.Boolean) {
	fmt.Fprintf(s.b, boolColor.Paint("%t"), b.Value)
}
func (s *singleLine) VisitCall(c *rules.Call) {
	s.writeOpenParen()
	c.Name.Visit(s)
	for _, arg := range c.Args {
		s.b.WriteRune(' ')
		arg.Visit(s)
	}
	s.writeCloseParen()
}
func (s *singleLine) VisitIdentifier(i *rules.Identifier) {
	s.b.WriteString(identColor.Paint(i.Value))
}
func (s *singleLine) VisitNumber(n *rules.Number) {
	fmt.Fprintf(s.b, numColor.Paint("%.0f"), n.Value)
}
func (s *singleLine) VisitString(r *rules.String) {
	s.b.WriteString(strColor.Paint(r.Value))
}
func (s *singleLine) VisitSubscript(r *rules.Subscript) {
	s.writeOpenParen()
	s.b.WriteString("[] ")
	r.Target.Visit(s)
	s.b.WriteRune(' ')
	r.Index.Visit(s)
	s.writeCloseParen()
}
func (s *singleLine) VisitTuple(t *rules.Tuple) {
	s.writeOpenParen()
	t.Elems[0].Visit(s)
	for _, arg := range t.Elems[1:] {
		s.b.WriteRune(' ')
		arg.Visit(s)
	}
	s.writeCloseParen()
}
func (s *singleLine) VisitUnary(u *rules.UnaryOp) {
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

func (s singleLine) bracketColor() ioutil.ForegroundColor {
	return bracketColors[s.depth%len(bracketColors)]
}
