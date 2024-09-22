package saburo

import (
	"slices"
	"strings"
	"sudonters/zootler/icearrow/analysis"
	"sudonters/zootler/icearrow/ast"
	parsing "sudonters/zootler/icearrow/parser"
	"sudonters/zootler/icearrow/zasm"
	"sudonters/zootler/internal"
	"sudonters/zootler/internal/app"
	"sudonters/zootler/internal/components"
	"sudonters/zootler/internal/entities"

	"github.com/etc-sudonters/substrate/peruse"
	"github.com/etc-sudonters/substrate/skelly/graph"
	"github.com/etc-sudonters/substrate/slipup"
)

type RuleAssembler struct {
	ScriptPath string
}

func (rc RuleAssembler) Setup(z *app.Zootlr) error {
	assembler := rc.createAssembler()
	locations := app.GetResource[entities.Locations](z).Res
	edges := app.GetResource[entities.Edges](z)
	collected := slices.Collect(edges.Res.All)
	slices.SortFunc(collected, func(a, b entities.Edge) int {
		return strings.Compare(string(a.Name()), string(b.Name()))
	})

	var edge entities.Edge

	whileHandlingRule := func(err error, action string) error {
		return slipup.Describef(err, "while %s rule %q", action, edge.GetRawRule())
	}
	grammar := parsing.NewRulesGrammar()
	analysisCtx := analysis.NewAnalysis(edges.Res)
	loadExpansions(&analysisCtx, rc.ScriptPath)
	loadIdentifierTypes(&analysisCtx, app.GetResource[entities.Tokens](z).Res)
	assemblies := zasm.Assembly{}
	for _, edge = range collected {
		origin, dest := edge.Retrieve("origin").(string), edge.Retrieve("dest").(string)
		if dest == "" || origin == "" {
			panic(slipup.Createf("edge %+v does not have dest/origin", edge))
		}
		analysisCtx.SetCurrent(origin)
		parser := peruse.NewParser(grammar, parsing.NewRulesLexer(string(edge.GetRawRule())))
		pt, ptErr := parser.ParseAt(parsing.LOWEST)
		if ptErr != nil {
			return whileHandlingRule(ptErr, "parsing")
		}
		astNodes, astErr := parsing.Transform(&ast.Ast{}, pt)
		if astErr != nil {
			return whileHandlingRule(astErr, "lowering tree")
		}
		astNodes, astErr = analysis.Analyze(astNodes, &analysisCtx)
		if astErr != nil {
			return whileHandlingRule(astErr, "lowering tree")
		}
		asm, asmErr := assembler.Assemble(string(edge.Name()), astNodes)
		if asmErr != nil {
			return whileHandlingRule(asmErr, "assembling")
		}
		assemblies.Include(asm)

	}

	for current, lateXpn := range (&analysisCtx).LateExpanders {
		edge = lateXpn.Edge
		analysisCtx.SetCurrent(current)
		astNodes, astErr := analysis.Analyze(lateXpn.Rule, &analysisCtx)
		if astErr != nil {
			return whileHandlingRule(astErr, "analyzing late expansion")
		}
		asm, asmErr := assembler.Assemble(string(edge.Name()), astNodes)
		if asmErr != nil {
			return whileHandlingRule(asmErr, "assembling")
		}
		assemblies.Include(asm)
		origin, _ := locations.Entity(components.Name(current))
		dest, _ := locations.Entity(components.Name(edge.Retrieve("dest").(string)))
		edge.Stash("originId", graph.Origination(origin.Id()))
		edge.Stash("destId", graph.Destination(dest.Id()))
	}

	tbls := assembler.CreateDataTables()
	assemblies.AttachDataTables(tbls)
	z.AddResource(assemblies)

	return nil
}

func (rc RuleAssembler) createAssembler() zasm.Assembler {
	return zasm.Assembler{
		Data: zasm.NewDataBuilder(),
	}
}

func loadExpansions(ctx *analysis.AnalysisContext, path string) {
	contents, err := internal.ReadJsonFileStringMap(path)
	if err != nil {
		panic(err)
	}
	grammar := parsing.NewRulesGrammar()
	for decl, body := range contents {
		name, params := decl, []string(nil)
		if strings.Contains(name, "(") {
			name, params = parseDecl(name)
		}
		parser := peruse.NewParser(grammar, parsing.NewRulesLexer(body))
		pt, _ := parser.ParseAt(parsing.LOWEST)
		ast, _ := parsing.Transform(&ast.Ast{}, pt)
		ctx.AddExpansion(name, params, ast)
	}
}

func loadIdentifierTypes(ac *analysis.AnalysisContext, tokens entities.Tokens) {
	for token := range tokens.All {
		ac.NameToken(string(token.Name()))
	}

	settingNames := []string{
		"logic_rules",
		"adult_trade_shuffle",
		"big_poe_count",
		"bridge",
		"complete_mask_quest",
		"damage_multiplier",
		"deadly_bonks",
		"disable_trade_revert",
		"dungeon_shortcuts",
		"entrance_shuffle",
		"fix_broken_drops",
		"free_bombchu_drops",
		"free_scarecrow",
		"ganon_bosskey_hearts",
		"ganon_bosskey_medallions",
		"ganon_bosskey_rewards",
		"ganon_bosskey_stones",
		"ganon_bosskey_tokens",
		"gerudo_fortress",
		"lacs_condition",
		"lacs_hearts",
		"lacs_medallions",
		"open_door_of_time",
		"open_forest",
		"open_kakariko",
		"plant_beans",
		"chicken_count",
		"selected_adult_trade_item",
		"shuffle_dungeon_entrances",
		"shuffle_empty_pots",
		"shuffle_expensive_merchants",
		"shuffle_ganon_bosskey",
		"shuffle_individual_ocarina_notes",
		"shuffle_interior_entrances",
		"shuffle_overworld_entrances",
		"shuffle_pots",
		"shuffle_scrubs",
		"shuffle_silver_rupees",
		"shuffle_tcgkeys",
		"skip_child_zelda",
		"skip_reward_from_rauru",
		"skipped_trials",
		"triforce_goal_per_world",
		"warp_songs",
		"zora_fountain",

		"bridge_hearts",
		"bridge_medallions",
		"bridge_rewards",
		"bridge_stones",
		"bridge_tokens",
		"ganon_bosskey_tokens_hearts",
		"ganon_bosskey_tokens_medallions",
		"ganon_bosskey_tokens_rewards",
		"ganon_bosskey_tokens_stones",
		"ganon_bosskey_tokens_tokens",
		"lacs_hearts",
		"lacs_medallions",
		"lacs_rewards",
		"lacs_stones",
		"lacs_tokens",
		"starting_age",
	}

	for _, setting := range settingNames {
		ac.NameSetting(setting)
	}

	builtIns := []string{
		"load_setting",
		"has_dungeon_shortcuts",
		"is_trial_skipped",
		"at",
		"at_dampe_time",
		"at_day",
		"at_night",
		"had_night_start",
		"has_bottle",
		"has_hearts",
		"has_stones",
		"here",
	}
	for _, builtIn := range builtIns {
		ac.NameBuiltIn(builtIn)
	}
}

func parseDecl(raw string) (string, []string) {
	parts := strings.Split(strings.TrimSuffix(raw, ")"), "(")
	args := strings.Split(parts[1], ",")
	return parts[0], args
}
