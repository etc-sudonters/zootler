package main

import (
	"fmt"
	"sudonters/zootler/internal/settings"
	"sudonters/zootler/mido"
	"sudonters/zootler/mido/ast"
	"sudonters/zootler/mido/symbols"
)

func wrongargcount(name string, expected, actual int) error {
	return fmt.Errorf("%q expects %d arguments, received %d", name, expected, actual)
}

func wrongargtype(name string, pos int, expected, actual ast.Kind) error {
	return fmt.Errorf("%q argument %d expected to be %q but was %q", name, pos, expected, actual)
}

type settingCompilerFuncs struct {
	settings *settings.Zootr
	symbols  *symbols.Table
}

func (this settingCompilerFuncs) IsTrickEnabled(nodes []ast.Node) (ast.Node, error) {
	if len(nodes) != 1 {
		return nil, wrongargcount("is_trick_enabled", 1, len(nodes))
	}

	node := nodes[0]
	name, isStr := node.(ast.String)
	if !isStr {
		return nil, wrongargtype("is_trick_enabled", 0, ast.KindString, node.Kind())
	}

	enabled := this.settings.Tricks.Enabled[string(name)]
	return ast.Boolean(enabled), nil
}

func (this settingCompilerFuncs) HadNightStart(nodes []ast.Node) (ast.Node, error) {
	if len(nodes) != 0 {
		return nil, wrongargcount("had_night_start", 0, len(nodes))
	}
	switch this.settings.Starting.TimeOfDay {
	case settings.StartingTimeOfDaySunset, settings.StartingTimeOfDayEvening, settings.StartingTimeOfDayMidnight, settings.StartingTimeOfDayWitching:
		return ast.Boolean(true), nil
	default:
		return ast.Boolean(false), nil
	}
}

func (this settingCompilerFuncs) loadsetting(nodes []ast.Node) (ast.Node, error) {
	if len(nodes) != 1 {
		return nil, wrongargcount("load_setting", 1, len(nodes))
	}
	node := nodes[0]
	symbol := ast.LookUpNodeInTable(this.symbols, node)
	if symbol == nil {
		return nil, fmt.Errorf("could not derefence setting: %#v", node)
	}
	panic("not implemented")
}

type compfuncs struct {
	constCompileFuncs
	settingCompilerFuncs
	mido.ConnectionGenerator
}

type constCompileFuncs ast.Boolean

func (ct constCompileFuncs) LoadSetting([]ast.Node) (ast.Node, error) {
	return ast.Boolean(ct), nil
}

func (ct constCompileFuncs) LoadSetting2([]ast.Node) (ast.Node, error) {
	return ast.Boolean(ct), nil
}

func (ct constCompileFuncs) CompareSetting([]ast.Node) (ast.Node, error) {
	return ast.Boolean(ct), nil
}

func (this constCompileFuncs) RegionHasShortcuts([]ast.Node) (ast.Node, error) {
	return ast.Boolean(this), nil
}
