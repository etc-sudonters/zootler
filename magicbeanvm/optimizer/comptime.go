package optimizer

import (
	"sudonters/zootler/magicbeanvm/ast"
	"sudonters/zootler/magicbeanvm/symbols"
)

var compileTimeNames = []string{
	"at",
	"compare_setting",
	"here",
	"is_trick_enabled",
	"load_setting",
	"load_setting_2",
	"had_night_start",
}

func CompileTimeNames() []string {
	return compileTimeNames[:]
}

type CompileTime interface {
	At([]ast.Node) (ast.Node, error)
	CompareSetting([]ast.Node) (ast.Node, error)
	HadNightStart([]ast.Node) (ast.Node, error)
	Here([]ast.Node) (ast.Node, error)
	IsTrickEnabled([]ast.Node) (ast.Node, error)
	LoadSetting([]ast.Node) (ast.Node, error)
	LoadSetting2([]ast.Node) (ast.Node, error)
}

func RunCompileTimeFuncs(symbols *symbols.Table, funcs CompileTime) ast.Rewriter {
	ct := comptime{symbols, funcs}
	return ast.Rewriter{
		Invoke: ct.Invoke,
	}
}

type comptime struct {
	symbols *symbols.Table
	funcs   CompileTime
}

func (this comptime) Invoke(node ast.Invoke, _ ast.Rewriting) (ast.Node, error) {
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
	default:
		return node, nil
	}
}
