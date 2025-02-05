package optimizer

import (
	"sudonters/zootler/internal"
	"sudonters/zootler/internal/settings"
	"sudonters/zootler/mido/ast"
)

type setting func(*settings.Zootr, string) ast.Node

func StrSetting(these *settings.Zootr, name string) ast.Node {
	value, err := these.String(name)
	internal.PanicOnError(err)
	return ast.String(value)
}

func Float64Setting(these *settings.Zootr, name string) ast.Node {
	value, err := these.Float64(name)
	internal.PanicOnError(err)
	return ast.Number(value)
}

func BoolSetting(these *settings.Zootr, name string) ast.Node {
	value, err := these.Bool(name)
	internal.PanicOnError(err)
	return ast.Boolean(value)
}

type settinginline struct {
	these *settings.Zootr
	fs    map[string]setting
}
