package analysis

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"sudonters/zootler/icearrow/ast"
	"sudonters/zootler/internal"
	"sudonters/zootler/internal/components"
	"sudonters/zootler/internal/entities"
	"sudonters/zootler/internal/table"

	"github.com/etc-sudonters/substrate/slipup"
)

func NewAnalysis(edges entities.Edges) AnalysisContext {
	var ac AnalysisContext
	ac.edges = edges
	ac.lateExpansions = make(map[string][]LateExpansion)
	ac.expansions = make(map[string]replacement)
	ac.names = make(map[internal.NormalizedStr]ast.AstIdentifierKind)

	ac.names[internal.NormalizedStr("age")] = ast.AST_IDENT_VAR
	symbolic := []string{
		"adult", "both", "child", "either", "fire", "forest", "light",
		"shadow", "spirit", "water",
	}

	for _, symbol := range symbolic {
		ac.names[internal.Normalize(symbol)] = ast.AST_IDENT_SYM
	}

	return ac
}

func Analyze(node ast.Node, ctx *AnalysisContext) (ast.Node, error) {
	ctx.expandTokenLike = false
	report := analyze(node, ctx)
	if report.lateExpansions {
		node, _ = yankLateExpansions(node, ctx)
	}
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

type LateExpansion struct {
	Name, Parent string
	Rule         ast.Node
	Edge         entities.Edge
}

type AnalysisContext struct {
	names           map[internal.NormalizedStr]ast.AstIdentifierKind
	edges           entities.Edges
	current         string
	expansions      map[string]replacement
	lateExpansions  map[string][]LateExpansion
	expandTokenLike bool
}

func (a *AnalysisContext) SetCurrent(location string) {
	a.current = location
}

func (ctx *AnalysisContext) NameToken(name string) {
	ctx.names[internal.Normalize(name)] = ast.AST_IDENT_TOK
}

func (ctx *AnalysisContext) NameSetting(name string) {
	ctx.names[internal.Normalize(name)] = ast.AST_IDENT_SET
}

func (ctx *AnalysisContext) NameBuiltIn(name string) {
	ctx.names[internal.Normalize(name)] = ast.AST_IDENT_BIF
}

func (ctx *AnalysisContext) AddExpansion(name string, params []string, body ast.Node) {
	ctx.expansions[name] = replacement{name, params, body}
}

func (ctx *AnalysisContext) LateExpanders(yield func(string, LateExpansion) bool) {
	var lateXpns map[string][]LateExpansion
	for len(ctx.lateExpansions) > 0 {
		lateXpns, ctx.lateExpansions = ctx.lateExpansions, make(map[string][]LateExpansion)
		for current, xpns := range lateXpns {
			for _, xpn := range xpns {
				if !yield(current, xpn) {
					return
				}
			}
		}
	}
}

func (ctx *AnalysisContext) lateExpander(call *ast.Call) (ast.Node, error) {
	var target string
	var rule ast.Node

	switch call.Callee {
	case "here":
		target = ctx.current
		rule = call.Args[0]
		break
	case "at":
		trgt := ast.MustAssertAs[*ast.Literal](call.Args[0])
		target = trgt.Value.(string)
		rule = call.Args[1]
	default:
		panic(slipup.Createf("unknown late expansions %q", call.Callee))
	}

	if target == "" {
		panic(slipup.Createf(
			"at/here expansion w/o target available:\ncurrent: %q\n%#v",
			ctx.current, call,
		))
	}

	rules := ctx.lateExpansions[target]
	name := fmt.Sprintf("%s Reachability %d", target, len(rules)+1)
	edge, err := ctx.edges.Entity(components.Name(name + " edge"))
	if err != nil {
		return nil, slipup.Describef(err, "could not create token for subrule %q", name)
	}
	addErr := edge.AddComponents(table.Values{components.EventEdge{}, components.AnonymousEvent{}})
	// this edge and its destination are synthetic and we'll produce them later
	edge.Stash("origin", target)
	edge.Stash("dest", name)
	if addErr != nil {
		return nil, slipup.Describef(addErr, "could not describe token for subrule %q", name)
	}

	lateXpns := LateExpansion{
		Name:   name,
		Parent: target,
		Rule:   rule,
		Edge:   edge,
	}
	rules = append(rules, lateXpns)
	ctx.lateExpansions[target] = rules
	return &ast.Call{
		Callee: "has",
		Args: []ast.Node{
			&ast.Identifier{
				Name: name,
				Kind: ast.AST_IDENT_EVT,
			},
			&ast.Literal{Value: float64(1), Kind: ast.AST_LIT_NUM},
		},
	}, nil
}

func (ctx *AnalysisContext) tagIdentifier(ident *ast.Identifier) {
	name := ident.Name
	if ident.Kind != ast.AST_IDENT_UNK {
		return
	} else if kind, assigned := ctx.names[internal.Normalize(name)]; assigned {
		ident.Kind = kind
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
		panic(slipup.Createf("unknown identifier %q", ident.Name))
	}
}

func (ctx *AnalysisContext) isToken(name string) bool {
	ident, exists := ctx.names[internal.Normalize(name)]
	return exists && ident == ast.AST_IDENT_TOK
}

func (ctx *AnalysisContext) isSetting(name string) bool {
	ident, exists := ctx.names[internal.Normalize(name)]
	return exists && ident == ast.AST_IDENT_SET
}

func (ctx *AnalysisContext) isBuiltIn(name string) bool {
	ident, exists := ctx.names[internal.Normalize(name)]
	return exists && ident == ast.AST_IDENT_BIF
}

func (ctx *AnalysisContext) isVar(name string) bool {
	ident, exists := ctx.names[internal.Normalize(name)]
	return exists && ident == ast.AST_IDENT_VAR
}

func (ctx *AnalysisContext) isExpandable(name string) bool {
	if ctx.looksLikeToken(name) && !ctx.expandTokenLike {
		return false
	}
	_, exists := ctx.expansions[name]
	return exists
}

func (_ *AnalysisContext) looksLikeToken(str string) bool {
	return looksLikeToken.MatchString(str)
}

func chaseIdentifier(node ast.Node, ctx *AnalysisContext) (*ast.Identifier, bool) {
	ident := func(i *ast.Identifier) (*ast.Identifier, error) {
		return i, nil
	}
	literal := func(l *ast.Literal) (*ast.Identifier, error) {
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

func writeCompareToSetting(setting *ast.Identifier, comperand ast.Node) (*ast.Call, error) {
	return &ast.Call{
		Callee: "compare_to_setting",
		Args:   []ast.Node{setting, comperand},
	}, nil
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
		if c.Callee == "load_setting_2" {
			ident, isIdent := ast.AssertAs[*ast.Identifier](c.Args[0])
			if isIdent {
				switch ident.Name {
				case "dungeon_shortcuts":
					c.Callee = "has_dungeon_shortcuts"
					c.Args = c.Args[1:]
					break
				case "skipped_trials":
					c.Callee = "is_trial_skipped"
					c.Args = c.Args[1:]
					break
				}

			}
		}
		return c, nil
	}

	isSettingIdent := func(node ast.Node) (*ast.Identifier, bool) {
		ident, isIdent := ast.AssertAs[*ast.Identifier](node)
		return ident, isIdent && ident.Kind == ast.AST_IDENT_SET
	}

	isVarIdent := func(node ast.Node) (*ast.Identifier, bool) {
		ident, isIdent := ast.AssertAs[*ast.Identifier](node)
		return ident, isIdent && ident.Kind == ast.AST_IDENT_VAR
	}

	isSettingLoad := func(node ast.Node) (*ast.Call, bool) {
		call, isCall := ast.AssertAs[*ast.Call](node)
		return call, isCall && (call.Callee == "load_setting" || call.Callee == "load_setting_2" || call.Callee == "is_trial_skipped" || call.Callee == "has_dungeon_shortcuts")
	}

	re.compares = func(r *rewriter, compare *ast.Comparison) (ast.Node, error) {
		if ident, isIdent := isSettingIdent(compare.LHS); isIdent {
			return writeCompareToSetting(ident, compare.RHS)
		}
		if ident, isIdent := isSettingIdent(compare.RHS); isIdent {
			return writeCompareToSetting(ident, compare.LHS)
		}

		if ident, isIdent := isVarIdent(compare.LHS); isIdent {
			if ident.Name == "age" {
				return &ast.Call{
					Callee: "check_age",
					Args:   []ast.Node{compare.RHS},
				}, nil
			}
		}

		if ident, isIdent := isVarIdent(compare.RHS); isIdent {
			if ident.Name == "age" {
				return &ast.Call{
					Callee: "check_age",
					Args:   []ast.Node{compare.LHS},
				}, nil
			}
		}

		newOp := new(ast.Comparison)
		newOp.Op = compare.Op
		newOp.LHS, _ = ast.Transform(r, compare.LHS)
		newOp.RHS, _ = ast.Transform(r, compare.RHS)
		return newOp, nil

	}

	re.booleans = func(r *rewriter, boolop *ast.BooleanOp) (ast.Node, error) {
		lhs, _ := ast.Transform(r, boolop.LHS)
		if boolop.Op == ast.AST_BOOL_NEGATE {
			if setting, isSettingIdent := isSettingIdent(boolop.LHS); isSettingIdent {
				return &ast.Call{
					Callee: "invert_load_setting",
					Args:   []ast.Node{setting},
				}, nil
			}

			if call, isCall := isSettingLoad(boolop.LHS); isCall {
				return &ast.Call{
					Callee: fmt.Sprintf("invert_%s", call.Callee),
					Args:   call.Args,
				}, nil
			}
		}

		newOp := new(ast.BooleanOp)
		newOp.Op = boolop.Op
		newOp.LHS = lhs
		newOp.RHS, _ = ast.Transform(r, boolop.RHS)
		return newOp, nil
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

func yankLateExpansions(node ast.Node, ctx *AnalysisContext) (ast.Node, error) {
	var re rewriter
	re.call = func(_ *rewriter, call *ast.Call) (ast.Node, error) {
		if call.Callee != "at" && call.Callee != "here" {
			return call, nil
		}
		return ctx.lateExpander(call)
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
