package vm

import "sudonters/zootler/magicbeanvm/objects"

func GlobalNames() []string {
	return globalNames[:]
}

var globalNames = []string{
	"Fire",
	"Forest",
	"Light",
	"Shadow",
	"Spirit",
	"Water",
	"adult",
	"age",
	"both",
	"either",
	"child",
}

type Runner struct {
	builtins objects.BuiltInFunctionTable
}

func New(builtins objects.BuiltInFunctionTable) Runner {
	var r Runner
	r.builtins = builtins
	return r
}
