package optimizer

import (
	"sudonters/zootler/magicbeanvm/ast"
	"sudonters/zootler/magicbeanvm/symbols"
)

var compilerFuncNames = []string{
	"at",
	"compare_setting",
	"here",
	"is_trick_enabled",
	"load_setting",
	"load_setting_2",
	"had_night_start",
	"region_has_shortcuts",
}

func CompilerFuncNames() []string {
	return compilerFuncNames[:]
}

type CompilerFunctions interface {
	At([]ast.Node) (ast.Node, error)
	CompareSetting([]ast.Node) (ast.Node, error)
	HadNightStart([]ast.Node) (ast.Node, error)
	Here([]ast.Node) (ast.Node, error)
	IsTrickEnabled([]ast.Node) (ast.Node, error)
	LoadSetting([]ast.Node) (ast.Node, error)
	LoadSetting2([]ast.Node) (ast.Node, error)
	RegionHasShortcuts([]ast.Node) (ast.Node, error)
}

func RunCompilerFunctions(symbols *symbols.Table, funcs CompilerFunctions) ast.Rewriter {
	ct := compilerfuncs{symbols, funcs}
	return ast.Rewriter{
		Invoke: ct.Invoke,
	}
}

type compilerfuncs struct {
	symbols *symbols.Table
	funcs   CompilerFunctions
}

func (this compilerfuncs) Invoke(node ast.Invoke, _ ast.Rewriting) (ast.Node, error) {
	symbol := ast.LookUpNodeInTable(this.symbols, node.Target)
	if symbol == nil {
		return node, nil
	}

	switch symbol.Name {
	case "at":
		return this.funcs.At(node.Args)
	case "compare_setting":
		return this.funcs.CompareSetting(node.Args)
	case "had_night_start":
		return this.funcs.HadNightStart(node.Args)
	case "here":
		return this.funcs.Here(node.Args)
	case "is_trick_enabled":
		return this.funcs.IsTrickEnabled(node.Args)
	case "load_setting":
		return this.funcs.LoadSetting(node.Args)
	case "load_setting_2":
		return this.funcs.LoadSetting2(node.Args)
	case "region_has_shortcuts":
		return this.funcs.RegionHasShortcuts(node.Args)
	default:
		return node, nil
	}
}
