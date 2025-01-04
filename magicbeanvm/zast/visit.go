package zast

import (
	"fmt"
	"strings"

	"github.com/etc-sudonters/substrate/stageleft"
)

func Render(a Ast) string {
	var sb strings.Builder
	r := Renderer(&sb)
	r.Visit(a)
	return sb.String()
}

func Renderer(sb *strings.Builder) Visitor {
	renderer := renderer{sb}
	return Visitor{
		Boolean:    renderer.Boolean,
		Comparison: renderer.Comparison,
		Identifier: renderer.Identifier,
		Invoke:     renderer.Invoke,
		Value:      renderer.Value,
	}

}

type visit func(Ast) error

type renderer struct {
	*strings.Builder
}

func (r renderer) Boolean(b Boolean, visit visit) error {
	r.WriteRune('(')
	switch b.Op {
	case BoolAnd:
		r.WriteString("and ")
		visit(b.LHS)
		r.WriteRune(' ')
		visit(b.RHS)
	case BoolOr:
		r.WriteString("or ")
		visit(b.LHS)
		r.WriteRune(' ')
		visit(b.RHS)
	case BoolInvert:
		r.WriteString("not ")
		visit(b.LHS)
	default:
		panic("unrecognized")
	}
	r.WriteRune(')')
	return nil
}

func (r renderer) Comparison(c Comparison, visit visit) error {
	r.WriteRune('(')
	switch c.Op {
	case CompareEqual:
		fmt.Fprint(r, "== ")
		visit(c.LHS)
		r.WriteRune(' ')
		visit(c.RHS)
	case CompareNotEqual:
		fmt.Fprint(r, "!= ")
		visit(c.LHS)
		r.WriteRune(' ')
		visit(c.RHS)
	case CompareLessThan:
		r.WriteRune('<')
		r.WriteRune(' ')
		visit(c.LHS)
		r.WriteRune(' ')
		visit(c.RHS)
	}
	r.WriteRune(')')
	return nil
}

func (r renderer) Identifier(i Identifier, visit visit) error {
	fmt.Fprintf(r, "(<LOAD_%02X> %q)", i.Kind, i.Name)
	return nil
}

func (r renderer) Invoke(i Invoke, visit visit) error {
	r.WriteString("(invoke ")
	visit(i.Target)

	if n := len(i.Args); n > 0 {
		for idx := range n {
			r.WriteRune(' ')
			visit(i.Args[idx])
		}
	}
	r.WriteRune(')')
	return nil

}

func (r renderer) Value(v Value, visit visit) error {
	fmt.Fprintf(r, "%v", v.any)
	return nil
}

type VisitFunc[T Ast] func(T, visit) error

type Visitor struct {
	Boolean    VisitFunc[Boolean]
	Comparison VisitFunc[Comparison]
	Identifier VisitFunc[Identifier]
	Invoke     VisitFunc[Invoke]
	Value      VisitFunc[Value]
	Hole       VisitFunc[Hole]
}

func (v Visitor) Visit(a Ast) error {
	var visit func(Ast) error
	visit = func(a Ast) error {
		switch ast := a.(type) {
		case Boolean:
			if v.Boolean == nil {
				lhsErr := visit(ast.LHS)
				if lhsErr != nil {
					panic(fmt.Errorf("error handling not impled: %w", lhsErr))
				}
				rhsErr := visit(ast.RHS)
				if rhsErr != nil {
					panic(fmt.Errorf("error handling not impled: %w", rhsErr))
				}

				return nil
			}
			return v.Boolean(ast, visit)
		case Comparison:
			if v.Comparison == nil {
				lhsErr := visit(ast.LHS)
				if lhsErr != nil {
					panic(fmt.Errorf("error handling not impled: %w", lhsErr))
				}
				rhsErr := visit(ast.RHS)
				if rhsErr != nil {
					panic(fmt.Errorf("error handling not impled: %w", rhsErr))
				}

				return nil
			}
			return v.Comparison(ast, visit)
		case Identifier:
			if v.Identifier == nil {
				return nil
			}
			return v.Identifier(ast, visit)
		case Invoke:
			if v.Invoke == nil {
				targetErr := visit(ast.Target)
				if targetErr != nil {
				}
				for i := range ast.Args {
					argErr := visit(ast.Args[i])
					if argErr != nil {
						panic(fmt.Errorf("error handling not impled: %w", argErr))
					}
				}

				return nil
			}
			return v.Invoke(ast, visit)
		case Value:
			if v.Value == nil {
				return nil
			}
			return v.Value(ast, visit)
		case Hole:
			if v.Hole == nil {
				return nil
			}
			return v.Hole(ast, visit)
		default:
			panic(stageleft.AttachExitCode(
				fmt.Errorf("unknown node type %T", ast),
				stageleft.ExitCode(90),
			))
		}
	}
	return visit(a)
}
