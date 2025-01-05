package main

import (
	"fmt"
	"sudonters/zootler/internal/ruleparser"
	"sudonters/zootler/magicbeanvm/ast"
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
	funcTable, funcTableErr := ast.BuildFunctionTable(&symbolTable, grammar, ReadHelpers(".data/logic/helpers.json"))
	if funcTableErr != nil {
		panic(funcTableErr)
	}

	loadingRules, loadRulesErr := loaddir(".data/logic/glitchless")
	if loadRulesErr != nil {
		panic(loadRulesErr)
	}

	ctx := optimizer.NewCtx()
	rw := []ast.Rewriter{
		optimizer.InlineCalls(&ctx, &symbolTable, &funcTable),
		optimizer.FoldConstants(&symbolTable),
		optimizer.EnsureFuncs(&symbolTable, &funcTable),
		optimizer.CollapseHas(&symbolTable),
		optimizer.ExtractLateExpansions(&ctx, &symbolTable),
		optimizer.PromoteTokens(&symbolTable),
		optimizer.RunCompileTimeFuncs(&symbolTable, boolcomptime(ast.Bool(true))),
	}

	findstrings := make(findstrings)
	stringfinder := ast.Visitor{
		Invoke: ast.DontVisit[ast.Invoke](),
		String: findstrings.String,
	}

	// handle forward declaration for events, exits and locations
	for _, rule := range loadingRules {
		switch rule.kind {
		case symbols.TRANSIT, symbols.LOCATION:
			rule.name = fmt.Sprintf("%s -> %s", rule.parent, rule.name)
		}
		symbolTable.Declare(rule.name, rule.kind)
	}

	for _, rule := range loadingRules {
		where, logic := rule.parent, rule.body
		nodes, astErr := ast.Parse(logic, &symbolTable, grammar)
		if astErr != nil {
			panic(astErr)
		}

		ctx.Store(optimizer.CurrentLocationKey, where)

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

		stringfinder.Visit(nodes)
		fmt.Printf("%s -> %s: %s\n", where, rule.name, ast.Render(nodes))
	}

	reexpand(&ctx, &symbolTable, rw)
	fmt.Println()

	var events []string
	for symbol := range symbolTable.All {
		switch symbol.Kind {
		case symbols.EVENT:
			events = append(events, symbol.Name)
		}
	}
	if eventCount := len(events); eventCount != 0 {
		fmt.Printf("Generated %d events:\n", eventCount)
		for _, name := range events {
			fmt.Printf("\t%s\n", name)
		}
	}
}

func reexpand(ctx *optimizer.Context, symbolTable *symbols.Table, rw []ast.Rewriter) {
	expansions := ctx.Swap(optimizer.LateExpansionKey, make(optimizer.LateExpansions))
	if expansions, casted := expansions.(optimizer.LateExpansions); casted {
		for expansions.Size() != 0 {
			fmt.Printf("\n\nFound %04d expansions\n", expansions.Size())
			for _, expand := range expansionsToRules(expansions) {
				var rewriteErr error
				var nodes ast.Node
				token := symbolTable.LookUpByIndex(expand.token)
				where := symbolTable.LookUpByIndex(expand.where)
				ctx.Store(optimizer.CurrentLocationKey, where.Name)
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
			}

			expansions = ctx.Swap(optimizer.LateExpansionKey, make(optimizer.LateExpansions)).(optimizer.LateExpansions)
		}
	}
}

func expansionsToRules(expansions optimizer.LateExpansions) []partialRule {
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

func (ct boolcomptime) Here([]ast.Node) (ast.Node, error) {
	panic("not implemented")
}

func (ct boolcomptime) At([]ast.Node) (ast.Node, error) {
	panic("not implemented")
}
