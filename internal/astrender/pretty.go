package astrender

import (
	"fmt"
	"strings"

	"sudonters/zootler/pkg/rules/ast"

	"github.com/etc-sudonters/substrate/dontio"
)

func NewPretty() *PrettyPrint {
	return &PrettyPrint{
		b:      &strings.Builder{},
		indent: 0,
		scheme: LipglossColorScheme(),
	}
}

type PrettyPrint struct {
	b      *strings.Builder
	indent int
	scheme ColorScheme
}

func (p PrettyPrint) String() string {
	return p.b.String()
}

func (w *PrettyPrint) VisitAttrAccess(a *ast.AttrAccess) error {
	w.writeObject(a)
	w.writeProperty("Target", a.Target)
	w.writeProperty("Attr", a.Attr)
	w.writeObjectEnd()
	return nil
}
func (w *PrettyPrint) VisitBinOp(b *ast.BinOp) error {
	w.writeObject(b)
	w.writeKeywordProperty("Op", string(b.Op))
	w.writeProperty("Left", b.Left)
	w.writeProperty("Right", b.Right)
	w.writeObjectEnd()
	return nil
}

func (w *PrettyPrint) VisitBoolOp(b *ast.BoolOp) error {
	w.writeObject(b)
	w.writeKeywordProperty("Op", string(b.Op))
	w.writeProperty("Left", b.Left)
	w.writeProperty("Right", b.Right)
	w.writeObjectEnd()
	return nil
}

func (w *PrettyPrint) VisitBoolean(b *ast.Boolean) error {
	w.writeObject(b)
	w.writeBoolProperty("Value", b.Value)
	w.writeObjectEnd()
	return nil
}
func (w *PrettyPrint) VisitCall(c *ast.Call) error {
	w.writeObject(c)
	w.writeProperty("Fn", c.Callee)
	w.writePropertyName("Args")
	w.writeArrStart()
	for _, node := range c.Args {
		w.writeArrElem(node)
	}
	w.writeArrEnd()
	w.writeObjectEnd()
	return nil
}
func (w *PrettyPrint) VisitIdentifier(i *ast.Identifier) error {
	w.writeObject(i)
	w.writeIdentifierProperty("Value", i.Value)
	w.writeObjectEnd()
	return nil
}

func (w *PrettyPrint) VisitNumber(n *ast.Number) error {
	w.writeObject(n)
	w.writeNumProperty("Value", n.Value)
	w.writeObjectEnd()
	return nil
}
func (w *PrettyPrint) VisitString(s *ast.String) error {
	w.writeObject(s)
	w.writeStrProperty("Value", s.Value)
	w.writeObjectEnd()
	return nil
}
func (w *PrettyPrint) VisitSubscript(s *ast.Subscript) error {
	w.writeObject(s)
	w.writeProperty("Target", s.Target)
	w.writeProperty("Index", s.Index)
	w.writeObjectEnd()
	return nil
}
func (w *PrettyPrint) VisitTuple(t *ast.Tuple) error {
	w.writeObject(t)
	w.writePropertyName("Elems")
	w.writeArrStart()
	for _, e := range t.Elems {
		w.writeArrElem(e)
	}
	w.writeArrEnd()
	w.writeObjectEnd()
	return nil
}

func (w *PrettyPrint) VisitUnary(u *ast.UnaryOp) error {
	w.writeObject(u)
	w.writeKeywordProperty("Op", string(u.Op))
	w.writeProperty("Target", u.Target)
	w.writeObjectEnd()
	return nil
}

func (a PrettyPrint) writeIndent() {
	if a.indent > 0 {
		fmt.Fprint(a.b, strings.Repeat("  ", a.indent))
	}
}

func (a *PrettyPrint) writeObjectType(o interface{}) {
	fmt.Fprintf(a.b, a.scheme.Node.Paint("%T"), o)
}

func (a PrettyPrint) bracketColor() dontio.Painter {
	return a.scheme.BracketFor(a.indent)
}

func (a *PrettyPrint) writeColoredBracket(s string) {
	a.b.WriteString(a.bracketColor().Paint(s))
}

func (a *PrettyPrint) writeObject(o interface{}) {
	a.writeObjectType(o)
	a.b.WriteRune(' ')
	a.writeColoredBracket("{")
	a.b.WriteRune('\n')
	a.indent += 1
}
func (a *PrettyPrint) writeObjectEnd() {
	a.indent -= 1
	a.writeIndent()
	a.writeColoredBracket("}")
}

func (a *PrettyPrint) writePropertyName(name string) {
	a.writeIndent()
	a.b.WriteString(a.scheme.Property.Paint(name))
	a.b.WriteString(":  ")
}

func (a *PrettyPrint) writePropertyEnd() {
	a.b.Write([]byte(",\n"))
}

func (a *PrettyPrint) writeStr(s string) {
	fmt.Fprintf(a.b, a.scheme.String.Paint("%s"), s)
}

func (a *PrettyPrint) writeStrProperty(name, value string) {
	a.writePropertyName(name)
	a.writeStr(value)
	a.writePropertyEnd()
}

func (a *PrettyPrint) writeNumber(f float64) {
	fmt.Fprintf(a.b, a.scheme.Number.Paint("%.0f"), f)
}

func (a *PrettyPrint) writeNumProperty(name string, value float64) {
	a.writePropertyName(name)
	a.writeNumber(value)
	a.writePropertyEnd()
}

func (a *PrettyPrint) writeBool(b bool) {
	fmt.Fprintf(a.b, a.scheme.Boolean.Paint("%t"), b)
}

func (a *PrettyPrint) writeBoolProperty(name string, value bool) {
	a.writePropertyName(name)
	a.writeBool(value)
	a.writePropertyEnd()
}

func (a *PrettyPrint) writeIdentifier(i string) {
	fmt.Fprint(a.b, a.scheme.Identifier.Paint(i))
}

func (a *PrettyPrint) writeIdentifierProperty(name, i string) {
	a.writePropertyName(name)
	a.writeIdentifier(i)
	a.writePropertyEnd()
}

func (a *PrettyPrint) writeKeyword(kw string) {
	fmt.Fprint(a.b, a.scheme.Keyword.Paint(kw))
}

func (a *PrettyPrint) writeKeywordProperty(name, kw string) {
	a.writePropertyName(name)
	a.writeKeyword(kw)
	a.writePropertyEnd()
}

func (a *PrettyPrint) writeProperty(name string, value ast.Expression) error {
	a.writePropertyName(name)
	ast.Visit(a, value)
	a.writePropertyEnd()
	return nil
}

func (a *PrettyPrint) writeArrStart() {
	a.writeColoredBracket("[")
	a.b.WriteRune('\n')
	a.indent += 1
}

func (a *PrettyPrint) writeArrElem(v ast.Expression) {
	a.writeIndent()
	ast.Visit(a, v)
	fmt.Fprintf(a.b, ",\n")
}

func (a *PrettyPrint) writeArrEnd() {
	a.indent -= 1
	a.writeIndent()
	a.writeColoredBracket("]")
	fmt.Fprint(a.b, ",\n")
}

func (a *PrettyPrint) finish() string {
	return a.b.String()
}
