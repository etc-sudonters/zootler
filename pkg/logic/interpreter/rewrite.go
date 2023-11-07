package interpreter

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"strings"
	"sudonters/zootler/internal/astrender"
	"sudonters/zootler/pkg/logic"
	"sudonters/zootler/pkg/rules/ast"
	"sudonters/zootler/pkg/world"
	"sudonters/zootler/pkg/worldloader"
)

type eq int

const (
	_ eq = iota
	definitelyEq
	definitelyNotEq
	// we have to resolve this at runtime
	notSure
)

var _ Evaluation[ast.Expression] = Rewriter{}

func NewRewriter() *Rewriter {
	return &Rewriter{}
}

type Rewriter struct {
	Settings      map[string]any
	Tricks        map[string]bool
	SkippedTrials map[string]bool
	Builder       *world.Builder
	regionName    string
}

func (rw *Rewriter) SetRegion(r string) {
	rw.regionName = r
}

func (rw Rewriter) Rewrite(expr ast.Expression, env Environment) ast.Expression {
	return Evaluate(rw, expr, env)
}

func (rw Rewriter) areEq(left, right ast.Expression, env Environment) eq {
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
func (rw Rewriter) resolveToValue(expr ast.Expression, env Environment) Value {
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

func (rw Rewriter) fromEnv(i ast.Expression, env Environment) (Value, bool) {
	if i.Type() == ast.ExprIdentifier {
		return env.Get(i.(*ast.Identifier).Value)
	}

	return nil, false
}

func (rw Rewriter) entityLiteral(lit string, env Environment) (*ast.Identifier, bool) {
	escaped := worldloader.EscapeName(lit)
	v, ok := env.Get(escaped)
	if !ok {
		return nil, false
	}

	if v.Type() != ENT_TYPE {
		return nil, false
	}

	return &ast.Identifier{Value: escaped}, true

}

func (rw Rewriter) EvalLiteral(literal *ast.Literal, _ Environment) ast.Expression {
	return literal
}

func (rw Rewriter) EvalBinOp(op *ast.BinOp, env Environment) ast.Expression {
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
		return rw.Rewrite(&ast.Subscript{Target: op.Right, Index: op.Left}, env)
	}

	return &ast.BinOp{
		Left:  left,
		Op:    op.Op,
		Right: right,
	}
}

func (rw Rewriter) isFuncPointer(expr ast.Expression, env Environment) (*ast.Identifier, bool) {
	v := rw.resolveToValue(expr, env)
	if v == nil {
		return nil, false
	}

	if v.Type() != CALL_TYPE {
		return nil, false
	}

	switch v := v.(type) {
	case Fn:
		return v.Name, true
	case PartiallyEvaluatedFn:
		return &ast.Identifier{Value: v.Name}, false
	default:
		panic(errors.New("missed a callable type"))
	}
}

func (rw Rewriter) make0ArityFnCall(expr ast.Expression, env Environment) (ast.Expression, bool) {
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

func (rw Rewriter) EvalBoolOp(op *ast.BoolOp, env Environment) ast.Expression {
	left := rw.Rewrite(op.Left, env)
	if call, ok := rw.make0ArityFnCall(left, env); ok {
		left = call
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
	if call, ok := rw.make0ArityFnCall(right, env); ok {
		right = call
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

func (rw Rewriter) EvalUnary(unary *ast.UnaryOp, env Environment) ast.Expression {
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
			case Callable, Entity:
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

func (rw Rewriter) EvalCall(call *ast.Call, env Environment) ast.Expression {
	if ident, ok := call.Callee.(*ast.Identifier); ok && (ident.Value == "here" || ident.Value == "at") {
		var name string
		var body ast.Expression

		switch ident.Value {
		case "here":
			name = rw.regionName
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
	case *ast.Literal, *ast.Identifier:
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
		return &ast.Identifier{Value: newName}
	}
}

func (rw Rewriter) identInEnv(expr ast.Expression, env Environment) bool {
	ident, ok := expr.(*ast.Identifier)
	if !ok {
		return false
	}

	_, ok = env.Get(ident.Value)
	return ok
}

func (rw Rewriter) expandMacro(where string, rule ast.Expression, env Environment) *ast.Identifier {
	addr := contentAddress(rule)
	eventName := fmt.Sprintf("%s@%s", where, addr)

	// creating and pushing an entity to the graph are idempotent
	// because we might be referencing a place that doesn't exist yet
	origin, err := rw.Builder.Entity(world.Name(where))
	if err != nil {
		panic(err)
	}
	rw.Builder.Node(origin)

	event, err := rw.Builder.Entity(world.Name(eventName))
	if err != nil {
		panic(err)
	}
	rw.Builder.Node(event)
	event.Add(logic.Token{})
	event.Add(logic.Event{})

	edge, err := rw.Builder.Edge(origin, event)
	if err != nil {
		panic(err)
	}

	rule = rw.Rewrite(rule, env)
	edge.Add(rule)

	env.Set(eventName, Box(event.Model()))

	return &ast.Identifier{Value: eventName}
}

func contentAddress(expr ast.Expression) string {
	s := astrender.NewSexpr(astrender.DontTheme())
	ast.Visit(s, expr)
	hash := sha256.New()
	hash.Write([]byte(s.String()))
	return fmt.Sprintf("sha256:%x", hash.Sum(nil))
}

func (rw Rewriter) EvalIdentifier(ident *ast.Identifier, env Environment) ast.Expression {
	/*
		if in rw.Settings rewrite(&ast.Subscript)
		otherwise just return it
	*/

	if v, ok := rw.fromEnv(ident, env); ok {
		if CanLiteralfy(v) {
			return Literalify(v)
		}
		return ident
	}

	name := ident.Value
	if setting, ok := rw.Settings[name]; ok {
		env.Set(name, Box(setting))
		return Literalify(setting)
	}

	if strings.HasPrefix(name, "logic_") {
		v := rw.Tricks[strings.TrimPrefix(name, "logic_")]
		env.Set(name, Box(v))
		return Literalify(v)
	}

	return ident
}

// lowers to a boolean from a passed settings dict
func (rw Rewriter) EvalSubscript(subscript *ast.Subscript, env Environment) ast.Expression {
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
	default:
		panic(parseError("unknown subscript target %s[%s]", settings.Value, value))
	}
}

func (rw Rewriter) EvalTuple(tup *ast.Tuple, env Environment) ast.Expression {
	if len(tup.Elems) != 2 {
		panic(BadTupleErr)
	}

	item := tup.Elems[0]
	want := tup.Elems[1]
	var ident *ast.Identifier
	var qty float64
	var ok bool

	switch item := item.(type) {
	case *ast.Identifier:
		ident = item
	case *ast.Literal:
		if item.Kind == ast.LiteralStr {
			ident, ok = rw.entityLiteral(ident.Value, env)
			if ok {
				break
			}
		}
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

	return &ast.Call{
		Callee: &ast.Identifier{Value: "has"},
		Args:   []ast.Expression{ident, Literalify(qty)},
	}
}

var BadTupleErr = errors.New("tuple must be (Ident, Number)")
