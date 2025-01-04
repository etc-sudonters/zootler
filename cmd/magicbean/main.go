package main

import (
	"fmt"
	"sudonters/zootler/internal/ruleparser"
	"sudonters/zootler/magicbeanvm/ast"
	"sudonters/zootler/magicbeanvm/optimizer"
	"sudonters/zootler/magicbeanvm/symbols"
)

type rule struct{ where, logic string }
type partialRule struct {
	where, token symbols.Index
	body         ast.Node
}

var rules = []rule{
	{"nested-and", "item and adult and dance"},
	{"nested-or", "item or adult or dance"},
	{"late-expand-at", "at('Forest Temple Outside Upper Ledge', True)"},
	{"late-expand-here", "here(logic_forest_mq_hallway_switch_boomerang and can_use(Boomerang))"},
	{"is-trick-enabled", "logic_forest_mq_hallway_switch_jumpslash"},
	{"call-func", "can_use(Hover_Boots)"},
	{"true", "True"},
	{"true-or", "True or at('Forest Temple Outside Upper Ledge', False)"},
	{"true-and", "True and at('Forest Temple Outside Upper Ledge', True)"},
	{"or-true", "at('Forest Temple Outside Upper Ledge', False) or True)"},
	{"and-true", "at('Forest Temple Outside Upper Ledge', True) and True"},
	{"false", "False"},
	{"false-or", "False or at('Forest Temple Outside Upper Ledge', False)"},
	{"false-and", "False and at('Forest Temple Outside Upper Ledge', False)"},
	{"or-false", "at('Forest Temple Outside Upper Ledge', False) or False)"},
	{"and-false", "at('Forest Temple Outside Upper Ledge', False) and False"},
	{"compare-eq-same-is-true", "same == same"},
	{"compare-nq-diff-is-true", "same != diff"},
	{"compare-eq-diff-is-false", "same == diff"},
	{"compare-nq-same-is-false", "same != same"},
	{"compare-lt", "chicken_count < 7"},
	{"compare-setting", "deadly_bonks == 'ohko'"},
	{"uses-setting", "('Triforce Piece', victory_goal_count)"},
	{"contains", "'Deku Tree' in dungeon_shortcuts"},
	{"subscript", "skipped_trials[Forest]"},
	{"float", "can_live_dmg(0.5)"},
	{"promote-standalone-token", "Progressive_Hookshot"},
	{"has-all", "has(taco, 1) and has(burrito, 1) and has(taquito, 1)"},
	{"has-any", "has(taco, 1) or has(burrito, 1) or has(taquito, 1)"},
	{"has-all-mix", "has(taco, 2) and has(burrito, 1) and has(taquito, 1) and is_adult"},
	{"has-any-mix", "has(taco, 2) or has(burrito, 1) or has(taquito, 1) or is_child"},
	{"call-helper", "can_use(Dins_Fire)"},
	{"can-use-hookshot", "can_use(Hookshot)"},
	{"can-use-goron-tunic", "can_use(Goron_Tunic)"},
	{"promote-standalone-func", "is_adult"},
	{"implicit-has", "(Spirit_Temple_Small_Key, 15)"},
	{"really-implicit-has", "Dins_Fire"},
	{"really-really-implicit-has", "'Goron Tunic'"},
	{"goron-tunic", "is_adult and ('Goron Tunic' or Buy_Goron_Tunic)"},
	{"goron-tunic", "is_adult or ('Goron Tunic' or Buy_Goron_Tunic)"},
	{"subscripts", "(skipped_trials[Forest] or 'Forest Trial Clear') and (skipped_trials[Fire] or 'Fire Trial Clear') and (skipped_trials[Water] or 'Water Trial Clear') and (skipped_trials[Shadow] or 'Shadow Trial Clear') and (skipped_trials[Spirit] or 'Spirit Trial Clear') and (skipped_trials[Light] or 'Light Trial Clear')"},
	{"logic_rules", "logic_rules == 'glitched'"},
	{"recursive-macro", "here(at('dance hall', dance))"},
}

func main() {
	rawTokens, tokenErr := loadTokensNames(".data/data/items.json")
	if tokenErr != nil {
		panic(tokenErr)
	}

	symbolTable := symbols.NewTable()
	symbolTable.DeclareMany(symbols.COMP_TIME, compTime)
	symbolTable.DeclareMany(symbols.BUILT_IN, builtIns)
	symbolTable.DeclareMany(symbols.GLOBAL, globals)
	symbolTable.DeclareMany(symbols.SETTING, settings)
	symbolTable.DeclareMany(symbols.TOKEN, tokens)
	symbolTable.DeclareMany(symbols.TOKEN, rawTokens)

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
		optimizer.InlineCalls(&ctx, &symbolTable, funcTable),
		optimizer.FoldConstants(&symbolTable),
		optimizer.EnsureFuncs(&symbolTable, funcTable),
		optimizer.PromoteTokens(&symbolTable),
		optimizer.CollapseHas(&symbolTable),
		optimizer.ExtractLateExpansions(&ctx, &symbolTable),
	}

	for _, rule := range loadingRules {
		where, logic := rule.parent, rule.body
		nodes, astErr := ast.Parse(logic, &symbolTable, grammar)
		ctx.Store(optimizer.CurrentLocationKey, where)
		if astErr != nil {
			panic(astErr)
		}

		for range 5 {
			var rewriteErr error
			nodes, rewriteErr = ast.RewriteWithEvery(nodes, rw)
			if rewriteErr != nil {
				panic(rewriteErr)
			}
		}

		fmt.Printf("%s -> %s: %s\n", where, rule.name, ast.Render(nodes))
	}

	var expansions optimizer.LateExpansions
	expansions = ctx.Swap(optimizer.LateExpansionKey, make(optimizer.LateExpansions)).(optimizer.LateExpansions)
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
