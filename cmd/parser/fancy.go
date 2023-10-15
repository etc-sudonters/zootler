package main

import (
	"fmt"
	"strings"

	"sudonters/zootler/internal/ioutil"
	"sudonters/zootler/internal/rules"
)

func newFancy() *FancyAstWriter {
	return &FancyAstWriter{
		b:      &strings.Builder{},
		indent: 0,
	}
}

var (
	nodeColor    ioutil.ForegroundColor = 244
	propColor    ioutil.ForegroundColor = 252
	strColor     ioutil.ForegroundColor = 112
	numColor     ioutil.ForegroundColor = 33
	boolColor    ioutil.ForegroundColor = 160
	identColor   ioutil.ForegroundColor = 99
	keywordColor ioutil.ForegroundColor = 208
	fnColor      ioutil.ForegroundColor = 228
)

var bracketColors []ioutil.ForegroundColor = []ioutil.ForegroundColor{
	1, 2, 3, 4, 5, 6, 7,
}

type FancyAstWriter struct {
	b      *strings.Builder
	indent int
}

func (w *FancyAstWriter) VisitAttrAccess(a *rules.AttrAccess) {
	w.writeObject(a)
	w.writeProperty("Target", a.Target)
	w.writeProperty("Attr", a.Attr)
	w.writeObjectEnd()
}
func (w *FancyAstWriter) VisitBinOp(b *rules.BinOp) {
	w.writeObject(b)
	w.writeKeywordProperty("Op", string(b.Op))
	w.writeProperty("Left", b.Left)
	w.writeProperty("Right", b.Right)
	w.writeObjectEnd()
}

func (w *FancyAstWriter) VisitBoolOp(b *rules.BoolOp) {
	w.writeObject(b)
	w.writeKeywordProperty("Op", string(b.Op))
	w.writeProperty("Left", b.Left)
	w.writeProperty("Right", b.Right)
	w.writeObjectEnd()
}

func (w *FancyAstWriter) VisitBoolean(b *rules.Boolean) {
	w.writeObject(b)
	w.writeBoolProperty("Value", b.Value)
	w.writeObjectEnd()
}
func (w *FancyAstWriter) VisitCall(c *rules.Call) {
	w.writeObject(c)
	w.writeProperty("Fn", c.Name)
	w.writePropertyName("Args")
	w.writeArrStart()
	for _, e := range c.Args {
		w.writeArrElem(e)
	}
	w.writeArrEnd()
	w.writeObjectEnd()
}
func (w *FancyAstWriter) VisitIdentifier(i *rules.Identifier) {
	w.writeObject(i)
	w.writeIdentifierProperty("Value", i.Value)
	w.writeObjectEnd()
}

func (w *FancyAstWriter) VisitNumber(n *rules.Number) {
	w.writeObject(n)
	w.writeNumProperty("Value", n.Value)
	w.writeObjectEnd()
}
func (w *FancyAstWriter) VisitString(s *rules.String) {
	w.writeObject(s)
	w.writeStrProperty("Value", s.Value)
	w.writeObjectEnd()
}
func (w *FancyAstWriter) VisitSubscript(s *rules.Subscript) {
	w.writeObject(s)
	w.writeProperty("Target", s.Target)
	w.writeProperty("Index", s.Index)
	w.writeObjectEnd()
}
func (w *FancyAstWriter) VisitTuple(t *rules.Tuple) {
	w.writeObject(t)
	w.writePropertyName("Elems")
	w.writeArrStart()
	for _, e := range t.Elems {
		w.writeArrElem(e)
	}
	w.writeArrEnd()
	w.writeObjectEnd()
}

func (w *FancyAstWriter) VisitUnary(u *rules.UnaryOp) {
	w.writeObject(u)
	w.writeKeywordProperty("Op", string(u.Op))
	w.writeProperty("Target", u.Target)
	w.writeObjectEnd()
}

func (a FancyAstWriter) writeIndent() {
	if a.indent > 0 {
		fmt.Fprint(a.b, strings.Repeat("  ", a.indent))
	}
}

func (a *FancyAstWriter) writeObjectType(o interface{}) {
	fmt.Fprintf(a.b, nodeColor.Paint("%T"), o)
}

func (a FancyAstWriter) bracketColor() ioutil.ForegroundColor {
	return bracketColors[a.indent%len(bracketColors)]
}

func (a *FancyAstWriter) writeColoredBracket(s string) {
	a.b.WriteString(a.bracketColor().Paint(s))
}

func (a *FancyAstWriter) writeObject(o interface{}) {
	a.writeObjectType(o)
	a.b.WriteRune(' ')
	a.writeColoredBracket("{")
	a.b.WriteRune('\n')
	a.indent += 1
}
func (a *FancyAstWriter) writeObjectEnd() {
	a.indent -= 1
	a.writeIndent()
	a.writeColoredBracket("}")
}

func (a *FancyAstWriter) writePropertyName(name string) {
	a.writeIndent()
	a.b.WriteString(propColor.Paint(name))
	a.b.WriteString(":  ")
}

func (a *FancyAstWriter) writePropertyEnd() {
	a.b.Write([]byte(",\n"))
}

func (a *FancyAstWriter) writeStr(s string) {
	fmt.Fprintf(a.b, strColor.Paint("%s"), s)
}

func (a *FancyAstWriter) writeStrProperty(name, value string) {
	a.writePropertyName(name)
	a.writeStr(value)
	a.writePropertyEnd()
}

func (a *FancyAstWriter) writeNumber(f float64) {
	fmt.Fprintf(a.b, numColor.Paint("%.0f"), f)
}

func (a *FancyAstWriter) writeNumProperty(name string, value float64) {
	a.writePropertyName(name)
	a.writeNumber(value)
	a.writePropertyEnd()
}

func (a *FancyAstWriter) writeBool(b bool) {
	fmt.Fprintf(a.b, boolColor.Paint("%t"), b)
}

func (a *FancyAstWriter) writeBoolProperty(name string, value bool) {
	a.writePropertyName(name)
	a.writeBool(value)
	a.writePropertyEnd()
}

func (a *FancyAstWriter) writeIdentifier(i string) {
	fmt.Fprint(a.b, identColor.Paint(i))
}

func (a *FancyAstWriter) writeIdentifierProperty(name, i string) {
	a.writePropertyName(name)
	a.writeIdentifier(i)
	a.writePropertyEnd()
}

func (a *FancyAstWriter) writeKeyword(kw string) {
	fmt.Fprint(a.b, keywordColor.Paint(kw))
}

func (a *FancyAstWriter) writeKeywordProperty(name, kw string) {
	a.writePropertyName(name)
	a.writeKeyword(kw)
	a.writePropertyEnd()
}

func (a *FancyAstWriter) writeProperty(name string, value rules.Expression) {
	a.writePropertyName(name)
	value.Visit(a)
	a.writePropertyEnd()
}

func (a *FancyAstWriter) writeArrStart() {
	a.writeColoredBracket("[")
	a.b.WriteRune('\n')
	a.indent += 1
}

func (a *FancyAstWriter) writeArrElem(v rules.Expression) {
	a.writeIndent()
	v.Visit(a)
	fmt.Fprintf(a.b, ",\n")
}

func (a *FancyAstWriter) writeArrEnd() {
	a.indent -= 1
	a.writeIndent()
	a.writeColoredBracket("]")
	fmt.Fprint(a.b, ",\n")
}

func (a *FancyAstWriter) finish() string {
	return a.b.String()
}
