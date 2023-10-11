package main

import (
	"fmt"
	"strings"

	"github.com/etc-sudonters/zootler/internal/console"
	"github.com/etc-sudonters/zootler/internal/rules"
)

func newFancy() *FancyAstWriter {
	return &FancyAstWriter{
		b:      &strings.Builder{},
		indent: 0,
	}
}

var (
	nodeColor  console.ForegroundColor = 57
	propColor  console.ForegroundColor = 69
	strColor   console.ForegroundColor = 106
	numColor   console.ForegroundColor = 159
	boolColor  console.ForegroundColor = 208
	identColor console.ForegroundColor = 212
	fnColor    console.ForegroundColor = 160
)

var bracketColors []console.ForegroundColor = []console.ForegroundColor{
	196,
	201,
	100,
	92,
	123,
	220,
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
	w.writeProperty("Left", b.Left)
	w.writeStrProperty("Op", string(b.Op))
	w.writeProperty("Right", b.Right)
	w.writeObjectEnd()
}

func (w *FancyAstWriter) VisitBoolOp(b *rules.BoolOp) {
	w.writeObject(b)
	w.writeProperty("Left", b.Left)
	w.writeStrProperty("Op", string(b.Op))
	w.writeProperty("Right", b.Right)
	w.writeObjectEnd()
}

func (w *FancyAstWriter) VisitBoolean(b *rules.Boolean) {
	w.writeObject(b)
	w.writeBool(b.Value)
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
	w.writeStrProperty("Value", i.Value)
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

func (a FancyAstWriter) writeIndent() {
	if a.indent > 0 {
		fmt.Fprint(a.b, strings.Repeat("  ", a.indent))
	}
}

func (a *FancyAstWriter) writeObjectType(o interface{}) {
	fmt.Fprintf(a.b, nodeColor.Paint("%T"), o)
}

func (a FancyAstWriter) bracketColor() console.ForegroundColor {
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
	fmt.Fprintf(a.b, strColor.Paint("%q"), s)
}

func (a *FancyAstWriter) writeStrProperty(name, value string) {
	a.writePropertyName(name)
	a.writeStr(value)
	a.writePropertyEnd()
}

func (a *FancyAstWriter) writeNumber(f float64) {
	fmt.Fprintf(a.b, numColor.Paint("%f"), f)
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
