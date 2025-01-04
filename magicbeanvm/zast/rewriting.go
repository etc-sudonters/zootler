package zast

import (
	"fmt"

	"github.com/etc-sudonters/substrate/stageleft"
)

func NewRewriteContext() RewriteContext {
	var ctx RewriteContext
	ctx.values = make(map[any]any)
	ctx.names = make(map[string]IdentKind)
	ctx.rewriters = make([]Rewriter, 0, 16)
	return ctx
}

type RewriteContext struct {
	values    map[any]any
	names     map[string]IdentKind
	rewriters []Rewriter
}

func (ctx *RewriteContext) KindFor(who string) (kind IdentKind, exists bool) {
	kind, exists = ctx.names[who]
	return kind, exists
}

func (ctx *RewriteContext) AttachRewriter(rw Rewriter) {
	ctx.rewriters = append(ctx.rewriters, rw)
}

func (ctx *RewriteContext) Declare(name string, kind IdentKind) {
	ctx.names[name] = kind
}

func (ctx *RewriteContext) Store(key any, value any) {
	ctx.values[key] = value
}

func (ctx *RewriteContext) Retrieve(key any) (any, bool) {
	value, exists := ctx.values[key]
	return value, exists
}

func (ctx *RewriteContext) StoreIfNotPresent(key any, f func() any) {
	ctx.RetrieveOrStore(key, f)
}

func (ctx *RewriteContext) RetrieveOrStore(key any, f func() any) (value any, existed bool) {
	value, existed = ctx.Retrieve(key)
	if !existed {
		value = f()
		ctx.Store(key, value)
	}

	return
}

func Analyze(ctx *RewriteContext, ast Ast) (Ast, error) {
	var astErr error
	for i := range ctx.rewriters {
		ast, astErr = ctx.rewriters[i].Rewrite(ast)
		if astErr != nil {
			break
		}
	}
	return ast, astErr
}

func RewriteMany(rewrite rewrite, ast ...Ast) ([]Ast, error) {
	var result []Ast
	if n := len(ast); n == 0 {
		return nil, nil
	} else {
		result = make([]Ast, n)
	}

	var err error
	for i := range ast {
		result[i], err = rewrite(ast[i])
		if err != nil {
			break
		}
	}

	return result, err
}

type Rewriter struct {
	Boolean    RewriteFunc[Boolean]
	Comparison RewriteFunc[Comparison]
	Identifier RewriteFunc[Identifier]
	Invoke     RewriteFunc[Invoke]
	Value      RewriteFunc[Value]
}

func (r Rewriter) Func() rewrite {
	var rewrite rewrite
	rewrite = func(a Ast) (Ast, error) {
		switch ast := a.(type) {
		case Boolean:
			if r.Boolean == nil {
				lhs, lhsErr := rewrite(ast.LHS)
				if lhsErr != nil {
					panic(fmt.Errorf("error handling not impled: %w", lhsErr))
				}
				rhs, rhsErr := rewrite(ast.RHS)
				if rhsErr != nil {
					panic(fmt.Errorf("error handling not impled: %w", rhsErr))
				}

				return Boolean{lhs, rhs, ast.Op}, nil
			}
			return r.Boolean(ast, rewrite)
		case Comparison:
			if r.Comparison == nil {
				lhs, lhsErr := rewrite(ast.LHS)
				if lhsErr != nil {
					panic(fmt.Errorf("error handling not impled: %w", lhsErr))
				}
				rhs, rhsErr := rewrite(ast.RHS)
				if rhsErr != nil {
					panic(fmt.Errorf("error handling not impled: %w", rhsErr))
				}

				return Comparison{lhs, rhs, ast.Op}, nil
			}
			return r.Comparison(ast, rewrite)
		case Identifier:
			if r.Identifier == nil {
				return ast, nil
			}
			return r.Identifier(ast, rewrite)
		case Invoke:
			if r.Invoke == nil {
				target, targetErr := rewrite(ast.Target)
				if targetErr != nil {
				}
				ident, isIdent := target.(Identifier)
				if !isIdent {
					panic(fmt.Errorf("expected identifer for function target, received: %T", target))
				}

				invoke := Invoke{
					Target: ident,
					Args:   make([]Ast, len(ast.Args)),
				}

				var argErr error
				for i := range ast.Args {
					invoke.Args[i], argErr = rewrite(ast.Args[i])
					if argErr != nil {
						panic(fmt.Errorf("error handling not impled: %w", argErr))
					}
				}

				return invoke, nil
			}
			return r.Invoke(ast, rewrite)
		case Value:
			if r.Value == nil {
				return ast, nil
			}
			return r.Value(ast, rewrite)
		case Hole:
			return ast, nil
		default:
			panic(stageleft.AttachExitCode(
				fmt.Errorf("unknown node type %T", ast),
				stageleft.ExitCode(90),
			))
		}
	}

	return rewrite
}

func (r Rewriter) Rewrite(node Ast) (Ast, error) {
	visit := r.Func()
	return visit(node)
}

type rewrite = func(Ast) (Ast, error)
type RewriteFunc[T Ast] func(node T, visit rewrite) (Ast, error)
