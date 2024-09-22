package main

import (
	"sudonters/zootler/carpenters/shiro"
	"sudonters/zootler/icearrow/compiler"
	"sudonters/zootler/icearrow/runtime"
	"sudonters/zootler/internal/app"

	"github.com/etc-sudonters/substrate/dontio"
)

type FakeVMState struct{}

func (_ FakeVMState) HasQty(uint32, uint8) bool { return true }
func (_ FakeVMState) HasAny(...uint32) bool     { return true }
func (_ FakeVMState) HasAll(...uint32) bool     { return true }
func (_ FakeVMState) HasBottle() bool           { return true }
func (_ FakeVMState) IsAdult() bool             { return true }
func (_ FakeVMState) IsChild() bool             { return true }
func (_ FakeVMState) AtTod(uint8) bool          { return true }

func RunVM(z *app.Zootlr) error {
	program := app.GetResource[shiro.CompiledWorldRules](z)
	state := FakeVMState{}
	vm := runtime.VM{}
	symbols := &program.Res.Symbols

	for name, rule := range program.Res.Rules {
		dontio.WriteLineOut(z.Ctx(), "running %q", name)
		dontio.WriteLineOut(z.Ctx(), compiler.ReadTape(&rule, symbols))
		vm.Execute(&rule, state, symbols)
	}
	panic("aaahhh!!!!")
	//return nil
}
