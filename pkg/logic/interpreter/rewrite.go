package interpreter

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"sudonters/zootler/internal/astrender"
	"sudonters/zootler/pkg/logic"
	"sudonters/zootler/pkg/rules/ast"
	"sudonters/zootler/pkg/world"
	"sudonters/zootler/pkg/world/components"
)

var _entityName = regexp.MustCompile("^[A-Z][A-Za-z_]+$")

type eq int

const (
	_ eq = iota
	definitelyEq
	definitelyNotEq
	// we have to resolve this at runtime
	notSure
)

var _ Evaluation[ast.Expression] = Inliner{}

func NewInliner(globals Environment) *Inliner {
	return &Inliner{Globals: globals}
}

// does compile time execution to resolve and inline as many things as possible
// allows us to do stuff like cleave branches that are always false now
// this includes recursing into function calls and either replacing it with a constant
// or storing the partially executed function into the environment and replaces the general
// call to the optimized call
type Inliner struct {
	Globals          Environment
	Settings         map[string]any
	Tricks           map[string]bool
	SkippedTrials    map[string]bool
	DungeonShortcuts map[string]bool
	Builder          *world.Builder
	RegionName       string
}

func (rw Inliner) Rewrite(expr ast.Expression, env Environment) ast.Expression {
	return Evaluate(rw, expr, env)
}

func (rw Inliner) areEq(left, right ast.Expression, env Environment) eq {
	if left.Type() == ast.ExprCall && right.Type() == ast.ExprCall {
		lFn, rFn := left.(*ast.Call), right.(*ast.Call)

		if sameFunc := rw.areEq(lFn.Callee, rFn.Callee, env); sameFunc != definitelyEq {
			return sameFunc
		}

		for i := range lFn.Args {
			if sameArg := rw.areEq(lFn.Args[i], rFn.Args[i], env); sameArg != definitelyEq {
				return sameArg
			}
		}

		return definitelyEq
	}

	l, r := rw.resolveToValue(left, env), rw.resolveToValue(right, env)
	if l == nil || r == nil {
		return notSure
	}

	if l.Eq(r) {
		return definitelyEq
	}
	return definitelyNotEq
}

// expr is already rewritten
func (rw Inliner) resolveToValue(expr ast.Expression, env Environment) Value {
	switch expr := expr.(type) {
	case *ast.Literal:
		return ReifyLiteral(expr)
	case *ast.Identifier:
		if v, ok := rw.fromEnv(expr, env); ok {
			return v
		}
	}

	return nil
}

func IsTruthy(v Value) bool {
	switch v := v.(type) {
	case Boolean:
		return v.Value
	case Number:
		return v.Value != 0
	case String:
		return v.Value != ""
	default:
		panic(errors.New("unknown truthiness kind"))
	}
}

func (rw Inliner) fromEnv(i ast.Expression, env Environment) (Value, bool) {
	if i.Type() == ast.ExprIdentifier {
		return env.Get(i.(*ast.Identifier).Value)
	}

	return nil, false
}

func (rw Inliner) EvalLiteral(literal *ast.Literal, env Environment) ast.Expression {
	if literal.Kind == ast.LiteralStr {
		ident := &ast.Identifier{Value: literal.Value.(string)}
		_, ok := env.Get(ident.Value)
		if ok {
			return ident
		}

		typ := rw.Builder.TypedStrs.Typed(ident.Value)
		rw.Builder.Components.RowOf(typ)
		env.Set(ident.Value, Token{Component: typ, Literal: ident.Value})
		return ident
	}
	return literal
}

func (rw Inliner) EvalBinOp(op *ast.BinOp, env Environment) ast.Expression {
	left := rw.Rewrite(op.Left, env)
	right := rw.Rewrite(op.Right, env)

	switch op.Op {
	case ast.BinOpEq:
		if eq := rw.areEq(left, right, env); eq != notSure {
			return Literalify(eq == definitelyEq)
		}
		break
	case ast.BinOpNotEq:
		if eq := rw.areEq(left, right, env); eq != notSure {
			return Literalify(eq != definitelyEq)
		}
		break
	case ast.BinOpLt:
		r, ok := right.(*ast.Literal)
		if !ok || r.Kind != ast.LiteralNum {
			panic(parseError("cmp(<) only between numbers"))
		}

		var l float64

		switch left.Type() {

		case ast.ExprIdentifier:
			lv, ok := rw.fromEnv(left, env)
			if !ok {
				panic(parseError("expected %q to be available at compile time", left.(*ast.Identifier).Value))
			}

			if lv.Type() != NUM_TYPE {
				panic(parseError("cmp(<) only between numbers"))
			}

			l = lv.(Number).Value
			break
		case ast.ExprLiteral:
			lv := left.(*ast.Literal)
			if lv.Kind != ast.LiteralNum {
				panic(parseError("cmp(<) only between numbers"))
			}

			l = lv.Value.(float64)
		}

		return Literalify(l < r.Value.(float64))
	case ast.BinOpContains:
		if op.Left.Type() == ast.ExprLiteral {
			// subscript assumes identifiers only
			op.Left = &ast.Identifier{Value: op.Left.(*ast.Literal).Value.(string)}
		}
		return rw.Rewrite(&ast.Subscript{Target: op.Right, Index: op.Left}, env)
	}

	return &ast.BinOp{
		Left:  left,
		Op:    op.Op,
		Right: right,
	}
}

func (rw Inliner) isFuncPointer(expr ast.Expression, env Environment) (*ast.Identifier, bool) {
	v, ok := rw.fromEnv(expr, env)

	if !ok {
		return nil, false
	}

	if v.Type() != CALL_TYPE {
		return nil, false
	}

	return expr.(*ast.Identifier), true
}

func (rw Inliner) Make0ArityFnCall(expr ast.Expression, env Environment) (ast.Expression, bool) {
	if ident, ok := rw.isFuncPointer(expr, env); ok {
		fn, _ := rw.fromEnv(ident, env)
		if fn.(Callable).Arity() != 0 {
			panic(errors.New("must be 0 arity func"))
		}

		return &ast.Call{
			Callee: ident,
			Args:   nil,
		}, true
	}

	return expr, false
}

func (rw Inliner) EvalBoolOp(op *ast.BoolOp, env Environment) ast.Expression {
	left := rw.Rewrite(op.Left, env)
	if call, ok := rw.Make0ArityFnCall(left, env); ok {
		left = rw.Rewrite(call, env)
	}

	if literal, ok := left.(*ast.Literal); ok && literal.Kind == ast.LiteralBool {
		l := ReifyLiteral(literal).(Boolean).Value

		switch op.Op {
		case ast.BoolOpOr:
			if l {
				return Literalify(l)
			}
			return rw.Rewrite(op.Right, env)
		case ast.BoolOpAnd:
			if !l {
				return Literalify(l)
			}
			return rw.Rewrite(op.Right, env)
		}
	}

	right := rw.Rewrite(op.Right, env)
	if call, ok := rw.Make0ArityFnCall(right, env); ok {
		right = rw.Rewrite(call, env)
	}
	if literal, ok := right.(*ast.Literal); ok && literal.Kind == ast.LiteralBool {
		r := ReifyLiteral(literal).(Boolean).Value
		switch op.Op {
		case ast.BoolOpOr:
			if !r {
				return left
			}
			return Literalify(r)
		case ast.BoolOpAnd:
			if r {
				return left
			}
			return Literalify(r)
		}
	}

	return &ast.BoolOp{
		Left:  left,
		Op:    op.Op,
		Right: right,
	}
}

func (rw Inliner) EvalUnary(unary *ast.UnaryOp, env Environment) ast.Expression {
	target := rw.Rewrite(unary.Target, env)
	switch unary.Op {
	case ast.UnaryNot:
		if target.Type() == ast.ExprLiteral {
			b := target.(*ast.Literal)
			if b.Kind != ast.LiteralBool {
				panic(parseError("can only negate literal bools"))
			}
			return Literalify(!b.Value.(bool))
		}
		t, _ := rw.fromEnv(target, env)
		if t != nil {
			switch t := t.(type) {
			case Boolean:
				return Literalify(!t.Value)
			case Callable:
				break
			default:
				panic(parseError("can only negate literal bools"))
			}
		}
	default:
		panic(parseError("unknown unary op: %q", unary.Op))
	}

	return &ast.UnaryOp{
		Op:     unary.Op,
		Target: target,
	}
}

func (rw Inliner) EvalCall(call *ast.Call, env Environment) ast.Expression {
	if ident, ok := call.Callee.(*ast.Identifier); ok && (ident.Value == "here" || ident.Value == "at") {
		var name string
		var body ast.Expression

		switch ident.Value {
		case "here":
			name = rw.RegionName
			body = call.Args[0]
		case "at":
			name = call.Args[0].(*ast.Literal).Value.(string)
			body = call.Args[1]
		}

		return rw.expandMacro(name, body, env)
	}

	newCall := new(ast.Call)
	newCall.Callee = rw.Rewrite(call.Callee, env)
	newCall.Args = make([]ast.Expression, len(call.Args))
	for i := range newCall.Args {
		newCall.Args[i] = rw.Rewrite(call.Args[i], env)
	}

	v, ok := rw.fromEnv(newCall.Callee, env)
	if !ok || v.Type() != CALL_TYPE {
		return newCall
	}

	fn, ok := v.(Fn) // specifically
	if !ok {
		return newCall
	}

	if fn.Arity() != len(newCall.Args) {
		panic(parseError("mismatch arg count: wanted %d but got %d", fn.Arity(), len(newCall.Args)))
	}

	enclosed := env.Enclosed()
	for i, arg := range newCall.Args {
		a := rw.resolveToValue(arg, env)
		if a == nil {
			return newCall
		}

		enclosed.Set(fn.Params[i], a)
	}

	body := rw.Rewrite(fn.Body, enclosed)
	switch body.(type) {
	case *ast.Literal:
		return body
	default:
		addr := contentAddress(body)
		newName := fmt.Sprintf("%s@%s", fn.Name.Value, addr)
		partialFn := PartiallyEvaluatedFn{
			Body: body,
			Env:  enclosed,
			Name: newName,
		}
		env.Set(newName, partialFn)
		return &ast.Call{
			Callee: &ast.Identifier{Value: newName},
		}
	}
}

func (rw Inliner) identInEnv(expr ast.Expression, env Environment) bool {
	ident, ok := expr.(*ast.Identifier)
	if !ok {
		return false
	}

	_, ok = env.Get(ident.Value)
	return ok
}

func (rw Inliner) expandMacro(where string, rule ast.Expression, env Environment) *ast.Identifier {
	addr := contentAddress(rule)
	eventName := fmt.Sprintf("%s@%s", where, addr)

	// creating and pushing an entity to the graph are idempotent
	// because we might be referencing a place that doesn't exist yet
	origin, err := rw.Builder.Entity(components.Name(where))
	if err != nil {
		panic(err)
	}
	rw.Builder.Node(origin)

	event, err := rw.Builder.Entity(components.Name(eventName))
	if err != nil {
		panic(err)
	}

	rw.Builder.Node(event)

	arch := components.EventArchetype{}
	if err := arch.Apply(event); err != nil {
		panic(err)
	}

	event.Add(rw.Builder.TypedStrs.InstanceOf(eventName))
	rw.Globals.Set(eventName, Token{
		Literal:   eventName,
		Component: rw.Builder.TypedStrs.Typed(eventName),
	})

	edge, err := rw.Builder.Edge(origin, event)
	if err != nil {
		panic(err)
	}

	rule = rw.Rewrite(rule, env)
	edge.Add(rule)
	return &ast.Identifier{Value: eventName}
}

func contentAddress(expr ast.Expression) string {
	s := astrender.NewSexpr(astrender.DontTheme())
	ast.Visit(s, expr)
	hash := sha256.New()
	hash.Write([]byte(s.String()))
	return fmt.Sprintf("sha256:%x", hash.Sum(nil))
}

func (rw Inliner) EvalIdentifier(ident *ast.Identifier, env Environment) ast.Expression {
	if v, ok := rw.fromEnv(ident, env); ok {
		if CanLiteralfy(v) {
			return Literalify(v)
		}
		return ident
	}

	name := ident.Value
	if setting, ok := rw.Settings[name]; ok {
		rw.Globals.Set(name, Box(setting))
		return Literalify(setting)
	}

	if strings.HasPrefix(name, "logic_") {
		v := rw.Tricks[strings.TrimPrefix(name, "logic_")]
		rw.Globals.Set(name, Box(v))
		return Literalify(v)
	}

	if _entityName.MatchString(name) {
		entity, err := rw.Builder.Entity(components.Name(name))
		if err != nil {
			panic(err)
		}

		entity.Add(rw.Builder.TypedStrs.InstanceOf(logic.EscapeName(name)))
		rw.Globals.Set(name, Token{
			Literal:   name,
			Component: rw.Builder.TypedStrs.Typed(logic.EscapeName(name)),
		})
	}

	return ident
}

// lowers to a boolean from a passed settings dict
func (rw Inliner) EvalSubscript(subscript *ast.Subscript, env Environment) ast.Expression {
	if subscript.Target.Type() != ast.ExprIdentifier || subscript.Index.Type() != ast.ExprIdentifier {
		panic("subscript only with identifiers")
	}

	settings := subscript.Target.(*ast.Identifier)
	value := subscript.Index.(*ast.Identifier).Value

	switch settings.Value {
	case "tricks":
		return Literalify(rw.Tricks[value])
	case "skipped_trials":
		return Literalify(rw.SkippedTrials[value])
	case "settings":
		return Literalify(rw.Settings[value])
	case "dungeon_shortcuts":
		return Literalify(rw.DungeonShortcuts[value])
	default:
		panic(parseError("unknown subscript target %s[%s]", settings.Value, value))
	}
}

func (rw Inliner) EvalTuple(tup *ast.Tuple, env Environment) ast.Expression {
	if len(tup.Elems) != 2 {
		panic(BadTupleErr)
	}

	item := tup.Elems[0]
	want := tup.Elems[1]
	var ident *ast.Identifier
	var qty float64

	switch item := item.(type) {
	case *ast.Identifier:
		ident = item
	case *ast.Literal:
		panic(BadTupleErr)
	default:
		panic(BadTupleErr)
	}
	switch want := want.(type) {
	case *ast.Identifier:
		v := rw.resolveToValue(want, env)
		qty = v.(Number).Value
		break
	case *ast.Literal:
		if want.Kind == ast.LiteralNum {
			qty = want.Value.(float64)
			break
		}
		panic(BadTupleErr)
	default:
		panic(BadTupleErr)
	}

	// ensure that we create the entity in the global env
	rw.EvalIdentifier(ident, env)

	return &ast.Call{
		Callee: &ast.Identifier{Value: "has"},
		Args:   []ast.Expression{ident, Literalify(qty)},
	}
}

var BadTupleErr = errors.New("tuple must be (Ident, Number)")
