package main

import (
	"fmt"
	"sudonters/zootler/internal/ruleparser"
	"sudonters/zootler/magicbeanvm"
	"sudonters/zootler/magicbeanvm/ast"
	"sudonters/zootler/magicbeanvm/code"
	"sudonters/zootler/magicbeanvm/compiler"
	"sudonters/zootler/magicbeanvm/optimizer"
	"sudonters/zootler/magicbeanvm/symbols"
	"sudonters/zootler/magicbeanvm/vm"
)

func main() {
	rawTokens, tokenErr := loadTokensNames(".data/data/items.json")
	if tokenErr != nil {
		panic(tokenErr)
	}

	symbolTable := symbols.NewTable()
	symbolTable.DeclareMany(symbols.COMP_TIME, optimizer.CompileTimeNames())
	symbolTable.DeclareMany(symbols.BUILT_IN, vm.BuiltInFunctionNames())
	symbolTable.DeclareMany(symbols.GLOBAL, vm.GlobalNames())
	symbolTable.DeclareMany(symbols.SETTING, settings)
	symbolTable.DeclareMany(symbols.TOKEN, rawTokens)
	aliasTokens(&symbolTable, rawTokens)

	grammar := ruleparser.NewRulesGrammar()
	funcTable, funcTableErr := ast.BuildCompilingFunctionTable(&symbolTable, grammar, ReadHelpers(".data/logic/helpers.json"))
	if funcTableErr != nil {
		panic(funcTableErr)
	}

	loadingRules, loadRulesErr := loaddir(".data/logic/glitchless")
	if loadRulesErr != nil {
		panic(loadRulesErr)
	}

	ctx := optimizer.NewCtx()
	comptime := optimizer.RunCompileTimeFuncs(
		&symbolTable,
		comptime{
			boolcomptime(ast.Bool(true)),
			magicbeanvm.ExtractLateExpansions(&ctx, &symbolTable),
		})
	_ = comptime

	rw := []ast.Rewriter{
		optimizer.InlineCalls(&ctx, &symbolTable, &funcTable),
		optimizer.FoldConstants(&symbolTable),
		optimizer.EnsureFuncs(&symbolTable, &funcTable),
		optimizer.CollapseHas(&symbolTable),
		optimizer.PromoteTokens(&symbolTable),
		comptime,
	}

	nodecounter := make(nodecounter)
	invokeCounter := findinvokes{&symbolTable, make(map[string]symcount), make(map[string]int)}
	countInvokes := ast.Visitor{Invoke: invokeCounter.Invoke}
	countnodes := countnodes(nodecounter)

	// handle forward declaration for events, exits and locations
	for _, rule := range loadingRules {
		switch rule.kind {
		case symbols.TRANSIT, symbols.LOCATION:
			rule.name = fmt.Sprintf("%s -> %s", rule.parent, rule.name)
		}
		symbol := symbolTable.Declare(rule.name, rule.kind)

		if rule.kind == symbols.EVENT {
			symbolTable.Alias(symbol, escape(symbol.Name))
		}
	}

	for _, rule := range loadingRules {
		where, logic := rule.parent, rule.body
		nodes, astErr := ast.Parse(logic, &symbolTable, grammar)
		if astErr != nil {
			panic(astErr)
		}

		magicbeanvm.SetCurrentLocation(&ctx, where)

		for range 5 {
			var rewriteErr error
			nodes, rewriteErr = ast.RewriteWithEvery(nodes, rw)
			if nodes == nil {
				panic(fmt.Errorf("%s -> %s produced nil rule", where, rule.name))
			}
			if rewriteErr != nil {
				panic(rewriteErr)
			}
		}

		fmt.Printf("%s -> %s: %s\n", where, rule.name, ast.Render(nodes))
		countInvokes.Visit(nodes)
		countnodes.Visit(nodes)
	}

	ctx.Store("invokeCounter", &countInvokes)
	ctx.Store("countnodes", &countnodes)
	reexpand(&ctx, &symbolTable, rw)
	fmt.Println()

	fmt.Println("INVOKE TOTALS")
	for name, item := range invokeCounter.counting {
		fmt.Printf("%06d\t%s\t\t%s\n", item.count, item.kind, name)
	}
	fmt.Println()

	fmt.Println("NODE TOTALS")
	for kind, count := range nodecounter {
		fmt.Printf("%06d\t%s\n", count, kind)
	}
	fmt.Println()
	size, total, aliased := symbolTable.Size(), symbolTable.RawSize(), symbolTable.AliasCount()
	fmt.Printf("ALIAS: %04d %04X\n", aliased, aliased)
	fmt.Printf("COUNT: %04d %04X\n", size, size)
	fmt.Printf("TOTAL: %04d %04X\n", total, total)
	fmt.Println()

	var tape compiler.Tape
	tape.Write(code.Make(code.BEAN_NOP))
	tape.Write(code.Make(code.BEAN_CHK_QTY, 0x7FF7, 0xAB))
	tape.Write(code.Make(code.BEAN_PUSH_CONST, 0xBEEF))
	tape.Write(code.Make(code.BEAN_PUSH_PTR, 0xDEAD))
	tape.Write(code.Make(code.BEAN_PUSH_FUNC, 0xCAFE))
	tape.Write(code.Make(code.BEAN_CALL, 3))
	tape.Write(code.Make(code.BEAN_NEED_ALL, 2))
	tape.Write(code.Make(code.BEAN_PUSH_T))
	tape.Write(code.Make(code.BEAN_PUSH_F))
	tape.Write(code.Make(code.BEAN_NEED_ANY, 3))
	tape.Write(code.Make(code.BEAN_ERR))
	fmt.Print(code.Disassemble(tape.Read()))
}

func reexpand(ctx *optimizer.Context, symbolTable *symbols.Table, rw []ast.Rewriter) {
	countInvokes := ctx.Retrieve("invokeCounter").(*ast.Visitor)
	countnodes := ctx.Retrieve("countnodes").(*ast.Visitor)
	expansions := magicbeanvm.SwapLateExpansions(ctx)
	for expansions.Size() != 0 {
		fmt.Printf("\n\nFound %04d expansions\n", expansions.Size())
		for _, expand := range expansionsToRules(expansions) {
			var rewriteErr error
			var nodes ast.Node
			token := symbolTable.LookUpByIndex(expand.token)
			where := symbolTable.LookUpByIndex(expand.where)
			magicbeanvm.SetCurrentLocation(ctx, where.Name)
			nodes = expand.body
			for range 5 {
				nodes, rewriteErr = ast.RewriteWithEvery(nodes, rw)
				if nodes == nil {
					panic(fmt.Errorf("%q produced nil rule", token.Name))
				}
				if rewriteErr != nil {
					panic(rewriteErr)
				}
			}
			fmt.Printf("%s %s\n", token.Name, ast.Render(nodes))
			countInvokes.Visit(nodes)
			countnodes.Visit(nodes)
		}

		expansions = magicbeanvm.SwapLateExpansions(ctx)
	}
}

func expansionsToRules(expansions magicbeanvm.LateExpansions) []partialRule {
	rules := make([]partialRule, 0, expansions.Size())

	for _, partials := range expansions {
		for _, partial := range partials {
			rules = append(rules, partialRule{
				where: partial.AttachedTo,
				token: partial.Token,
				body:  partial.Rule,
			})
		}
	}

	return rules
}

type rest interface {
	At([]ast.Node) (ast.Node, error)
	Here([]ast.Node) (ast.Node, error)
}

type comptime struct {
	boolcomptime
	rest
}

type boolcomptime ast.Bool

func (ct boolcomptime) LoadSetting([]ast.Node) (ast.Node, error) {
	return ast.Bool(ct), nil
}

func (ct boolcomptime) LoadSetting2([]ast.Node) (ast.Node, error) {
	return ast.Bool(ct), nil
}

func (ct boolcomptime) CompareSetting([]ast.Node) (ast.Node, error) {
	return ast.Bool(ct), nil
}

func (ct boolcomptime) IsTrickEnabled([]ast.Node) (ast.Node, error) {
	return ast.Bool(ct), nil
}

func (ct boolcomptime) HadNightStart([]ast.Node) (ast.Node, error) {
	return ast.Bool(ct), nil
}
