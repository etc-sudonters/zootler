package debug

import (
	"fmt"
	"strings"
	"sudonters/zootler/icearrow/ast"
	"sudonters/zootler/icearrow/parser"
)

func PtSexpr(ast parser.Expression) string {
	var pt ptsexpr
	parser.Visit(&pt, ast)
	return pt.sb.String()
}

func AstSexpr(node ast.Node) string {
	var s astsexpr
	ast.Visit(&s, node)
	return s.sb.String()
}

type astsexpr struct {
	sb strings.Builder
}

type ptsexpr struct {
	sb strings.Builder
}

func (pt *ptsexpr) VisitBinOp(ast *parser.BinOp) error {
	pt.sb.WriteRune('(')
	pt.sb.WriteString(string(ast.Op))
	pt.sb.WriteRune(' ')
	parser.Visit(pt, ast.Left)
	pt.sb.WriteRune(' ')
	parser.Visit(pt, ast.Right)
	pt.sb.WriteRune(')')
	return nil
}
func (pt *ptsexpr) VisitBoolOp(ast *parser.BoolOp) error {
	pt.sb.WriteRune('(')
	pt.sb.WriteString(string(ast.Op))
	pt.sb.WriteRune(' ')
	parser.Visit(pt, ast.Left)
	pt.sb.WriteRune(' ')
	parser.Visit(pt, ast.Right)
	pt.sb.WriteRune(')')
	return nil
}
func (pt *ptsexpr) VisitCall(ast *parser.Call) error {
	pt.sb.WriteString("([")
	parser.Visit(pt, ast.Callee)
	pt.sb.WriteRune(' ')
	for i, arg := range ast.Args {
		if i != 0 {
			pt.sb.WriteRune(' ')
		}
		parser.Visit(pt, arg)
	}
	pt.sb.WriteString("])")
	return nil
}
func (pt *ptsexpr) VisitIdentifier(ast *parser.Identifier) error {
	pt.sb.WriteString(ast.Value)
	return nil
}

func (pt *ptsexpr) VisitSubscript(ast *parser.Subscript) error {
	pt.sb.WriteString("(load_setting_2 ")
	parser.Visit(pt, ast.Target)
	pt.sb.WriteRune(' ')
	parser.Visit(pt, ast.Index)
	pt.sb.WriteRune(')')
	return nil
}
func (pt *ptsexpr) VisitTuple(ast *parser.Tuple) error {
	pt.sb.WriteRune('(')
	pt.sb.WriteRune(' ')
	for idx, elm := range ast.Elems {
		if idx != 0 {
			pt.sb.WriteRune(' ')
		}
		parser.Visit(pt, elm)
	}
	return nil
}
func (pt *ptsexpr) VisitUnary(ast *parser.UnaryOp) error {
	fmt.Fprintf(&pt.sb, "(%s ", ast.Op)
	parser.Visit(pt, ast.Target)
	pt.sb.WriteRune(')')
	return nil
}
func (pt *ptsexpr) VisitLiteral(ast *parser.Literal) error {
	fmt.Fprintf(&pt.sb, "%+v", ast.Value)
	return nil
}

func (a *astsexpr) Comparison(node *ast.Comparison) error {
	a.sb.WriteRune('(')
	switch node.Op {
	case ast.AST_CMP_EQ:
		a.sb.WriteString("== ")
		break
	case ast.AST_CMP_NQ:
		a.sb.WriteString("!= ")
		break
	case ast.AST_CMP_LT:
		a.sb.WriteString("< ")
		break
	}

	ast.Visit(a, node.LHS)
	a.sb.WriteRune(' ')
	ast.Visit(a, node.RHS)
	a.sb.WriteRune(')')
	return nil
}
func (a *astsexpr) BooleanOp(node *ast.BooleanOp) error {
	a.sb.WriteRune('(')
	switch node.Op {
	case ast.AST_BOOL_AND:
		a.sb.WriteString("and ")
		break
	case ast.AST_BOOL_OR:
		a.sb.WriteString("or ")
		break
	case ast.AST_BOOL_NEGATE:
		a.sb.WriteString("not ")
		ast.Visit(a, node.LHS)
		a.sb.WriteRune(')')
		return nil
	}

	ast.Visit(a, node.LHS)
	a.sb.WriteRune(' ')
	ast.Visit(a, node.RHS)
	a.sb.WriteRune(')')
	return nil
}
func (a *astsexpr) Call(node *ast.Call) error {
	fmt.Fprintf(&a.sb, "([%s", node.Callee)
	a.sb.WriteRune(' ')
	for i, arg := range node.Args {
		if i != 0 {
			a.sb.WriteRune(' ')
		}
		ast.Visit(a, arg)
	}
	a.sb.WriteString("])")
	return nil
}
func (a *astsexpr) Identifier(node *ast.Identifier) error {
	switch node.Kind {
	case ast.AST_IDENT_EXP, ast.AST_IDENT_BIF:
		fmt.Fprintf(&a.sb, "<>%s", node.Name)
		break
	case ast.AST_IDENT_TOK, ast.AST_IDENT_EVT:
		fmt.Fprintf(&a.sb, "@%s", node.Name)
		break
	case ast.AST_IDENT_SET, ast.AST_IDENT_TRK:
		fmt.Fprintf(&a.sb, "&%s", node.Name)
		break
	case ast.AST_IDENT_VAR:
		fmt.Fprintf(&a.sb, "$%s", node.Name)
		break
	case ast.AST_IDENT_UNK, ast.AST_IDENT_UNP:
		fallthrough
	default:
		fmt.Fprintf(&a.sb, "??%s", node.Name)
		break
	}
	return nil
}
func (a *astsexpr) Literal(node *ast.Literal) error {
	fmt.Fprintf(&a.sb, "%+v", node.Value)
	return nil
}
func (a *astsexpr) Empty(node *ast.Empty) error {
	return nil
}
