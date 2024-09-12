package analysis

import (
	"regexp"
	"strings"
	"sudonters/zootler/icearrow/ast"
	"sudonters/zootler/internal"

	"github.com/etc-sudonters/substrate/slipup"
)

var looksLikeToken = regexp.MustCompile("[A-Z][a-zA-Z_]")

type AnalysisContext struct {
	Expansions map[string]replacement
	Promotions map[string]ast.AstIdentifierKind

	tokenNames  map[internal.NormalizedStr]struct{}
	settingName map[internal.NormalizedStr]struct{}
	varName     map[internal.NormalizedStr]struct{}
	builtInName map[internal.NormalizedStr]struct{}
	unpromoted  map[string]struct{}
}

func (a *AnalysisContext) Unpromoted() []string {
	names := make([]string, 0, len(a.unpromoted))
	for k := range a.unpromoted {
		names = append(names, k)
	}

	return names
}

func (a *AnalysisContext) NameToken(name string) {
	a.tokenNames[internal.Normalize(name)] = struct{}{}
}

func (a *AnalysisContext) isToken(name string) bool {
	_, exists := a.tokenNames[internal.Normalize(name)]
	return exists
}

func (a *AnalysisContext) NameSetting(name string) {
	a.settingName[internal.Normalize(name)] = struct{}{}
}

func (a *AnalysisContext) isSetting(name string) bool {
	_, exists := a.settingName[internal.Normalize(name)]
	return exists
}

func (a *AnalysisContext) NameBuiltIn(name string) {
	a.builtInName[internal.Normalize(name)] = struct{}{}
}

func (a *AnalysisContext) isBuiltIn(name string) bool {
	_, exists := a.builtInName[internal.Normalize(name)]
	return exists
}

func (a *AnalysisContext) AddExpansion(name string, params []string, body ast.Node) {
	a.Expansions[name] = replacement{name, params, body}
}

func NewAnalysis() AnalysisContext {
	var ac AnalysisContext
	ac.Expansions = make(map[string]replacement)
	ac.settingName = make(map[internal.NormalizedStr]struct{})
	ac.tokenNames = make(map[internal.NormalizedStr]struct{})
	ac.builtInName = make(map[internal.NormalizedStr]struct{})
	ac.unpromoted = make(map[string]struct{})

	ac.varName = make(map[internal.NormalizedStr]struct{})
	ac.varName[internal.NormalizedStr("age")] = struct{}{}
	ac.varName[internal.NormalizedStr("either")] = struct{}{}
	ac.varName[internal.NormalizedStr("both")] = struct{}{}
	ac.varName[internal.NormalizedStr("adult")] = struct{}{}
	ac.varName[internal.NormalizedStr("child")] = struct{}{}
	ac.varName[internal.NormalizedStr("fire")] = struct{}{}
	ac.varName[internal.NormalizedStr("forest")] = struct{}{}
	ac.varName[internal.NormalizedStr("light")] = struct{}{}
	ac.varName[internal.NormalizedStr("shadow")] = struct{}{}
	ac.varName[internal.NormalizedStr("spirit")] = struct{}{}
	ac.varName[internal.NormalizedStr("water")] = struct{}{}

	return ac
}

func (ac *AnalysisContext) isExpandable(name string) bool {
	_, exists := ac.Expansions[name]
	return exists
}

func Analyze(node ast.Node, ctx *AnalysisContext) (ast.Node, error) {
	report := analyze(node, ctx)
	node, _ = tagIdentifiers(node, ctx)
	if report.compares {
		node, _ = constCompares(node)
	}

	if report.branches {
		node, _ = constBranches(node)
	}

	if report.expansions {
		node, _ = expand(node, ctx)
	}

	node, _ = promotions(node, ctx)

	return node, nil
}

func tagIdentifiers(node ast.Node, ctx *AnalysisContext) (ast.Node, error) {
	var re rewriter
	re.identifier = func(_ *rewriter, ident *ast.Identifier) (ast.Node, error) {
		name := ident.Name
		normaled := internal.Normalize(name)
		if ident.Kind != ast.AST_IDENT_UNK {
			return ident, nil
		} else if strings.HasPrefix(ident.Name, "logic_") {
			ident.Kind = ast.AST_IDENT_TRK
		} else if ctx.isToken(ident.Name) {
			ident.Kind = ast.AST_IDENT_TOK
		} else if ctx.isSetting(ident.Name) {
			ident.Kind = ast.AST_IDENT_SET
		} else if ctx.isExpandable(name) {
			ident.Kind = ast.AST_IDENT_EXP
		} else if ctx.isBuiltIn(name) {
			ident.Kind = ast.AST_IDENT_BIF
		} else if _, isVar := ctx.varName[normaled]; isVar {
			ident.Kind = ast.AST_IDENT_VAR
		} else if looksLikeToken.MatchString(name) {
			ident.Kind = ast.AST_IDENT_TOK
		} else {
			ident.Kind = ast.AST_IDENT_UNP
			ctx.unpromoted[name] = struct{}{}
		}
		return ident, nil
	}

	return ast.Transform(&re, node)
}

func constCompares(node ast.Node) (ast.Node, error) {
	var re rewriter
	re.compares = func(re *rewriter, compare *ast.Comparison) (ast.Node, error) {
		if compare.Op != ast.AST_CMP_EQ && compare.Op != ast.AST_CMP_NQ {
			return compare, nil
		}
		lhs, _ := ast.Transform(re, compare.LHS)
		rhs, _ := ast.Transform(re, compare.RHS)

		lhIdent, lhFail := ast.Unify(lhs,
			func(node *ast.Identifier) (*ast.Identifier, error) {
				return node, nil
			},
			func(node *ast.Call) (*ast.Identifier, error) {
				if node.Callee != "has" {
					return nil, slipup.Createf("no identifier to extract")
				}
				return ast.MustAssertAs[*ast.Identifier](node.Args[0]), nil
			})

		rhIdent, rhFail := ast.Unify(rhs,
			func(node *ast.Identifier) (*ast.Identifier, error) {
				return node, nil
			},
			func(node *ast.Call) (*ast.Identifier, error) {
				if node.Callee != "has" {
					return nil, slipup.Createf("no identifier to extract")
				}
				return ast.MustAssertAs[*ast.Identifier](node.Args[0]), nil
			})

		if lhFail != nil || rhFail != nil {
			return compare, nil
		}

		if lhIdent.Kind != rhIdent.Kind {
			return compare, nil
		}
		isSameIdent := internal.Normalize(lhIdent.Name) == internal.Normalize(rhIdent.Name)
		wantsSameIdent := compare.Op == ast.AST_CMP_EQ
		return ast.LiteralBool(isSameIdent && wantsSameIdent), nil
	}
	return ast.Transform(&re, node)
}

func constBranches(node ast.Node) (ast.Node, error) {
	var re rewriter
	re.booleans = func(trans *rewriter, boolop *ast.BooleanOp) (ast.Node, error) {
		op := boolop.Op
		if op == ast.AST_BOOL_NEGATE {
			return boolop, nil
		}
		lhs, _ := ast.Transform(trans, boolop.LHS)
		lhLit, lhIsLit := ast.AssertAs[*ast.Literal](lhs)
		if lhIsLit {
			lhBool, lhIsBool := lhLit.Value.(bool)
			if lhIsBool {
				if op == ast.AST_BOOL_AND && lhBool {
					return ast.Transform(trans, boolop.RHS)
				}
				if op == ast.AST_BOOL_AND && !lhBool {
					return ast.LiteralBool(false), nil
				}
				if op == ast.AST_BOOL_OR && lhBool {
					return ast.LiteralBool(true), nil
				}
				if op == ast.AST_BOOL_OR && !lhBool {
					return ast.Transform(trans, boolop.RHS)
				}
			}
		}

		rhs, _ := ast.Transform(trans, boolop.RHS)
		rhLit, rhIsLit := ast.AssertAs[*ast.Literal](rhs)
		if !rhIsLit {
			return &ast.BooleanOp{
				LHS: lhs,
				RHS: rhs,
				Op:  boolop.Op,
			}, nil
		}
		rhBool, rhIsBool := rhLit.Value.(bool)

		if !rhIsBool {
			return &ast.BooleanOp{
				LHS: lhs,
				RHS: rhs,
				Op:  boolop.Op,
			}, nil

		}
		if op == ast.AST_BOOL_AND && rhBool {
			return ast.Transform(trans, boolop.LHS)
		}
		if op == ast.AST_BOOL_AND && !rhBool {
			return ast.LiteralBool(false), nil
		}
		if op == ast.AST_BOOL_OR && rhBool {
			return ast.LiteralBool(true), nil
		}
		if op == ast.AST_BOOL_OR && !rhBool {
			return ast.Transform(trans, boolop.LHS)
		}

		panic("unreachable")
	}
	return ast.Transform(&re, node)
}

func replaceIdentifiers(fragment ast.Node, repls map[string]ast.Node) ast.Node {
	var re rewriter
	re.identifier = func(_ *rewriter, ident *ast.Identifier) (ast.Node, error) {
		newNodes, replace := repls[ident.Name]
		if replace {
			return newNodes, nil
		}
		return ident, nil
	}
	nodes, _ := ast.Transform(&re, fragment)
	return nodes
}

func expand(node ast.Node, ctx *AnalysisContext) (ast.Node, error) {
	var re rewriter
	re.call = func(trans *rewriter, call *ast.Call) (ast.Node, error) {
		expansion, has := ctx.Expansions[call.Callee]
		if !has {
			return call, nil
		}
		replacements := map[string]ast.Node{}
		for idx, name := range expansion.Params {
			replacements[name] = call.Args[idx]
		}

		body, _ := copyfragment(expansion.Body)
		return replaceIdentifiers(body, replacements), nil
	}

	re.identifier = func(_ *rewriter, ident *ast.Identifier) (ast.Node, error) {
		expansion, has := ctx.Expansions[ident.Name]
		if !has {
			return ident, nil
		}

		if len(expansion.Params) != 0 {
			return nil, slipup.Createf("expected 0 parameter expansion for %s", ident.Name)
		}

		return copyfragment(expansion.Body)
	}
	expansion, _ := ast.Transform(&re, node)
	return Analyze(expansion, ctx)
}

func promotions(node ast.Node, _ *AnalysisContext) (ast.Node, error) {
	var re rewriter
	re.call = func(r *rewriter, c *ast.Call) (ast.Node, error) {
		return c, nil
	}
	re.identifier = func(r *rewriter, ident *ast.Identifier) (ast.Node, error) {
		switch ident.Kind {
		case ast.AST_IDENT_TOK, ast.AST_IDENT_EVT:
			return &ast.Call{
				Callee: "has",
				Args:   []ast.Node{ident, ast.LiteralNumber(1)},
			}, nil
		case ast.AST_IDENT_SET:
			return &ast.Call{
				Callee: "load_setting",
				Args:   []ast.Node{ident},
			}, nil
		}
		return ident, nil
	}
	return ast.Transform(&re, node)
}

type rewriter struct {
	compares   func(*rewriter, *ast.Comparison) (ast.Node, error)
	booleans   func(*rewriter, *ast.BooleanOp) (ast.Node, error)
	call       func(*rewriter, *ast.Call) (ast.Node, error)
	identifier func(*rewriter, *ast.Identifier) (ast.Node, error)
	literal    func(*rewriter, *ast.Literal) (ast.Node, error)
}

func (r *rewriter) Comparison(node *ast.Comparison) (ast.Node, error) {
	if r.compares != nil {
		return r.compares(r, node)
	}

	n := ast.Comparison{
		Op: node.Op,
	}

	n.LHS, _ = ast.Transform(r, node.LHS)
	n.RHS, _ = ast.Transform(r, node.RHS)
	return &n, nil
}

func (r *rewriter) BooleanOp(node *ast.BooleanOp) (ast.Node, error) {
	if r.booleans != nil {
		return r.booleans(r, node)
	}
	n := ast.BooleanOp{
		Op: node.Op,
	}

	n.LHS, _ = ast.Transform(r, node.LHS)
	n.RHS, _ = ast.Transform(r, node.RHS)
	return &n, nil

}

func (r *rewriter) Call(node *ast.Call) (ast.Node, error) {
	if r.call != nil {
		return r.call(r, node)
	}

	call := ast.Call{
		Callee: node.Callee,
		Args:   make([]ast.Node, len(node.Args)),
	}

	for idx, arg := range node.Args {
		call.Args[idx], _ = ast.Transform(r, arg)
	}

	return &call, nil
}

func (r *rewriter) Identifier(node *ast.Identifier) (ast.Node, error) {
	if r.identifier != nil {
		return r.identifier(r, node)
	}
	return node, nil
}

func (r *rewriter) Literal(node *ast.Literal) (ast.Node, error) {
	if r.literal != nil {
		return r.literal(r, node)
	}
	return node, nil
}

func (r *rewriter) Empty(node *ast.Empty) (ast.Node, error) {
	return node, nil
}

type replacement struct {
	Name   string
	Params []string
	Body   ast.Node
}
