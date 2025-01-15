package main

import (
	"fmt"
	"sudonters/zootler/mido"
	"sudonters/zootler/mido/ast"
	"sudonters/zootler/mido/code"
	"sudonters/zootler/mido/objects"
	"sudonters/zootler/mido/vm"
)

func ExectuteAll(env *mido.CompileEnv, compiled []mido.CompiledSource) {
	objTable := objects.NewTable(
		objects.BuildTableFrom(env.Objects),
		objects.TableWithBuiltIns(objects.BuiltInFunctions{
			{Name: "has", Params: 2, Fn: constBuiltInFunc},
			{Name: "has_anyof", Params: -1, Fn: constBuiltInFunc},
			{Name: "has_every", Params: -1, Fn: constBuiltInFunc},
			{Name: "is_adult", Params: 0, Fn: constBuiltInFunc},
			{Name: "is_child", Params: 0, Fn: constBuiltInFunc},
			{Name: "has_bottle", Params: 0, Fn: constBuiltInFunc},
			{Name: "has_dungeon_rewards", Params: 1, Fn: constBuiltInFunc},
			{Name: "has_hearts", Params: 1, Fn: constBuiltInFunc},
			{Name: "has_medallions", Params: 1, Fn: constBuiltInFunc},
			{Name: "has_stones", Params: 1, Fn: constBuiltInFunc},
			{Name: "is_starting_age", Params: 0, Fn: constBuiltInFunc},
		}),
	)
	engine := vm.VM{&objTable}

	for i := range compiled {
		compiled := compiled[i]
		result, err := engine.Execute(compiled.ByteCode)

		if err == nil && result != nil {
			continue
		}
		fmt.Println()
		fmt.Printf("%s -> %s\n", compiled.OriginatingRegion, compiled.Destination)
		if compiled.String != "" {
			fmt.Println()
			fmt.Println("raw")
			fmt.Println(compiled.String)
		}
		if compiled.Ast != nil {
			fmt.Println()
			fmt.Println("ast")
			fmt.Println(ast.Render(compiled.Ast))
		}

		fmt.Println()
		fmt.Println(code.DisassembleToString(compiled.ByteCode.Tape))

		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("produced nil result")
		}
		fmt.Println()
	}
}

type constbuiltins bool

func (this constbuiltins) AtDampeTime([]objects.Object) (objects.Object, error) {
	return objects.Boolean(this), nil
}
func (this constbuiltins) AtDay([]objects.Object) (objects.Object, error) {
	return objects.Boolean(this), nil
}
func (this constbuiltins) AtNight([]objects.Object) (objects.Object, error) {
	return objects.Boolean(this), nil
}
func (this constbuiltins) Has([]objects.Object) (objects.Object, error) {
	return objects.Boolean(this), nil
}
func (this constbuiltins) HasAllNotesForSong([]objects.Object) (objects.Object, error) {
	return objects.Boolean(this), nil
}
func (this constbuiltins) HasAnyOf([]objects.Object) (objects.Object, error) {
	return objects.Boolean(this), nil
}
func (this constbuiltins) HasBottle([]objects.Object) (objects.Object, error) {
	return objects.Boolean(this), nil
}
func (this constbuiltins) HasDungeonRewards([]objects.Object) (objects.Object, error) {
	return objects.Boolean(this), nil
}
func (this constbuiltins) HasEvery([]objects.Object) (objects.Object, error) {
	return objects.Boolean(this), nil
}
func (this constbuiltins) HasHearts([]objects.Object) (objects.Object, error) {
	return objects.Boolean(this), nil
}
func (this constbuiltins) HasMedallions([]objects.Object) (objects.Object, error) {
	return objects.Boolean(this), nil
}
func (this constbuiltins) HasStones([]objects.Object) (objects.Object, error) {
	return objects.Boolean(this), nil
}
func (this constbuiltins) IsAdult([]objects.Object) (objects.Object, error) {
	return objects.Boolean(this), nil
}
func (this constbuiltins) IsChild([]objects.Object) (objects.Object, error) {
	return objects.Boolean(this), nil
}
func (this constbuiltins) IsStartingAge([]objects.Object) (objects.Object, error) {
	return objects.Boolean(this), nil
}
func (this constbuiltins) RegionHasShortcuts([]objects.Object) (objects.Object, error) {
	return objects.Boolean(this), nil
}
