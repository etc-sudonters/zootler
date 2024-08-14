package main

import (
	"context"
	"strings"
	"sudonters/zootler/internal/app"
	"sudonters/zootler/internal/components"
	"sudonters/zootler/internal/query"
	"sudonters/zootler/internal/rules/parser"
	"sudonters/zootler/internal/rules/runtime"
	"sudonters/zootler/internal/slipup"
)

func example(z *app.Zootlr) error {
	// hardcoded to progressive hookshot right now
	name := "has(amount)"
	code := name
	ast, parseErr := parser.Parse(code)
	if parseErr != nil {
		return slipup.TraceMsg(parseErr, "while parsing rule '%s'", code)
	}

	c, compileErr := runtime.Compile(ast)
	if compileErr != nil {
		return slipup.TraceMsg(parseErr, "while compiling rule '%s'", code)
	}

	env := runtime.NewEnv()
	env.Set("amount", runtime.StackValueOrPanic(2))
	memory := runtime.NewVmMem()
	memory.AddFunc("has", &HasQtyOf{z.Engine()})

	WriteLineOut(z.Ctx(), c.Disassemble(name))

	vm := runtime.CreateVM(env, memory)
	vm.Debug(true)
	result, runErr := vm.Run(z.Ctx(), c)
	if runErr == nil {
		WriteLineOut(z.Ctx(), "result:\t%#v", result.Unwrap())
	}
	WriteLineOut(z.Ctx(), "vm dump:\n%#v", vm)
	return runErr
}

type HasQtyOf struct {
	storage query.Engine
}

func (h HasQtyOf) Arity() int {
	return 1
}

func (h HasQtyOf) Run(ctx context.Context, _ *runtime.VM, values runtime.Values) (runtime.Value, error) {
	if len(values) != 1 {
		return runtime.NullValue(), slipup.Create("expected 1 arguments, received: %d", len(values))
	}

	qty := values[0].Unwrap().(int)

	has, err := h.qtyOf(ctx, "Progressive Hookshot")
	if err != nil {
		return runtime.NullValue(), err
	}

	return runtime.ValueFromBool(has >= qty), nil
}

func (h *HasQtyOf) qtyOf(ctx context.Context, needle string) (int, error) {
	q := h.storage.CreateQuery()
	q.Exists(T[components.CollectableGameToken]())
	q.Exists(T[components.Advancement]())
	q.Load(T[components.Name]())

	haystack, err := h.storage.Retrieve(q)
	if err != nil {
		return -1, slipup.Trace(err, "while looking up collected items")
	}

	qty := 0
	for haystack.MoveNext() {
		name := haystack.Current().Values[0].(components.Name)
		if strings.ToLower(needle) == strings.ToLower(string(name)) {
			qty++
		}
	}

	WriteLineOut(ctx, "Found %d %s", qty, needle)
	return qty, nil
}
