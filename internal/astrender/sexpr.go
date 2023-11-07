package astrender

import (
	"fmt"
	"strings"

	ast "sudonters/zootler/pkg/rules/ast"

	"github.com/etc-sudonters/substrate/dontio"
)

type IdentifierCallback func(ident *ast.Identifier) ast.Expression

func NewSexpr(s ColorScheme) *sexprFormatter {
	f := new(sexprFormatter)
	f.scheme = s
	f.Clear()
	return f
}

type sexprFormatter struct {
	b          *strings.Builder
	depth      int
	scheme     ColorScheme
	identifier IdentifierCallback
}

func (s sexprFormatter) String() string {
	return s.b.String()
}

func (s *sexprFormatter) Clear() {
	s.b = new(strings.Builder)
	s.depth = 0
}

func (s *sexprFormatter) SetColorScheme(scheme ColorScheme) {
	s.scheme = scheme
}

func (s *sexprFormatter) SetIdentifier(cb IdentifierCallback) {
	s.identifier = cb
}

func (s *sexprFormatter) VisitBinOp(b *ast.BinOp) error {
	s.writeOpenParen()
	s.b.WriteString(s.scheme.Keyword.Paint(string(b.Op)))
	s.b.WriteRune(' ')
	ast.Visit(s, b.Left)
	s.b.WriteRune(' ')
	ast.Visit(s, b.Right)
	s.writeCloseParen()
	return nil
}
func (s *sexprFormatter) VisitBoolOp(b *ast.BoolOp) error {
	s.writeOpenParen()
	s.b.WriteString(s.scheme.Keyword.Paint(string(b.Op)))
	s.b.WriteRune(' ')
	ast.Visit(s, b.Left)
	s.b.WriteRune(' ')
	ast.Visit(s, b.Right)
	s.writeCloseParen()
	return nil
}

func (s *sexprFormatter) VisitLiteral(l *ast.Literal) error {
	switch l.Kind {
	case ast.LiteralBool:
		fmt.Fprintf(s.b, s.scheme.Boolean.Paint("%t"), l.Value)
	case ast.LiteralNum:
		fmt.Fprintf(s.b, s.scheme.Number.Paint("%.0f"), l.Value)
	case ast.LiteralStr:
		s.b.WriteString(s.scheme.String.Paint(l.Value.(string)))
	default:
		fmt.Fprintf(s.b, s.scheme.Property.Paint("%+v"), l.Value)
	}
	return nil
}

func (s *sexprFormatter) VisitCall(c *ast.Call) error {
	s.writeOpenParen()
	ast.Visit(s, c.Callee)
	for _, arg := range c.Args {
		s.b.WriteRune(' ')
		ast.Visit(s, arg)
	}

	if len(c.Args) == 0 {
		s.b.WriteString(" @")
	}

	s.writeCloseParen()
	return nil
}
func (s *sexprFormatter) VisitIdentifier(i *ast.Identifier) error {
	if s.identifier != nil {
		expr := s.identifier(i)
		if expr != nil && expr.Type() != ast.ExprIdentifier {
			return ast.Visit(s, expr)
		}
	}

	s.b.WriteString(s.scheme.Identifier.Paint(i.Value))
	return nil
}
func (s *sexprFormatter) VisitSubscript(r *ast.Subscript) error {
	s.writeOpenParen()
	s.b.WriteString("[] ")
	ast.Visit(s, r.Target)
	s.b.WriteRune(' ')
	ast.Visit(s, r.Index)
	s.writeCloseParen()
	return nil
}
func (s *sexprFormatter) VisitTuple(t *ast.Tuple) error {
	s.writeOpenParen()
	ast.Visit(s, t.Elems[0])
	for _, arg := range t.Elems[1:] {
		s.b.WriteRune(' ')
		ast.Visit(s, arg)
	}
	s.writeCloseParen()
	return nil
}
func (s *sexprFormatter) VisitUnary(u *ast.UnaryOp) error {
	s.writeOpenParen()
	s.b.WriteString(s.scheme.Keyword.Paint(string(u.Op)))
	s.b.WriteRune(' ')
	ast.Visit(s, u.Target)
	s.writeCloseParen()
	return nil
}

func (s *sexprFormatter) writeOpenParen() {
	s.b.WriteString(s.bracketColor().Paint("("))
	s.depth += 1
}

func (s *sexprFormatter) writeCloseParen() {
	s.depth -= 1
	s.b.WriteString(s.bracketColor().Paint(")"))
}

func (s sexprFormatter) bracketColor() dontio.Painter {
	return s.scheme.BracketFor(s.depth)
}
