package analysis

import (
	"errors"
	"regexp"
	"strings"
	"sudonters/zootler/icearrow/ast"
	"sudonters/zootler/internal"

	"github.com/etc-sudonters/substrate/slipup"
)

func NewAnalysis() AnalysisContext {
	var ac AnalysisContext
	ac.expansions = make(map[string]replacement)
	ac.settingName = make(map[internal.NormalizedStr]struct{})
	ac.tokenNames = make(map[internal.NormalizedStr]struct{})
	ac.builtInName = make(map[internal.NormalizedStr]struct{})
	ac.unpromoted = make(map[string]struct{})

	ac.varName = make(map[internal.NormalizedStr]struct{})
	varNames := []string{
		"adult", "age", "both", "child", "either", "fire", "forest", "light",
		"shadow", "spirit", "water",
	}
	for _, name := range varNames {
		ac.varName[internal.Normalize(name)] = struct{}{}
	}

	return ac
}

func Analyze(node ast.Node, ctx *AnalysisContext) (ast.Node, error) {
	ctx.expandTokenLike = false
	report := analyze(node, ctx)
	if report.expansions {
		node, _ = expand(node, ctx)
	}
	node, _ = promotions(node, ctx)
	node, _ = constCompares(node, ctx)
	node, _ = constBranches(node)
	if report.expansions {
		// run compares before and after identifier expansion
		// some symbols -- Progressive_Hookshot -- are expandable but
		// occasionally used as symbols
		ctx.expandTokenLike = true
		defer func() { ctx.expandTokenLike = false }()
		node, _ = expand(node, ctx)
		node, _ = constCompares(node, ctx)
		node, _ = constBranches(node)
	}
	return node, nil
}

type AnalysisContext struct {
	expansions      map[string]replacement
	tokenNames      map[internal.NormalizedStr]struct{}
	settingName     map[internal.NormalizedStr]struct{}
	varName         map[internal.NormalizedStr]struct{}
	builtInName     map[internal.NormalizedStr]struct{}
	unpromoted      map[string]struct{}
	expandTokenLike bool
}

func (a *AnalysisContext) Unpromoted() []string {
	names := make([]string, 0, len(a.unpromoted))
	for k := range a.unpromoted {
		names = append(names, k)
	}

	return names
}

func (ctx *AnalysisContext) NameToken(name string) {
	ctx.tokenNames[internal.Normalize(name)] = struct{}{}
}

func (ctx *AnalysisContext) isToken(name string) bool {
	_, exists := ctx.tokenNames[internal.Normalize(name)]
	return exists
}

func (ctx *AnalysisContext) NameSetting(name string) {
	ctx.settingName[internal.Normalize(name)] = struct{}{}
}

func (ctx *AnalysisContext) isSetting(name string) bool {
	_, exists := ctx.settingName[internal.Normalize(name)]
	return exists
}

func (ctx *AnalysisContext) NameBuiltIn(name string) {
	ctx.builtInName[internal.Normalize(name)] = struct{}{}
}

func (ctx *AnalysisContext) isBuiltIn(name string) bool {
	_, exists := ctx.builtInName[internal.Normalize(name)]
	return exists
}

func (ctx *AnalysisContext) isVar(name string) bool {
	_, exists := ctx.varName[internal.Normalize(name)]
	return exists
}

func (ctx *AnalysisContext) isExpandable(name string) bool {
	if ctx.looksLikeToken(name) && !ctx.expandTokenLike {
		return false
	}
	_, exists := ctx.expansions[name]
	return exists
}

func (ctx *AnalysisContext) tagIdentifier(ident *ast.Identifier) {
	name := ident.Name
	if ident.Kind != ast.AST_IDENT_UNK {
		return
	} else if strings.HasPrefix(ident.Name, "logic_") {
		ident.Kind = ast.AST_IDENT_TRK
	} else if ctx.isToken(name) || ctx.looksLikeToken(name) {
		ident.Kind = ast.AST_IDENT_TOK
	} else if ctx.isSetting(ident.Name) {
		ident.Kind = ast.AST_IDENT_SET
	} else if ctx.isExpandable(name) {
		ident.Kind = ast.AST_IDENT_EXP
	} else if ctx.isBuiltIn(name) {
		ident.Kind = ast.AST_IDENT_BIF
	} else if ctx.isVar(name) {
		ident.Kind = ast.AST_IDENT_VAR
	} else {
		ident.Kind = ast.AST_IDENT_UNP
		ctx.unpromoted[name] = struct{}{}
	}
}

func (_ *AnalysisContext) looksLikeToken(str string) bool {
	return looksLikeToken.MatchString(str)
}

func (ctx *AnalysisContext) AddExpansion(name string, params []string, body ast.Node) {
	ctx.expansions[name] = replacement{name, params, body}
}

func chaseIdentifier(node ast.Node, ctx *AnalysisContext) (*ast.Identifier, bool) {
	ident := func(i *ast.Identifier) (*ast.Identifier, error) {
		return i, nil
	}
	literal :=
		func(l *ast.Literal) (*ast.Identifier, error) {
			str, isStr := l.Value.(string)
			if isStr && ctx.looksLikeToken(str) {
				i := new(ast.Identifier)
				i.Kind = ast.AST_IDENT_TOK
				i.Name = str
				return i, nil
			}
			return nil, errors.New("nope")
		}

	res, err := ast.Unify3(node,
		ident,
		literal,
		func(c *ast.Call) (*ast.Identifier, error) {
			if c.Callee != "has" {
				return nil, errors.New("nope")
			}

			return ast.Unify(c.Args[0], ident, literal)
		})
	return res, err == nil

}

func constCompares(node ast.Node, ctx *AnalysisContext) (ast.Node, error) {
	var re rewriter
	re.compares = func(trans *rewriter, compare *ast.Comparison) (ast.Node, error) {
		if compare.Op != ast.AST_CMP_LT {
			wantsSame := compare.Op == ast.AST_CMP_EQ
			lhsIdent, lhsIsIdent := chaseIdentifier(compare.LHS, ctx)
			rhsIdent, rhsIsIdent := chaseIdentifier(compare.RHS, ctx)
			if lhsIsIdent && rhsIsIdent {
				ctx.tagIdentifier(lhsIdent)
				ctx.tagIdentifier(rhsIdent)
				if lhsIdent.Kind == rhsIdent.Kind {
					areSame := internal.Normalize(lhsIdent.Name) == internal.Normalize(rhsIdent.Name)
					return ast.LiteralBool(wantsSame && areSame), nil
				}
			}
		}
		newOp := new(ast.Comparison)
		newOp.Op = compare.Op
		newOp.LHS, _ = ast.Transform(trans, compare.LHS)
		newOp.RHS, _ = ast.Transform(trans, compare.RHS)
		return newOp, nil
	}
	return ast.Transform(&re, node)
}

func assertIsBool(node ast.Node) (bool, bool) {
	lit, isLit := ast.AssertAs[*ast.Literal](node)
	if !isLit {
		return false, false
	}
	b, isBool := lit.Value.(bool)
	return b, isBool
}

func constBranches(node ast.Node) (ast.Node, error) {
	var re rewriter
	re.booleans = func(trans *rewriter, boolop *ast.BooleanOp) (ast.Node, error) {
		op := boolop.Op
		if op != ast.AST_BOOL_NEGATE {
			lhBool, lhIsBool := assertIsBool(boolop.LHS)
			rhBool, rhIsBool := assertIsBool(boolop.RHS)

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

			if rhIsBool {
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
			}
		}

		newOp := new(ast.BooleanOp)
		newOp.Op = boolop.Op
		newOp.LHS, _ = ast.Transform(trans, boolop.LHS)
		newOp.RHS, _ = ast.Transform(trans, boolop.RHS)
		return newOp, nil
	}
	return ast.Transform(&re, node)
}

func replaceIdentifierArgs(fragment ast.Node, repls map[string]ast.Node) ast.Node {
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
		expansion, has := ctx.expansions[call.Callee]
		if !has {
			return call, nil
		}
		replacements := map[string]ast.Node{}
		for idx, name := range expansion.Params {
			replacements[name] = call.Args[idx]
		}

		body, _ := copyfragment(expansion.Body)
		return replaceIdentifierArgs(body, replacements), nil
	}

	re.identifier = func(_ *rewriter, ident *ast.Identifier) (ast.Node, error) {
		if !ctx.isExpandable(ident.Name) {
			return ident, nil
		}
		expansion := ctx.expansions[ident.Name]

		if len(expansion.Params) != 0 {
			return nil, slipup.Createf("expected 0 parameter expansion for %s", ident.Name)
		}

		return copyfragment(expansion.Body)
	}
	expansion, _ := ast.Transform(&re, node)
	return Analyze(expansion, ctx)
}

func promotions(node ast.Node, ctx *AnalysisContext) (ast.Node, error) {
	var re rewriter
	re.call = func(r *rewriter, c *ast.Call) (ast.Node, error) {
		if c.Callee == "has" {
			lit, isLit := ast.AssertAs[*ast.Literal](c.Args[0])
			if isLit {
				c.Args[0] = &ast.Identifier{
					Name: lit.Value.(string),
					Kind: ast.AST_IDENT_TOK,
				}
			}
		}
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
	re.literal = func(_ *rewriter, lit *ast.Literal) (ast.Node, error) {
		str, isStr := lit.Value.(string)
		if isStr && ctx.isToken(str) || ctx.looksLikeToken(str) {
			return &ast.Identifier{
				Name: lit.Value.(string),
				Kind: ast.AST_IDENT_TOK,
			}, nil
		}
		return lit, nil
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

var looksLikeToken = regexp.MustCompile("[A-Z][a-zA-Z_ ]")
