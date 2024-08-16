package main

import (
	"context"
	"strings"
	"sudonters/zootler/internal"
	"sudonters/zootler/internal/app"
	"sudonters/zootler/internal/components"
	"sudonters/zootler/internal/query"
	"sudonters/zootler/internal/rules/parser"
	"sudonters/zootler/internal/rules/runtime"
	"sudonters/zootler/internal/slipup"
)

func example(z *app.Zootlr) error {
	env := runtime.NewEnv()
	env.Set("amount", runtime.ValueFromInt(2))
	functions := runtime.NewFuncNamespace()
	functions.AddFunc("has", &HasQtyOf{z.Engine()})
	compiler := runtime.CompilerUsing(functions, env)
	code := "has('Progressive Hookshot', 2)"
	ast, parseErr := parser.Parse(code)
	if parseErr != nil {
		return slipup.Describef(parseErr, "while parsing rule '%s'", code)
	}

	c, compileErr := compiler.CompileEdgeRule(ast)
	if compileErr != nil {
		return slipup.Describef(compileErr, "while compiling rule '%s'", code)
	}

	internal.WriteLineOut(z.Ctx(), c.Disassemble(code))
	vm := runtime.CreateVM(env, functions)
	result, runErr := vm.Run(z.Ctx(), c)
	if runErr == nil {
		internal.WriteLineOut(z.Ctx(), "result:\t%#v", result.Unwrap())
	}
	return runErr
}

type HasQtyOf struct {
	storage query.Engine
}

func (h HasQtyOf) Arity() int {
	return 2
}

func (h HasQtyOf) Run(ctx context.Context, _ *runtime.VM, values runtime.Values) (runtime.Value, error) {
	if len(values) != 2 {
		return runtime.NullValue(), slipup.Createf("expected 2 arguments, received: %d", len(values))
	}

	name, castErr := values[0].AsStr()
	if castErr != nil {
		return runtime.NullValue(), slipup.Describe(castErr, "expected qty to be string")
	}

	qty, castErr := values[1].AsInt()
	if castErr != nil {
		return runtime.NullValue(), slipup.Describe(castErr, "expected qty to be number")
	}

	has, err := h.qtyOf(ctx, name)
	if err != nil {
		return runtime.NullValue(), err
	}

	return runtime.ValueFromBool(has >= qty), nil
}

func (h *HasQtyOf) qtyOf(_ context.Context, needle string) (int, error) {
	eng := h.storage
	q := eng.CreateQuery()
	q.Exists(query.MustAsColumnId[components.CollectableGameToken](eng))
	q.Exists(query.MustAsColumnId[components.Advancement](eng))
	q.Load(query.MustAsColumnId[components.Name](eng))

	haystack, err := h.storage.Retrieve(q)
	if err != nil {
		return -1, slipup.Describe(err, "while looking up collected items")
	}

	qty := 0
	for haystack.MoveNext() {
		name := haystack.Current().Values[0].(components.Name)
		if strings.ToLower(needle) == strings.ToLower(string(name)) {
			qty++
		}
	}

	return qty, nil
}
