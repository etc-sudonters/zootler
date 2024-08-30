package worldloader

import (
	"sudonters/zootler/internal/components"
	"sudonters/zootler/internal/rules/parser"
	"sudonters/zootler/internal/rules/preprocessor"
	"sudonters/zootler/internal/rules/runtime"
	"sudonters/zootler/internal/table"

	"github.com/etc-sudonters/substrate/slipup"
)

type LogicCompiler struct {
	P *preprocessor.P
	C *runtime.Compiler
}

func (lc *LogicCompiler) CompileAll(locations *Locations, tokens *Tokens) error {
	var rule components.CompiledRule
	for name, token := range tokens.item {
		lc.P.Env[name] = parser.TokenLiteral(token)
	}

	for _, edge := range locations.edges {
		ast, err := parser.Parse(edge.rule)
		if err != nil {
			return slipup.Describef(err, "while parsing edge %s", edge.name)
		}
		ast, err = lc.P.Process(string(edge.origin.name), ast)
		if err != nil {
			return slipup.Describef(err, "while processing edge %s", edge.name)
		}
		rule.Bytecode, err = runtime.CompileEdgeRule(lc.C, ast)
		err = edge.Attach(table.Values{rule})
		if err != nil {
			return slipup.Describef(err, "while compiling edge %s", edge.name)
		}
	}

	// delayed rules _could_ create more
	for len(lc.P.Delayed) != 0 {
		delayedRules := lc.P.Delayed
		lc.P.Delayed = make(preprocessor.DelayedRules, 0)
		for target, delayeds := range delayedRules {
			origin, err := locations.Build(components.Name(target))
			if err != nil {
				return slipup.Describef(err, "while retrieving %s for delayed edge processing", target)
			}
			for _, delayed := range delayeds {
				ast, err := lc.P.Process(target, delayed.Rule)
				if err != nil {
					return slipup.Describef(err, "while processing delayed edge %s", delayed.Name)
				}
				rule.Bytecode, err = runtime.CompileEdgeRule(lc.C, ast)
				if err != nil {
					return slipup.Describef(err, "while compiling delayed edge %s", delayed.Name)
				}

				dest, err := locations.Build(components.Name(delayed.Name))
				if err != nil {
					return slipup.Describef(err, "while compiling delayed edge %s", delayed.Name)
				}
				edge, err := locations.Connect(origin, dest, "")
				if err != nil {
					return slipup.Describef(err, "while compiling delayed edge %s", delayed.Name)
				}
				err = edge.Attach(table.Values{
					rule,
					components.EventEdge{},
					components.Advancement{},
					components.CollectableGameToken{},
				})
				if err != nil {
					return slipup.Describef(err, "while compiling delayed edge %s", delayed.Name)
				}
			}
		}
	}

	return nil
}
