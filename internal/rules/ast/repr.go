package ast

import (
	"fmt"
	"strings"
)

type AstRender struct {
	s strings.Builder
}

func (r AstRender) String() string {
	return r.s.String()
}

func (r *AstRender) start() {
	r.s.WriteRune('(')
}

func (r *AstRender) end() error {
	r.s.WriteRune(')')
	return nil
}

func (r *AstRender) Comparison(ast *Comparison) error {
	r.start()
	switch ast.Op {
	case AST_CMP_EQ:
		r.s.WriteString("EQ ")
		break
	case AST_CMP_NQ:
		r.s.WriteString("NQ ")
		break
	case AST_CMP_LT:
		r.s.WriteString("LT ")
		break
	}

	Visit(r, ast.LHS)
	Visit(r, ast.RHS)
	return r.end()
}

func (r *AstRender) BooleanOp(ast *BooleanOp) error {
	r.start()
	switch ast.Op {
	case AST_BOOL_NEGATE:
		r.s.WriteString("NO! ")
		Visit(r, ast.LHS)
		return r.end()
	case AST_BOOL_AND:
		r.s.WriteString("BOTH ")
		break
	case AST_BOOL_OR:
		r.s.WriteString("EITHER ")
		break
	}
	Visit(r, ast.LHS)
	r.s.WriteRune(' ')
	Visit(r, ast.RHS)
	return r.end()
}

func (r *AstRender) Call(ast *Call) error {
	r.start()

	if ast.Macro {
		fmt.Fprintf(&r.s, "{MACRO %s}", ast.Callee)
	} else {
		fmt.Fprintf(&r.s, "%s", ast.Callee)
	}

	for _, arg := range ast.Args {
		r.s.WriteRune(' ')
		Visit(r, arg)
	}

	return r.end()
}

func (r *AstRender) Identifier(ast *Identifier) error {
	switch ast.Kind {
	case AST_IDENT_UNK:
		fmt.Fprintf(&r.s, "UNK{%s}", ast.Name)
		break
	case AST_IDENT_SET:
		fmt.Fprintf(&r.s, "SET{%s}", ast.Name)
		break
	case AST_IDENT_TOK:
		fmt.Fprintf(&r.s, "TOK{%s}", ast.Name)
		break
	case AST_IDENT_TRK:
		fmt.Fprintf(&r.s, "TRK{%s}", ast.Name)
		break
	default:
		panic("unreachable")
	}

	return nil
}

func (r *AstRender) Literal(ast *Literal) error {
	switch ast.Kind {
	case AST_LIT_NUM:
		fmt.Fprintf(&r.s, "%f", ast.Value.(float64))
		return nil
	case AST_LIT_BOOL:
		fmt.Fprintf(&r.s, "%t", ast.Value.(bool))
		return nil
	case AST_LIT_STR:
		fmt.Fprintf(&r.s, "'%s'", ast.Value.(string))
		return nil
	default:
		panic("unreachable")
	}
}

func (r *AstRender) Empty(ast *Empty) error { return nil }
