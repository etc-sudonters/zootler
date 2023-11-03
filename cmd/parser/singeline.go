package main

import (
	"fmt"
	"strings"

	rulesparser "sudonters/zootler/pkg/rules/parser"

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
	rulesparser.Visit(s, a.Target)
	s.b.WriteRune('.')
	rulesparser.Visit(s, a.Attr)
}
func (s *singleLine) VisitBinOp(b *rulesparser.BinOp) {
	s.writeOpenParen()
	s.b.WriteString(keywordColor.Paint(string(b.Op)))
	s.b.WriteRune(' ')
	rulesparser.Visit(s, b.Left)
	s.b.WriteRune(' ')
	rulesparser.Visit(s, b.Right)
	s.writeCloseParen()
}
func (s *singleLine) VisitBoolOp(b *rulesparser.BoolOp) {
	s.writeOpenParen()
	s.b.WriteString(keywordColor.Paint(string(b.Op)))
	s.b.WriteRune(' ')
	rulesparser.Visit(s, b.Left)
	s.b.WriteRune(' ')
	rulesparser.Visit(s, b.Right)
	s.writeCloseParen()
}
func (s *singleLine) VisitBoolean(b *rulesparser.Boolean) {
	fmt.Fprintf(s.b, boolColor.Paint("%t"), b.Value)
}
func (s *singleLine) VisitCall(c *rulesparser.Call) {
	s.writeOpenParen()
	rulesparser.Visit(s, c.Name)
	for _, arg := range c.Args {
		s.b.WriteRune(' ')
		rulesparser.Visit(s, arg)
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
	rulesparser.Visit(s, r.Target)
	s.b.WriteRune(' ')
	rulesparser.Visit(s, r.Index)
	s.writeCloseParen()
}
func (s *singleLine) VisitTuple(t *rulesparser.Tuple) {
	s.writeOpenParen()
	rulesparser.Visit(s, t.Elems[0])
	for _, arg := range t.Elems[1:] {
		s.b.WriteRune(' ')
		rulesparser.Visit(s, arg)
	}
	s.writeCloseParen()
}
func (s *singleLine) VisitUnary(u *rulesparser.UnaryOp) {
	s.writeOpenParen()
	s.b.WriteString(keywordColor.Paint(string(u.Op)))
	s.b.WriteRune(' ')
	rulesparser.Visit(s, u.Target)
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
