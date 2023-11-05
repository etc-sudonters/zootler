package interpreter

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"strings"
	"sudonters/zootler/internal/astrender"
	"sudonters/zootler/internal/entity"
	"sudonters/zootler/pkg/rules/ast"
	"sudonters/zootler/pkg/worldloader"

	"github.com/etc-sudonters/substrate/stageleft"
)

/*

Notes:

1. convert item/event references into entity.Model
	("item name", 2) -> (entity.Model(69), 2) -> Zoot_HasQuantityOf{ entity.Model(69), 2 }
2. inline functions and settings
	// NOTE at_day is a builtin
	is_child and at_day and (can_break_crate or chicken_count < 7)
	// descend leftest branch first
	age == 'child' and at_day and ((can_bonk or can_blast_or_smash) or (chicken_count < 7))
	// dead_bonks is a setting, the default is 'none'
	age == 'child' and at_day and (((deadly_bonks != 'ohko' or Fairy or can_use(Nayrus_Love)) or (has_explosives or can_use(Megaton_Hammer))) or (chicken_count < 7))
	// 'none' != 'ohko' is something we can determine right now
	age == 'child' and at_day and ((('none' != 'ohko' or Fairy or can_use(Nayrus_Love)) or (has_explosives or can_use(Megaton_Hammer))) or (chicken_count < 7))
	// or short circuits on the first true
	age == 'child' and at_day and (((true or Fairy or can_use(Nayrus_Love)) or (has_explosives or can_use(Megaton_Hammer))) or (chicken_count < 7))
	// reduce
	age == 'child' and at_day and (true or (chicken_count < 7))
	age == 'child' and at_day and true
	age == 'child' and at_day

3. at and here need to be expanded in this pass, these are metarules, essentially a macro that creates things
	they both work similiar, except at is remote and here is the current region
	creates an event at the specified REGION with the specified rule, THIS rule now depends on THAT event:

	at('SFM Entryway', Scarecrow)
		* locate this entity
		* does this event rule exist already? cool, let's just use that
		* IMPL: sha256(stringify(expr)) so it's content addressable, why not bring other things into this :shrug:
		* create a new event with a generated name -- 'SFM Entryway Subrule 1'
		* create edge 'SFM Entryway -> SFM Entryway Subrule 1'
		* rewrite the rule as much as possible
		* attach rewritten rule + identifier to SFM Entryway
		* Spit out a HasQuanityOf{ "SFM Entryway Subrule 1", 1 }

	here(Scarecrow) -> at(???, Scarecrow)
		okay but how do we know where "HERE" is? easy, we're here.
		who's on first, what's on second and iunno's on third
		the rewriter is being run region by region so it knows exactly where its at
		here(has_bottle) -> at('ZR Fairy Grotto', has_bottle) <- could fall into a at_fast(rw.region, has_bottle) since we have the region
*/

func Rewrite(old ast.Expression, env Environment) ast.Expression {
	r := Rewriter{}
	return r.Rewrite(old, env)
}

func NewRewriter(tricks map[string]bool) *Rewriter {
	return &Rewriter{
		tricks:           tricks,
		skippedTrials:    make(map[string]bool),
		dungeonShortcuts: make(map[string]bool),
	}
}

type Rewriter struct {
	testingEntityId  int
	tricks           map[string]bool
	skippedTrials    map[string]bool
	dungeonShortcuts map[string]bool
	region           string
}

func (r *Rewriter) SetRegion(region string) {
	r.region = region
}

func (r *Rewriter) Rewrite(old ast.Expression, env Environment) ast.Expression {
	return Evaluate(r, old, env)
}

func (r *Rewriter) EvalAttrAccess(access *ast.AttrAccess, env Environment) ast.Expression {
	panic(stageleft.NotImplErr)
}

func (r *Rewriter) EvalBinOp(op *ast.BinOp, env Environment) ast.Expression {
	left := r.Rewrite(op.Left, env)
	right := r.Rewrite(op.Right, env)

	switch op.Op {
	case ast.BinOpEq:
		if left.Type() == right.Type() {
			switch left.Type() {
			case ast.ExprString:
				l := left.(*ast.String).Value
				r := right.(*ast.String).Value
				eq := l == r
				return Literalify(eq)
			case ast.ExprBoolean:
				l := left.(*ast.Boolean).Value
				r := right.(*ast.Boolean).Value
				eq := l == r
				return Literalify(eq)
			case ast.ExprNumber:
				l := left.(*ast.Number).Value
				r := right.(*ast.Number).Value
				eq := l == r
				return Literalify(eq)
			case ast.ExprIdentifier:
				l := left.(*ast.Identifier).Value
				r := right.(*ast.Identifier).Value
				eq := l == r
				if eq { // can't know now if two _different_ identifiers are eq
					return Literalify(eq)
				}
			}
		}
		break
	case ast.BinOpNotEq:
		if left.Type() == right.Type() {
			switch left.Type() {
			case ast.ExprString:
				l := left.(*ast.String).Value
				r := right.(*ast.String).Value
				eq := l == r
				return Literalify(!eq)
			case ast.ExprBoolean:
				l := left.(*ast.Boolean).Value
				r := right.(*ast.Boolean).Value
				eq := l == r
				return Literalify(!eq)
			case ast.ExprNumber:
				l := left.(*ast.Number).Value
				r := right.(*ast.Number).Value
				eq := l == r
				return Literalify(!eq)
			case ast.ExprIdentifier:
				l := left.(*ast.Identifier).Value
				r := right.(*ast.Identifier).Value
				eq := l == r
				if eq { // can't know now if two _different_ identifiers are eq
					return Literalify(!eq)
				}
			}
		}
		break
	case ast.BinOpLt:
		if left.Type() == right.Type() && left.Type() == ast.ExprNumber {
			return Literalify(left.(*ast.Number).Value != right.(*ast.Number).Value)
		}
		panic(parseError("comparisons can only happen between numbers, cannot %T < %T", left, right))
	case ast.BinOpContains:
		if left.Type() == ast.ExprString && right.Type() == ast.ExprIdentifier {
			sub := &ast.Subscript{
				Target: right,
				Index:  left,
			}
			return r.Rewrite(sub, env)
		}
		panic(parseError("invalid contains"))
	default:
		panic(parseError("unknown binary operator %q", op.Op))
	}

	return &ast.BinOp{
		Left:  left,
		Op:    op.Op,
		Right: right,
	}
}

func (r *Rewriter) EvalBoolOp(op *ast.BoolOp, env Environment) ast.Expression {
	left := r.Rewrite(op.Left, env)

	// discard left side of bool wrapper if possible
	if b, ok := left.(*ast.Boolean); ok {
		if (b.Value && op.Op == ast.BoolOpOr) || (!b.Value && op.Op == ast.BoolOpAnd) {
			return b
		}
		return r.Rewrite(op.Right, env)
	}

	right := r.Rewrite(op.Right, env)

	// discard right side of bool wrapper if possible
	if b, ok := right.(*ast.Boolean); ok {
		// discard partially evaulated left side
		if (b.Value && op.Op == ast.BoolOpOr) || (!b.Value && op.Op == ast.BoolOpAnd) {
			return b
		}

		// discard completely evaluated right side
		if (!b.Value && op.Op == ast.BoolOpOr) || (b.Value && op.Op == ast.BoolOpAnd) {
			return left
		}
	}

	newOp := &ast.BoolOp{
		Left:  left,
		Right: right,
		Op:    op.Op,
	}

	return newOp
}

func (r *Rewriter) EvalBoolean(bool *ast.Boolean) ast.Expression {
	return bool
}

func isMacro(call *ast.Call) bool {
	ident, ok := call.Callee.(*ast.Identifier)
	return ok && ident.Value == "at" || ident.Value == "here"
}

func (r Rewriter) doMacro(call *ast.Call, env Environment) ast.Expression {
	ident := call.Callee.(*ast.Identifier)
	var where string
	var rule ast.Expression

	if ident.Value == "at" {
		if len(call.Args) != 2 {
			panic(parseError("invalid at() arguments"))
		}

		if region, ok := call.Args[0].(*ast.String); !ok {
			panic(parseError("invalid at() arguments"))
		} else {
			where = region.Value
		}

		rule = call.Args[1]
	} else {
		if len(call.Args) != 1 {
			panic(parseError("invalid here() arguments"))
		}
		where = r.region
		rule = call.Args[0]
	}

	// hypothetically this can recurse
	return r.testingRunAt(where, r.Rewrite(rule, env), env)
}

func (r *Rewriter) EvalCall(call *ast.Call, env Environment) ast.Expression {
	// these are special, think macros
	if isMacro(call) {
		return r.doMacro(call, env)
	}

	c := &ast.Call{}
	c.Callee = r.Rewrite(call.Callee, env)
	c.Args = make([]ast.Expression, len(call.Args))
	for i := range c.Args {
		c.Args[i] = r.Rewrite(call.Args[i], env)
	}

	if fn, ok := canOptimize(c, env); ok {
		return r.doOptimize(fn, c, env)
	}

	return c
}

// at best we can compile time evaluate a function to a constant
// if we can't, we replace the function with our partially evaluated one
func (r *Rewriter) doOptimize(fn Fn, c *ast.Call, env Environment) ast.Expression {
	enclosed := env.Enclosed()
	for i := range fn.Params {
		arg := c.Args[i]
		switch {
		case CanReifyLiteral(arg):
			enclosed.Set(fn.Params[i], ReifyLiteral(arg))
		case arg.Type() == ast.ExprIdentifier:
			v, _ := env.Get((arg.(*ast.Identifier)).Value)
			enclosed.Set(fn.Params[i], v)
		}
	}

	body := r.Rewrite(fn.Body, enclosed)

	if CanReifyLiteral(body) || body.Type() == ast.ExprIdentifier {
		return body
	}

	s := astrender.NewSexpr(astrender.DontTheme())
	ast.Visit(s, body)
	hash := sha256.New()
	hash.Write([]byte(fn.Name))
	hash.Write([]byte(s.String()))
	name := fmt.Sprintf("%s:sha256:%x", fn.Name, hash.Sum(nil))
	partialFn := PartiallyEvaluatedFn{
		Name: name,
		Body: body,
		Env:  enclosed,
	}

	env.Set(name, partialFn)
	return &ast.Call{
		Callee: &ast.Identifier{Value: name},
	}

}

// if all the args compile evaluate to constants, then we can at least
// partially evaluate the function body at compile time
func canOptimize(c *ast.Call, env Environment) (Fn, bool) {
	var defaultFn Fn
	tryOptimize := true
	for i := range c.Args {
		tryOptimize = tryOptimize && CanReifyLiteral(c.Args[i]) || isIdentInCompileEnv(c.Args[i], env)
		if !tryOptimize {
			return defaultFn, false
		}
	}

	if c.Callee.Type() != ast.ExprIdentifier {
		return defaultFn, false
	}
	ident := c.Callee.(*ast.Identifier)
	v, ok := env.Get(ident.Value)
	if !ok {
		return defaultFn, false
	}
	fn, ok := v.(Fn) // specifically, not any callable
	if !ok {
		return defaultFn, false
	}
	return fn, true
}

func isIdentInCompileEnv(expr ast.Expression, env Environment) bool {
	ident, ok := expr.(*ast.Identifier)
	if !ok {
		return false
	}

	_, ok = env.Get(ident.Value)
	return ok
}

func (r *Rewriter) EvalIdentifier(ident *ast.Identifier, env Environment) ast.Expression {
	/*
		Zootr's Priority:
			1. at/here/TOD
			2. things defined in helpers.json
			3. escaped_item? HasQuanityOf{ Ident, 1 }
			4. names off World.World in World.py
			5. World Settings
			6. funcs off State -- heart count, can live dmg, etc
			7. maybe it's an event? if so, make event item and replace w/ HasQuanityOf{ EventName, 1 }
			8. Throw

		What we'll do:
			1. at/here -> panic, these are special special
			2. if it exists in the env
				a. and it's a Int|Str|Bool: just return that
				b. just return the identifier
			3. resolve "static" values: settings mostly
				* NOTE zootr runs ast.parse on the values pulled from the world and settings
			4. try unaliasing to entity.Model w/ component.Token, if so stick an Entity in the env and return the ident
			5. check if we can eventify it, if so do and stuff that into the env, return the ident
			6. panic? or make it a runtime problem?
	*/

	// settings are already loaded into the environment
	if strings.HasPrefix(ident.Value, "logic_") {
		return Literalify(r.resolveTrick(ident.Value))
	}

	if v, ok := env.Get(ident.Value); ok && CanLiteralfy(v) {
		return Literalify(v)
	}

	if literal, _ := r.tryLiteralFromEnv(ident, env); literal != nil {
		return literal
	}

	return ident
}

func (r *Rewriter) resolveTrick(name string) bool {
	name = strings.TrimSuffix(name, "logic_")
	return r.tricks[name] // if we don't fine it it's not on :brain:
}

// try to embed a value as a literal for evaluation now
// if we can, we return an expression and true
func (r *Rewriter) tryLiteralFromEnv(name *ast.Identifier, env Environment) (expr ast.Expression, present bool) {
	var v Value
	v, present = env.Get(name.Value)
	if !present {
		v, present = env.Get(worldloader.EscapeName(name.Value))
		if !present {
			return
		}
	}

	if !CanLiteralfy(v) {
		return
	}

	expr = Literalify(v)
	return
}

func (r *Rewriter) EvalNumber(num *ast.Number) ast.Expression {
	return num
}

func (r *Rewriter) EvalString(str *ast.String) ast.Expression {
	return str
}

func (r *Rewriter) hasAtLeastOneOf(name string, env Environment) ast.Expression {
	/*
		some strings are actual literals for has(name, 1)

		```zoot
		Goron_Tunic: 'Goron_Tunic' or can_buy_goron_tunic
		```

		When we process this, we don't recurse into `Goron_Tunic` instead
		`'Goron_Tunic'` is replaced with `has(entity.Model, 1)` the
		`entity.Model` is either retrieved from the world's pool or a new event
		is created in the pool
	*/

	return &ast.Call{
		Callee: r.Rewrite(&ast.Identifier{Value: "has"}, env),
		Args:   []ast.Expression{r.resolveEntity(name, env), Literalify(1)},
	}
}

func (r *Rewriter) EvalSubscript(subscript *ast.Subscript, env Environment) ast.Expression {
	/*

		```zoot
		Ganons Castle Tower: |
			(skipped_trials[Forest] or 'Forest Trial Clear') and
			(skipped_trials[Fire] or 'Fire Trial Clear') and
			(skipped_trials[Water] or 'Water Trial Clear') and
			(skipped_trials[Shadow] or 'Shadow Trial Clear') and
			(skipped_trials[Spirit] or 'Spirit Trial Clear') and
			(skipped_trials[Light] or 'Light Trial Clear')
		```

		skipped_trials is a setting that's map[string]bool we'll convert this
		into an identifier `skipped_trials_Forest` and insert the boolean into
		the environment
	*/

	target, ok := subscript.Target.(*ast.Identifier)
	if !ok {
		panic(parseError("subscription can only be done with identifiers"))
	}

	key, ok := subscript.Index.(*ast.Identifier)
	if !ok {
		if raw, ok := subscript.Index.(*ast.String); ok {
			key = &ast.Identifier{Value: raw.Value}
		} else {
			panic(parseError("subscription can only be done with identifiers or strings"))
		}
	}

	setting, ok := r.resolveSetting(target.Value).(map[string]bool)
	if !ok { // surely this is always true
		panic(parseError("subscription %q did not resolve to a map", target.Value))
	}

	return Literalify(setting[key.Value])
}

func (r *Rewriter) EvalTuple(tup *ast.Tuple, env Environment) ast.Expression {
	/*
		arity 2
		item = _.0 <- ident of item
		count = _.1 <- number or settings ref
	*/

	s := astrender.NewSexpr(astrender.DontTheme())
	ast.Visit(s, tup)

	if len(tup.Elems) != 2 {
		panic(parseError("tup must be 2 elements, got: %d\n%s", len(tup.Elems), s.String()))
	}

	var token *ast.Identifier
	var qty *ast.Number

	{
		t := r.Rewrite(tup.Elems[0], env)
		q := r.Rewrite(tup.Elems[1], env)
		var ok bool

		if token, ok = t.(*ast.Identifier); !ok {
			panic(parseError("tup must be (Ident, Num), got (%T, %T)\n%s", t, q, s.String()))
		}

		if qty, ok = q.(*ast.Number); !ok {
			// we'll let runtime worry about this one
			if ident, ok := q.(*ast.Identifier); ok {
				return &ast.Call{
					Callee: &ast.Identifier{Value: "has"},
					Args:   []ast.Expression{token, ident},
				}
			}

			panic(parseError("tup must be (Ident, Num), got (%T, %T)\n%s", t, q, s.String()))
		}
	}

	return &ast.Call{
		Callee: &ast.Identifier{Value: "has"},
		Args:   []ast.Expression{token, qty},
	}
}

func (r *Rewriter) EvalUnary(unary *ast.UnaryOp, env Environment) ast.Expression {
	switch unary.Op {
	case ast.UnaryNot:
		res := r.Rewrite(unary.Target, env)
		if b, ok := res.(*ast.Boolean); ok {
			return &ast.Boolean{Value: !b.Value}
		}
		return &ast.UnaryOp{
			Op:     unary.Op,
			Target: res,
		}
	default:
		panic(parseError("unknown unary operator: %q", unary.Op))
	}
}

func (r *Rewriter) testingRunAt(where string, rule ast.Expression, env Environment) ast.Expression {
	// hash the raw unrewritten rule first
	s := astrender.NewSexpr(astrender.DontTheme())
	ast.Visit(s, rule)
	h := sha256.New()
	h.Write([]byte(s.String()))
	name := fmt.Sprintf("%q:sha256:%x", where, h.Sum(nil))
	ident := &ast.Identifier{Value: name}

	if _, exists := env.Get(name); exists {
		return ident
	}

	r.testingEntityId++
	env.Set(name, Entity{entity.Model(r.testingEntityId)})

	return ident
}

func (r *Rewriter) resolveSetting(block string) any {
	switch block {
	case "skipped_trials":
		return r.skippedTrials
	case "dungeon_shortcuts":
		return r.dungeonShortcuts
	default:
		panic(parseError("setting block %q doesn't exist", block))
	}
}

func (r *Rewriter) resolveEntity(name string, env Environment) *ast.Identifier {
	/*
		escape name
		if it exists in the env, just return the identifier now
		look up name in entity pool
		no entity? :think: zootr makes an event and inserts it
		put escaped named + entity into environment
		return escaped name
	*/
	name = worldloader.EscapeName(name)
	ident := &ast.Identifier{Value: name}

	ent, ok := env.Get(name)
	if ok {
		if _, ok = ent.(Entity); ok {
			return ident
		} else {
			panic(fmt.Errorf("expected to resolve entity from %q but resolved %T", name, ent))
		}
	}

	if testingIsEntity(name) {
		r.testingEntityId++
		env.Set(name, Box(entity.Model(r.testingEntityId)))
		return ident
	}

	return &ast.Identifier{Value: name}
}

func testingIsEntity(name string) bool {
	return strings.Contains(name, "_")
}

func parseError(reason string, v ...any) error {
	reason = fmt.Sprintf(reason, v...)
	return fmt.Errorf("%w: %s", parseErr, reason)
}

var parseErr = errors.New("parse error")
